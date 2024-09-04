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

type LibconfigT struct {
	ConfigStruct []SettingT
	configMap    map[string]any
}

type SettingT struct {
	Name  string
	Value SettingValueT
}

type SettingValueT struct {
	Primitive string
	Array     []string
	// Group     GroupT
	Group []SettingT
	List  ListT
}

type GroupT struct {
}

type ListT struct {
}

// ----------------------------------------------------------------
// Decode/Encode LIBCONFIG data structure
// ----------------------------------------------------------------

// Decode functions

func (e *LibconfigT) DecodeConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = e.DecodeConfigBytes(configBytes)
	return err
}

func (e *LibconfigT) DecodeConfigBytes(configBytes []byte) (err error) {
	// Remove one line comments in file
	configStr := string(configBytes)
	configStr = regexp.MustCompile(`#[^\n]*`).ReplaceAllString(configStr, "")

	err = e.parseLibconfigString(configStr)
	return err
}

func isWhitespace(b byte) bool {
	return slices.Contains([]byte{'\n', '\t', '\r', ' '}, b)
}

func isOpenScope(b byte) bool {
	return slices.Contains([]byte{'{', '(', '['}, b)
}

func isCloseScope(b byte) bool {
	return slices.Contains([]byte{']', ')', '}'}, b)
}

func isCloseValue(b byte) bool {
	return slices.Contains([]byte{',', ';'}, b)
}

func isEqual(b byte) bool {
	return slices.Contains([]byte{'=', ':'}, b)
}

func (e *LibconfigT) parseLibconfigString(config string) (err error) {
	e.configMap, err = parseSettings(config)
	repr.Println(e.configMap, repr.Indent("  "), repr.OmitEmpty(true))
	return err
}

func parseSettings(config string) (settings map[string]any, err error) {
	settings = map[string]any{}

	configLen := len(config)
	for i := 0; i < configLen; i++ {
		spacesDiff := skipParseSpaces(config[i:])
		i += spacesDiff
		if i >= configLen {
			break
		}

		// parse name
		name, nameDiff := parseSettingName(config[i:])
		settings[name] = nil
		i += nameDiff
		if i >= configLen {
			break
		}

		spacesDiff = skipParseSpaces(config[i:])
		i += spacesDiff
		if i >= configLen {
			break
		}

		// parse value
		value, valueDiff, err := parseSettingValue(config[i:])
		if err != nil {
			return settings, err
		}

		settings[name] = value
		i += valueDiff
		if i >= configLen {
			break
		}

		spacesDiff = skipParseSpaces(config[i:])
		i += spacesDiff
		if i >= configLen {
			break
		}

		if isCloseValue(config[i]) && i+1 < configLen {
			i++
		}
	}

	return settings, err
}

func parseSettingValue(config string) (value any, diff int, err error) {
	switch config[0] {
	case '[':
		value, diff, err = parseSettingValueArray(config)
	case '{':
		value, diff, err = parseSettingValueGroup(config)
	case '(':
		value, diff, err = parseSettingValueList(config)
	default:
		value, diff = parseSettingValuePrimitive(config)
	}

	return value, diff, err
}

func skipParseSpaces(config string) (diff int) {
	for diff = 0; diff < len(config) && isWhitespace(config[diff]); diff++ {
	}
	if diff > 0 {
		diff--
	}

	return diff
}

func parseSettingName(config string) (name string, diff int) {
	configLen := len(config)
	for diff = 0; diff < configLen; diff++ {
		if isEqual(config[diff]) || isWhitespace(config[diff]) {
			break
		}
	}
	name = config[:diff]

	config = config[diff:]
	configLen = len(config)
	for i := 0; i < configLen; i++ {
		if isEqual(config[i]) {
			diff += i + 1
			break
		}
	}

	return name, diff
}

func parseSettingValuePrimitive(config string) (value string, diff int) {
	configLen := len(config)
	if config[0] == '"' {
		for diff := 1; diff < configLen; diff++ {
			if config[diff] == '"' && config[diff-1] != '\\' {
				diff++
				value = config[:diff]
				return value, diff
			}
		}
	}

	for diff := 1; diff < configLen; diff++ {
		if isWhitespace(config[diff]) || isCloseValue(config[diff]) || isCloseScope(config[diff]) {
			value = config[:diff]
			return value, diff
		}
	}

	return value, diff
}

func parseSettingValueArray(config string) (value []string, diff int, err error) {
	configLen := len(config)
	count := 1
	for diff = 1; diff < configLen; diff++ {
		if config[diff] == '[' {
			count++
		}

		if config[diff] == ']' {
			count--
			if count <= 0 {
				diff++
				break
			}
		}

	}

	if count > 0 {
		err = fmt.Errorf("unclose array")
		return value, diff, err
	}

	arrayConfigStr := regexp.MustCompile(`[\s]`).ReplaceAllString(config[:diff], "")
	arrayConfigStrLen := len(arrayConfigStr)
	for i := 0; i < arrayConfigStrLen; i++ {
		if isCloseValue(arrayConfigStr[i]) || isOpenScope(arrayConfigStr[i]) || isWhitespace(arrayConfigStr[i]) {
			continue
		}

		if i >= arrayConfigStrLen {
			break
		}
		primitive, pDiff := parseSettingValuePrimitive(arrayConfigStr[i:])
		value = append(value, primitive)
		i += pDiff
		if i >= arrayConfigStrLen {
			break
		}
	}

	return value, diff, err
}

func parseSettingValueGroup(config string) (groupSettings map[string]any, diff int, err error) {
	configLen := len(config)
	count := 1
	for diff = 1; diff < configLen; diff++ {
		if config[diff] == '{' {
			count++
		}

		if config[diff] == '}' {
			count--
			if count <= 0 {
				// diff++
				break
			}
		}
	}

	if count > 0 {
		err = fmt.Errorf("unclose group")
		return groupSettings, diff, err
	}

	groupStr := strings.TrimSuffix(strings.TrimPrefix(config[:diff], "{"), "}")
	fmt.Println("--------------------------------------------------------------------")
	fmt.Println("group:", groupStr)
	fmt.Println("--------------------------------------------------------------------")
	groupSettings, err = parseSettings(groupStr)

	return groupSettings, diff, err
}

func parseSettingValueList(config string) (valueList []any, diff int, err error) {
	configLen := len(config)
	count := 1
	for diff = 1; diff < configLen && count <= 0; diff++ {
		if config[diff] == '(' {
			count++
		}

		if config[diff] == ')' {
			count--
		}
	}

	if count > 0 {
		err = fmt.Errorf("unclose list")
		return valueList, diff, err
	}

	listConfigStr := strings.TrimSuffix(strings.TrimPrefix(config[:diff], "("), ")")

	listConfigLen := len(listConfigStr)
	for i := 0; i < listConfigLen; i++ {
		spacesDiff := skipParseSpaces(listConfigStr[i:])
		i += spacesDiff
		if i >= listConfigLen {
			break
		}

		var value any
		var valueDiff int
		value, valueDiff, err = parseSettingValue(listConfigStr[i:])
		if err != nil {
			return valueList, diff, err
		}
		valueList = append(valueList, value)
		i += valueDiff
		if i >= listConfigLen {
			break
		}

		spacesDiff = skipParseSpaces(listConfigStr[i:])
		i += spacesDiff
		if i >= listConfigLen {
			break
		}
	}

	return valueList, diff, err
}

func main() {
	globals.InitLogger(globals.DEBUG, nil)
	program := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		globals.Logger.Fatalf("file as argument not provided (usage: %s <filepath>)", program)
	}
	filepath := os.Args[1]

	// ----------------------------------------------------------------
	// LIBCONFIG file parser
	// ----------------------------------------------------------------
	libconfig := LibconfigT{}
	err := libconfig.DecodeConfig(filepath)
	if err != nil {
		globals.Logger.Fatalf("unable to parse file %s: %s", filepath, err.Error())
	}
}
