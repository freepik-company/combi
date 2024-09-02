package daemon

import (
	"context"

	"combi/api/v1alpha2"
	"combi/internal/conditions"
	"combi/internal/encoding"
)

type daemonT struct {
	context context.Context
	flags   daemonFlagsT
	config  v1alpha2.CombiConfigT
}

func (d *daemonT) mergeConfigurations(globalConfigBytes, configBytes []byte) (globalEncoder encoding.EncoderT, err error) {
	globalEncoder = encoding.GetEncoders()[d.config.Kind]
	configEncoder := encoding.GetEncoders()[d.config.Kind]

	err = globalEncoder.DecodeConfigBytes(globalConfigBytes)
	if err != nil {
		return globalEncoder, err
	}

	err = configEncoder.DecodeConfigBytes(configBytes)
	if err != nil {
		return globalEncoder, err
	}

	globalEncoder.MergeConfigs(configEncoder.GetConfigStruct())

	return globalEncoder, err
}

func (d *daemonT) evaluateConditions(targetConfigMap map[string]interface{}) (result bool, err error) {
	// Check global config conditions
	evalGlobalSuccess, err := conditions.EvalConditions(&d.config.Global.Conditions, &targetConfigMap)
	if err != nil {
		return result, err
	}

	// Check local config conditions
	result, err = conditions.EvalConditions(&d.config.Config.Conditions, &targetConfigMap)
	if err != nil {
		return result, err
	}

	result = result && evalGlobalSuccess

	return result, err
}
