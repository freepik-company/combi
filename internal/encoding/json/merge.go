package json

// ----------------------------------------------------------------
// Merge JSON data structure
// ----------------------------------------------------------------
func (e *JsonT) GetConfigStruct() (config interface{}) {
	return e.ConfigStruct
}

func (e *JsonT) MergeConfigs(source interface{}) {
	mergeMaps(e.ConfigStruct.(map[string]interface{}), source.(map[string]interface{}))
}

func mergeMaps(destination map[string]interface{}, source map[string]interface{}) {

}
