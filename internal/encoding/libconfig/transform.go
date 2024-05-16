package libconfig

func ConfigToMap(config *LIBCONFIG) (libconfigMap map[string]interface{}) {
	libconfigMap = make(map[string]interface{})
	configSettingsToMap(&config.Settings, &libconfigMap)
	return libconfigMap
}

func configSettingsToMap(settings *[]*SettingT, configMap *map[string]interface{}) {
	for _, setting := range *settings {
		if setting.SettingValue.Primitive != nil {
			(*configMap)[setting.SetingName] = setting.SettingValue.Primitive.Value
		}

		if setting.SettingValue.Array != nil {
			tmpArray := []string{}
			for _, item := range setting.SettingValue.Array.Primitives {
				tmpArray = append(tmpArray, item.Value)
			}
			(*configMap)[setting.SetingName] = tmpArray
		}

		if setting.SettingValue.Group != nil {
			tmpConfigMap := map[string]interface{}{}
			configSettingsToMap(&setting.SettingValue.Group.Settings, &tmpConfigMap)
			(*configMap)[setting.SetingName] = tmpConfigMap
		}

		if setting.SettingValue.List != nil {
			var tmpValues []interface{}
			configListToMap(&setting.SettingValue.List.SettingValues, &tmpValues)
			(*configMap)[setting.SetingName] = tmpValues
		}
	}
}

func configListToMap(settingValues *[]*SettingValueT, values *[]interface{}) {
	for _, sv := range *settingValues {
		if sv.Primitive != nil {
			*values = append(*values, sv.Primitive.Value)
		}

		if sv.Array != nil {
			tmpArray := []string{}
			for _, item := range sv.Array.Primitives {
				tmpArray = append(tmpArray, item.Value)
			}
			*values = append(*values, tmpArray)
		}

		if sv.Group != nil {
			tmpConfigMap := map[string]interface{}{}
			configSettingsToMap(&sv.Group.Settings, &tmpConfigMap)
			*values = append(*values, tmpConfigMap)
		}

		if sv.List != nil {
			var tmpValues []interface{}
			configListToMap(&sv.List.SettingValues, &tmpValues)
			*values = append(*values, tmpValues)
		}
	}
}
