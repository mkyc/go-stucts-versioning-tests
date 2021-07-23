package v0

import (
	"encoding/json"
	"errors"

	awsbi "github.com/epiphany-platform/e-structures/awsbi/v0"
	azks "github.com/epiphany-platform/e-structures/azks/v0"
	hi "github.com/epiphany-platform/e-structures/hi/v0"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
	maps "github.com/mitchellh/mapstructure"
)

type Status string

const (
	kind    = "state"
	version = "v0.0.5"

	Initialized Status = "initialized"
	Applied     Status = "applied"
	Destroyed   Status = "destroyed"
)

type AwsBIState struct {
	Status Status        `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *awsbi.Config `json:"config" validate:"omitempty"`
	Output *awsbi.Output `json:"output" validate:"omitempty"`
}

type HiState struct {
	Status Status     `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *hi.Config `json:"config" validate:"omitempty"`
}

func (s *HiState) GetConfig() *hi.Config {
	if s == nil {
		return nil
	}
	return s.Config
}

type AzKSState struct {
	Status Status       `json:"status" validate:"required,eq=initialized|eq=applied|eq=destroyed"`
	Config *azks.Config `json:"config" validate:"omitempty"`
	Output *azks.Output `json:"output" validate:"omitempty"`
}

func (s *AzKSState) GetConfig() *azks.Config {
	if s == nil {
		return nil
	}
	return s.Config
}

func (s *AzKSState) GetOutput() *azks.Output {
	if s == nil {
		return nil
	}
	return s.Output
}

// TODO change into Modules

type State struct {
	Kind    *string     `json:"kind" validate:"required,eq=state"`
	Version *string     `json:"version" validate:"required,version=~0"`
	Unused  []string    `json:"-"`
	AzKS    *AzKSState  `json:"azks" validate:"omitempty"`
	Hi      *HiState    `json:"hi" validate:"omitempty"`
	AwsBI   *AwsBIState `json:"awsbi" validate:"omitempty"`
}

func (s *State) GetAzKSState() *AzKSState {
	if s == nil {
		return nil
	}
	return s.AzKS
}

func (s *State) GetHiState() *HiState {
	if s == nil {
		return nil
	}
	return s.Hi
}

//TODO test
func NewState() *State {
	return &State{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Unused:  []string{},
	}
}

func (s *State) Marshal() ([]byte, error) {
	err := s.isValid()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(s, "", "\t")
}

func (s *State) Unmarshal(b []byte) (err error) {
	var input map[string]interface{}
	if err = json.Unmarshal(b, &input); err != nil {
		return
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &s,
	})
	if err != nil {
		return
	}
	err = d.Decode(input)
	if err != nil {
		return
	}
	s.Unused = md.Unused
	err = s.isValid()
	return
}

func (s *State) isValid() error {
	if s == nil {
		return errors.New("state is nil")
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

// DO NOT USE!!!
// This is temporary function used to fix existing issue (https://github.com/epiphany-platform/e-structures/issues/10)
// in some modules and will be removed shortly after issue is resolved in all modules
func (s *State) UnmarshalDoNotUse(b []byte) error {
	var input map[string]interface{}
	if err := json.Unmarshal(b, &input); err != nil {
		return err
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &s,
	})
	if err != nil {
		return err
	}
	err = d.Decode(input)
	if err != nil {
		return err
	}
	s.Unused = md.Unused
	return nil
}

// DO NOT USE!!!
// This is temporary function used to fix existing issue (https://github.com/epiphany-platform/e-structures/issues/10)
// in some modules and will be removed shortly after issue is resolved in all modules
func (s *State) IsValidDoNotUse() error {
	return s.isValid()
}
