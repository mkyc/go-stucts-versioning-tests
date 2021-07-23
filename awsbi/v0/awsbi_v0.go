package v0

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/epiphany-platform/e-structures/utils/validators"
	"github.com/go-playground/validator/v10"
	maps "github.com/mitchellh/mapstructure"
)

const (
	kind    = "awsbi"
	version = "v0.0.1"
)

type DataDisk struct {
	DeviceName *string `json:"device_name" validate:"required,min=1"` // https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/volume_attachment#device_name
	GbSize     *int    `json:"disk_size_gb" validate:"required,min=1"`
	Type       *string `json:"type" validate:"required,eq=standard|eq=gp2|eq=gp3|eq=io1|eq=io2|eq=sc1|eq=st1"` // https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/ebs_volume#type
}

type VmImage struct {
	AMI   *string `json:"ami" validate:"required,min=1"`
	Owner *string `json:"owner" validate:"required,min=1"`
}

type VmGroup struct {
	Name               *string    `json:"name" validate:"required,min=1"`
	VmCount            *int       `json:"vm_count" validate:"required,min=1"`
	VmSize             *string    `json:"vm_size" validate:"required,min=1"`
	UsePublicIp        *bool      `json:"use_public_ip" validate:"required"`
	SubnetNames        []string   `json:"subnet_names" validate:"omitempty,min=1,dive,required"`
	SecurityGroupNames []string   `json:"sg_names" validate:"omitempty,min=1,dive,required"`
	VmImage            *VmImage   `json:"vm_image" validate:"required,dive"`
	RootVolumeGbSize   *int       `json:"root_volume_size" validate:"required,min=1"`
	DataDisks          []DataDisk `json:"data_disks" validate:"omitempty,dive"`
}

type SecurityRule struct {
	Protocol   *string  `json:"protocol" validate:"required,min=1"`
	FromPort   *int     `json:"from_port" validate:"required,min=0"`
	ToPort     *int     `json:"to_port" validate:"required,min=0"`
	CidrBlocks []string `json:"cidr_blocks" validate:"omitempty,min=1,dive,required,cidr"`
}

type Rules struct {
	Ingress []SecurityRule `json:"ingress" validate:"omitempty,min=1,dive,required"`
	Egress  []SecurityRule `json:"egress" validate:"omitempty,min=1,dive,required"`
}

type SecurityGroup struct {
	Name  *string `json:"name" validate:"required,min=1"`
	Rules *Rules  `json:"rules" validate:"required,dive"`
}

type Subnet struct {
	Name             *string `json:"name" validate:"required,min=1"`
	AvailabilityZone *string `json:"availability_zone" validate:"required,min=1"`
	AddressPrefixes  *string `json:"address_prefixes" validate:"required,min=1,cidr"`
}

type Subnets struct {
	Private []Subnet `json:"private" validate:"required_without=Public"`
	Public  []Subnet `json:"public" validate:"required_without=Private"`
}

type Params struct {
	Name                  *string `json:"name" validate:"required,min=1"`
	Region                *string `json:"region" validate:"required,min=1"`
	NatGatewayCount       *int    `json:"nat_gateway_count" validate:"required,min=0"`
	VirtualPrivateGateway *bool   `json:"virtual_private_gateway" validate:"required"`

	RsaPublicKeyPath *string `json:"rsa_pub_path" validate:"required,min=1"`

	VpcAddressSpace *string         `json:"vpc_address_space" validate:"required,min=1,cidr"`
	Subnets         *Subnets        `json:"subnets" validate:"required,dive,omitempty"`
	SecurityGroups  []SecurityGroup `json:"security_groups" validate:"required,dive"`
	VmGroups        []VmGroup       `json:"vm_groups" validate:"required,dive"`
}

type Config struct {
	Kind    *string  `json:"kind" validate:"required,eq=awsbi"`
	Version *string  `json:"version" validate:"required,version=~0"`
	Params  *Params  `json:"params" validate:"required,dive"`
	Unused  []string `json:"-"`
}

//TODO test
func NewConfig() *Config {
	return &Config{
		Kind:    to.StrPtr(kind),
		Version: to.StrPtr(version),
		Params: &Params{
			Name:                  to.StrPtr("epiphany"),
			Region:                to.StrPtr("eu-central-1"),
			NatGatewayCount:       to.IntPtr(1),
			VirtualPrivateGateway: to.BoolPtr(false),
			RsaPublicKeyPath:      to.StrPtr("/shared/vms_rsa.pub"),
			VpcAddressSpace:       to.StrPtr("10.1.0.0/20"),
			Subnets: &Subnets{
				Private: []Subnet{
					{
						Name:             to.StrPtr("first_private_subnet"),
						AvailabilityZone: to.StrPtr("any"),
						AddressPrefixes:  to.StrPtr("10.1.1.0/24"),
					},
				},
				Public: []Subnet{
					{
						Name:             to.StrPtr("first_public_subnet"),
						AvailabilityZone: to.StrPtr("any"),
						AddressPrefixes:  to.StrPtr("10.1.2.0/24"),
					},
				},
			},
			SecurityGroups: []SecurityGroup{
				{
					Name: to.StrPtr("default_sg"),
					Rules: &Rules{
						Ingress: []SecurityRule{
							{
								Protocol:   to.StrPtr("-1"),
								FromPort:   to.IntPtr(0),
								ToPort:     to.IntPtr(0),
								CidrBlocks: []string{"10.1.0.0/20"},
							},
							{
								Protocol:   to.StrPtr("tcp"),
								FromPort:   to.IntPtr(22),
								ToPort:     to.IntPtr(22),
								CidrBlocks: []string{"0.0.0.0/0"},
							},
						},
						Egress: []SecurityRule{
							{
								Protocol:   to.StrPtr("-1"),
								FromPort:   to.IntPtr(0),
								ToPort:     to.IntPtr(0),
								CidrBlocks: []string{"0.0.0.0/0"},
							},
						},
					},
				},
			},
			VmGroups: []VmGroup{
				{
					Name:               to.StrPtr("vm-group0"),
					VmCount:            to.IntPtr(1),
					VmSize:             to.StrPtr("t3.medium"),
					UsePublicIp:        to.BoolPtr(false),
					SubnetNames:        []string{"first_private_subnet"},
					SecurityGroupNames: []string{"default_sg"},
					VmImage: &VmImage{
						AMI:   to.StrPtr("RHEL-7.8_HVM_GA-20200225-x86_64-1-Hourly2-GP2"),
						Owner: to.StrPtr("309956199498"),
					},
					RootVolumeGbSize: to.IntPtr(30),
					DataDisks: []DataDisk{
						{
							DeviceName: to.StrPtr("/dev/sdf"),
							GbSize:     to.IntPtr(16),
							Type:       to.StrPtr("gp2"),
						},
					},
				},
			},
		},
		Unused: []string{},
	}
}

func (c *Config) Marshal() ([]byte, error) {
	err := c.isValid()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(c, "", "\t")
}

func (c *Config) Unmarshal(b []byte) (err error) {
	var input map[string]interface{}
	if err = json.Unmarshal(b, &input); err != nil {
		return
	}
	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{
		Metadata: &md,
		TagName:  "json",
		Result:   &c,
	})
	if err != nil {
		return
	}
	err = d.Decode(input)
	if err != nil {
		return
	}
	c.Unused = md.Unused
	err = c.isValid()
	return
}

func (c *Config) isValid() error {
	if c == nil {
		return errors.New("azbi config is nil")
	}
	validate := validator.New()

	err := validate.RegisterValidation("version", validators.HasVersion)
	if err != nil {
		return err
	}
	validate.RegisterStructValidation(AwsBIParamsValidation, Params{})
	err = validate.Struct(c)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		return err
	}
	return nil
}

type OutputDataDisk struct {
	Size       *int    `json:"size"`
	DeviceName *string `json:"device_name"`
}

type OutputVm struct {
	Name      *string          `json:"name"`
	PublicIp  *string          `json:"public_ip"`
	PrivateIp *string          `json:"private_ip"`
	DataDisks []OutputDataDisk `json:"data_disks"`
}

type OutputVmGroup struct {
	Name *string    `json:"name"`
	Vms  []OutputVm `json:"vms"`
}

type Output struct {
	VpcId             *string         `json:"vpc_id"`
	PrivateSubnetIds  []string        `json:"private_subnet_ids"`
	PublicSubnetIds   []string        `json:"public_subnet_ids"`
	PrivateRouteTable *string         `json:"private_route_table"`
	VmGroups          []OutputVmGroup `json:"vm_groups"`
}

func AwsBIParamsValidation(sl validator.StructLevel) {
	params := sl.Current().Interface().(Params)
	if len(params.VmGroups) > 0 {
		for i, vmGroup := range params.VmGroups {
			for j, sn := range vmGroup.SubnetNames {
				if sn == "" {
					sl.ReportError(
						params.VmGroups[i].SubnetNames[j],
						fmt.Sprintf("VmGroups[%d].SubnetNames[%d]", i, j),
						fmt.Sprintf("SubnetNames[%d]", j),
						"required",
						"")
				}
				found := false
				if params.Subnets != nil {
					for _, s := range params.Subnets.Private {
						if s.Name != nil && sn == *s.Name {
							found = true
						}
					}
					for _, s := range params.Subnets.Public {
						if s.Name != nil && sn == *s.Name {
							found = true
						}
					}
				}
				if !found {
					sl.ReportError(
						params.VmGroups[i].SubnetNames[j],
						fmt.Sprintf("VmGroups[%d].SubnetNames[%d]", i, j),
						fmt.Sprintf("SubnetNames[%d]", j),
						"insubnets",
						"")
				}
			}
			for j, sg := range vmGroup.SecurityGroupNames {
				if sg == "" {
					sl.ReportError(
						params.VmGroups[i].SecurityGroupNames[j],
						fmt.Sprintf("VmGroups[%d].SecurityGroupNames[%d]", i, j),
						fmt.Sprintf("SecurityGroupNames[%d]", j),
						"required",
						"")
				}
				found := false
				if params.SecurityGroups != nil {
					for _, s := range params.SecurityGroups {
						if s.Name != nil && sg == *s.Name {
							found = true
						}
					}
				}
				if !found {
					sl.ReportError(
						params.VmGroups[i].SecurityGroupNames[j],
						fmt.Sprintf("VmGroups[%d].SecurityGroupNames[%d]", i, j),
						fmt.Sprintf("SecurityGroupNames[%d]", j),
						"insecuritygroups",
						"")
				}
			}
		}
	}
	if params.Subnets != nil {
		if len(params.Subnets.Private) == 0 && len(params.Subnets.Public) == 0 {
			sl.ReportError(
				params.Subnets,
				"Subnets",
				"Subnets",
				"private_or_public",
				"")
		}
		if len(params.Subnets.Private) > 0 {
			for i, s := range params.Subnets.Private {
				validate := validator.New()
				err := validate.Struct(s)
				if err != nil {
					if e, ok := err.(validator.ValidationErrors); ok {
						namespace := fmt.Sprintf("Subnets.Private[%d].", i)
						sl.ReportValidationErrors(namespace, namespace, e)
					} else {
						sl.ReportError(
							params.Subnets,
							"Subnets.Private",
							"Private",
							"fatal",
							"")
					}
				}
			}
		}
		if len(params.Subnets.Public) > 0 {
			for i, s := range params.Subnets.Public {
				validate := validator.New()
				err := validate.Struct(s)
				if err != nil {
					if e, ok := err.(validator.ValidationErrors); ok {
						namespace := fmt.Sprintf("Subnets.Public[%d].", i)
						sl.ReportValidationErrors(namespace, namespace, e)
					} else {
						sl.ReportError(
							params.Subnets,
							"Subnets.Public",
							"Public",
							"fatal",
							"")
					}
				}
			}
		}
	}
}
