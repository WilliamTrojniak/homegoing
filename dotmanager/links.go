
package dotmanager

import (
	"errors"
	"os"
	"path"
	"strings"
)

type SymLink struct {
  // Path to file that acts as the symlink
  dest string
  // Path to the actual file
  src string
}

func (link SymLink) GetValue() string {
  return path.Base(link.dest);
}

func (link SymLink) GetId() string {
  return link.dest + " -> " + link.src;
}


func GetSymLinksInDir(absLinkDestDirPath string, absLinkSrcDirPath string) ([]SymLink, error) {
  if !path.IsAbs(absLinkDestDirPath) {
    return nil, errors.New("Link dest directory path must be absolute");
  }

  if !path.IsAbs(absLinkSrcDirPath) {
    return nil, errors.New("Link src directory path must be absolute.");
  }

  absLinkDestDirPath = path.Clean(absLinkDestDirPath);
  absLinkSrcDirPath = path.Clean(absLinkSrcDirPath);
  
  readRes, err := os.ReadDir(absLinkDestDirPath);
  if err != nil {
    return nil, err;
  }
  
  links := make([]SymLink, 0)

  for _, file := range readRes {
    
    // Check if the file is a symlink
    fileInfo, _ := file.Info();
    if fileInfo.Mode() & os.ModeSymlink != os.ModeSymlink {
      continue;
    }

    // Check if its linked from the src directory
    filePath := path.Join(absLinkDestDirPath, file.Name());
    srcPath, readLinkErr := os.Readlink(filePath);
    if readLinkErr == nil && strings.HasPrefix(srcPath, absLinkSrcDirPath) {
      // Add it to the output
      links = append(links, SymLink{src: srcPath, dest: filePath})

    }
  }
  return links, nil;
}
