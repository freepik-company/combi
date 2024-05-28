package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

const (
	commentsRegex                    = `([#][^\n]*)`
	escapeCharsRegex                 = `(;|{|}|\s)`
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
// NGINX file parser
// ----------------------------------------------------------------

// type NGINX struct {
// 	Settings []*SettingT `parser:"@@*"`
// }

// type SettingT struct {
// 	SetingName   string         `parser:"@Name "`
// 	SettingValue *SettingValueT `parser:"@@"`
// }

// type SettingValueT struct {
// 	Primitive *PrimitiveT `parser:"( @@"`
// 	Block     *BlockT     `parser:" | @@ )"`
// }

// type PrimitiveT struct {
// 	Value string `parser:"@Value ';'"`
// }

// type BlockT struct {
// 	Settings []*SettingT `parser:"'{' @@* '}'"`
// }

// /////////////////////////////////////////////////////
// type Config struct {
// 	// Entries    []*Entry      `parser:"@@*"`
// 	Directives []*DirectiveT `parser:"@@*"`
// }

// type DirectiveT struct {
// 	DirectiveKey   string          `parser:"@Ident"`
// 	DirectiveValue DirectiveValueT `parser:"@@"`
// }

// type DirectiveValueT struct {
// 	DirectiveName *DirectiveNameT `parser:"( @Indent"`
// 	Value         *ValueT         `parser:"  | @@"`
// 	Block         *BlockT         `parser:"  | @@ )"`
// }

// type DirectiveNameT struct {
// 	Val string `parser:" @Ident (';'|'{')"`
// }

// type ValueT struct {
// 	Val string `parser:" @Ident ';'"`
// }

// type BlockT struct {
// 	Directives []*DirectiveT `parser:"'{' @@* '}'"`
// }

///////////////////////////////////////////////////////////////////////////

// type Entry struct {
// 	EntryKey  string `parser:"@Ident"`
// 	EntryName *Value `parser:"( @@"`
// 	// Value      *Value      `parser:"( @@"`
// 	DobleValue *DobleValue `parser:"  | @@"`
// 	Block      *Block      `parser:"  | @@ )"`
// }

// type Value struct {
// 	Val string `parser:" @Ident (@Ident)* (';'|'{')"`
// }

// type DobleValue struct {
// 	Val1 string `parser:" @Ident"`
// 	Val2 string `parser:" @Ident ';'"`
// }

// type Block struct {
// 	Entries []*Entry `parser:"'{' @@* '}'"`
// }

////////////////////////////////////////////////////////////////////////////////

type NGINX struct {
	Entries []*EntryT `parser:"@@*"`
}

type EntryT struct {
	Key        string        `parser:"@Key ' '*"`
	Directives []*DirectiveT `parser:"( @@ ';')"`
	// Blocks     []*BlockT     `parser:"| @@ )"`
}

type DirectiveT struct {
	Val string `parser:"@Val*"`
	// Val1 string `parser:"@Val ';'?"`
	// Val2 string `parser:"@Val ';'"`
	// Val []*ValT `parser:"@@ ';'"`
}

type ValT struct {
	Val string `parser:"@Val"`
}

type BlockT struct {
	Entries []*EntryT `parser:"'{' @@ '}'"`
}

func main() {
	// ----------------------------------------------------------------
	// NGINX file parser configuration
	// ----------------------------------------------------------------
	// nginxLexer := lexer.MustSimple([]lexer.SimpleRule{
	// 	{Name: `Name`, Pattern: settingNameRegex},
	// 	{Name: `Value`, Pattern: `[^;]+`},
	// 	{Name: "EscapeChars", Pattern: escapeCharsRegex},
	// 	{Name: "Comments", Pattern: commentsRegex},
	// 	{Name: "whitespace", Pattern: `(\s+)|(\n+)`},
	// })
	// nginxParser := participle.MustBuild[NGINX](
	// 	participle.Lexer(nginxLexer),
	// )

	// fmt.Printf("%s\n\n", nginxParser.String())

	// ----------------------------------------------------------------
	// NGINX parse file
	// ----------------------------------------------------------------
	nginxConfigBytes, err := os.ReadFile("nginx-primitive.conf")
	if err != nil {
		panic(err)
	}

	// nginx, err := nginxParser.ParseString("", string(nginxConfigBytes))
	// repr.Println(nginx, repr.Indent("  "), repr.OmitEmpty(true))
	// if err != nil {
	// 	panic(err)
	// }

	nginxLexer := lexer.MustSimple([]lexer.SimpleRule{
		// {Name: `Name`, Pattern: settingNameRegex},
		// {Name: `Key`, Pattern: `[A-Za-z*][-A-Za-z0-9_*]*`},
		{Name: `Key`, Pattern: `[^\s]+`},
		// {Name: `Val`, Pattern: `([\.\/A-Za-z0-9_*]+)`},
		{Name: `Val`, Pattern: `([^;]+)`},
		{Name: "Comments", Pattern: commentsRegex},
		{Name: "EscapeChars", Pattern: escapeCharsRegex},
		{Name: "whitespace", Pattern: `(\s+)`},
	})
	var parser = participle.MustBuild[NGINX](
		participle.Lexer(nginxLexer),
	)
	sep := "--------------------------------\n"
	fmt.Printf(
		sep+"%s\n"+sep,
		parser.String(),
	)
	expr, err := parser.ParseString("", string(nginxConfigBytes))
	repr.Println(expr, repr.Indent("  "), repr.OmitEmpty(true))
	fmt.Print(sep)
	if err != nil {
		panic(err)
	}
}
