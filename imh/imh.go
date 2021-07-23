package imh

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/e-structures/shared"
	"os"
	"path/filepath"
	"time"
)

type Modulator interface {
	shared.Initializer
	shared.Backupper
	shared.Loader
	shared.Saver
	shared.Printer
	shared.Validator
	shared.Upgrader
	shared.WithUnused
}

const (
	configFileName      = "config.json"
	stateFileName       = "state.json"
	backupDirectoryName = "backup"
)

var ncverr shared.NotCurrentVersionError

type InfrastructureModuleHelper struct {
	ModuleDirectoryPath string
	ModuleVersion       string
}

func (h InfrastructureModuleHelper) Initialize(config Modulator, state Modulator) (Modulator, Modulator, error) {
	// TODO test
	// check if required fields are set
	if h.ModuleDirectoryPath == "" {
		return nil, nil, fmt.Errorf("setup module directory path first")
	}
	if h.ModuleVersion == "" {
		return nil, nil, fmt.Errorf("setup module version first")
	}

	// ensure module directory
	err := os.MkdirAll(h.ModuleDirectoryPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	// create default files paths
	configFilePath := filepath.Join(h.ModuleDirectoryPath, configFileName)
	stateFilePath := filepath.Join(h.ModuleDirectoryPath, stateFileName)

	// load state file
	err = state.Load(stateFilePath)
	if os.IsNotExist(err) {
		// if no state loaded then init it
		state.Init(h.ModuleVersion)
	} else if errors.As(err, &ncverr) {
		// if old version was found try to upgrade it
		err2 := state.Upgrade(stateFilePath)
		if err2 != nil {
			return nil, nil, err2
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("load state failed: %v", err)
	}
	// load config file
	err = config.Load(configFilePath)
	if os.IsNotExist(err) {
		// if no config loaded then init it
		config.Init(h.ModuleVersion)
	} else if errors.As(err, &ncverr) {
		// if old version was found try to upgrade it
		err2 := config.Upgrade(configFilePath)
		if err2 != nil {
			return nil, nil, err2
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("load config failed: %v", err)
	}

	// backup
	err = backup(h.ModuleDirectoryPath, config, state)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func (h InfrastructureModuleHelper) Load(config Modulator, state Modulator) (Modulator, Modulator, error) {
	// TODO test
	// check if required fields are set
	if h.ModuleDirectoryPath == "" {
		return nil, nil, fmt.Errorf("setup module directory path first")
	}
	if h.ModuleVersion == "" {
		return nil, nil, fmt.Errorf("setup module version first")
	}

	// ensure module directory
	err := os.MkdirAll(h.ModuleDirectoryPath, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	// create default files paths
	configFilePath := filepath.Join(h.ModuleDirectoryPath, configFileName)
	stateFilePath := filepath.Join(h.ModuleDirectoryPath, stateFileName)

	// load state file
	err = state.Load(stateFilePath)
	if errors.As(err, &ncverr) {
		// if old version was found try to upgrade it
		err2 := state.Upgrade(stateFilePath)
		if err2 != nil {
			return nil, nil, err2
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("load state failed: %v", err)
	}
	// load config file
	err = config.Load(configFilePath)
	if errors.As(err, &ncverr) {
		// if old version was found try to upgrade it
		err2 := config.Upgrade(configFilePath)
		if err2 != nil {
			return nil, nil, err2
		}
	} else if err != nil {
		return nil, nil, fmt.Errorf("load config failed: %v", err)
	}

	// backup
	err = backup(h.ModuleDirectoryPath, config, state)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func (h InfrastructureModuleHelper) Save(config Modulator, state Modulator) error {
	// TODO test
	// check if required fields are set
	if h.ModuleDirectoryPath == "" {
		return fmt.Errorf("setup module directory path first")
	}
	if h.ModuleVersion == "" {
		return fmt.Errorf("setup module version first")
	}

	// ensure module directory
	err := os.MkdirAll(h.ModuleDirectoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	// create default files paths
	configFilePath := filepath.Join(h.ModuleDirectoryPath, configFileName)
	stateFilePath := filepath.Join(h.ModuleDirectoryPath, stateFileName)

	err = state.Save(stateFilePath)
	if err != nil {
		return err
	}
	err = config.Save(configFilePath)
	if err != nil {
		return err
	}
	return nil
}

func backup(moduleDirectoryPath string, config Modulator, state Modulator) error {
	backupDirectoryPath := filepath.Join(moduleDirectoryPath, backupDirectoryName)

	// ensure backups directory
	err := os.MkdirAll(backupDirectoryPath, os.ModePerm)
	if err != nil {
		return err
	}

	// prepare timestamp
	t := time.Now().Format("20060102-150405.999999")
	// backup config to
	err = config.Backup(filepath.Join(backupDirectoryPath, fmt.Sprintf("config-%s.json", t)))
	if err != nil {
		return fmt.Errorf("config backup failed: %v", err)
	}
	err = state.Backup(filepath.Join(backupDirectoryPath, fmt.Sprintf("state-%s.json", t)))
	if err != nil {
		return fmt.Errorf("state backup failed: %v", err)
	}

	return nil
}
