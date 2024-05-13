package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// ----------------------------------------------------------------
// INI file parser
// ----------------------------------------------------------------
type INI struct {
	Properties []*Property `@@*`
	Sections   []*Section  `@@*`
}

type Section struct {
	Identifier string      `"[" @Ident "]"`
	Properties []*Property `@@*`
}

type Property struct {
	Key   string `@Ident "="`
	Value Value  `@@`
}

type Value interface{ value() }

type String struct {
	String string `@String`
}

func (String) value() {}

type Number struct {
	Number float64 `@Float`
}

func (Number) value() {}

func main() {
	// ----------------------------------------------------------------
	// INI file parser
	// ----------------------------------------------------------------
	iniLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Ident`, Pattern: `[a-zA-Z][a-zA-Z_\d]*`},
		{Name: `String`, Pattern: `"(?:\\.|[^"])*"`},
		{Name: `Float`, Pattern: `\d+(?:\.\d+)?`},
		{Name: `Punct`, Pattern: `[][=]`},
		{Name: "comment", Pattern: `[#;][^\n]*`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	parser := participle.MustBuild[INI](
		participle.Lexer(iniLexer),
		participle.Unquote("String"),
		participle.Union[Value](String{}, Number{}),
	)

	fmt.Printf("%s\n\n", parser.String())

	iniConfigStr, err := os.ReadFile("config.ini")
	if err != nil {
		panic(err)
	}

	ini, err := parser.ParseString("", string(iniConfigStr))
	repr.Println(ini, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
