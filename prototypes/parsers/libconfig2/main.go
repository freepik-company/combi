package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"prototypes/globals"

	"github.com/alecthomas/repr"
)

type Libconfig2T struct {
	ConfigStruct SettingListT
	configMap    map[string]any
}

type SettingListT struct {
	PrimitiveList []Primitive2T
	ArrayList     []Array2T
	GroupList     []Group2T
}

type Array2T struct {
	Values []string
}

type Group2T struct {
	Values []string
}

type Primitive2T struct {
	Name  string
	Value string
}

// ----------------------------------------------------------------
// Decode/Encode NGINX data structure
// ----------------------------------------------------------------

// Decode functions

func (e *Libconfig2T) DecodeConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = e.DecodeConfigBytes(configBytes)
	return err
}

func (e *Libconfig2T) DecodeConfigBytes(configBytes []byte) (err error) {
	// Remove one line comments in file
	configStr := string(configBytes)
	configStr = regexp.MustCompile(`#[^\n]*`).ReplaceAllString(configStr, "")
	configStr = regexp.MustCompile(`[\s]*[=:][\s]*`).ReplaceAllStringFunc(configStr, func(match string) string {
		result := "="
		if strings.Contains(match, ":") {
			result = ":"
		}
		return result
	})

	// Parse formatted nginx configuration
	// configStrLines := strings.Split(configStr, "\n")

	// configStrLines = slices.DeleteFunc(configStrLines, func(str string) bool {
	// 	return str == ""
	// })

	// repr.Println(configStrLines, repr.Indent("  "), repr.OmitEmpty(true))
	// fmt.Println(configStr)
	fmt.Println("------------------------------------------------------------")
	e.configMap = map[string]any{}
	e.parseString(configStr)
	// err = parseNginxBlockContent(&e.ConfigStruct, configStrLines)
	return err
}

func isWhitespace(b byte) bool {
	return slices.Contains([]byte{'\n', '\t', '\r', ' '}, b)
}

func isScope(b byte) bool {
	return slices.Contains([]byte{'{', '(', '[', ']', ')', '}'}, b)
}

func isCloseValueChar(b byte) bool {
	return slices.Contains([]byte{',', ';'}, b)
}

func (e *Libconfig2T) parseString(config string) {
	configLen := len(config)
	for i := 0; i < configLen; i++ {
		if isWhitespace(config[i]) || isCloseValueChar(config[i]) {
			continue
		}

		// parse name
		name, nameDiff := parseSettingName(config[i:])
		e.configMap[name] = nil
		fmt.Printf("setting name: '%s'; start: %d; end: %d\n", name, i, i+nameDiff)
		i += nameDiff

		// parse primitive value
		if !isScope(config[i]) {
			value, valueDiff := parseSettingPrimitiveValue(config[i:])
			e.configMap[name] = value
			fmt.Printf("setting value '%s'; start: %d; end: %d\n", value, i, i+valueDiff)
			i += valueDiff
			continue
		}

		// parse primitive value

		break

	}

	fmt.Println("------------------------------------------------------------")
}

func parseSettingName(config string) (value string, diff int) {
	configLen := len(config)
	index := 0
	diff = index + 1
	found := false
	for ; diff < configLen && !found; diff++ {
		if config[diff] == '=' {
			found = true
		}
	}
	value = config[index : diff-1]
	return value, diff
}

func parseSettingPrimitiveValue(config string) (value string, diff int) {
	configLen := len(config)
	index := 0
	diff = index + 1
	found := false
	if config[index] == '"' {
		for ; diff < configLen && !found; diff++ {
			if config[diff] == '"' && config[diff-1] != '\\' {
				found = true
			}
		}
		value = config[index:diff]
		return value, diff
	}

	for ; diff < configLen && !found; diff++ {
		if isWhitespace(config[diff]) || isCloseValueChar(config[diff]) {
			found = true
		}
	}
	value = config[index : diff-1]

	return value, diff
}

// func parseNginxBlockContent(blockContent *BlockContentT, blockContentLines []string) (err error) {
// 	// Parse block content
// 	for blockLineIndex := 0; blockLineIndex < len(blockContentLines); blockLineIndex++ {
// 		line := strings.TrimSpace(blockContentLines[blockLineIndex])
// 		// Skip empty strings (only to be carefull)
// 		if len(line) == 0 {
// 			continue
// 		}

// 		// Parse nginx directives
// 		if strings.HasSuffix(line, ";") {
// 			directiveParts := strings.Fields(line)
// 			subDirective := DirectiveT{
// 				Name:  directiveParts[0],
// 				Value: strings.Join(directiveParts[1:len(directiveParts)-1], " "),
// 			}
// 			if len(directiveParts) > 3 {
// 				subDirective.Param = directiveParts[1]
// 				subDirective.Value = strings.Join(directiveParts[2:len(directiveParts)-1], " ")
// 			}
// 			blockContent.Directives = append(blockContent.Directives, subDirective)
// 		}

// 		// Parse ngix blocks
// 		if strings.HasSuffix(line, "{") {
// 			startBlockLines := blockContentLines[blockLineIndex:]
// 			i, open := 0, 0
// 			for ; i < len(startBlockLines); i++ {
// 				if strings.HasSuffix(startBlockLines[i], "{") {
// 					open += 1
// 				}
// 				if strings.HasSuffix(startBlockLines[i], "}") {
// 					open -= 1
// 				}
// 				if open <= 0 {
// 					break
// 				}
// 			}
// 			if open > 0 {
// 				err = fmt.Errorf("nginx block '%s' without closed bracket '}'", blockContentLines[blockLineIndex])
// 				return err
// 			}
// 			subBlock, err := parseNginxBlock(blockContentLines[blockLineIndex : blockLineIndex+i+1])
// 			if err != nil {
// 				return err
// 			}
// 			blockContent.Blocks = append(blockContent.Blocks, subBlock)
// 			blockLineIndex += i
// 			continue
// 		}

// 		if strings.HasSuffix(line, "}") {
// 			err = fmt.Errorf("over closed bracket '}' in nginx configuration")
// 			return err
// 		}
// 	}

// 	return err
// }

// func parseNginxBlock(blockLines []string) (block BlockT, err error) {
// 	// Parse name and parameters
// 	nameParamsParts := strings.Fields(blockLines[0])
// 	block.Name = nameParamsParts[0]
// 	block.Params = strings.Join(nameParamsParts[1:len(nameParamsParts)-1], " ")

// 	// Parse block content
// 	err = parseNginxBlockContent(&block.BlockContent, blockLines[1:len(blockLines)-1])
// 	return block, err
// }

// // Encode functions

// func (e *NginxT) EncodeConfigString() (configStr string) {
// 	configStr = encodeNginxBlockContent(e.ConfigStruct, 0)
// 	return configStr
// }

// func encodeNginxBlockContent(blockContent BlockContentT, indent int) (configStr string) {
// 	indentStr := ""
// 	for i := 0; i < indent; i++ {
// 		indentStr += "    "
// 	}

// 	for _, val := range blockContent.Directives {
// 		configStr += indentStr + val.Name + " " + val.Param + " " + val.Value + ";\n"
// 	}

// 	for _, val := range blockContent.Blocks {
// 		configStr += indentStr + val.Name + " " + val.Params + " {\n"
// 		configStr += encodeNginxBlockContent(val.BlockContent, indent+1)
// 		configStr += indentStr + "}\n"
// 	}

// 	return configStr
// }

func main() {
	globals.InitLogger(globals.DEBUG, nil)
	program := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		globals.Logger.Fatalf("file as argument not provided (usage: %s <filepath>)", program)
	}

	filepath := os.Args[1]
	// libconfigConfigBytes, err := os.ReadFile(filepath)
	// if err != nil {
	// 	globals.Logger.Fatalf("unable to read file %s: %s", filepath, err.Error())
	// }

	// ----------------------------------------------------------------
	// LIBCONFIG file parser
	// ----------------------------------------------------------------
	libconfig := Libconfig2T{}
	err := libconfig.DecodeConfig(filepath)
	repr.Println(libconfig.ConfigStruct, repr.Indent("  "), repr.OmitEmpty(true))
	if err != nil {
		globals.Logger.Fatalf("unable to parse file %s: %s", filepath, err.Error())
	}
}
