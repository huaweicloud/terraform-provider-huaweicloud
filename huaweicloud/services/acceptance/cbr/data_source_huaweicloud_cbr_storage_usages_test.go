package cbr

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStorageUsages_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbr_storage_usages.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byResourceID   = "data.huaweicloud_cbr_storage_usages.filter_by_resource_id"
		dcByResourceID = acceptance.InitDataSourceCheck(byResourceID)

		byResourceType   = "data.huaweicloud_cbr_storage_usages.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataStorageUsages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.#"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.backup_count"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.backup_size"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_usages.0.backup_size_multiaz"),

					dcByResourceID.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataStorageUsages_basic = `
data "huaweicloud_cbr_storage_usages" "test" {}

data "huaweicloud_cbr_storage_usages" "filter_by_resource_id" {
  resource_id = data.huaweicloud_cbr_storage_usages.test.storage_usages[0].resource_id
}

output "resource_id_filter_is_useful" {
  value = length(data.huaweicloud_cbr_storage_usages.filter_by_resource_id.storage_usages) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_storage_usages.filter_by_resource_id.storage_usages[*] :
  v.resource_id == data.huaweicloud_cbr_storage_usages.test.storage_usages[0].resource_id])
}

data "huaweicloud_cbr_storage_usages" "filter_by_resource_type" {
  resource_type = data.huaweicloud_cbr_storage_usages.test.storage_usages[0].resource_type
}

output "resource_type_filter_is_useful" {
  value = length(data.huaweicloud_cbr_storage_usages.filter_by_resource_type.storage_usages) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_storage_usages.filter_by_resource_type.storage_usages[*] :
  v.resource_type == data.huaweicloud_cbr_storage_usages.test.storage_usages[0].resource_type])
}
`
