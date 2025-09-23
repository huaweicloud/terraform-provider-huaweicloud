package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rms_tags.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.value.#"),
					resource.TestCheckOutput("key_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsTags_basic() string {
	return `
data "huaweicloud_rms_tags" "test" {}

data "huaweicloud_rms_tags" "key_filter" {
  key = data.huaweicloud_rms_tags.test.tags[0].key
}
locals {
  key = data.huaweicloud_rms_tags.test.tags[0].key
}
output "key_filter_is_useful" {
  value = length(data.huaweicloud_rms_tags.key_filter.tags) > 0 && alltrue(
  [for v in data.huaweicloud_rms_tags.key_filter.tags[*].key : v == local.key]
  )
}
`
}
