package models

type Config struct {
	Data struct {
		File    string `yaml:"file"`
		Headers bool   `yaml:"headers"`
	} `yaml:"data"`

	Processor struct {
		Request struct {
			Method  string            `yaml:"method"`
			BaseUrl string            `yaml:"baseUrl"`
			Uri     string            `yaml:"uri"`
			Body    string            `yaml:"body"`
			Headers map[string]string `yaml:"headers"`
		} `yaml:"request"`
	} `yaml:"processor"`

	Result struct {
		File         string `yaml:"file"`
		Append       bool   `yaml:"append_results"`
		WriteSuccess bool   `yaml:"write_success"`
		WriteFails   bool   `yaml:"write_fails"`
	} `yaml:"result"`

	Executor struct {
		Type     string `yaml:"type"`
		Routines int    `yaml:"routines"`
		Flush    int    `yaml:"flush_number"`
	} `yaml:"executor"`

	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
}
