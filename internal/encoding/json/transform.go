package json

func (e *JsonT) ConfigToMap() (configMap map[string]interface{}) {
	configMap = e.ConfigStruct.(map[string]interface{})
	return configMap
}
