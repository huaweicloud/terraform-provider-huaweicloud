package rocketmq

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dms_rocketmq_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckOutput("is_tags_exist", "true"),
					resource.TestCheckOutput("is_key_set", "true"),
					resource.TestCheckOutput("is_values_set", "true"),
				),
			},
		},
	})
}

func testAccTags_basic() string {
	return `
data "huaweicloud_dms_rocketmq_tags" "test" {}

locals {
  tags = data.huaweicloud_dms_rocketmq_tags.test.tags
}

output "is_tags_exist" {
  value = length(local.tags) >= 0
}

output "is_key_set" {
  value = length(local.tags) > 0 ? alltrue([for tag in local.tags : tag.key != ""]) : true
}

output "is_values_set" {
  value = length(local.tags) > 0 ? alltrue([for tag in local.tags : length(tag.values) >= 0]) : true
}
`
}
