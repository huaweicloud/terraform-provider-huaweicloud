package identitycenter

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentitycenterPermissionSets_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identitycenter_permission_sets.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceIdentitycenterPermissionSets_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.0.permission_set_id"),
					resource.TestCheckResourceAttrSet(dataSource, "permission_sets.0.name"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceIdentitycenterPermissionSets_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_identitycenter_permission_sets" "test" {
  instance_id = data.huaweicloud_identitycenter_instance.system.id
}

locals {
  permission_set_id = data.huaweicloud_identitycenter_permission_sets.test.permission_sets[0].permission_set_id
  name              = data.huaweicloud_identitycenter_permission_sets.test.permission_sets[0].name
}

data "huaweicloud_identitycenter_permission_sets" "filter_by_id" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = local.permission_set_id
}

data "huaweicloud_identitycenter_permission_sets" "filter_by_name" {
  instance_id = data.huaweicloud_identitycenter_instance.system.id
  name        = local.name
}

locals {
  list_by_id   = data.huaweicloud_identitycenter_permission_sets.filter_by_id.permission_sets
  list_by_name = data.huaweicloud_identitycenter_permission_sets.filter_by_name.permission_sets
}
  
output "is_id_filter_useful" {
  value = length(local.list_by_id) > 0 && alltrue(
    [for v in local.list_by_id[*].permission_set_id : v == local.permission_set_id]
  )
}

output "is_name_filter_useful" {
  value = length(local.list_by_name) > 0 && alltrue(
    [for v in local.list_by_name[*].name : v == local.name]
  )
}
`, testPermissionSet_basic(name))
}
