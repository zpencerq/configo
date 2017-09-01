package configo

type ConfigoChain struct {
	configos []Configo
}

func NewConfigoChain(configos ...Configo) *ConfigoChain {
	return &ConfigoChain{configos: configos}
}

func NewDefaultConfigoChain(file string) *ConfigoChain {
	configos := []Configo{
		NewDefaultsConfigo(),
		NewTomlConfigo(file),
		NewEnvConfigo(),
	}
	return &ConfigoChain{configos: configos}
}

func (chain *ConfigoChain) Load(v interface{}) error {
	var err error
	for _, c := range chain.configos {
		err = c.Load(v)
		if err != nil {
			return err
		}
	}
	return nil
}
