package json

import (
	"encoding/json"
	"os"
)

type JsonT struct {
	ConfigStruct interface{}
}

// ----------------------------------------------------------------
// Decode/Encode JSON data structure
// ----------------------------------------------------------------

// Decode functions

func (e *JsonT) DecodeConfig(filepath string) (err error) {
	configBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = e.DecodeConfigBytes(configBytes)
	return err
}

func (e *JsonT) DecodeConfigBytes(configBytes []byte) (err error) {
	err = json.Unmarshal(configBytes, &e.ConfigStruct)
	return err
}

// Encode functions

func (e *JsonT) EncodeConfigString() (configStr string) {
	configBytes, _ := json.MarshalIndent(e.ConfigStruct, "", "  ")
	configStr = string(configBytes)
	return configStr
}
