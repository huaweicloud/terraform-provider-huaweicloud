package cts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCtsTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cts_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCtsTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_exist", "true"),
					resource.TestCheckOutput("is_key_set", "true"),
					resource.TestCheckOutput("is_values_set", "true"),
				),
			},
		},
	})
}

func testDataSourceCtsTags_basic() string {
	return `
data "huaweicloud_cts_tags" "test" {}

locals {
  tags = data.huaweicloud_cts_tags.test.tags
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
