package v0

import (
	"errors"
	"fmt"
	"github.com/epiphany-platform/e-structures/shared"
	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_Init(t *testing.T) {
	tests := []struct {
		name          string
		moduleVersion string
		want          *Config
	}{
		{
			name:          "happy path",
			moduleVersion: "v1.1.1",
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v1.1.1"),
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := &Config{}
			got.Init(tt.moduleVersion)
			a.Equal(tt.want, got)
		})
	}
}

func TestConfig_Backup(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name:    "happy path",
			config:  &Config{},
			wantErr: nil,
		},
		{
			name:    "file already exists",
			config:  &Config{},
			wantErr: os.ErrExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p, err := createTempDirectory("azbi-config-backup")
			if errors.Is(tt.wantErr, os.ErrExist) {
				err = ioutil.WriteFile(filepath.Join(p, "backup-file.json"), []byte("content"), 0644)
				t.Logf("path: %s", filepath.Join(p, "backup-file.json"))
				a.NoError(err)
			}
			err = tt.config.Backup(filepath.Join(p, "backup-file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				a.Equal(tt.wantErr, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestConfig_Load(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [{
			"name": "main",
			"address_prefixes": [
				"10.0.1.0/24"
			]
		}],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 3,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}, 
			"data_disks": [
				{
					"disk_size_gb": 10, 
					"storage_type": "Premium_LRS"
				}
			]
		}],
		"admin_username": "operations", 
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
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
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"extra_inner_field" : "extra_inner_value",
		"location": "northeurope",
		"name": "epiphany",
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main", 
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 3,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			}, 
			"data_disks": [
				{
					"disk_size_gb": 10, 
					"storage_type": "Premium_LRS"
				}
			]
		}],
		"admin_username": "operations", 
		"rsa_pub_path": "/shared/vms_rsa.pub"
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:     to.StrPtr("northeurope"),
					Name:         to.StrPtr("epiphany"),
					AddressSpace: []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
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
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
				},
				Unused: []string{"params.extra_inner_field", "extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "ensure load is performing validation",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "minimal correct json",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name", 
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": false,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple subnets configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name",  
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "main",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			},
			{
				"name": "second",
				"address_prefixes": [
					"10.0.2.0/24"
				]
			}
		],
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 3,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"subnet_names": ["main"],
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"main"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "vm_group without networking",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 3,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": true,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location: to.StrPtr("northeurope"),
					Name:     to.StrPtr("epiphany"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple vm_groups configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name",  
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "first",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			}
		],
		"vm_groups": [
			{
				"name": "first",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["first"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			},
			{
				"name": "second",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["first"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			}
		]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple vm_groups and subnets configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name", 
		"address_space": [
			"10.0.0.0/16"
		], 
		"subnets": [
			{
				"name": "first",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			},
			{
				"name": "second",
				"address_prefixes": [
					"10.0.2.0/24"
				]
			}
		],
		"vm_groups": [
			{
				"name": "first",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["first"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			},
			{
				"name": "second",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["second"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			}
		]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"first"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"second"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "2 vm_groups and 3 subnets configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name",  
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "first",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			},
			{
				"name": "second",
				"address_prefixes": [
					"10.0.2.0/24"
				]
			},
			{
				"name": "third",
				"address_prefixes": [
					"10.0.3.0/24"
				]
			}
		],
		"vm_groups": [
			{
				"name": "first",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["first", "third"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			},
			{
				"name": "second",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["second", "third"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			}
		]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
						{
							Name:            to.StrPtr("third"),
							AddressPrefixes: []string{"10.0.3.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"first", "third"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"second", "third"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "multiple vm_groups and subnets and data disks configuration",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name",  
		"address_space": [
			"10.0.0.0/16"
		],
		"subnets": [
			{
				"name": "first",
				"address_prefixes": [
					"10.0.1.0/24"
				]
			},
			{
				"name": "second",
				"address_prefixes": [
					"10.0.2.0/24"
				]
			}
		],
		"vm_groups": [
			{
				"name": "first",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["first"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": [
					{
						"disk_size_gb": 10, 
						"storage_type": "Premium_LRS"
					},
					{
						"disk_size_gb": 20, 
						"storage_type": "Standard_LRS"
					}
				]
			},
			{
				"name": "second",
				"vm_count": 3,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": true,
				"subnet_names": ["second"],
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": [
					{
						"disk_size_gb": 30, 
						"storage_type": "StandardSSD_LRS"
					},
					{
						"disk_size_gb": 40, 
						"storage_type": "UltraSSD_LRS"
					}
				]
			}
		]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.2.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("first"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"first"},
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
								{
									GbSize:      to.IntPtr(20),
									StorageType: to.StrPtr("Standard_LRS"),
								},
							},
						},
						{
							Name:        to.StrPtr("second"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"second"},
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{
								{
									GbSize:      to.IntPtr(30),
									StorageType: to.StrPtr("StandardSSD_LRS"),
								},
								{
									GbSize:      to.IntPtr(40),
									StorageType: to.StrPtr("UltraSSD_LRS"),
								},
							},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "old version exception",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.0",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"rsa_pub_path": "some-file-name", 
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": false,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want:    nil,
			wantErr: shared.NotCurrentVersionError{Version: "v0.2.0"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azbi-config-load", tt.json)
			r.NoError(err)
			got := &Config{}
			err = got.Load(p)
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				wj, err2 := tt.want.Print()
				a.NoError(err2)
				gj, err2 := got.Print()
				a.NoError(err2)
				a.Equal(string(wj), string(gj))
			}
		})
	}
}

func TestConfig_Save(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"address_space": null,
		"subnets": null,
		"vm_groups": [
			{
				"name": "vm-group0",
				"vm_count": 1,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": false,
				"subnet_names": null,
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			}
		],
		"admin_username": "operations",
		"rsa_pub_path": "some-file-name"
	}
}`),
			wantErr: nil,
		},
		{
			name:   "invalid",
			config: &Config{},
			want:   nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDirectory("azbi-config-save")
			a.NoError(err)

			err = tt.config.Save(filepath.Join(p, "file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				a.FileExists(filepath.Join(p, "file.json"))
				got, err2 := ioutil.ReadFile(filepath.Join(p, "file.json"))
				a.NoError(err2)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestConfig_Print(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		want    []byte
		wantErr bool
	}{
		{
			name: "happy path",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"address_space": null,
		"subnets": null,
		"vm_groups": [
			{
				"name": "vm-group0",
				"vm_count": 1,
				"vm_size": "Standard_DS2_v2",
				"use_public_ip": false,
				"subnet_names": null,
				"vm_image": {
					"publisher": "Canonical",
					"offer": "UbuntuServer",
					"sku": "18.04-LTS",
					"version": "18.04.202006101"
				},
				"data_disks": []
			}
		],
		"admin_username": "operations",
		"rsa_pub_path": "some-file-name"
	}
}`),
			wantErr: false,
		},
		{
			name:    "invalid",
			config:  &Config{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got, err := tt.config.Print()
			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestConfig_Valid(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name: "minimal correct",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name:   "empty struct",
			config: &Config{},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params",
					Field: "Params",
					Tag:   "required",
				},
			},
		},
		{
			name: "meta missing",
			config: &Config{
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
			},
		},
		{
			name: "major version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v100.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Meta.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.100.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.1.100"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "just vm_groups in params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Location",
					Field: "Location",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "required",
				},
			},
		},
		{
			name: "missing requested subnets list",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							SubnetNames: []string{"main"},
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "VmGroups[0].SubnetNames[0]",
					Tag:   "insubnets",
				},
			},
		},
		{
			name: "empty subnets list",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets:          []Subnet{},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "min",
				},
			},
		},
		{
			name: "subnets but no address space",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "excluded_without",
				},
			},
		},
		{
			name: "address space but no subnets",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets",
					Field: "Subnets",
					Tag:   "required_with",
				},
			},
		},
		{
			name: "missing subnet params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty subnet params",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr(""),
							AddressPrefixes: []string{},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes",
					Field: "AddressPrefixes",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty subnet address prefixes element and not cidr",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"10.0.0.0/16"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("first"),
							AddressPrefixes: []string{""},
						},
						{
							Name:            to.StrPtr("second"),
							AddressPrefixes: []string{"10.0.1.0"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.Subnets[0].AddressPrefixes[0]",
					Field: "AddressPrefixes[0]",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.Subnets[1].AddressPrefixes[0]",
					Field: "AddressPrefixes[0]",
					Tag:   "cidr",
				},
			},
		},
		{
			name: "emtpy address_space",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace",
					Field: "AddressSpace",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty address_space element or not cidr",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					AddressSpace:     []string{"", "10.0.1.0"},
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace[0]",
					Field: "AddressSpace[0]",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AddressSpace[1]",
					Field: "AddressSpace[1]",
					Tag:   "cidr",
				},
			},
		},
		{
			name: "empty name location admin_username and rsa_pub_path ",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.1.0"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr(""),
					Name:             to.StrPtr(""),
					AdminUsername:    to.StrPtr(""),
					RsaPublicKeyPath: to.StrPtr(""),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.Location",
					Field: "Location",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing vm groups",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups",
					Field: "VmGroups",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm groups",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AddressSpace:     []string{"10.0.0.0/16"},
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
					Subnets: []Subnet{
						{
							Name:            to.StrPtr("main"),
							AddressPrefixes: []string{"10.0.1.0/24"},
						},
					},
					VmGroups: []VmGroup{},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "missing vm_groups parameters",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
						{},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].Name",
					Field: "Name",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmSize",
					Field: "VmSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].UsePublicIP",
					Field: "UsePublicIP",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage",
					Field: "VmImage",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks",
					Field: "DataDisks",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups parameters",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr(""),
							VmCount:     to.IntPtr(0),
							VmSize:      to.StrPtr(""),
							SubnetNames: []string{},
							UsePublicIP: nil,
							VmImage:     &VmImage{},
							DataDisks: []DataDisk{
								{},
							},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].Name",
					Field: "Name",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmSize",
					Field: "VmSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].UsePublicIP",
					Field: "UsePublicIP",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames",
					Field: "SubnetNames",
					Tag:   "min",
				},
			},
		},
		{
			name: "vm_groups negative vm_count and subnet_names list empty value",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(-1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{""},
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
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmCount",
					Field: "VmCount",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "SubnetNames[0]",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].SubnetNames[0]",
					Field: "VmGroups[0].SubnetNames[0]",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.data_disks list value",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
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
								{},
							},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "required",
				},
			},
		},
		{
			name: "incorrect vm_groups.data_disks list value",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
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
									GbSize:      to.IntPtr(0),
									StorageType: to.StrPtr("incorrect"),
								},
								{
									GbSize:      to.IntPtr(-1),
									StorageType: to.StrPtr(""),
								},
							},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[0].StorageType",
					Field: "StorageType",
					Tag:   "eq=Standard_LRS|eq=Premium_LRS|eq=StandardSSD_LRS|eq=UltraSSD_LRS",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[1].GbSize",
					Field: "GbSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].DataDisks[1].StorageType",
					Field: "StorageType",
					Tag:   "eq=Standard_LRS|eq=Premium_LRS|eq=StandardSSD_LRS|eq=UltraSSD_LRS",
				},
			},
		},
		{
			name: "missing vm_groups.vm_image parameters",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"main"},
							VmImage:     &VmImage{},
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
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty vm_groups.vm_image parameters",
			config: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
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
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(3),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(true),
							SubnetNames: []string{"main"},
							VmImage: &VmImage{
								Publisher: to.StrPtr(""),
								Offer:     to.StrPtr(""),
								Sku:       to.StrPtr(""),
								Version:   to.StrPtr(""),
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
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Publisher",
					Field: "Publisher",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Offer",
					Field: "Offer",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Sku",
					Field: "Sku",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VmGroups[0].VmImage.Version",
					Field: "Version",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			err := tt.config.Validate()
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				_, ok = err.(validator.ValidationErrors)
				r.Equal(true, ok)
				errs := err.(validator.ValidationErrors)
				a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))

				for _, e := range errs {
					found := false
					for _, we := range tt.wantErr.(test.TestValidationErrors) {
						if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
					}
				}
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestConfig_Upgrade(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path nothing to upgrade",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"admin_username": "operations", 
		"rsa_pub_path": "some-file-name", 
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": false,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "ensure that validation is also performed in upgrade",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"rsa_pub_path": "some-file-name", 
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": false,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "required",
				},
			},
		},
		{
			name: "upgrade v0.2.0 to v0.2.1",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.2.0",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany",
		"rsa_pub_path": "some-file-name", 
		"vm_groups": [{
			"name": "vm-group0",
			"vm_count": 1,
			"vm_size": "Standard_DS2_v2",
			"use_public_ip": false,
			"vm_image": {
				"publisher": "Canonical",
				"offer": "UbuntuServer",
				"sku": "18.04-LTS",
				"version": "18.04.202006101"
			},
			"data_disks": []
		}]
	}
}
`),
			want: &Config{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiConfig"),
					Version:       to.StrPtr("v0.2.1"),
					ModuleVersion: to.StrPtr("v0.0.1"),
				},
				Params: &Params{
					Location:         to.StrPtr("northeurope"),
					Name:             to.StrPtr("epiphany"),
					AdminUsername:    to.StrPtr("operations"),
					RsaPublicKeyPath: to.StrPtr("some-file-name"),
					VmGroups: []VmGroup{
						{
							Name:        to.StrPtr("vm-group0"),
							VmCount:     to.IntPtr(1),
							VmSize:      to.StrPtr("Standard_DS2_v2"),
							UsePublicIP: to.BoolPtr(false),
							VmImage: &VmImage{
								Publisher: to.StrPtr("Canonical"),
								Offer:     to.StrPtr("UbuntuServer"),
								Sku:       to.StrPtr("18.04-LTS"),
								Version:   to.StrPtr("18.04.202006101"),
							},
							DataDisks: []DataDisk{},
						},
					},
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "some unknown version",
			json: []byte(`{
	"meta": {
		"kind": "azbiConfig",
		"version": "v0.0.1",
		"module_version": "v0.0.1"
	},
	"params": {
		"location": "northeurope",
		"name": "epiphany"
	}
}
`),
			want:    nil,
			wantErr: errors.New("unknown version to upgrade"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azbi-config-load", tt.json)
			r.NoError(err)
			got := &Config{}
			err = got.Upgrade(p)
			if tt.wantErr != nil {
				r.Error(err)
				_, ok := err.(*validator.InvalidValidationError)
				r.Equal(false, ok)
				errs, ok := err.(validator.ValidationErrors)
				if ok {
					for _, e := range errs {
						found := false
						for _, we := range tt.wantErr.(test.TestValidationErrors) {
							if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
								found = true
								break
							}
						}
						if !found {
							t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), tt.wantErr.Error())
						}
					}
					a.Equal(len(tt.wantErr.(test.TestValidationErrors)), len(errs))
				} else {
					a.Equal(tt.wantErr, err)
				}
			} else {
				a.NoError(err)
				wj, err2 := tt.want.Print()
				a.NoError(err2)
				gj, err2 := got.Print()
				a.NoError(err2)
				a.Equal(string(wj), string(gj))
			}
		})
	}
}

func TestParams_ExtractEmptySubnets(t *testing.T) {
	tests := []struct {
		name   string
		params *Params
		want   []Subnet
	}{
		{
			name: "happy path",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet2"),
					AddressPrefixes: []string{"2.2.2.2/24"},
				},
			},
		},
		{
			name:   "nil params",
			params: nil,
			want:   nil,
		},
		{
			name: "nil subnets",
			params: &Params{
				Subnets: nil,
			},
			want: nil,
		},
		{
			name: "empty subnets",
			params: &Params{
				Subnets: []Subnet{},
			},
			want: nil,
		},
		{
			name: "nil vm_groups",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: nil,
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
			},
		},
		{
			name: "empty vm_groups",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
			},
		},
		{
			name: "no empty subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups no empty subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
					{
						SubnetNames: []string{"subnet2"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups reuse one subnet",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet1"},
					},
					{
						SubnetNames: []string{"subnet1"},
					},
				},
			},
			want: []Subnet{},
		},
		{
			name: "multiple vm_groups some free subnets",
			params: &Params{
				Subnets: []Subnet{
					{
						Name:            to.StrPtr("subnet1"),
						AddressPrefixes: []string{"1.1.1.1/24"},
					},
					{
						Name:            to.StrPtr("subnet2"),
						AddressPrefixes: []string{"2.2.2.2/24"},
					},
					{
						Name:            to.StrPtr("subnet3"),
						AddressPrefixes: []string{"3.3.3.3/24"},
					},
					{
						Name:            to.StrPtr("subnet4"),
						AddressPrefixes: []string{"4.4.4.4/24"},
					},
					{
						Name:            to.StrPtr("subnet5"),
						AddressPrefixes: []string{"5.5.5.5/24"},
					},
				},
				VmGroups: []VmGroup{
					{
						SubnetNames: []string{"subnet2", "subnet5"},
					},
					{
						SubnetNames: []string{"subnet2", "subnet4"},
					},
				},
			},
			want: []Subnet{
				{
					Name:            to.StrPtr("subnet1"),
					AddressPrefixes: []string{"1.1.1.1/24"},
				},
				{
					Name:            to.StrPtr("subnet3"),
					AddressPrefixes: []string{"3.3.3.3/24"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := tt.params.ExtractEmptySubnets()
			a.Equal(tt.want, got)
		})
	}
}

func createTempDocumentFile(name string, document []byte) (string, error) {
	p, err := ioutil.TempDir("", fmt.Sprintf("e-structures-%s-*", name))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(filepath.Join(p, "file.json"), document, 0644)
	return filepath.Join(p, "file.json"), err
}

func createTempDirectory(name string) (string, error) {
	return ioutil.TempDir("", fmt.Sprintf("e-structures-%s-*", name))
}
