package dotmanager

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

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

func ReadConfig(filepath string) (DotConfig, error) {

  if !path.IsAbs(filepath) {
    wd, _ := os.Getwd();
    filepath = path.Join(wd, filepath);
  }

  var config DotConfig;
  _, err := toml.DecodeFile(filepath, &config.data);

  srcPath := path.Dir(filepath);

  for i, mod := range config.data.Modules {
    mod.rootSrc = &srcPath;
    mod.rootDst = &config.data.Dest 
    config.data.Modules[i] = mod;
  }

  return config, err;
} 
