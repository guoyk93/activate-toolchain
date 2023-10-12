package activate_toolchain

import (
	"context"
	"errors"
	"github.com/guoyk93/activate-toolchain/pkg/unarchive"
	"log"
	"os"
	"path/filepath"
)

type InstallArchiveOptions struct {
	URLs           []string
	Filename       string
	Name           string
	StripDirectory bool
}

// InstallArchive installs an archive to a directory
func InstallArchive(ctx context.Context, opts InstallArchiveOptions) (dir string, err error) {
	var home string
	if home, err = os.UserHomeDir(); err != nil {
		return
	}
	base := filepath.Join(home, ".atc")
	if err = os.MkdirAll(base, 0755); err != nil {
		return
	}

	dir = filepath.Join(base, opts.Name)

	// check if already installed
	{
		var stat os.FileInfo
		if stat, err = os.Stat(dir); err == nil {
			if stat.IsDir() {
				return
			}
			err = errors.New("target directory is not a directory: " + dir)
			return
		} else {
			if !os.IsNotExist(err) {
				return
			}
			err = nil
		}
	}

	// temp directory for atomic installation
	dirTemp := filepath.Join(base, opts.Name+".tmp")
	os.RemoveAll(dirTemp)
	if err = os.MkdirAll(dirTemp, 0755); err != nil {
		return
	}
	defer os.RemoveAll(dirTemp)

	file := filepath.Join(base, opts.Filename)

	// ensure file exists
	{
		var stat os.FileInfo
		if stat, err = os.Stat(file); err == nil {
			if stat.IsDir() {
				err = errors.New("target file is a directory: " + file)
				return
			}
		} else {
			if !os.IsNotExist(err) {
				return
			}
			err = nil

			log.Println("downloading:", opts.Filename)

			if err = AdvancedFetchFile(ctx, opts.URLs, file); err != nil {
				return
			}
		}
	}

	log.Println("extracting:", opts.Filename)

	// extract file
	{
		var f *os.File
		if f, err = os.OpenFile(file, os.O_RDONLY, 0644); err != nil {
			return
		}

		if err = unarchive.Unarchive(f, dirTemp); err != nil {
			_ = f.Close()
			return
		}

		_ = f.Close()
	}

	dirInstall := dirTemp

	if opts.StripDirectory {
		var dirs []os.DirEntry
		if dirs, err = os.ReadDir(dirTemp); err != nil {
			return
		}
		for _, dir := range dirs {
			if dir.IsDir() {
				dirInstall = filepath.Join(dirTemp, dir.Name())
				break
			}
		}
	}

	os.RemoveAll(dir)

	if err = os.Rename(dirInstall, dir); err != nil {
		return
	}

	log.Println("installed:", opts.Name)

	os.RemoveAll(file)

	return
}
