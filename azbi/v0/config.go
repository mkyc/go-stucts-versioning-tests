package v0

import (
	"errors"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Meta   *Meta    `json:"meta" validate:"required"`
	Params *Params  `json:"params" validate:"required"`
	Unused []string `json:"-"`
}

func (c *Config) Init(moduleVersion string) {
	*c = Config{
		Meta: &Meta{
			Kind:          to.StrPtr(configKind),
			Version:       to.StrPtr(configVersion),
			ModuleVersion: to.StrPtr(moduleVersion),
		},
		Params: &Params{
			Name:             to.StrPtr("unknown"),
			Location:         to.StrPtr("northeurope"),
			AddressSpace:     []string{"10.0.0.0/16"},
			AdminUsername:    to.StrPtr("operations"),
			RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
			Subnets: []Subnet{
				{
					Name:            to.StrPtr("main"),
					AddressPrefixes: []string{"10.0.1.0/24"},
				},
			},
			VmGroups: []VmGroup{
				{
					Name:        to.StrPtr("vm-group-0"),
					VmCount:     to.IntPtr(1),
					VmSize:      to.StrPtr("Standard_DS2_v2"),
					UsePublicIP: to.BoolPtr(true),
					SubnetNames: []string{"main"},
					VmImage: &VmImage{
						Publisher: to.StrPtr("Canonical"),
						Offer:     to.StrPtr("UbuntuServer"),
						Sku:       to.StrPtr("18.04-LTS"),
						Version:   to.StrPtr("18.04.202006101"),
					},
					DataDisks: []DataDisk{
						{
							GbSize:      to.IntPtr(10),
							StorageType: to.StrPtr("Premium_LRS"),
						},
					},
				},
			},
		},
		Unused: []string{},
	}
	// TODO consider if we should call Validate() here
}

func (c *Config) Backup(path string) error {
	return shared.Backup(c, path)
}

func (c *Config) Load(path string) error {
	i, err := shared.Load(c, path, configVersion)
	if err != nil {
		return err
	}
	config, ok := i.(*Config)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = config.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*c = *config
	return nil
}

func (c *Config) Save(path string) error {
	return shared.Save(c, path)
}

func (c *Config) Print() ([]byte, error) {
	return shared.Print(c)
}

func (c *Config) Validate() error {
	if c == nil {
		return errors.New("expected config is nil")
	}
	validate := validator.New()

	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}
	validate.RegisterStructValidation(AzBISubnetsValidation, Params{})
	err = validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

func (c *Config) Upgrade(path string) error {
	i, err := shared.Upgrade(c, path)
	if err != nil {
		return err
	}
	config, ok := i.(*Config)
	if !ok {
		return errors.New("incorrect casting")
	}
	err = config.Validate() // TODO rethink if validation should be done here
	if err != nil {
		return err
	}
	*c = *config
	return nil
}

func (c *Config) UpgradeFunc(input map[string]interface{}) error {
	upgraded := false
	for !upgraded {
		v, err := shared.GetVersion(input)
		if err != nil {
			return err
		}
		switch v {
		case "v0.2.0":
			meta, ok := input["meta"].(map[string]interface{})
			if !ok {
				return errors.New("incorrect casting")
			}
			meta["version"] = "v0.2.1"
			input["meta"] = meta

			params, ok := input["params"].(map[string]interface{})
			if !ok {
				return errors.New("incorrect casting")
			}
			params["admin_username"] = "operations"
			input["params"] = params
		default:
			v, err2 := shared.GetVersion(input)
			if err2 != nil {
				return err2
			}
			if v != configVersion {
				return errors.New("unknown version to upgrade")
			}
			upgraded = true
		}
	}
	return nil
}

func (c *Config) SetUnused(unused []string) {
	c.Unused = unused
}

type Meta struct {
	Kind          *string `json:"kind" validate:"required,eq=azbiConfig|eq=azbiState"`
	Version       *string `json:"version" validate:"required,version=~0"`
	ModuleVersion *string `json:"module_version" validate:"required"`
}

type Params struct {
	Name             *string   `json:"name" validate:"required,min=1"`
	Location         *string   `json:"location" validate:"required,min=1"`
	AddressSpace     []string  `json:"address_space" validate:"omitempty,min=1,dive,min=1,cidr"`
	Subnets          []Subnet  `json:"subnets" validate:"required_with=AddressSpace,excluded_without=AddressSpace,omitempty,min=1,dive,required"` // TODO custom validator that subnets are in AddressSpaces
	VmGroups         []VmGroup `json:"vm_groups" validate:"required,dive"`
	AdminUsername    *string   `json:"admin_username" validate:"required,min=1"`
	RsaPublicKeyPath *string   `json:"rsa_pub_path" validate:"required,min=1"`
}

// ExtractEmptySubnets extracts list of Subnet unassigned to any of VmGroup
func (p *Params) ExtractEmptySubnets() []Subnet {
	if p == nil {
		return nil
	}
	if p.Subnets == nil || len(p.Subnets) == 0 {
		return nil
	}
	if p.VmGroups == nil || len(p.VmGroups) == 0 {
		return p.Subnets
	}
	m := make(map[string]Subnet)
	for _, subnet := range p.Subnets {
		m[*subnet.Name] = subnet
	}
	for _, vmGroup := range p.VmGroups {
		for _, subnet := range p.Subnets {
			for _, subnetName := range vmGroup.SubnetNames {
				if *subnet.Name == subnetName {
					_, ok := m[subnetName]
					if ok {
						delete(m, subnetName)
						break
					}
				}
			}
		}
	}
	result := make([]Subnet, 0)
	for _, s := range p.Subnets {
		if v, ok := m[*s.Name]; ok {
			result = append(result, v)
		}
	}
	return result
}

type Subnet struct {
	Name            *string  `json:"name" validate:"required,min=1"`
	AddressPrefixes []string `json:"address_prefixes" validate:"required,min=1,dive,required,cidr"`
}

type VmGroup struct {
	Name        *string    `json:"name" validate:"required,min=1"`
	VmCount     *int       `json:"vm_count" validate:"required,min=1"`
	VmSize      *string    `json:"vm_size" validate:"required,min=1"`
	UsePublicIP *bool      `json:"use_public_ip" validate:"required"`
	SubnetNames []string   `json:"subnet_names" validate:"omitempty,min=1,dive,required"`
	VmImage     *VmImage   `json:"vm_image" validate:"required,dive"`
	DataDisks   []DataDisk `json:"data_disks" validate:"required,dive"`
}

type VmImage struct {
	Publisher *string `json:"publisher" validate:"required,min=1"`
	Offer     *string `json:"offer" validate:"required,min=1"`
	Sku       *string `json:"sku" validate:"required,min=1"`
	Version   *string `json:"version" validate:"required,min=1"`
}

type DataDisk struct {
	GbSize      *int    `json:"disk_size_gb" validate:"required,min=1"`
	StorageType *string `json:"storage_type" validate:"required,eq=Standard_LRS|eq=Premium_LRS|eq=StandardSSD_LRS|eq=UltraSSD_LRS"` // https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/managed_disk#storage_account_type
}
