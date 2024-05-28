// nolint: golint, dupl
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/repr"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

const (
	jsonCommentRegex     = `(\/\/[^\n]*)`
	jsonEscapeCharsRegex = `((=|:)|(;|,)|({|})|(\[|\])|(\(|\)))`
	jsonStringRegex      = `("(\\"|[^"])*")`
	jsonNumberRegex      = `([-+]?(\d*\.)?\d+)`
	jsonPunctRegex       = `([-[!@#$%^&*()+_={}\|:;"'<,>.?/]|])`
	jsonNullRegex        = `(null)`
	jsonTrueRegex        = `(true)`
	jsonFalseRegex       = `(false)`
	jsonPrimitiveRegex   = `(` +
		jsonStringRegex + `|` +
		jsonNumberRegex + `|` +
		jsonNullRegex + `|` +
		jsonTrueRegex + `|` +
		jsonFalseRegex + `)`
)

type JSON struct {
	Value *JsonValueT `parser:"@@"`
}

type JsonValueT struct {
	Object    *ObjectPairT `parser:"('{' @@ '}') |"`
	Array     *ArrayT      `parser:"@@ |"`
	Primitive *PrimitiveT  `parser:"@@"`
}

type ObjectT struct {
	Pair []*ObjectPairT `parser:"'{' @@ '}'"`
}

type ArrayT struct {
	Items []*JsonValueT `parser:"'[' @@ (',' @@)* ']'"`
}

type ObjectPairT struct {
	Key   string      `parser:"@String ':'"`
	Value *JsonValueT `parser:"@@"`
}

type PrimitiveT struct {
	Value string `parser:"@Primitive"`
}

func main() {
	libconfigConfigBytes, err := os.ReadFile("test.json")
	if err != nil {
		panic(err)
	}

	// ----------------------------------------------------------------
	// JSON file parser
	// ----------------------------------------------------------------
	jsonLexer := lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Comment", Pattern: jsonCommentRegex},
		{Name: "Primitive", Pattern: jsonPrimitiveRegex},
		{Name: "String", Pattern: jsonStringRegex},
		{Name: "Punct", Pattern: jsonPunctRegex},
		{Name: "EOL", Pattern: `[\n\r]+`},
		{Name: "Whitespace", Pattern: `[ \t]+`},
	})
	jsonParser := participle.MustBuild[JSON](
		participle.Lexer(jsonLexer),
	)

	fmt.Printf("%s\n\n", jsonParser.String())

	json, err := jsonParser.ParseString("", string(libconfigConfigBytes))
	repr.Println(json, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		panic(err)
	}
}
