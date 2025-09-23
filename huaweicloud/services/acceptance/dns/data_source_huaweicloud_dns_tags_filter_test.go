package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDnsTagsFilter_basic(t *testing.T) {
	var (
		name   = fmt.Sprintf("acpttest-zone-%s.com", acctest.RandString(5))
		dcName = "data.huaweicloud_dns_tags_filter.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDnsTagsFilter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "resources.#"),
					resource.TestCheckResourceAttrSet(dcName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dcName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dcName, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dcName, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dcName, "resources.0.tags.0.value"),

					resource.TestCheckOutput("is_all_tags_result_useful", "true"),
					resource.TestCheckOutput("is_any_tags_result_useful", "true"),
					resource.TestCheckOutput("is_without_all_tags_result_useful", "true"),
					resource.TestCheckOutput("is_without_any_tags_result_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDnsTagsFilter_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  tag_key   = keys(huaweicloud_dns_zone.test.tags)[0]
  tag_value = huaweicloud_dns_zone.test.tags[local.tag_key]
}

data "huaweicloud_dns_tags_filter" "test" {
  depends_on = [huaweicloud_dns_zone.test]

  resource_type = "DNS-public_zone"
}

# filter by all tags 
data "huaweicloud_dns_tags_filter" "filter_by_all_tags" {
  resource_type = "DNS-public_zone"

  tags {
    key    = local.tag_key
    values = [local.tag_value]
  }
}

locals {
  all_tag_keys = flatten([
    for v in data.huaweicloud_dns_tags_filter.filter_by_all_tags.resources : [
      for tag in v.tags : tag.key
    ]
  ])
}

output "is_all_tags_result_useful" {
  value = contains(local.all_tag_keys, local.tag_key)
}

# filter by any tags
data "huaweicloud_dns_tags_filter" "filter_by_any_tags" {
  resource_type = "DNS-public_zone"

  tags_any {
    key    = local.tag_key
    values = [local.tag_value]
  }
}

locals {
  any_tag_keys = flatten([
    for v in data.huaweicloud_dns_tags_filter.filter_by_any_tags.resources : [
      for tag in v.tags : tag.key
    ]
  ])
}

output "is_any_tags_result_useful" {
  value = contains(local.any_tag_keys, local.tag_key)
}

# filter by without all tags
data "huaweicloud_dns_tags_filter" "filter_by_without_all_tags" {
  resource_type = "DNS-public_zone"

  not_tags {
    key    = local.tag_key
    values = [local.tag_value]
  }
}

locals {
  without_all_tag_keys = flatten([
    for v in data.huaweicloud_dns_tags_filter.filter_by_without_all_tags.resources : [
      for tag in v.tags : tag.key
    ]
  ])
}

output "is_without_all_tags_result_useful" {
  value = !contains(local.without_all_tag_keys, local.tag_key)
}

# filter by without any tags
data "huaweicloud_dns_tags_filter" "filter_by_without_any_tags" {
  resource_type = "DNS-public_zone"

  not_tags_any {
    key    = local.tag_key
    values = [local.tag_value]
  }
}

locals {
  without_any_tag_keys = flatten([
    for v in data.huaweicloud_dns_tags_filter.filter_by_without_any_tags.resources : [
      for tag in v.tags : tag.key
    ]
  ])
}

output "is_without_any_tags_result_useful" {
  value = !contains(local.without_any_tag_keys, local.tag_key)
}
`, testAccDNSZone_basic(name))
}
