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
			var tmpValue interface{}
			configListToMap(&setting.SettingValue.List.SettingValues, &tmpValue)
			(*configMap)[setting.SetingName] = tmpValue
		}
	}
}

func configListToMap(settingValues *[]*SettingValueT, value *interface{}) {
	for _, sv := range *settingValues {
		if sv.Primitive != nil {
			*value = sv.Primitive.Value
		}

		if sv.Array != nil {
			tmpArray := []string{}
			for _, item := range sv.Array.Primitives {
				tmpArray = append(tmpArray, item.Value)
			}
			*value = tmpArray
		}

		if sv.Group != nil {
			tmpConfigMap := map[string]interface{}{}
			configSettingsToMap(&sv.Group.Settings, &tmpConfigMap)
			*value = tmpConfigMap
		}

		if sv.List != nil {
			var tmpValue interface{}
			configListToMap(&sv.List.SettingValues, &tmpValue)
			*value = tmpValue
		}
	}
}
