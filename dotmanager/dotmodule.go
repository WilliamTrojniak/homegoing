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

type dotModule struct {
	data dotModuleData
}

func (mod *dotModule) GetDest() string {
	dest := path.Join(os.ExpandEnv(mod.data.Dest), mod.getTarget())
	if path.IsAbs(dest) {
		return dest
	}
	return path.Join(os.ExpandEnv(*mod.data.rootDst), dest)

}

func (mod *dotModule) GetSrc() string {
	src := os.ExpandEnv(mod.data.Src)
	if path.IsAbs(src) {
		return path.Clean(src)
	}
	return path.Join(*mod.data.rootSrc, src)
}

func (mod *dotModule) GetName() string {
	if len(mod.data.Name) == 0 {
		return path.Base(mod.GetSrc())
	}

	return mod.data.Name
}

func (mod *dotModule) getTarget() string {
	if len(mod.data.Target) == 0 {
		return path.Base(mod.GetSrc())
	}
	return mod.data.Target
}
