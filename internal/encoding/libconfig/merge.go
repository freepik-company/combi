package libconfig

// ----------------------------------------------------------------
// Merge LIBCONFIG data structure
// ----------------------------------------------------------------

func MergeConfigs(destination *LIBCONFIG, source *LIBCONFIG) {
	mergeSettings(&destination.Settings, &source.Settings)
}

func mergeSettings(destination *[]*SettingT, source *[]*SettingT) {
	for _, sSetting := range *source {
		// Merge settings with primitive
		foundSetting := false
		if sSetting.SettingValue.Primitive != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					dSetting.SettingValue.Primitive.Value = sSetting.SettingValue.Primitive.Value
				}
			}
		}

		// Merge settings with group
		if sSetting.SettingValue.Group != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueGroups(dSetting.SettingValue.Group, sSetting.SettingValue.Group)
				}
			}
		}

		// Merge settings with array
		if sSetting.SettingValue.Array != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueArrays(dSetting.SettingValue.Array, sSetting.SettingValue.Array)
				}
			}
		}

		// Merge settings with list
		if sSetting.SettingValue.List != nil {
			for _, dSetting := range *destination {
				if dSetting.SetingName == sSetting.SetingName {
					foundSetting = true
					mergeSettingValueLists(dSetting.SettingValue.List, sSetting.SettingValue.List)
				}
			}
		}

		// Append not found setting
		if !foundSetting {
			*destination = append(*destination, sSetting)
		}
	}
}

func mergeSettingValueArrays(destination *ArrayT, source *ArrayT) {
	if len(source.Primitives) > 0 {
		for _, sPrimitive := range source.Primitives {
			found := false
			for _, dPrimitive := range destination.Primitives {
				if dPrimitive.Value == sPrimitive.Value {
					found = true
				}
			}
			if !found {
				destination.Primitives = append(destination.Primitives, sPrimitive)
			}
		}
	}
}

func mergeSettingValueGroups(destination *GroupT, source *GroupT) {
	if len(source.Settings) > 0 {
		mergeSettings(&destination.Settings, &source.Settings)
	}
}

func mergeSettingValueLists(destination *ListT, source *ListT) {
	if len(source.SettingValues) > 0 {
		for _, sSettingValue := range source.SettingValues {
			found := false
			if sSettingValue.Primitive != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Primitive != nil && sSettingValue.Primitive.Value != dSettingValue.Primitive.Value {
						found = true
						dSettingValue.Primitive.Value = sSettingValue.Primitive.Value
					}
				}
			}

			if sSettingValue.Array != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Array != nil {
						found = true
						mergeSettingValueArrays(dSettingValue.Array, sSettingValue.Array)
					}
				}
			}

			if sSettingValue.Group != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.Group != nil {
						found = true
						mergeSettingValueGroups(dSettingValue.Group, sSettingValue.Group)
					}
				}
			}

			if sSettingValue.List != nil {
				for _, dSettingValue := range destination.SettingValues {
					if dSettingValue.List != nil {
						found = true
						mergeSettingValueLists(dSettingValue.List, sSettingValue.List)
					}
				}
			}

			if !found {
				destination.SettingValues = append(destination.SettingValues, sSettingValue)
			}
		}
	}
}
