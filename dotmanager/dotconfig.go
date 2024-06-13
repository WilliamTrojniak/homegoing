package dotmanager

import (
	"errors"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
)

func (config *DotConfig) unmarshalFolderHelper(parentSrc string, parentDest string, data interface{}) error {

  var currentDest, currentSrc string;
  folderData, _ := data.(map[string]interface{});

  var ok bool;
  var folderDest string;
  if folderDest, ok = folderData["dest"].(string); !ok {
    folderDest = "";
  }
  folderDest = os.ExpandEnv(folderDest);

  if path.IsAbs(folderDest) {
    currentDest = path.Clean(folderDest);
  } else {
    currentDest = path.Join(parentDest, folderDest);
  }

  var folderSrc string;
  if folderSrc, ok = folderData["src"].(string); !ok {
    folderSrc = "";
  }
  folderSrc = os.ExpandEnv(folderSrc);


  if path.IsAbs(folderSrc) {
    currentSrc = path.Clean(folderSrc);
  } else {
    currentSrc = path.Join(parentSrc, folderSrc);
  }

  modules, ok := folderData["modules"].([]map[string]interface{});

  for _, module := range modules {
    var src, dest, name, target string;
    if src, ok = module["src"].(string); !ok {
      return errors.New("Failed to load module data: missing \"src\"");
    }
    src = os.ExpandEnv(src);

    if path.IsAbs(src) {
      src = path.Clean(src);
    } else {
      src = path.Join(currentSrc, src);
    }

    if name, ok = module["name"].(string); !ok {
      name = path.Base(src);
    }

    if target, ok = module["target"].(string); !ok {
      target = path.Base(src);
    }

    if dest, ok = module["dest"].(string); !ok {
      dest = "";
    }
    dest = os.ExpandEnv(dest);

    if path.IsAbs(dest) {
      dest = path.Join(dest, target);
    } else {
      dest = path.Join(currentDest, dest, target);
    }

    config.modules = append(config.modules, &DotModule{src: src, dest: dest, name: name, target: target});

  }

  var folders []map[string]interface{};
  if folders, ok = folderData["folders"].([]map[string]interface{}); !ok {
    return nil;
  }

  for _, data := range folders {
    if err := config.unmarshalFolderHelper(currentSrc, strings.Clone(currentDest), data); err != nil {
      return err;
    }
  }

  return nil;
}

type DotConfig struct {
  modules []*DotModule
}

type dotConfigData struct {
  Src string
  Dest string
  Groups []dotConfigData
  Modules []dotModuleData
}

func (config *DotConfig) GetModules() []*DotModule {
  out := make([]*DotModule, len(config.modules));
  copy(out, config.modules);
  return out;
}

func loadModulesFromConfigData(data *dotConfigData, parentSrc string, parentDest string) ([]*DotModule, error) {

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

  configData := dotConfigData{};
  _, err := toml.DecodeFile(filepath, &configData);
  if err != nil {
    return nil, err;
  }

  srcDir := path.Dir(filepath);
  modules, err := loadModulesFromConfigData(&configData, srcDir, srcDir);
  config := DotConfig{modules: modules};

  return &config, err;
} 
