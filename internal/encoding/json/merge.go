package json

import (
	"combi/internal/globals"
)

// ----------------------------------------------------------------
// Merge JSON data structure
// ----------------------------------------------------------------
func (e *JsonT) GetConfigStruct() (config interface{}) {
	return e.ConfigStruct
}

func (e *JsonT) MergeConfigs(source interface{}) {
	mergeJsonObjects(e.ConfigStruct.(map[string]any), source.(map[string]any))
}

func mergeJsonObjects(destination, source map[string]any) {
	for srcKey, srcVal := range source {

		if _, ok := destination[srcKey]; !ok {
			destination[srcKey] = srcVal
			continue
		}

		switch destination[srcKey].(type) {
		case float64, string, bool, nil:
			destination[srcKey] = srcVal
		case []any:
			mergeJsonArray(destination[srcKey].([]any), srcVal.([]any))
		case map[string]any:
			mergeJsonObjects(destination[srcKey].(map[string]any), srcVal.(map[string]any))
		default:
			globals.ExecContext.Logger.Debugf("invalid json type\n")
		}
	}
}

func mergeJsonArray(destination, source []interface{}) {
	gap := len(source) - len(destination)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			destination = append(destination, nil)
		}
	}
	for srcIndex, srcVal := range source {
		switch srcVal.(type) {
		case float64, string, bool, nil:
			destination[srcIndex] = srcVal
		case []any:
			{
				if destination[srcIndex] == nil {
					destination[srcIndex] = []any{}
				}
				mergeJsonArray(destination[srcIndex].([]any), srcVal.([]any))
			}
		case map[string]any:
			{
				if destination[srcIndex] == nil {
					destination[srcIndex] = map[string]any{}
				}
				mergeJsonObjects(destination[srcIndex].(map[string]any), srcVal.(map[string]any))
			}
		default:
			globals.ExecContext.Logger.Debugf("invalid json type\n")
		}
	}
}
