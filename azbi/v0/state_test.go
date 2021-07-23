package v0

import (
	"errors"
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

func TestState_Init(t *testing.T) {
	tests := []struct {
		name          string
		moduleVersion string
		want          *State
	}{
		{
			name:          "happy path",
			moduleVersion: "v1.1.1",
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			got := &State{}
			got.Init(tt.moduleVersion)
			a.Equal(tt.want, got)
		})
	}
}

func TestState_Backup(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr error
	}{
		{
			name:    "happy path",
			state:   &State{},
			wantErr: nil,
		},
		{
			name:    "file already exists",
			state:   &State{},
			wantErr: os.ErrExist,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			p, err := createTempDirectory("azbi-state-backup")
			if errors.Is(tt.wantErr, os.ErrExist) {
				err = ioutil.WriteFile(filepath.Join(p, "backup-file.json"), []byte("content"), 0644)
				t.Logf("path: %s", filepath.Join(p, "backup-file.json"))
				a.NoError(err)
			}
			err = tt.state.Backup(filepath.Join(p, "backup-file.json"))
			if tt.wantErr != nil {
				a.Error(err)
				a.Equal(tt.wantErr, err)
			} else {
				a.NoError(err)
			}
		})
	}
}

func TestState_Load(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *State
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"unknown_key_1": "unknown_value_1",
	"config": {
		"meta": {
			"kind": "azbiConfig",
			"version": "v0.1.1",
			"module_version": "v0.0.1"
		},
		"params": {
			"unknown_key_2": "unknown_value_2",
			"name": "epiphany",
			"location": "northeurope",
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
				"unknown_key_3": "unknown_value_3",
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
			}],
			"admin_username": "operations", 
			"rsa_pub_path": "/shared/vms_rsa.pub"
		}
	},
	"output": {
		"unknown_key_4": "unknown_value_4",
		"rg_name": "epiphany-rg",
		"vm_groups": [
			{
				"vm_group_name": "vm-group0",
				"unknown_key_5": "unknown_value_5",
				"vms": [
					{
						"private_ips": [
							"10.0.1.4"
						],
						"public_ip": "123.234.345.456",
						"vm_name": "epiphany-vm-group0-0"
					},
					{
						"private_ips": [
							"10.0.1.5"
						],
						"public_ip": "123.234.345.457",
						"vm_name": "epiphany-vm-group0-1"
					},
					{
						"private_ips": [
							"10.0.1.6"
						],
						"public_ip": "123.234.345.458",
						"vm_name": "epiphany-vm-group0-2"
					}
				]
			}
		],
		"vnet_name": "epiphany-vnet"
	}
}`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azbiConfig"),
						Version:       to.StrPtr("v0.1.1"),
						ModuleVersion: to.StrPtr("v0.0.1"),
					},
					Params: &Params{
						Name:         to.StrPtr("epiphany"),
						Location:     to.StrPtr("northeurope"),
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
								DataDisks: []DataDisk{},
							},
						},
						AdminUsername:    to.StrPtr("operations"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
					},
					Unused: nil,
				},
				Output: &Output{
					RgName:   to.StrPtr("epiphany-rg"),
					VnetName: to.StrPtr("epiphany-vnet"),
					VmGroups: []OutputVmGroup{
						{
							Name: to.StrPtr("vm-group0"),
							Vms: []OutputVm{
								{
									Name:       to.StrPtr("epiphany-vm-group0-0"),
									PrivateIps: []string{"10.0.1.4"},
									PublicIp:   to.StrPtr("123.234.345.456"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-1"),
									PrivateIps: []string{"10.0.1.5"},
									PublicIp:   to.StrPtr("123.234.345.457"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-2"),
									PrivateIps: []string{"10.0.1.6"},
									PublicIp:   to.StrPtr("123.234.345.458"),
								},
							},
						},
					},
				},
				Unused: []string{
					"config.params.vm_groups[0].unknown_key_3",
					"config.params.unknown_key_2",
					"output.vm_groups[0].unknown_key_5",
					"output.unknown_key_4",
					"unknown_key_1",
				},
			},
			wantErr: nil,
		},
		{
			name: "minimal state",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "dev"
	},
	"status": "initialized"
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "future version mismatch",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v100.0.0",
		"module_version": "dev"
	},
	"status": "initialized"
}
`),
			want:    nil,
			wantErr: shared.NotCurrentVersionError{Version: "v100.0.0"},
		},
		{
			name: "old version",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.1",
		"module_version": "dev"
	},
	"status": "initialized"
}
`),
			want:    nil,
			wantErr: shared.NotCurrentVersionError{Version: "v0.0.1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDocumentFile("azbi-state-load", tt.json)
			r.NoError(err)
			got := &State{}
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
				a.Equal(tt.want, got)
			}
		})
	}
}

func TestState_Save(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}`),
			wantErr: nil,
		},
		{
			name:  "invalid",
			state: &State{},
			want:  nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			p, err := createTempDirectory("azbi-state-save")
			a.NoError(err)

			err = tt.state.Save(filepath.Join(p, "file.json"))
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

func TestState_Print(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		want    []byte
		wantErr error
	}{
		{
			name: "happy path",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azbiConfig"),
						Version:       to.StrPtr("v0.1.1"),
						ModuleVersion: to.StrPtr("v0.0.1"),
					},
					Params: &Params{
						Name:         to.StrPtr("epiphany"),
						Location:     to.StrPtr("northeurope"),
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
								DataDisks: []DataDisk{},
							},
						},
						AdminUsername:    to.StrPtr("operations"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
					},
					Unused: nil,
				},
				Output: &Output{
					RgName:   to.StrPtr("epiphany-rg"),
					VnetName: to.StrPtr("epiphany-vnet"),
					VmGroups: []OutputVmGroup{
						{
							Name: to.StrPtr("vm-group0"),
							Vms: []OutputVm{
								{
									Name:       to.StrPtr("epiphany-vm-group0-0"),
									PrivateIps: []string{"10.0.1.4"},
									PublicIp:   to.StrPtr("123.234.345.456"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-1"),
									PrivateIps: []string{"10.0.1.5"},
									PublicIp:   to.StrPtr("123.234.345.457"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-2"),
									PrivateIps: []string{"10.0.1.6"},
									PublicIp:   to.StrPtr("123.234.345.458"),
								},
							},
						},
					},
				},
				Unused: []string{},
			},
			want: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"config": {
		"meta": {
			"kind": "azbiConfig",
			"version": "v0.1.1",
			"module_version": "v0.0.1"
		},
		"params": {
			"name": "epiphany",
			"location": "northeurope",
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
			"vm_groups": [
				{
					"name": "vm-group0",
					"vm_count": 3,
					"vm_size": "Standard_DS2_v2",
					"use_public_ip": true,
					"subnet_names": [
						"main"
					],
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
			"rsa_pub_path": "/shared/vms_rsa.pub"
		}
	},
	"output": {
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"vm_groups": [
			{
				"vm_group_name": "vm-group0",
				"vms": [
					{
						"vm_name": "epiphany-vm-group0-0",
						"private_ips": [
							"10.0.1.4"
						],
						"public_ip": "123.234.345.456",
						"data_disks": null
					},
					{
						"vm_name": "epiphany-vm-group0-1",
						"private_ips": [
							"10.0.1.5"
						],
						"public_ip": "123.234.345.457",
						"data_disks": null
					},
					{
						"vm_name": "epiphany-vm-group0-2",
						"private_ips": [
							"10.0.1.6"
						],
						"public_ip": "123.234.345.458",
						"data_disks": null
					}
				]
			}
		]
	}
}`),
			wantErr: nil,
		},
		{
			name:  "invalid",
			state: &State{},
			want:  nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := require.New(t)
			got, err := tt.state.Print()
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
				a.Equal(string(tt.want), string(got))
			}
		})
	}
}

func TestState_Valid(t *testing.T) {
	tests := []struct {
		name    string
		state   *State
		wantErr error
	}{
		{
			name: "minimal correct",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name:  "empty struct",
			state: &State{},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "required",
				},
			},
		},
		{
			name: "meta missing",
			state: &State{
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta",
					Field: "Meta",
					Tag:   "required",
				},
			},
		},
		{
			name: "major version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v100.0.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "minor version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.100.0"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "patch version mismatch",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.100"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "incorrect status",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: "incorrect",
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Status",
					Field: "Status",
					Tag:   "eq=initialized|eq=applied|eq=destroyed",
				},
			},
		},
		{
			name: "empty config and output",
			state: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{},
				Output: &Output{},
				Unused: []string{},
			},
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Config.Meta",
					Field: "Meta",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "State.Config.Params",
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
			err := tt.state.Validate()
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

func TestState_Upgrade(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *State
		wantErr error
	}{
		{
			name: "happy path nothing to upgrade state without config",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("dev"),
				},
				Status: shared.Initialized,
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "happy path nothing to upgrade state with config and output",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.2",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"config": {
		"meta": {
			"kind": "azbiConfig",
			"version": "v0.1.1",
			"module_version": "v0.0.1"
		},
		"params": {
			"name": "epiphany",
			"location": "northeurope",
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
				"data_disks": []
			}],
			"admin_username": "operations", 
			"rsa_pub_path": "/shared/vms_rsa.pub"
		}
	},
	"output": {
		"rg_name": "epiphany-rg",
		"vm_groups": [
			{
				"vm_group_name": "vm-group0",
				"vms": [
					{
						"private_ips": [
							"10.0.1.4"
						],
						"public_ip": "123.234.345.456",
						"vm_name": "epiphany-vm-group0-0"
					},
					{
						"private_ips": [
							"10.0.1.5"
						],
						"public_ip": "123.234.345.457",
						"vm_name": "epiphany-vm-group0-1"
					},
					{
						"private_ips": [
							"10.0.1.6"
						],
						"public_ip": "123.234.345.458",
						"vm_name": "epiphany-vm-group0-2"
					}
				]
			}
		],
		"vnet_name": "epiphany-vnet"
	}
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azbiConfig"),
						Version:       to.StrPtr("v0.1.1"),
						ModuleVersion: to.StrPtr("v0.0.1"),
					},
					Params: &Params{
						Name:         to.StrPtr("epiphany"),
						Location:     to.StrPtr("northeurope"),
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
								DataDisks: []DataDisk{},
							},
						},
						AdminUsername:    to.StrPtr("operations"),
						RsaPublicKeyPath: to.StrPtr("/shared/vms_rsa.pub"),
					},
					Unused: nil,
				},
				Output: &Output{
					RgName:   to.StrPtr("epiphany-rg"),
					VnetName: to.StrPtr("epiphany-vnet"),
					VmGroups: []OutputVmGroup{
						{
							Name: to.StrPtr("vm-group0"),
							Vms: []OutputVm{
								{
									Name:       to.StrPtr("epiphany-vm-group0-0"),
									PrivateIps: []string{"10.0.1.4"},
									PublicIp:   to.StrPtr("123.234.345.456"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-1"),
									PrivateIps: []string{"10.0.1.5"},
									PublicIp:   to.StrPtr("123.234.345.457"),
								},
								{
									Name:       to.StrPtr("epiphany-vm-group0-2"),
									PrivateIps: []string{"10.0.1.6"},
									PublicIp:   to.StrPtr("123.234.345.458"),
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
			name: "ensure that validation is also performed in upgrade",
			json: []byte(`{
	"meta": {
		"version": "v0.0.2",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
}
`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "State.Meta.Kind",
					Field: "Kind",
					Tag:   "required",
				},
			},
		},
		{
			name: "upgrade v0.0.1 to v0.0.2",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.1",
		"module_version": "v1.1.1"
	},
	"status": "initialized",
	"config": {
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
	},
	"output": {
		"rg_name": "epiphany-rg",
		"vm_groups": [
			{
				"vm_group_name": "vm-group0",
				"vms": [
					{
						"private_ips": [
							"10.0.1.4"
						],
						"public_ip": "123.234.345.456",
						"vm_name": "epiphany-vm-group0-0"
					}
				]
			}
		],
		"vnet_name": "epiphany-vnet"
	}
}
`),
			want: &State{
				Meta: &Meta{
					Kind:          to.StrPtr("azbiState"),
					Version:       to.StrPtr("v0.0.2"),
					ModuleVersion: to.StrPtr("v1.1.1"),
				},
				Status: shared.Initialized,
				Config: &Config{
					Meta: &Meta{
						Kind:          to.StrPtr("azbiConfig"),
						Version:       to.StrPtr("v0.2.1"),
						ModuleVersion: to.StrPtr("v0.0.1"),
					},
					Params: &Params{
						Name:     to.StrPtr("epiphany"),
						Location: to.StrPtr("northeurope"),
						VmGroups: []VmGroup{
							{
								Name:        to.StrPtr("vm-group0"),
								VmCount:     to.IntPtr(1),
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
						RsaPublicKeyPath: to.StrPtr("some-file-name"),
					},
					Unused: nil,
				},
				Output: &Output{
					RgName:   to.StrPtr("epiphany-rg"),
					VnetName: to.StrPtr("epiphany-vnet"),
					VmGroups: []OutputVmGroup{
						{
							Name: to.StrPtr("vm-group0"),
							Vms: []OutputVm{
								{
									Name:       to.StrPtr("epiphany-vm-group0-0"),
									PrivateIps: []string{"10.0.1.4"},
									PublicIp:   to.StrPtr("123.234.345.456"),
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
			name: "some unknown version",
			json: []byte(`{
	"meta": {
		"kind": "azbiState",
		"version": "v0.0.0",
		"module_version": "dev"
	},
	"status": "initialized",
	"config": null,
	"output": null
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
			p, err := createTempDocumentFile("azbi-state-load", tt.json)
			r.NoError(err)
			got := &State{}
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
