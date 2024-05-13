package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/repr"
)

// ----------------------------------------------------------------
// LIBCONFIG file parser
// ----------------------------------------------------------------

type LIBCONFIG struct {
	Entries []*SettingT `@@*`
}

type SettingT struct {
	Name  string  `@Name ("="|":")`
	Value *ValueT `( @@`
	Group *GroupT ` | @@`
	Array *ArrayT ` | @@`
	List  *ListT  ` | @@ )`
}

type GroupT struct {
	Settings []*SettingT `"{" @@* "}"`
}

type ArrayT struct {
	Value []*ValueT `"[" @@* "]"`
}

type ValueT struct {
	Value string `@Value (";"?","?)`
}

type ListT struct {
	Settings []*SettingListT `"(" @@* ")"`
}

type SettingListT struct {
	Value *ValueT `( @@`
	Group *GroupT ` | @@ (","?)`
	Array *ArrayT ` | @@ (","?)`
	List  *ListT  ` | @@ (","?))`
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
		{Name: `Name`, Pattern: `[A-Za-z*][-A-Za-z0-9_*]*`},
		{Name: `Value`, Pattern: `(\"([^\"\\]|\\.)*\")|(([-+]?([0-9]*)?\.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(\.[0-9]*)?[eE][-+]?[0-9]+))|(0[Xx][0-9A-Fa-f]+(L(L)?)?)|([-+]?[0-9]+(L(L)?)?)`},
		// {Name: `Value`, Pattern: `(\"([^\"\\]|\\.)*\")`},
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
