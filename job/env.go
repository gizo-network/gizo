package job

import (
	"encoding/json"

	"github.com/kpango/glg"
)

//EnvironmentVariables stores key and value of env variables
type EnvironmentVariable struct {
	Key   string
	Value string
}

func NewEnv(key, value string) *EnvironmentVariable {
	return &EnvironmentVariable{
		Key:   key,
		Value: value,
	}
}

func (env EnvironmentVariable) GetKey() string {
	return env.Key
}

func (env EnvironmentVariable) GetValue() string {
	return env.Value
}

type EnvironmentVariables []EnvironmentVariable

func (e EnvironmentVariables) Serialize() []byte {
	b, err := json.Marshal(e)
	if err != nil {
		glg.Fatal(err)
	}
	return b
}

func NewEnvVariables(variables ...EnvironmentVariable) EnvironmentVariables {
	return variables
}

func DeserializeEnvs(b []byte) EnvironmentVariables {
	var temp EnvironmentVariables
	err := json.Unmarshal(b, &temp)
	if err != nil {
		glg.Fatal(err)
	}
	return temp
}
