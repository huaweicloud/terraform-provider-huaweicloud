package cts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceTags_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cts_resource_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "region"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckOutput("is_tags_exist", "true"),
					resource.TestCheckOutput("is_key_set", "true"),
					resource.TestCheckOutput("is_values_set", "true"),
				),
			},
		},
	})
}

func testAccDataSourceResourceTags_base() string {
	return `
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "tf-test-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "test" {
  bucket_name = huaweicloud_obs_bucket.test.bucket
  file_prefix = "cts"

  tags = {
    foo1 = "bar1",
    foo2 = "bar2"
  }
}
`
}
func testAccDataSourceResourceTags_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cts_resource_tags" "test" {
  resource_id   = huaweicloud_cts_tracker.test.id
  resource_type = "cts-tracker"
}

locals {
  tags = data.huaweicloud_cts_resource_tags.test.tags
}

output "is_tags_exist" {
  value = length(local.tags) > 0
}

output "is_key_set" {
  value = alltrue([for v in local.tags[*].key : v != ""])
}

output "is_values_set" {
  value = alltrue([for v in local.tags[*].value : v != ""])
}
`, testAccDataSourceResourceTags_base())
}
