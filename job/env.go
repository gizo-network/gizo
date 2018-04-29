package job

//EnvironmentVariables stores key and value of env variables
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

func NewEnvVariables(variables ...EnvironmentVariable) EnvironmentVariables {
	return variables
}
