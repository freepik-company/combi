package daemon

import (
	"combi/api/v1alpha1"
	"combi/internal/conditions"
	"combi/internal/encoding"
)

func mergeConfigurations(combiConfig v1alpha1.GCMerge, configField string) (targetEncoder encoding.EncoderT, err error) {
	targetEncoder = encoding.GetEncoders()[combiConfig.Kind]
	localRawEncoder := encoding.GetEncoders()[combiConfig.Kind]
	globalRawEncoder := encoding.GetEncoders()[combiConfig.Kind]

	err = targetEncoder.DecodeConfig(combiConfig.Configs[configField].TargetConfig)
	if err != nil {
		return targetEncoder, err
	}
	err = localRawEncoder.DecodeConfigBytes([]byte(combiConfig.Configs[configField].RawConfig))
	if err != nil {
		return targetEncoder, err
	}
	targetEncoder.MergeConfigs(localRawEncoder.GetConfigStruct())

	err = globalRawEncoder.DecodeConfigBytes([]byte(combiConfig.Global.RawConfig))
	if err != nil {
		return targetEncoder, err
	}
	targetEncoder.MergeConfigs(globalRawEncoder.GetConfigStruct())

	return targetEncoder, err
}

func evaluateConditions(combiConfig v1alpha1.GCMerge, configField string, targetConfigMap map[string]interface{}) (result bool, err error) {
	// Check local config conditions
	localConds := combiConfig.Configs[configField].Conditions
	result, err = conditions.EvalConditions(&localConds, &targetConfigMap)
	if err != nil {
		return result, err
	}

	// Check global config conditions
	evalGlobalSuccess, err := conditions.EvalConditions(&combiConfig.Global.Conditions, &targetConfigMap)
	if err != nil {
		return result, err
	}
	result = result && evalGlobalSuccess

	return result, err
}
