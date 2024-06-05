package main

import (
	"fmt"
	"os"
	"path/filepath"
	"prototypes/globals"

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
	Value string `parser:"@Value"`
}

type GroupT struct {
	Settings []*SettingT `parser:"'{' @@* '}'"`
}

type ArrayT struct {
	Primitives []*PrimitiveT `parser:"'[' @@* ']'"`
}

type ListT struct {
	List []*SettingValueT `parser:"'(' @@* ')'"`
}

func main() {
	globals.InitLogger(globals.DEBUG, nil)
	program := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		globals.Logger.Fatalf("file as argument not provided (usage: %s <filepath>)", program)
	}

	filepath := os.Args[1]
	libconfigConfigBytes, err := os.ReadFile(filepath)
	if err != nil {
		globals.Logger.Fatalf("unable to read file %s: %s", filepath, err.Error())
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
	repr.Println(libconfig, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		globals.Logger.Fatalf("unable to parse file %s: %s", filepath, err.Error())
	}
}
