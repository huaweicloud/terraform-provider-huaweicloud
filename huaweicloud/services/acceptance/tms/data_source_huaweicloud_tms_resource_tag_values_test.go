package tms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsTagValues_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_resource_tag_values.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsTagValues_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "values.#"),
					resource.TestCheckResourceAttrSet(dataSource, "values.0"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceTmsTagValues_basic() string {
	return `
data "huaweicloud_tms_resource_tag_keys" "test" {}

data "huaweicloud_tms_resource_tag_values" "test" {
  key = data.huaweicloud_tms_resource_tag_keys.test.keys[0]
}

data "huaweicloud_tms_resource_tag_values" "region_id_filter" {
  key       = data.huaweicloud_tms_resource_tag_keys.test.keys[0]
  region_id = "cn-north-4"
}
output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_tms_resource_tag_values.region_id_filter.values) > 0
}
`
}
