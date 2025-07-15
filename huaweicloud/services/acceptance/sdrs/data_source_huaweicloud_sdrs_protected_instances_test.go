package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSdrsProtectedInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sdrs_protected_instances.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSdrsProtectedInstances_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.attachment.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.priority_station"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.server_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.source_server"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.target_server"),
					resource.TestCheckResourceAttrSet(dataSource, "protected_instances.0.updated_at"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_server_group_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSdrsProtectedInstances_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_protected_instances" "test" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]
}

# Filter by name
locals {
  name = data.huaweicloud_sdrs_protected_instances.test.protected_instances[0].name
}
data "huaweicloud_sdrs_protected_instances" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_sdrs_protected_instances.filter_by_name.protected_instances[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by server_group_id
locals {
  server_group_id = data.huaweicloud_sdrs_protected_instances.test.protected_instances[0].server_group_id
}
data "huaweicloud_sdrs_protected_instances" "filter_by_server_group_id" {
  server_group_id = local.server_group_id
}

locals {
  server_group_id_filter_result = [
    for v in data.huaweicloud_sdrs_protected_instances.filter_by_server_group_id.protected_instances[*].server_group_id :
	v == local.server_group_id
  ]
}

output "is_server_group_id_filter_useful" {
  value = length(local.server_group_id_filter_result) > 0 && alltrue(local.server_group_id_filter_result)
}

# Filter by status
locals {
  status = data.huaweicloud_sdrs_protected_instances.test.protected_instances[0].status
}
data "huaweicloud_sdrs_protected_instances" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_sdrs_protected_instances.filter_by_status.protected_instances[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testProtectedInstance_basic(name))
}
