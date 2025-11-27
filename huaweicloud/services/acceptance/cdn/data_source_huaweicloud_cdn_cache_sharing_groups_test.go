package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCacheSharingGroups_basic(t *testing.T) {
	var (
		rName     = "data.huaweicloud_cdn_cache_sharing_groups.test"
		dc        = acceptance.InitDataSourceCheck(rName)
		groupName = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdntDomainNames(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCacheSharingGroups_basic(groupName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "groups.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(rName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.group_name"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.primary_domain"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.share_cache_records.#"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.create_time"),
					resource.TestMatchResourceAttr(rName, "groups.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckOutput("is_cache_sharing_group_found", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCacheSharingGroups_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cdn_cache_sharing_group" "test" {
  name           = "%[1]s"
  primary_domain = try(element(split(",", "%[2]s"), 0), "")

  dynamic "share_cache_records" {
    for_each = split(",", "%[2]s")

    content {
      domain_name = share_cache_records.value
    }
  }
}
`, name, acceptance.HW_CDN_DOMAIN_NAMES)
}

func testAccDataSourceCacheSharingGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cdn_cache_sharing_groups" "test" {
  depends_on = [
    huaweicloud_cdn_cache_sharing_group.test
  ]
}

locals {
  cache_sharing_group_id           = huaweicloud_cdn_cache_sharing_group.test.id
  cache_sharing_group_query_result = try([
    for v in data.huaweicloud_cdn_cache_sharing_groups.test.groups : v if v.id == local.cache_sharing_group_id
  ][0], null)
}

output "is_cache_sharing_group_found" {
  value = local.cache_sharing_group_query_result != null && local.cache_sharing_group_query_result.id == local.cache_sharing_group_id
}
`, testAccDataSourceCacheSharingGroups_base(name))
}
