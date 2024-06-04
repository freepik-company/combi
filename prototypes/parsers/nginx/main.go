package main

import (
	"fmt"
	"os"
	"path/filepath"
	"prototypes/globals"
	"regexp"
	"strings"

	"github.com/alecthomas/repr"

	"github.com/tufanbarisyildirim/gonginx/parser"
)

// Config representa la configuración completa de NGINX
type NGINX struct {
	Directives []DirectiveT
	Blocks     []BlockT
}

// Directive representa una directiva simple de NGINX
type DirectiveT struct {
	Name  string
	Value string
}

// Block representa un bloque de configuración de NGINX (por ejemplo, server, location)
type BlockT struct {
	Name       string
	Parameters string
	Directives []DirectiveT
	Blocks     []BlockT
}

// ParseConfig analiza el contenido del archivo de configuración
func ParseNginxConfig(configStr string) (config NGINX) {
	configStrLines := strings.Split(configStr, "\n")

	for lineIndex := 0; lineIndex < len(configStrLines); lineIndex++ {
		line := strings.TrimSpace(configStrLines[lineIndex])
		// Skip Commets and empty strings
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		// parse ngix blocks
		if strings.HasSuffix(line, "{") {
			open := 0
			endBlockIndex := lineIndex
			for k, v := range configStrLines[lineIndex:] {
				if strings.HasSuffix(v, "}") {
					open -= 1
				} else if strings.HasSuffix(v, "{") {
					open += 1
				}
				if open <= 0 {
					endBlockIndex += k
					break
				}
			}
			fmt.Printf("block from '%d' to '%d'\n", lineIndex+1, endBlockIndex+1)
			config.Blocks = append(config.Blocks, parseNginxBlock(configStrLines[lineIndex:endBlockIndex]))
			lineIndex = endBlockIndex
			continue
		}

		// parse nginx directives
		endDirectiveIndex := lineIndex
		for k, v := range configStrLines[lineIndex:] {
			if strings.HasSuffix(v, ";") {
				endDirectiveIndex += k
				break
			}
		}
		fmt.Printf("directive from '%d' to '%d'\n", lineIndex+1, endDirectiveIndex+1)
		config.Directives = append(config.Directives, parseNginxDirective(configStrLines[lineIndex:endDirectiveIndex+1]))
	}

	return config
}

func parseNginxDirective(directiveLines []string) (directive DirectiveT) {
	if len(directiveLines) == 1 {
		parts := strings.Fields(directiveLines[0])
		directive.Name = parts[0]
		directive.Value = strings.TrimSuffix(strings.Join(parts[1:], " "), ";")
		repr.Println(directive, repr.Indent("  "), repr.OmitEmpty(true))
		return directive
	}

	directiveLines[len(directiveLines)-1] = strings.TrimSuffix(directiveLines[len(directiveLines)-1], ";")

	firstLineParts := strings.Fields(directiveLines[0])
	directive.Name = firstLineParts[0]
	directiveLines[0] = strings.Join(firstLineParts[1:len(firstLineParts)-1], " ")

	directive.Value = strings.Join(directiveLines, " ")

	repr.Println(directive, repr.Indent("  "), repr.OmitEmpty(true))
	return directive
}

func parseNginxBlock(blockLines []string) (block BlockT) {
	repr.Println(blockLines, repr.Indent("  "), repr.OmitEmpty(true))
	return block
}

func main() {
	globals.InitLogger(globals.DEBUG)

	program := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		globals.Logger.Error(fmt.Sprintf("file as argument not provided (usage: %s <filepath>)", program))
		os.Exit(1)
	}

	filepath := os.Args[1]
	nginxConfigBytes, err := os.ReadFile(filepath)
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("unable to read file %s: %s", filepath, err.Error()))
		os.Exit(1)
	}

	// ----------------------------------------------------------------
	// NGINX file parser configuration
	// ----------------------------------------------------------------

	config := ParseNginxConfig(string(nginxConfigBytes))
	_ = config

	// Imprimir la configuración analizada
	// repr.Println(config, repr.Indent("  "), repr.OmitEmpty(true))
	os.Exit(0)

	//--------------------------------------------
	p := parser.NewStringParser(string(nginxConfigBytes))
	c, err := p.Parse()
	if err != nil {
		globals.Logger.Error(fmt.Sprintf("unable to parse config '%s'", filepath))
		os.Exit(1)
	}
	// fmt.Print(dumper.DumpConfig(c, dumper.IndentedStyle))
	repr.Println(c, repr.Indent("  "), repr.OmitEmpty(true))
	os.Exit(0)

	nginxDirectiveRegex := "([a-zA-Z0-9_]*([ ]+))([a-zA-Z0-9_/.]*([ ]+))?([a-zA-Z0-9_/.]*);"
	re, _ := regexp.Compile(nginxDirectiveRegex)
	directiveMatches := re.FindAll(nginxConfigBytes, -1)
	if directiveMatches == nil {
		globals.Logger.Error(fmt.Sprintf("unable to match directives in file %s", filepath))
		os.Exit(1)
	}

	nginxContextRegex := "([a-zA-Z0-9_]*([ \n]+)){([ \n]+)(.*)([ \n]+)}"
	re, _ = regexp.Compile(nginxContextRegex)
	contextMatches := re.FindAll(nginxConfigBytes, -1)
	if contextMatches == nil {
		globals.Logger.Error(fmt.Sprintf("unable to match contexts in file %s", filepath))
		os.Exit(1)
	}
	repr.Println(contextMatches, repr.Indent("  "), repr.OmitEmpty(true))
}
