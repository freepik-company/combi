package libconfig

import (
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	commentsRegex                    = `([#][^\n]*)|(\/\/[^\n]*)|(\/\*.*[\n]\*\/)`
	escapeCharsRegex                 = `((=|:)|(;|,)|({|})|(\[|\])|(\(|\)))`
	settingNameRegex                 = `[A-Za-z*][-A-Za-z0-9_*]*`
	settingValuePrimitiveStringRegex = `(\"([^\"\\]|\\.)*\")`
	settingValuePrimitiveFloatRegex  = `(([-+]?([0-9]*)?\.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(\.[0-9]*)?[eE][-+]?[0-9]+))`
	settingValuePrimitiveHexRegex    = `(0[Xx][0-9A-Fa-f]+(L{1,2})?)`
	settingValuePrimitiveIntRegex    = `([-+]?[0-9]+(L{1,2})?)`
	settingValuePrimitiveRegex       = `(` +
		settingValuePrimitiveStringRegex + `|` +
		settingValuePrimitiveFloatRegex + `|` +
		settingValuePrimitiveHexRegex + `|` +
		settingValuePrimitiveIntRegex + `)`
)

// ----------------------------------------------------------------
// LIBCONFIG data structure
// ----------------------------------------------------------------

type LIBCONFIG struct {
	Settings []*SettingT `@@*`
}

type SettingT struct {
	SetingName   string         `@Name ("="|":")`
	SettingValue *SettingValueT `@@`
}

type SettingValueT struct {
	Primitive *PrimitiveT `( @@ (";"?","?)`
	Group     *GroupT     ` | @@ (","?)`
	Array     *ArrayT     ` | @@ (","?)`
	List      *ListT      ` | @@ (","?))`
}

type PrimitiveT struct {
	Value string `@Value (","?)`
}

type ArrayT struct {
	Primitives []*PrimitiveT `"[" @@* "]"`
}

type GroupT struct {
	Settings []*SettingT `"{" @@* "}"`
}

type ListT struct {
	SettingValues []*SettingValueT `"(" @@* ")"`
}

// ----------------------------------------------------------------
// Decode/Encode LIBCONFIG data structure
// ----------------------------------------------------------------

// Decode functions

func DecodeConfig(filepath string) (libconfig *LIBCONFIG, err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return libconfig, err
	}

	libconfig, err = DecodeConfigBytes(configBytes)
	return libconfig, err
}

func DecodeConfigBytes(configBytes []byte) (libconfig *LIBCONFIG, err error) {
	configLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Name`, Pattern: settingNameRegex},
		{Name: `Value`, Pattern: settingValuePrimitiveRegex},
		{Name: "EscapeChars", Pattern: escapeCharsRegex},
		{Name: "Comments", Pattern: commentsRegex},
		{Name: "whitespace", Pattern: `(\s+)`},
	})
	configParser := participle.MustBuild[LIBCONFIG](
		participle.Lexer(configLexer),
	)

	libconfig, err = configParser.ParseBytes("", configBytes)
	return libconfig, err
}

// Encode functions

func EncodeConfigString(config *LIBCONFIG) (configStr string) {
	configStr += encodeConfigSettingString(config.Settings, 0)
	return configStr
}

func encodeConfigSettingString(settings []*SettingT, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	// Encode settings with primitive
	for _, setting := range settings {
		if setting.SettingValue.Primitive != nil {
			configStr += indentStr + setting.SetingName + "=" + setting.SettingValue.Primitive.Value + ",\n"
		}
	}

	// Encode settings with Array
	for _, setting := range settings {
		if setting.SettingValue.Array != nil {
			configStr += setting.SetingName + "=" + "\n"
			configStr += encodeConfigArrayString(setting.SettingValue.Array, indent)
			configStr += ",\n"
		}
	}

	// Encode settings with Group
	for _, setting := range settings {
		if setting.SettingValue.Group != nil {
			configStr += setting.SetingName + "=" + "\n"
			configStr += encodeConfigGroupString(setting.SettingValue.Group, indent)
			configStr += ",\n"
		}
	}

	// Encode settings with List
	for _, setting := range settings {
		if setting.SettingValue.List != nil {
			configStr += setting.SetingName + "=" + "\n"
			configStr += encodeConfigListString(setting.SettingValue.List, indent)
			configStr += ",\n"
		}
	}

	return configStr
}

func encodeConfigArrayString(array *ArrayT, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "[\n" + indentStr + "  "
	for _, primitive := range array.Primitives {
		configStr += primitive.Value + ", "
	}
	configStr += "\n" + indentStr + "]"
	return configStr
}

func encodeConfigGroupString(group *GroupT, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "{\n"
	configStr += encodeConfigSettingString(group.Settings, indent+1)
	configStr += indentStr + "}"
	return configStr
}

func encodeConfigListString(list *ListT, indent int) (configStr string) {
	var indentStr string
	for i := 0; i < indent; i++ {
		indentStr += "  "
	}

	configStr += indentStr + "(\n"
	for _, settingValue := range list.SettingValues {
		if settingValue.Primitive != nil {
			configStr += "  " + settingValue.Primitive.Value + ",\n"
		}

		if settingValue.Array != nil {
			configStr += encodeConfigArrayString(settingValue.Array, indent+1)
			configStr += ",\n"
		}

		if settingValue.Group != nil {
			configStr += encodeConfigGroupString(settingValue.Group, indent+1)
			configStr += ",\n"
		}

		if settingValue.List != nil {
			configStr += encodeConfigListString(settingValue.List, indent+1)
			configStr += ",\n"
		}
	}
	configStr += indentStr + ")"
	return configStr
}

// ----------------------------------------------------------------
// Merge LIBCONFIG data structure
// ----------------------------------------------------------------

func MergeConfigs(destination *LIBCONFIG, source *LIBCONFIG) {
	mergeSettings(&destination.Settings, &source.Settings)
}

func mergeSettings(destination *[]*SettingT, source *[]*SettingT) {
	for _, sSetting := range *source {
		// Merge settings with primitive
		foundSetting := false
		if sSetting.SettingValue.Primitive != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					dSetting.SettingValue.Primitive.Value = sSetting.SettingValue.Primitive.Value
				}
			}
		}

		// Merge settings with group
		if sSetting.SettingValue.Group != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueGroups(dSetting.SettingValue.Group, sSetting.SettingValue.Group)
				}
			}
		}

		// Merge settings with array
		if sSetting.SettingValue.Array != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueArrays(dSetting.SettingValue.Array, sSetting.SettingValue.Array)
				}
			}
		}

		// Merge settings with list
		if sSetting.SettingValue.List != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueLists(dSetting.SettingValue.List, sSetting.SettingValue.List)
				}
			}
		}

		// Append not found setting
		if !foundSetting {
			*destination = append(*destination, sSetting)
		}
	}
}

func mergeSettingValueArrays(destination *ArrayT, source *ArrayT) {
	if len(source.Primitives) > 0 {
		for _, sPrimitive := range source.Primitives {
			found := false
			for _, dPrimitive := range destination.Primitives {
				if dPrimitive.Value == sPrimitive.Value {
					found = true
				}
			}
			if !found {
				destination.Primitives = append(destination.Primitives, sPrimitive)
			}
		}
	}
}

func mergeSettingValueGroups(destination *GroupT, source *GroupT) {
	if len(source.Settings) > 0 {
		mergeSettings(&destination.Settings, &source.Settings)
	}
}

func mergeSettingValueLists(destination *ListT, source *ListT) {
	if len(source.SettingValues) > 0 {
		for _, sSettingValue := range source.SettingValues {
			found := false
			if sSettingValue.Primitive != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Primitive != nil && sSettingValue.Primitive.Value != dSettingValue.Primitive.Value {
						found = true
						dSettingValue.Primitive.Value = sSettingValue.Primitive.Value
					}
				}
			}

			if sSettingValue.Array != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Array != nil {
						found = true
						mergeSettingValueArrays(dSettingValue.Array, sSettingValue.Array)
					}
				}
			}

			if sSettingValue.Group != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Group != nil {
						found = true
						mergeSettingValueGroups(dSettingValue.Group, sSettingValue.Group)
					}
				}
			}

			if sSettingValue.List != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.List != nil {
						found = true
						mergeSettingValueLists(dSettingValue.List, sSettingValue.List)
					}
				}
			}

			if !found {
				destination.SettingValues = append(destination.SettingValues, sSettingValue)
			}
		}
	}
}
