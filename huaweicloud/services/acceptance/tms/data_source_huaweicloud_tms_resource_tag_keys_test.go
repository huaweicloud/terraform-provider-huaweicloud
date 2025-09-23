package tms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsTagKeys_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_resource_tag_keys.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsTagKeys_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keys.#"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTmsTagKeys_basic() string {
	return `
data "huaweicloud_tms_resource_tag_keys" "test" {}

data "huaweicloud_tms_resource_tag_keys" "region_id_filter" {
  region_id = "cn-north-4"
}
output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_tms_resource_tag_keys.region_id_filter.keys) > 0
}
`
}
