package dotmanager

import (
	"os"
	"path"
)

type dotModuleData struct {
  Dest string
  Src string
  Name string
  Target string
}

type DotModule struct {
	// The path to the directory where the config should be linked to.
	// Can be an absolute path or relative to the dest defined by the dotmanager
	dest string

	// The path to the file to be linked. Relative to the dotfiles.toml config file
	src string

	// The name of the module, optionally specified in the config file
	name string

	// The filename of the created link, optionally specified in the config file
	target string
}

type LinkStatus int8

const (
  LINK_STATUS_UNLINKED LinkStatus = iota
  LINK_STATUS_EXISTS_CONFLICT
  LINK_STATUS_TARGET_CONFLICT
  LINK_STATUS_LINKED
)

func (status LinkStatus) String() string {
  switch status {
  case LINK_STATUS_UNLINKED:
    return "Unlinked";
  case LINK_STATUS_EXISTS_CONFLICT:
    return "A file or directory at the destination already exists";
  case LINK_STATUS_TARGET_CONFLICT:
    return "A symbolic link at the destination already exists";
  case LINK_STATUS_LINKED:
    return "Linked";
  }
  return "Unkown";
}

func (mod *DotModule) GetDest() string {
	return mod.dest;
}

func (mod *DotModule) GetSrc() string {
  return mod.src;
}

func (mod *DotModule) GetName() string {
  return mod.name;
}

func (mod *DotModule) GetTarget() string {
  return mod.target;
}

func (mod *DotModule) GetLinkStatus() (LinkStatus, bool) {
  destPath := mod.GetDest();
  srcPath := mod.GetSrc();

  destFile, destErr := os.Lstat(destPath);
  if destErr != nil {
    return LINK_STATUS_UNLINKED, false;
  }

  if destFile.Mode() & os.ModeSymlink != os.ModeSymlink {
    return LINK_STATUS_EXISTS_CONFLICT, false;
  }
  
  linkPath, linkErr := os.Readlink(destPath);
  if linkErr != nil || linkPath != srcPath {
    return LINK_STATUS_TARGET_CONFLICT, false;
  }

  return LINK_STATUS_LINKED, true;
}

func (mod *DotModule) LinkModule(force bool) error {
  destPath := mod.GetDest();
  srcPath := mod.GetSrc();

  if err := os.MkdirAll(path.Dir(destPath), 0700); err != nil {
    return err;
  }
  
  err := os.Symlink(srcPath, destPath);

  if !force || err == nil {
    return err;
  }

  if !os.IsExist(err) {
    return err;
  }

  os.RemoveAll(destPath);
  return os.Symlink(srcPath, destPath)

}

func (mod *DotModule) UnlinkModule() error {
  if _, isLinked := mod.GetLinkStatus(); !isLinked {
    return nil;
  }
  return os.Remove(mod.GetDest());
}
