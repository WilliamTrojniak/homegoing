package dotmanager

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)


type dotModuleData struct {
  // The path to the directory where the config should be linked to.
  // Can be an absolute path or relative to the dest defined by the dotmanager
  Dest string

  // The path to the file to be linked. Relative to the dotfiles.toml config file
  Src string

  // The path to the root to link, optionally specified in the config file
  rootDst *string

  // The path to the directory of where the config file was loaded from 
  rootSrc *string

  // The name of the module, optionally specified in the config file
  Name string

  // The filename of the created link, optionally specified in the config file
  Target string

}

type dotModule struct {
  data dotModuleData
}

func (mod *dotModule) GetDest() string {
  dest := path.Join(os.ExpandEnv(mod.data.Dest), mod.getTarget());
  if path.IsAbs(dest) {
    return dest;
  }
  return path.Join(os.ExpandEnv(*mod.data.rootDst), dest);

}

func (mod *dotModule) GetSrc() string {
  src := os.ExpandEnv(mod.data.Src);
  if path.IsAbs(src) {
    return path.Clean(src);
  }
  return path.Join(*mod.data.rootSrc, src);
}

func (mod *dotModule) GetName() string {
  if len(mod.data.Name) == 0 {
    return path.Base(mod.GetSrc());
  } 

  return mod.data.Name;
}

func (mod *dotModule) getTarget() string {
  if len(mod.data.Target) == 0 {
    return path.Base(mod.GetSrc());
  }
  return mod.data.Target;
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


func ReadConfig(filepath string) (DotConfig, error) {

  // TODO Add check that path is absolute

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
