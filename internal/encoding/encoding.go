package encoding

import (
	"combi/internal/encoding/json"
	"combi/internal/encoding/libconfig"
	"combi/internal/encoding/nginx"
)

const (
	jsonEncoderKey      = `json`
	nginxEncoderKey     = `nginx`
	libconfigEncoderKey = `libconfig`
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
		jsonEncoderKey:      &json.JsonT{},
		nginxEncoderKey:     &nginx.NginxT{},
		libconfigEncoderKey: &libconfig.LibconfigT{},
	}
	return encoders
}
