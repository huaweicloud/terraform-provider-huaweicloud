package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceVPCEPServicePermissions_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_vpcep_service_permissions.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpcepServicePermissions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.permission_id"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.permission"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.permission_type"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.created_at"),

					resource.TestCheckOutput("permission_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceVpcepServicePermissions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_service_permissions" "test" {
  service_id = huaweicloud_vpcep_service.test.id
}

data "huaweicloud_vpcep_service_permissions" "permission_filter" {
  service_id = huaweicloud_vpcep_service.test.id
  permission = data.huaweicloud_vpcep_service_permissions.test.permissions.0.permission
}

locals {
  permission = data.huaweicloud_vpcep_service_permissions.test.permissions.0.permission
}

output "permission_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_service_permissions.permission_filter.permissions) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_service_permissions.permission_filter.permissions[*].permission : v == local.permission]
  )  
}
`, testAccVPCEPService_Basic(name))
}
