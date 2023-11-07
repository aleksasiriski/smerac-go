package config

import (
	"os"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func (c *Config) Load(dataDirPath string, logDirPath string) {
	// Use "." as the key path delimiter. This can be "/" or any character.
	k := koanf.New(".")

	// Load default values using the structs provider.
	// We provide a struct along with the struct tag `koanf` to the
	// provider.
	if err := k.Load(structs.Provider(&c, "koanf"), nil); err != nil {
		log.Panic().Err(err).Msg("failed loading default values")
	}

	// Load YAML config
	yamlPath := path.Join(dataDirPath, "smerac.yaml")
	if _, err := os.Stat(yamlPath); err != nil {
		log.Trace().Msgf("no yaml config present at path: %v, looking for .yml", yamlPath)
		yamlPath = path.Join(dataDirPath, "smerac.yml")
		if _, errr := os.Stat(yamlPath); errr != nil {
			log.Trace().Msgf("no yaml config present at path: %v", yamlPath)
		} else if errr := k.Load(file.Provider(yamlPath), yaml.Parser()); errr != nil {
			log.Panic().Err(err).Msg("error loading yaml config")
		}
	} else if err := k.Load(file.Provider(yamlPath), yaml.Parser()); err != nil {
		log.Panic().Err(err).Msg("error loading yaml config")
	}

	// Load ENV config
	if err := k.Load(env.Provider("SMERAC_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "SMERAC_")), "_", ".", -1)
	}), nil); err != nil {
		log.Panic().Err(err).Msg("error loading env config")
	}

	// Unmarshal config into struct
	if err := k.Unmarshal("", &c); err != nil {
		log.Panic().Err(err).Msg("failed unmarshaling koanf config")
	}
}
