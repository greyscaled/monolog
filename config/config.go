package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"time"
)

// configFileName is the name of the monolog configuration file as stored on
// disk
const configFileName = ".monologconfig"

// File represents the monolog data persisted in a configuration file. It
// is read from or created if one does not exist on each run of the program.
type File struct {
	Latest time.Time
	Path   string
}

// New returns a ConfigFile with values as stored on disk.
func New() (*File, error) {
	pth, err := buildCfgPath()
	if err != nil {
		return nil, err
	}

	exists, err := fileExists(pth)
	if err != nil {
		return nil, err
	}

	if exists {
		f, err := os.Open(pth)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		cfgRaw, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		var cf File
		err = json.Unmarshal(cfgRaw, &cf)
		if err != nil {
			return nil, err
		}
		return &cf, nil
	}

	cf := File{
		Latest: time.Now().AddDate(0, 0, -1),
	}
	return &cf, nil
}

// Save persists the File to disk
func (cf *File) Save() error {
	cfgRaw, err := json.Marshal(&cf)
	if err != nil {
		return err
	}

	pth, err := buildCfgPath()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(pth, cfgRaw, 0644)

	return err
}

// buildCfgPath returns a string representing the disk path to configFileName
func buildCfgPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	// ~/.monologconfig
	pth := path.Join(u.HomeDir, configFileName)

	return pth, nil
}

// fileExists returns true if the pth exists, false if it doesn't and a
// non-nil error if existence cannot be determined
func fileExists(pth string) (bool, error) {
	_, err := os.Stat(pth)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
