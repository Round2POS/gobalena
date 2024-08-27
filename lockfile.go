package gobalena

import (
	"os"
	"path/filepath"
)

const (
	BalenaLockFile = "/tmp/balena/updates.lock"
)

func Unlock(file string) error {
	if _, err := os.Stat(file); err == nil {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}

func Lock(file string) error {
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, os.ModeExclusive)
	if err != nil {
		return err
	}

	defer f.Close()

	return nil
}
