package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceTmsTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.update_time"),
					resource.TestCheckOutput("key_filter_is_useful", "true"),
					resource.TestCheckOutput("value_filter_is_useful", "true"),
					resource.TestCheckOutput("order_field_filter_is_useful", "true"),
					resource.TestCheckOutput("order_method_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceTmsTags_base() string {
	return `
resource "huaweicloud_tms_tags" "test" {
  tags {
    key   = "key_1"
    value = "value_1"
  }
  tags {
    key   = "key_2"
    value = "value_2"
  }
  tags {
    key   = "key_2"
    value = "value_22"
  }
}
`
}

func testDataSourceDataSourceTmsTags_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_tms_tags" "test" {}

locals {
  key = tolist(huaweicloud_tms_tags.test.tags)[0].key
}
data "huaweicloud_tms_tags" "key_filter" {
  key = tolist(huaweicloud_tms_tags.test.tags)[0].key
}
output "key_filter_is_useful" {
  value = length(data.huaweicloud_tms_tags.key_filter.tags) > 0 && alltrue(
  [for v in data.huaweicloud_tms_tags.key_filter.tags[*].key : v == local.key]
  )  
}

locals {
  value = tolist(huaweicloud_tms_tags.test.tags)[0].value
}
data "huaweicloud_tms_tags" "value_filter" {
  value = tolist(huaweicloud_tms_tags.test.tags)[0].value
}
output "value_filter_is_useful" {
  value = length(data.huaweicloud_tms_tags.value_filter.tags) > 0 && alltrue(
  [for v in data.huaweicloud_tms_tags.value_filter.tags[*].value : v == local.value]
  )  
}

locals {
  order_field = "value"
}
data "huaweicloud_tms_tags" "order_field_filter" {
  order_field = "value"
}
output "order_field_filter_is_useful" {
  value = length(data.huaweicloud_tms_tags.order_field_filter.tags) > 0 
}

locals {
  order_method = "asc"
}
data "huaweicloud_tms_tags" "order_method_filter" {
  order_method = "asc"
}
output "order_method_filter_is_useful" {
  value = length(data.huaweicloud_tms_tags.order_method_filter.tags) > 0 
}
`, testDataSourceDataSourceTmsTags_base())
}
