package v0

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func AzBISubnetsValidation(sl validator.StructLevel) {
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
					return
				}
				found := false
				for _, s := range params.Subnets {
					if s.Name != nil && sn == *s.Name {
						found = true
					}
				}
				if !found {
					sl.ReportError(
						params.VmGroups[i].SubnetNames[j],
						fmt.Sprintf("VmGroups[%d].SubnetNames[%d]", i, j),
						fmt.Sprintf("SubnetNames[%d]", j),
						"insubnets",
						"")
					return
				}
			}
		}
	}
}
