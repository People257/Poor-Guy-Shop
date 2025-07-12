package conf

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

func MustLoad[T any](path string) (*koanf.Koanf, T) {
	k := koanf.New(".")
	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		panic(err)
	}
	var conf T
	if err := k.Unmarshal("", &conf); err != nil {
		panic(err)
	}
	return k, conf
}
