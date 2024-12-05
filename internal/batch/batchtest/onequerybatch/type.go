package onequerybatch

type OneQueryConfig struct {
	Type    string           `yaml:"type"`
	Request *OneQueryRequest `yaml:"request"`
}

type OneQueryRequest struct {
	EndpointType  string              `yaml:"endpointType"`
	QueryParam    map[string][]string `yaml:"queryParam"`
	PathVariables map[string]string   `yaml:"pathVariables"`
	Outputs       []OneQueryVariable  `yaml:"outputs"`
}

type OneQueryVariable struct {
	ID       string `yaml:"id"`
	JMESPath string `yaml:"jmesPath"`
	OnError  string `yaml:"onError"`
}
