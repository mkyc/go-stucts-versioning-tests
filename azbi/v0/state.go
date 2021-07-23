package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
)

type State struct {
	Meta   *Meta         `json:"meta" validate:"required"`
	Status shared.Status `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *Config       `json:"config" validate:"omitempty"`
	Output *Output       `json:"output" validate:"omitempty"`
	Unused []string      `json:"-"`
}

func (s *State) Init(moduleVersion string) {
	*s = State{
		Meta: &Meta{
			Kind:          to.StrPtr(stateKind),
			Version:       to.StrPtr(stateVersion),
			ModuleVersion: to.StrPtr(moduleVersion),
		},
		Status: shared.Initialized,
		Config: nil, // TODO should it be nil?
		Output: nil, // TODO should it be nil?
		Unused: []string{},
	}
}

func (s *State) Backup(path string) error {
	return shared.Backup(s, path)
}

func (s *State) Load(path string) error {
	i, err := shared.Load(s, path, stateVersion)
	if err != nil {
		return err
	}
	state, ok := i.(*State)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = state.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*s = *state
	return nil
}

func (s *State) Save(path string) error {
	return shared.Save(s, path)
}

func (s *State) Print() ([]byte, error) {
	return shared.Print(s)
}

func (s *State) Validate() error {
	if s == nil {
		return errors.New("expected state is nil")
	}
	validate := validator.New()
	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}
	err = validate.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

func (s *State) Upgrade(path string) error {
	i, err := shared.Upgrade(s, path)
	if err != nil {
		return err
	}
	state, ok := i.(*State)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = state.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*s = *state
	return nil
}

func (s *State) UpgradeFunc(input map[string]interface{}) error {
	upgraded := false
	for !upgraded {
		v, err := shared.GetVersion(input)
		if err != nil {
			return err
		}
		switch v {
		case "v0.0.1":
			meta, ok := input["meta"].(map[string]interface{})
			if !ok {
				return errors.New("incorrect casting")
			}
			meta["version"] = "v0.0.2"
			input["meta"] = meta

			configSubtree, ok := input["config"].(map[string]interface{})
			if !ok {
				return errors.New("incorrect casting")
			}
			c := Config{}
			err = c.UpgradeFunc(configSubtree)
			if err != nil {
				return err
			}
			input["config"] = configSubtree
		default:
			v, err2 := shared.GetVersion(input)
			if err2 != nil {
				return err2
			}
			if v != stateVersion {
				return errors.New("unknown version to upgrade")
			}
			upgraded = true
		}
	}
	return nil
}

func (s *State) SetUnused(unused []string) {
	s.Unused = unused
}

// TODO consider validation in output ... but really think about it hard. It might not be desired.

type Output struct {
	RgName   *string         `json:"rg_name"`
	VnetName *string         `json:"vnet_name"`
	VmGroups []OutputVmGroup `json:"vm_groups"`
}

type OutputVmGroup struct {
	Name *string    `json:"vm_group_name"`
	Vms  []OutputVm `json:"vms"`
}

type OutputVm struct {
	Name       *string          `json:"vm_name"`
	PrivateIps []string         `json:"private_ips"`
	PublicIp   *string          `json:"public_ip"`
	DataDisks  []OutputDataDisk `json:"data_disks"`
}

type OutputDataDisk struct {
	Size *int `json:"size"`
	Lun  *int `json:"lun"`
}
