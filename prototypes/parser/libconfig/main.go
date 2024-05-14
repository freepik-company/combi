package main

import (
	"fmt"
	"gcmerge/internal/conditions"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
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
// LIBCONFIG file parser
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
	Value string `@Value`
}

type GroupT struct {
	Settings []*SettingT `"{" @@* "}"`
}

type ArrayT struct {
	Primitives []*PrimitiveT `"[" @@* "]"`
}

type ListT struct {
	List []*SettingValueT `"(" @@* ")"`
}

func main() {
	libconfigConfigBytes, err := os.ReadFile("proxysql.cnf")
	if err != nil {
		panic(err)
	}

	// ----------------------------------------------------------------
	// LIBCONFIG file parser
	// ----------------------------------------------------------------
	libconfigLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Name`, Pattern: settingNameRegex},
		{Name: `Value`, Pattern: settingValuePrimitiveRegex},
		{Name: "EscapeChars", Pattern: escapeCharsRegex},
		{Name: "Comments", Pattern: commentsRegex},
		{Name: "whitespace", Pattern: `(\s+)`},
	})
	libconfigParser := participle.MustBuild[LIBCONFIG](
		participle.Lexer(libconfigLexer),
	)

	fmt.Printf("%s\n\n", libconfigParser.String())

	libconfig, err := libconfigParser.ParseString("", string(libconfigConfigBytes))
	configStruct := conditions.DeepCopy(libconfig)
	repr.Println(configStruct, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
