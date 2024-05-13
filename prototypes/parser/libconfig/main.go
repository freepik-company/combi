package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

const (
	settingNameRegex                 = `[A-Za-z*][-A-Za-z0-9_*]*`
	settingValuePrimitiveStringRegex = `(\"([^\"\\]|\\.)*\")`
	settingValuePrimitiveFloatRegex  = `(([-+]?([0-9]*)?\.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(\.[0-9]*)?[eE][-+]?[0-9]+))`
	settingValuePrimitiveHexRegex    = `(0[Xx][0-9A-Fa-f]+(L+)?)`
	settingValuePrimitiveIntRegex    = `([-+]?[0-9]+(L+)?)`
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
	Entries []*SettingT `@@*`
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
	Value []*PrimitiveT `"[" @@* "]"`
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
		{Name: `Equal`, Pattern: `=|:`},
		{Name: `EndSetting`, Pattern: `;|,`},
		{Name: `Keys`, Pattern: `{|}`},
		{Name: `Brackets`, Pattern: `\[|\]`},
		{Name: `Parentesis`, Pattern: `\(|\)`},
		{Name: "comments", Pattern: `([#][^\n]*)|(\/\/[^\n]*)|(\/\*.*[\n]\*\/)`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	libconfigParser := participle.MustBuild[LIBCONFIG](
		participle.Lexer(libconfigLexer),
	)

	fmt.Printf("%s\n\n", libconfigParser.String())

	libconfig, err := libconfigParser.ParseString("", string(libconfigConfigBytes))
	repr.Println(libconfig, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
