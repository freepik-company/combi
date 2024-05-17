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
	Settings []*SettingT `parser:"@@*"`
}

type SettingT struct {
	SetingName   string         `parser:"@Name ('='|':')"`
	SettingValue *SettingValueT `parser:"@@"`
}

type SettingValueT struct {
	Primitive *PrimitiveT `parser:"( @@ (';'?','?)"`
	Group     *GroupT     `parser:" | @@ (','?)"`
	Array     *ArrayT     `parser:" | @@ (','?)"`
	List      *ListT      `parser:" | @@ (','?))"`
}

type PrimitiveT struct {
	Value string `parser:"@Value (','?)"`
}

type ArrayT struct {
	Primitives []*PrimitiveT `parser:"'[' @@* ']'"`
}

type GroupT struct {
	Settings []*SettingT `parser:"'{' @@* '}'"`
}

type ListT struct {
	SettingValues []*SettingValueT `parser:"'(' @@* ')'"`
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
