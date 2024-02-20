package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDdsStorageTypes_basic(t *testing.T) {
	rName := "data.huaweicloud_dds_storage_types.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceStoragetype_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.name"),
					resource.TestCheckResourceAttrSet(rName, "storage_types.0.az_status.%"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceStoragetype_basic() string {
	return `

data "huaweicloud_dds_storage_types" "test" {}

data "huaweicloud_dds_storage_types" "name_filter" {
  engine_name = "DDS-Community"
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_dds_storage_types.name_filter.storage_types) > 0 
}
`
}
