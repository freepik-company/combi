package json

func (e *JsonT) ConfigToMap() (jsonMap map[string]interface{}) {
	jsonMap = e.ConfigStruct.(map[string]interface{})
	return jsonMap
}
