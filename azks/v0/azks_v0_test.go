package v0

import (
	"testing"

	"github.com/epiphany-platform/e-structures/utils/test"
	"github.com/epiphany-platform/e-structures/utils/to"
	"github.com/go-playground/validator/v10"
	"github.com/google/go-cmp/cmp"
)

// TestConfig_Load_general contains all general types of scenarios: happy path, unknown fields,
// kind and version validation, minimal correct and full json.
func TestConfig_Load_general(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "happy path",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr("Standard_DS2_v2"),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr("VirtualMachineScaleSets"),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr("600"),
						ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
						ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
						ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
						ScanInterval:                  to.StrPtr("10s"),
						ScaleDownUnneeded:             to.StrPtr("10m"),
						ScaleDownUnready:              to.StrPtr("10m"),
						ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
					},
					AzureAd:       nil,
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "unknown fields in multiple places",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"extra_outer_field" : "extra_outer_value",
	"params": {
		"name": "epiphany",
		"extra_inner_field" : "extra_inner_value",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null,
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr("Standard_DS2_v2"),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr("VirtualMachineScaleSets"),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr("600"),
						ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
						ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
						ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
						ScanInterval:                  to.StrPtr("10s"),
						ScaleDownUnneeded:             to.StrPtr("10m"),
						ScaleDownUnready:              to.StrPtr("10m"),
						ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
					},
					AzureAd:       nil,
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{"params.extra_inner_field", "extra_outer_field"},
			},
			wantErr: nil,
		},
		{
			name: "empty json",
			json: []byte(`{}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Kind",
					Field: "Kind",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
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
			name: "major version mismatch",
			json: []byte(`{
"kind": "azks",
"version": "100.0.0",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "some-name",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Version",
					Field: "Version",
					Tag:   "version",
				},
			},
		},
		{
			name: "full json",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr("Standard_DS2_v2"),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr("VirtualMachineScaleSets"),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr("600"),
						ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
						ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
						ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
						ScanInterval:                  to.StrPtr("10s"),
						ScaleDownUnneeded:             to.StrPtr("10m"),
						ScaleDownUnready:              to.StrPtr("10m"),
						ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
					},
					AzureAd: &AzureAd{
						Managed:             to.BoolPtr(true),
						TenantId:            to.StrPtr("123123123123"),
						AdminGroupObjectIds: []string{"123123123123"},
					},
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

// TestConfig_Load_params contains all scenarios related to validation of values stored directly in Params structure.
func TestConfig_Load_Params(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "missing params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}
	}
}`),
			want: nil,
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
					Key:   "Config.Params.RgName",
					Field: "RgName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.VnetName",
					Field: "VnetName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.SubnetName",
					Field: "SubnetName",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.KubernetesVersion",
					Field: "KubernetesVersion",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.EnableNodePublicIp",
					Field: "EnableNodePublicIp",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.EnableRbac",
					Field: "EnableRbac",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.IdentityType",
					Field: "IdentityType",
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
			name: "empty params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "",
		"location": "",
		"rsa_pub_path": "",
		"rg_name": "",
		"vnet_name": "",
		"subnet_name": "",
		"kubernetes_version": "",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "",
		"admin_username": ""
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
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
					Key:   "Config.Params.RsaPublicKeyPath",
					Field: "RsaPublicKeyPath",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.RgName",
					Field: "RgName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.VnetName",
					Field: "VnetName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.SubnetName",
					Field: "SubnetName",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.KubernetesVersion",
					Field: "KubernetesVersion",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.IdentityType",
					Field: "IdentityType",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AdminUsername",
					Field: "AdminUsername",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

// TestConfig_Load_DefaultNodePool contains scenarios related to validation of values stored in DefaultNodePool structure.
func TestConfig_Load_DefaultNodePool(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "missing default_node_pool",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool",
					Field: "DefaultNodePool",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty default_node_pool aka missing params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.VmSize",
					Field: "VmSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.DiskGbSize",
					Field: "DiskGbSize",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.AutoScaling",
					Field: "AutoScaling",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Type",
					Field: "Type",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty default_node_pool params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2, 
			"min": 2,
			"max": 5,
			"vm_size": "",
			"disk_gb_size": 0,
			"auto_scaling": true,
			"type": ""
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.VmSize",
					Field: "VmSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.DiskGbSize",
					Field: "DiskGbSize",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Type",
					Field: "Type",
					Tag:   "min",
				},
			},
		},
		{
			name: "missing default_node_pool.min",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "gtefield",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "gtefield",
				},
			},
		},
		{
			name: "missing default_node_pool.max",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "ltefield",
				},
			},
		},
		{
			name: "default_node_pool min > max",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 1, 
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "gtefield",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "ltefield",
				},
			},
		},
		{
			name: "default_node_pool size < min",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 1,
			"min": 2,
			"max": 3, 
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "gtefield",
				},
			},
		},
		{
			name: "default_node_pool size > max",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 4,
			"min": 2,
			"max": 3, 
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "ltefield",
				},
			},
		},
		{
			name: "default_node_pool negative sizes",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": -1,
			"min": -1,
			"max": -1, 
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Min",
					Field: "Min",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Max",
					Field: "Max",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.DefaultNodePool.Size",
					Field: "Size",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

// TestConfig_Load_AutoScalerProfile contains scenarios related to validation of values stored in AutoScalerProfile structure.
func TestConfig_Load_AutoScalerProfile(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "missing auto_scaler_profile",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile",
					Field: "AutoScalerProfile",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty auto_scaler_profile aka missing params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.BalanceSimilarNodeGroups",
					Field: "BalanceSimilarNodeGroups",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.MaxGracefulTerminationSec",
					Field: "MaxGracefulTerminationSec",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterAdd",
					Field: "ScaleDownDelayAfterAdd",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterDelete",
					Field: "ScaleDownDelayAfterDelete",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterFailure",
					Field: "ScaleDownDelayAfterFailure",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScanInterval",
					Field: "ScanInterval",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnneeded",
					Field: "ScaleDownUnneeded",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnready",
					Field: "ScaleDownUnready",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUtilizationThreshold",
					Field: "ScaleDownUtilizationThreshold",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty auto_scaler_profile params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "",
			"scale_down_delay_after_add": "",
			"scale_down_delay_after_delete": "",
			"scale_down_delay_after_failure": "",
			"scan_interval": "",
			"scale_down_unneeded": "",
			"scale_down_unready": "",
			"scale_down_utilization_threshold": ""
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				"123123123123"
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.MaxGracefulTerminationSec",
					Field: "MaxGracefulTerminationSec",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterAdd",
					Field: "ScaleDownDelayAfterAdd",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterDelete",
					Field: "ScaleDownDelayAfterDelete",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownDelayAfterFailure",
					Field: "ScaleDownDelayAfterFailure",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScanInterval",
					Field: "ScanInterval",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnneeded",
					Field: "ScaleDownUnneeded",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUnready",
					Field: "ScaleDownUnready",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AutoScalerProfile.ScaleDownUtilizationThreshold",
					Field: "ScaleDownUtilizationThreshold",
					Tag:   "min",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

// TestConfig_Load_AzureAd contains scenarios related to validation of values stored in AzureAd structure.
func TestConfig_Load_AzureAd(t *testing.T) {
	tests := []struct {
		name    string
		json    []byte
		want    *Config
		wantErr error
	}{
		{
			name: "missing azure_ad",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr("Standard_DS2_v2"),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr("VirtualMachineScaleSets"),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr("600"),
						ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
						ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
						ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
						ScanInterval:                  to.StrPtr("10s"),
						ScaleDownUnneeded:             to.StrPtr("10m"),
						ScaleDownUnready:              to.StrPtr("10m"),
						ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
					},
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "null azure_ad",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": null, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: &Config{
				Kind:    to.StrPtr("azks"),
				Version: to.StrPtr("v0.0.1"),
				Params: &Params{
					Location:           to.StrPtr("northeurope"),
					Name:               to.StrPtr("epiphany"),
					RsaPublicKeyPath:   to.StrPtr("/shared/vms_rsa.pub"),
					RgName:             to.StrPtr("epiphany-rg"),
					VnetName:           to.StrPtr("epiphany-vnet"),
					SubnetName:         to.StrPtr("azks"),
					KubernetesVersion:  to.StrPtr("1.18.14"),
					EnableNodePublicIp: to.BoolPtr(false),
					EnableRbac:         to.BoolPtr(false),
					DefaultNodePool: &DefaultNodePool{
						Size:        to.IntPtr(2),
						Min:         to.IntPtr(2),
						Max:         to.IntPtr(5),
						VmSize:      to.StrPtr("Standard_DS2_v2"),
						DiskGbSize:  to.IntPtr(36),
						AutoScaling: to.BoolPtr(true),
						Type:        to.StrPtr("VirtualMachineScaleSets"),
					},
					AutoScalerProfile: &AutoScalerProfile{
						BalanceSimilarNodeGroups:      to.BoolPtr(false),
						MaxGracefulTerminationSec:     to.StrPtr("600"),
						ScaleDownDelayAfterAdd:        to.StrPtr("10m"),
						ScaleDownDelayAfterDelete:     to.StrPtr("10s"),
						ScaleDownDelayAfterFailure:    to.StrPtr("10m"),
						ScanInterval:                  to.StrPtr("10s"),
						ScaleDownUnneeded:             to.StrPtr("10m"),
						ScaleDownUnready:              to.StrPtr("10m"),
						ScaleDownUtilizationThreshold: to.StrPtr("0.5"),
					},
					IdentityType:  to.StrPtr("SystemAssigned"),
					AdminUsername: to.StrPtr("operations"),
				},
				Unused: []string{},
			},
			wantErr: nil,
		},
		{
			name: "empty azure_ad aka missing params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.Managed",
					Field: "Managed",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.TenantId",
					Field: "TenantId",
					Tag:   "required",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds",
					Field: "AdminGroupObjectIds",
					Tag:   "required",
				},
			},
		},
		{
			name: "empty azure_ad params",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "",
			"admin_group_object_ids": []
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.TenantId",
					Field: "TenantId",
					Tag:   "min",
				},
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds",
					Field: "AdminGroupObjectIds",
					Tag:   "min",
				},
			},
		},
		{
			name: "empty azure_ad.admin_group_object_ids element",
			json: []byte(`{
	"kind": "azks",
	"version": "v0.0.1",
	"params": {
		"name": "epiphany",
		"location": "northeurope",
		"rsa_pub_path": "/shared/vms_rsa.pub",
		"rg_name": "epiphany-rg",
		"vnet_name": "epiphany-vnet",
		"subnet_name": "azks",
		"kubernetes_version": "1.18.14",
		"enable_node_public_ip": false,
		"enable_rbac": false,
		"default_node_pool": {
			"size": 2,
			"min": 2,
			"max": 5,
			"vm_size": "Standard_DS2_v2",
			"disk_gb_size": 36,
			"auto_scaling": true,
			"type": "VirtualMachineScaleSets"
		},
		"auto_scaler_profile": {
			"balance_similar_node_groups": false,
			"max_graceful_termination_sec": "600",
			"scale_down_delay_after_add": "10m",
			"scale_down_delay_after_delete": "10s",
			"scale_down_delay_after_failure": "10m",
			"scan_interval": "10s",
			"scale_down_unneeded": "10m",
			"scale_down_unready": "10m",
			"scale_down_utilization_threshold": "0.5"
		},
		"azure_ad": {
			"managed": true,
			"tenant_id": "123123123123",
			"admin_group_object_ids": [
				""
			]
		}, 
		"identity_type": "SystemAssigned",
		"admin_username": "operations"
	}
}`),
			want: nil,
			wantErr: test.TestValidationErrors{
				test.TestValidationError{
					Key:   "Config.Params.AzureAd.AdminGroupObjectIds[0]",
					Field: "AdminGroupObjectIds[0]",
					Tag:   "required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configLoadTestingBody(t, tt.json, tt.want, tt.wantErr)
		})
	}
}

func configLoadTestingBody(t *testing.T, json []byte, want *Config, wantErr error) {
	got := &Config{}
	err := got.Unmarshal(json)

	if wantErr != nil {

		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				t.Fatal(err)
			}
			errs := err.(validator.ValidationErrors)
			if len(errs) != len(wantErr.(test.TestValidationErrors)) {
				t.Fatalf("incorrect length of found errors. Got: \n%s\nExpected: \n%s", errs.Error(), wantErr.Error())
			}
			for _, e := range errs {
				found := false
				for _, we := range wantErr.(test.TestValidationErrors) {
					if we.Key == e.Namespace() && we.Tag == e.Tag() && we.Field == e.Field() {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Got unknown error:\n%s\nAll expected errors: \n%s", e.Error(), wantErr.Error())
				}
			}
		} else {
			t.Errorf("No errors got. All expected errors: \n%s", wantErr.Error())
		}
	} else {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Unmarshal() mismatch (-want +got):\n%s", diff)
		}
		if err != nil {
			t.Errorf("Unmarshal() unexpected error occured: %v", err)
		}
	}
}
