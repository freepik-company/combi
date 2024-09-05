package libconfig

import "combi/internal/logger"

// ----------------------------------------------------------------
// Merge LIBCONFIG data structure
// ----------------------------------------------------------------
func (e *LibconfigT) GetConfigStruct() (config any) {
	return e.ConfigStruct
}

func (e *LibconfigT) MergeConfigs(source any) {
	mergeSettings(e.ConfigStruct, source.(map[string]any))
}

func mergeSettings(destination, source map[string]any) {
	for srcKey, srcVal := range source {
		if _, ok := destination[srcKey]; !ok {
			destination[srcKey] = srcVal
			continue
		}

		switch destination[srcKey].(type) {
		case string, []string:
			destination[srcKey] = srcVal
		// case []string:
		// 	mergeSettingArray(destination[srcKey].([]string), srcVal.([]string))
		case []any:
			tmp := destination[srcKey].([]any)
			mergeSettingList(&tmp, srcVal.([]any))
			destination[srcKey] = tmp
		case map[string]any:
			mergeSettings(destination[srcKey].(map[string]any), srcVal.(map[string]any))
		default:
			logger.Log.Debugf("invalid libconfig type\n")
		}
	}
}

// func mergeSettingArray(destination, source []string) {
// 	gap := len(source) - len(destination)
// 	if gap > 0 {
// 		for i := 0; i < gap; i++ {
// 			destination = append(destination, "")
// 		}
// 	}
// 	for srcIndex, srcVal := range source {
// 		destination[srcIndex] = srcVal
// 	}
// }

func mergeSettingList(destination *[]any, source []any) {
	gap := len(source) - len(*destination)
	if gap > 0 {
		for i := 0; i < gap; i++ {
			*destination = append(*destination, nil)
		}
	}
	for srcIndex, srcVal := range source {
		switch srcVal.(type) {
		case string, []string, []any:
			(*destination)[srcIndex] = srcVal
		// case []any:
		// 	{
		// 		if (*destination)[srcIndex] == nil {
		// 			(*destination)[srcIndex] = []any{}
		// 		}
		// 		mergeSettingList((*destination)[srcIndex].([]any), srcVal.([]any))
		// 	}
		case map[string]any:
			{
				if (*destination)[srcIndex] == nil {
					(*destination)[srcIndex] = map[string]any{}
				}
				mergeSettings((*destination)[srcIndex].(map[string]any), srcVal.(map[string]any))
			}
		default:
			logger.Log.Debugf("invalid libconfig type\n")
		}
	}
}
