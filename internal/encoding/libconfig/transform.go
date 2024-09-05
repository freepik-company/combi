package libconfig

func (e *LibconfigT) ConfigToMap() (configMap map[string]interface{}) {
	return e.ConfigStruct
}
