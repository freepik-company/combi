package encoding

type EncoderT interface {
	// Encode/Decode configurations
	DecodeConfig()
	DecodeConfigBytes()
	EncodeConfigString()

	// Merge configurations
	MergeConfigs()

	// Transform configurations
	ConfigToMap()
}
