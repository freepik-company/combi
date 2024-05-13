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
	// Settings []*ConfigT `@Name ("="|":") (@String|@Float|@Hex32|@Hex64|@Integer32|@Integer64)`
	// Settings []*SettingT `@@*`
	Group []*GroupT `@@*`
	// Config []*ConfigT `@@*`
}

type ConfigT struct {
	// Name  string `@Name ("="|":")`
	// Value ValueT `@@ (";"?","?)`
}

type GroupT struct {
	Name     string      `@Name ("="|":")`
	Settings []*SettingT `"{" @@* "}"`
}

type SettingT struct {
	Name  string `@Name ("="|":")`
	Value string `@Value (";"?","?)`
	// Value ValueT `@@ (";"?","?)`
}

// type ValueT interface{ value() }

// type StringT struct {
// 	String string `@String`
// }

// func (StringT) value() {}

// type Float32T struct {
// 	Float float32 `@Float32`
// }

// func (Float32T) value() {}

// type Hex32T struct {
// 	Hex int32 `@Hex32`
// }

// func (Hex32T) value() {}

// type Hex64T struct {
// 	Hex string `@Hex64`
// }

// func (Hex64T) value() {}

// type Integer32T struct {
// 	Integer int32 `@Int32`
// }

// func (Integer32T) value() {}

// type Integer64T struct {
// 	Integer string `@Int64`
// }

// func (Integer64T) value() {}

// type BoolT struct {
// 	Boolean string `@Bool`
// }

// func (BoolT) value() {}

func main() {
	// ----------------------------------------------------------------
	// LIBCONFIG file parser
	// ----------------------------------------------------------------
	libconfigLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: `Name`, Pattern: `[A-Za-z*][-A-Za-z0-9_*]*`},
		{Name: `Value`, Pattern: `(\"([^\"\\]|\\.)*\")|(([-+]?([0-9]*)?\.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(\.[0-9]*)?[eE][-+]?[0-9]+))|(0[Xx][0-9A-Fa-f]+(L(L)?)?)|([-+]?[0-9]+(L(L)?)?)`},
		// {Name: `String`, Pattern: `\"([^\"\\]|\\.)*\"`},
		// {Name: `Float32`, Pattern: `([-+]?([0-9]*)?\.[0-9]*([eE][-+]?[0-9]+)?)|([-+]([0-9]+)(\.[0-9]*)?[eE][-+]?[0-9]+)`},
		// {Name: `Hex64`, Pattern: `0[Xx][0-9A-Fa-f]+L(L)?`},
		// {Name: `Hex32`, Pattern: `0[Xx][0-9A-Fa-f]+`},
		// {Name: `Int64`, Pattern: `[-+]?[0-9]+L(L)?`},
		// {Name: `Int32`, Pattern: `[-+]?[0-9]+`},
		// {Name: `Bool`, Pattern: `(([Tt][Rr][Uu][Ee])|([Ff][Aa][Ll][Ss][Ee]))`},
		{Name: `Equal`, Pattern: `=|:`},
		{Name: `EndSetting`, Pattern: `;|,`},
		{Name: `Brackets`, Pattern: `{|}`},
		{Name: "comments", Pattern: `([#][^\n]*)|(\/\/[^\n]*)|(\/\*.*[\n]\*\/)`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	libconfigParser := participle.MustBuild[LIBCONFIG](
		participle.Lexer(libconfigLexer),
		// participle.Unquote("String"),
		// participle.Union[ValueT](
		// 	StringT{},
		// 	Float32T{},
		// 	Hex32T{},
		// 	Hex64T{},
		// 	Integer32T{},
		// 	Integer64T{},
		// 	// BoolT{},
		// ),
	)

	fmt.Printf("%s\n\n", libconfigParser.String())

	libconfigConfigStr, err := os.ReadFile("proxysql-only-root-settings.cnf")
	if err != nil {
		panic(err)
	}

	libconfig, err := libconfigParser.ParseString("", string(libconfigConfigStr))
	repr.Println(libconfig, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
