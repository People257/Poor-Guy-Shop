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
	// 使用UnmarshalWithConf来支持mapstructure标签
	unmarshalConfig := koanf.UnmarshalConf{
		Tag:       "mapstructure",
		FlatPaths: false,
	}
	if err := k.UnmarshalWithConf("", &conf, unmarshalConfig); err != nil {
		panic(err)
	}
	return k, conf
}
