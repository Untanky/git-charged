package config

import "os"

type File interface {
	Has(key string) bool
	Get(key string) (string, bool)
}

type gitConfigFile struct {
	file      *os.File
	configMap map[string]string
}

func LoadFile(path string) (File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	configMap, err := ParseConfigFile(file)
	if err != nil {
		return nil, err
	}

	return gitConfigFile{
		file,
		configMap,
	}, nil
}

func (config gitConfigFile) Has(key string) bool {
	_, ok := config.configMap[key]
	return ok
}

func (config gitConfigFile) Get(key string) (string, bool) {
	value, ok := config.configMap[key]
	if !ok {
		return "", false
	}

	return value, true
}
