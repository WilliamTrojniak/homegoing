package dotmanager

import (
	"errors"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type DotConfig struct {
  modules []*DotModule
}

type dotConfigGroupData struct {
  Src string
  Dest string
  Groups []dotConfigGroupData
  Modules []dotConfigModuleData
}

type dotConfigModuleData struct {
  Dest string
  Src string
  Name string
  Target string
}

func (config *DotConfig) GetModules() []*DotModule {
  out := make([]*DotModule, len(config.modules));
  copy(out, config.modules);
  return out;
}

func (config *DotConfig) GetNumModules() int {
  return len(config.modules);
}

func loadModulesFromConfigData(data *dotConfigGroupData, parentSrc string, parentDest string) ([]*DotModule, error) {

  var groupSrc, groupDest string = os.ExpandEnv(data.Src), os.ExpandEnv(data.Dest);

  if path.IsAbs(groupSrc) {
    groupSrc = path.Clean(groupSrc);
  } else {
    groupSrc = path.Join(parentSrc, groupSrc);
  }

  if path.IsAbs(groupDest) {
    groupDest = path.Clean(groupDest);
  } else {
    groupDest = path.Join(parentDest, groupDest);
  }


  modules := make([]*DotModule, 0, len(data.Modules));
  for _, module := range data.Modules {

    // Module Validation
    if module.Src == "" {
      return modules, errors.New("Error loading modules from config: Missing field \"src\" on module object.");
    }

    var src, dest, name, target string = module.Src, module.Dest, module.Name, module.Target;
    src, dest = os.ExpandEnv(src), os.ExpandEnv(dest);

    if path.IsAbs(src) {
      src = path.Clean(src);
    } else {
      src = path.Join(groupSrc, src);
    }

    if name == "" {
      name = path.Base(src);
    }

    if target == "" {
      target = path.Base(src);
    }

    if path.IsAbs(dest) {
      dest = path.Join(dest, target);
    } else {
      dest = path.Join(groupDest, dest, target);
    }

    modules = append(modules, &DotModule{src: src, dest: dest, name: name, target: target});

  }

  for _, group := range data.Groups {
    submodules, err := loadModulesFromConfigData(&group, groupSrc, groupDest);
    modules = append(modules, submodules...);
    if err != nil {
      return modules, err;
    }
  }

  return modules, nil;

}

func LoadConfig(filepath string) (*DotConfig, error) {

  if !path.IsAbs(filepath) {
    wd, _ := os.Getwd();
    filepath = path.Join(wd, filepath);
  }

  configData := dotConfigGroupData{};
  _, err := toml.DecodeFile(filepath, &configData);
  if err != nil {
    return nil, err;
  }

  srcDir := path.Dir(filepath);
  modules, err := loadModulesFromConfigData(&configData, srcDir, srcDir);
  config := DotConfig{modules: modules};

  return &config, err;
} 
