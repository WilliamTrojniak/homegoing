package dotmanager

import (
	"os"
	"path"
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

type DotModule struct {
	data dotModuleData
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
	dest := path.Join(os.ExpandEnv(mod.data.Dest), mod.getTarget())
	if path.IsAbs(dest) {
		return dest
	}
	return path.Join(os.ExpandEnv(*mod.data.rootDst), dest)

}

func (mod *DotModule) GetSrc() string {
	src := os.ExpandEnv(mod.data.Src)
	if path.IsAbs(src) {
		return path.Clean(src)
	}
	return path.Join(*mod.data.rootSrc, src)
}

func (mod *DotModule) GetName() string {
	if len(mod.data.Name) == 0 {
		return path.Base(mod.GetSrc())
	}

	return mod.data.Name
}

func (mod *DotModule) getTarget() string {
	if len(mod.data.Target) == 0 {
		return path.Base(mod.GetSrc())
	}
	return mod.data.Target
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
