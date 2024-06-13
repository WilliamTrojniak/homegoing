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


func (config *DotConfig) UnmarshalTOML(data interface{}) error {

  return config.unmarshalFolderHelper(config.GetConfigRepoPath(), config.GetConfigRepoPath(), data);

}

type DotConfig struct {
  filepath string
  modules []*DotModule
}

func (config *DotConfig) GetConfigRepoPath() string {
  return path.Dir(config.filepath);
}

func (config *DotConfig) GetModules() []*DotModule {
  out := make([]*DotModule, len(config.modules));
  copy(out, config.modules);
  return out;
}

func ReadConfig(filepath string) (*DotConfig, error) {

  if !path.IsAbs(filepath) {
    wd, _ := os.Getwd();
    filepath = path.Join(wd, filepath);
  }

  config := DotConfig{filepath: filepath};
  _, err := toml.DecodeFile(filepath, &config);

  return &config, err;
} 
