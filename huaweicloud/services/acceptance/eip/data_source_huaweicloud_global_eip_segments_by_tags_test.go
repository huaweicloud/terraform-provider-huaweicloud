package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipSegmentsByTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eip_segments_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGlobalEipSegmentsByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.key", "key1"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.value", "value1"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("multiple_tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGlobalEipSegmentsByTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-south-guangzhou"
  name        = "bgp_segment_default"
}

resource "huaweicloud_global_eip_segment" "test" {
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  mask                  = 29
  name                  = "%[1]s"
  description           = "description test"
  enterprise_project_id = "%[2]s"

  tags {
    key   = "key1"
    value = "value1"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceGlobalEipSegmentsByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_eip_segments_by_tags" "test" {
  depends_on = [huaweicloud_global_eip_segment.test]

  tags {
    key   = "key1"
    value = "value1"
  }
}

data "huaweicloud_global_eip_segments_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_global_eip_segment.test]

  tags {
    key   = data.huaweicloud_global_eip_segments_by_tags.test.resources[0].tags[0].key
    value = data.huaweicloud_global_eip_segments_by_tags.test.resources[0].tags[0].value
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_global_eip_segments_by_tags.filter_by_tags.resources) > 0
}

data "huaweicloud_global_eip_segments_by_tags" "multiple_tags_filter" {
  depends_on = [huaweicloud_global_eip_segment.test]

  tags {
    key   = "key1"
    value = "value1"
  }

  tags {
    key   = "non-exist-key"
    value = "non-exist-value"
  }
}

output "multiple_tags_filter_is_useful" {
  value = length(data.huaweicloud_global_eip_segments_by_tags.multiple_tags_filter.resources) == 0
}
`, testDataSourceGlobalEipSegmentsByTags_base(name))
}
