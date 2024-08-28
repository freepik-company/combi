package main

import (
	"fmt"
	"os"
	"path/filepath"
	"prototypes/globals"
	"regexp"
	"strings"

	"github.com/alecthomas/repr"
)

type NGINX struct {
	RootContent BlockContentT
}

type BlockContentT struct {
	Directives []DirectiveT
	Blocks     []BlockT
}

type DirectiveT struct {
	Name  string
	Value string
}

type BlockT struct {
	Name         string
	Parameters   string
	BlockContent BlockContentT
}

// ParseConfig analiza el contenido del archivo de configuración
func ParseNginxConfig(configStr string) (config NGINX, err error) {
	// Remove one line comments in file
	re := regexp.MustCompile(`#[^\n]*?\n`)
	configStr = re.ReplaceAllString(configStr, "\n")

	// Format configuration to parse nginx config format by line
	configStr = strings.Join(strings.Fields(configStr), " ")
	configStr = strings.ReplaceAll(configStr, "\n", "")
	configStr = strings.ReplaceAll(configStr, "{", " {\n")
	configStr = strings.ReplaceAll(configStr, "}", "}\n")
	configStr = strings.ReplaceAll(configStr, ";", " ;\n")
	configStr = strings.ReplaceAll(configStr, "\n ", "\n")

	// Parse formatted nginx configuration
	configStrLines := strings.Split(configStr, "\n")
	err = parseNginxBlockContent(&config.RootContent, configStrLines)
	return config, err
}

func parseNginxBlockContent(blockContent *BlockContentT, blockContentLines []string) (err error) {
	// Parse block content
	for blockLineIndex := 0; blockLineIndex < len(blockContentLines); blockLineIndex++ {
		line := strings.TrimSpace(blockContentLines[blockLineIndex])
		// Skip empty strings (only to be carefull)
		if len(line) == 0 {
			continue
		}

		// Parse nginx directives
		if strings.HasSuffix(line, ";") {
			subDirective := parseNginxDirective(blockContentLines[blockLineIndex])
			blockContent.Directives = append(blockContent.Directives, subDirective)
		}

		// Parse ngix blocks
		if strings.HasSuffix(line, "{") {
			startBlockLines := blockContentLines[blockLineIndex:]
			i, open := 0, 0
			for ; i < len(startBlockLines); i++ {
				if strings.HasSuffix(startBlockLines[i], "{") {
					open += 1
				}
				if strings.HasSuffix(startBlockLines[i], "}") {
					open -= 1
				}
				if open <= 0 {
					break
				}
			}
			if open > 0 {
				err = fmt.Errorf("nginx block '%s' not close", blockContentLines[blockLineIndex])
				return err
			}
			subBlock, err := parseNginxBlock(blockContentLines[blockLineIndex : blockLineIndex+i+1])
			if err != nil {
				return err
			}
			blockContent.Blocks = append(blockContent.Blocks, subBlock)
			blockLineIndex += i
			continue
		}

		if strings.HasSuffix(line, "}") {
			err = fmt.Errorf("over closed bracket '}' in nginx configuration")
			return err
		}
	}

	return err
}

func parseNginxDirective(directiveLine string) (directive DirectiveT) {
	parts := strings.Fields(directiveLine)
	directive.Name = parts[0]
	directive.Value = strings.Join(parts[1:len(parts)-1], " ")
	return directive
}

func parseNginxBlock(blockLines []string) (block BlockT, err error) {
	// Parse name and parameters
	nameParamsParts := strings.Fields(blockLines[0])
	block.Name = nameParamsParts[0]
	block.Parameters = strings.Join(nameParamsParts[1:len(nameParamsParts)-1], " ")

	// Parse block content
	err = parseNginxBlockContent(&block.BlockContent, blockLines[1:len(blockLines)-1])
	return block, err
}

func main() {
	globals.InitLogger(globals.DEBUG, nil, "prototype", "nginx")

	program := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		globals.Logger.Fatalf("file as argument not provided (usage: %s <filepath>)", program)
	}

	filepath := os.Args[1]
	globals.Logger.Infof("reading file '%s'", filepath)
	nginxConfigBytes, err := os.ReadFile(filepath)
	if err != nil {
		globals.Logger.Fatalf("unable to read file %s: %s", filepath, err.Error())
	}

	// ----------------------------------------------------------------
	// NGINX file parser configuration
	// ----------------------------------------------------------------

	globals.Logger.Infof("parse configuration")
	config, err := ParseNginxConfig(string(nginxConfigBytes))
	if err != nil {
		globals.Logger.Fatalf("unable to parse file %s: %s", filepath, err.Error())
	}
	repr.Println(config, repr.Indent("  "), repr.OmitEmpty(true))
}
