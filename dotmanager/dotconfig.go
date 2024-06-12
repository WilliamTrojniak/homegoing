package dotmanager

import (
	"os"

	"github.com/BurntSushi/toml"
)


type dotModuleData struct {
  // The path to the directory where the config should be linked to.
  // Can be an absolute path or relative to the dest defined by the dotmanager
  Dest string

  // The path to the file to be linked. Relative to the dotfiles.toml config file
  Src string

  // The name of the created link
  Link string
}

type dotModule struct {
  data dotModuleData
}

func (mod *dotModule) GetDest() string {
  return os.ExpandEnv(mod.data.Dest);
}

type dotConfigData struct {
  // The path to the directory where the config should be linked to.
  Dest string `toml:"dest"`

  // Available configs that can be linked
  Modules []dotModuleData `toml:"modules"`

}

type DotConfig struct {
  data dotConfigData
}

func (config *DotConfig) GetRootDest() string {
  return os.ExpandEnv(config.data.Dest);
}

func (config *DotConfig) GetModules() []dotModule {
  data := make([]dotModule, 0, len(config.data.Modules));
  for _, mod := range config.data.Modules {
    data = append(data, dotModule{data: mod});
  }
  return data;
}


func ReadConfig(path string) (DotConfig, error) {
  var config DotConfig;
  _, err := toml.DecodeFile(path, &config.data);

  return config, err;
} 
