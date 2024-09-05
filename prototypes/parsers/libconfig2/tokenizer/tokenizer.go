package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"

	"prototypes/globals"

	"github.com/alecthomas/repr"
)

const (
	TOKEN_TYPE_DEFAULT int = iota
	TOKEN_TYPE_NAME
	TOKEN_TYPE_EQUAL
	TOKEN_TYPE_VALUE
	TOKEN_TYPE_OPEN_BRACKET
	TOKEN_TYPE_CLOSE_BRACKET
	TOKEN_TYPE_OPEN_PAREN
	TOKEN_TYPE_CLOSE_PAREN
	TOKEN_TYPE_OPEN_SQUARE_BRACKET
	TOKEN_TYPE_CLOSE_SQUARE_BRACKET
)

type LibconfigTokenT struct {
	Type  int
	Token string
}

type LibconfigT struct {
	// ConfigStruct []SettingT
	configMap map[string]any
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

func isSpace(b byte) bool {
	return slices.Contains([]byte{'\n', '\t', '\r', ' '}, b)
}

func isScope(b byte) bool {
	return slices.Contains([]byte{'{', '(', '[', ']', ')', '}'}, b)
}

func isCloseValue(b byte) bool {
	return slices.Contains([]byte{',', ';'}, b)
}

func isEqual(b byte) bool {
	return slices.Contains([]byte{'=', ':'}, b)
}

func getScopeType(b byte) (result int) {
	switch b {
	case '[':
		result = TOKEN_TYPE_OPEN_SQUARE_BRACKET
	case '{':
		result = TOKEN_TYPE_OPEN_BRACKET
	case '(':
		result = TOKEN_TYPE_OPEN_PAREN
	case ')':
		result = TOKEN_TYPE_CLOSE_PAREN
	case '}':
		result = TOKEN_TYPE_CLOSE_BRACKET
	case ']':
		result = TOKEN_TYPE_CLOSE_SQUARE_BRACKET
	}

	return result
}

func (e *LibconfigT) parseLibconfigString(config string) (err error) {
	tokens, err := parseTokens(config)
	if err != nil {
		return err
	}

	e.configMap, err = parseSettings(tokens)
	if err != nil {
		return err
	}
	repr.Println(e.configMap, repr.Indent("  "), repr.OmitEmpty(true))
	return err
}

func parseToken(config string) (value string, diff int, err error) {
	configLen := len(config)
	if config[0] == '"' {
		for diff = 1; diff < configLen; diff++ {
			if config[diff] == '"' && config[diff-1] != '\\' {
				diff++
				value = config[:diff]
				break
			}
		}

		if diff >= configLen {
			err = fmt.Errorf("unclosed string")
		}

		return value, diff, err
	}

	for diff = 1; diff < configLen; diff++ {
		if isSpace(config[diff]) || isEqual(config[diff]) || isScope(config[diff]) || isCloseValue(config[diff]) {
			value = config[:diff]
			diff--
			break
		}
	}

	return value, diff, err
}

func parseTokens(config string) (tokens []LibconfigTokenT, err error) {
	configLen := len(config)
	diff := 0
	for i := 0; i < configLen; i += 1 + diff {
		// fmt.Println("looop")
		diff = 0
		if isSpace(config[i]) || isCloseValue(config[i]) {
			continue
		}

		token := string(config[i])
		if isEqual(config[i]) {
			tokens = append(tokens, LibconfigTokenT{
				Type:  TOKEN_TYPE_EQUAL,
				Token: token,
			})
			continue
		}

		if isScope(config[i]) {
			tokens = append(tokens, LibconfigTokenT{
				Type:  getScopeType(config[i]),
				Token: token,
			})
			continue
		}

		// parse token
		token, diff, err = parseToken(config[i:])
		if err != nil {
			return tokens, err
		}

		TokenType := TOKEN_TYPE_NAME
		prevTokensIndex := len(tokens) - 1
		if prevTokensIndex >= 0 {
			if tokens[prevTokensIndex].Type == TOKEN_TYPE_EQUAL {
				TokenType = TOKEN_TYPE_VALUE
			}
		}

		tokens = append(tokens, LibconfigTokenT{
			Type:  TokenType,
			Token: token,
		})
	}

	isArrayValue := 0
	for i := range tokens {
		if tokens[i].Type == TOKEN_TYPE_OPEN_SQUARE_BRACKET {
			isArrayValue++
		}
		if tokens[i].Type == TOKEN_TYPE_CLOSE_SQUARE_BRACKET {
			isArrayValue--
		}

		if isArrayValue == 1 && tokens[i].Type == TOKEN_TYPE_NAME {
			tokens[i].Type = TOKEN_TYPE_VALUE
		}
	}

	return tokens, err
}

func parseSettings(tokens []LibconfigTokenT) (configMap map[string]any, err error) {
	configMap = map[string]any{}

	tokensLen := len(tokens)
	diff := 0
	for i := 0; i < tokensLen; i += 1 + diff {
		diff = 0
		if tokens[i].Type == TOKEN_TYPE_NAME {
			if i+2 < tokensLen && tokens[i+1].Type == TOKEN_TYPE_EQUAL {
				var value any
				value, diff, err = parseSettingValue(tokens[i+2:])
				configMap[tokens[i].Token] = value
				if err != nil {
					return configMap, err
				}
				continue
			}

			err = fmt.Errorf("setting '%s' without value", tokens[i].Token)
			return configMap, err
		}
	}

	return configMap, err
}

func parseSettingValueArray(tokens []LibconfigTokenT) (value []string, diff int, err error) {
	for _, token := range tokens {
		if token.Type == TOKEN_TYPE_CLOSE_SQUARE_BRACKET {
			break
		}
		if token.Type == TOKEN_TYPE_VALUE {
			value = append(value, token.Token)
		}
		diff++
	}
	return value, diff, err
}

func parseSettingValueGroup(tokens []LibconfigTokenT) (groupSettings map[string]any, diff int, err error) {
	tokensLen := len(tokens)
	count := 0
	for diff = 0; diff < tokensLen; diff++ {
		if tokens[diff].Type == TOKEN_TYPE_OPEN_BRACKET {
			count++
		}

		if tokens[diff].Type == TOKEN_TYPE_CLOSE_BRACKET {
			count--
			if count <= 0 {
				break
			}
		}
	}

	if count > 0 {
		err = fmt.Errorf("unclose group")
		return groupSettings, diff, err
	}

	groupSettings, err = parseSettings(tokens[1:diff])

	return groupSettings, diff, err
}

func parseSettingValueList(tokens []LibconfigTokenT) (valueList []any, diff int, err error) {
	tokensLen := len(tokens)
	count := 0
	for diff = 0; diff < tokensLen; diff++ {
		if tokens[diff].Type == TOKEN_TYPE_OPEN_PAREN {
			count++
		}

		if tokens[diff].Type == TOKEN_TYPE_CLOSE_PAREN {
			count--
			if count <= 0 {
				break
			}
		}
	}

	if count > 0 {
		err = fmt.Errorf("unclose list")
		return valueList, diff, err
	}

	settingValueList := tokens[1:diff]
	settingValueListLen := len(settingValueList)
	valueDiff := 0
	for i := 0; i < settingValueListLen; i += 1 + valueDiff {
		valueDiff = 0

		var value any
		value, valueDiff, err = parseSettingValue(settingValueList[i:])
		if err != nil {
			return valueList, diff, err
		}
		valueList = append(valueList, value)
	}
	repr.Println(settingValueList, repr.Indent("  "), repr.OmitEmpty(true))

	return valueList, diff, err
}

func parseSettingValue(tokens []LibconfigTokenT) (value any, diff int, err error) {
	// repr.Println(tokens, repr.Indent("  "), repr.OmitEmpty(true))
	switch tokens[0].Type {
	case TOKEN_TYPE_OPEN_SQUARE_BRACKET:
		value, diff, err = parseSettingValueArray(tokens)
	case TOKEN_TYPE_OPEN_BRACKET:
		value, diff, err = parseSettingValueGroup(tokens)
	case TOKEN_TYPE_OPEN_PAREN:
		value, diff, err = parseSettingValueList(tokens)
	case TOKEN_TYPE_VALUE:
		diff = 0
		value = tokens[0].Token
	default:
		err = fmt.Errorf("fail in parsing process")
	}

	return value, diff, err
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
