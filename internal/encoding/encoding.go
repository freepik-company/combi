package encoding

import (
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
)

const (
	libconfigEncoderKey = `libconfig`
	jsonEncoderKey      = `json`
)

type EncoderT interface {
	// Encode/Decode configurations
	DecodeConfig(filepath string) (err error)
	DecodeConfigBytes(configBytes []byte) (err error)
	EncodeConfigString() (configStr string)

	// Merge configurations
	MergeConfigs(source interface{})
	GetConfigStruct() (config interface{})

	// Transform configurations
	ConfigToMap() (configMap map[string]interface{})
}

func GetEncoders() (encoders map[string]EncoderT) {
	encoders = map[string]EncoderT{
		libconfigEncoderKey: &libconfig.LibconfigT{},
		jsonEncoderKey:      &json.JsonT{},
	}
	return encoders
}
