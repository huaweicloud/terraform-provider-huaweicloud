package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGlobalEipsByTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_global_eips_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGlobalEipsByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.key", "foo"),
					resource.TestCheckResourceAttr(dataSource, "resources.0.tags.0.value", "bar"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("multiple_tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGlobalEipsByTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_global_eip_pools" "all" {
  access_site = "cn-north-beijing"
  ip_version  = 4
}

resource "huaweicloud_global_internet_bandwidth" "test" {
  access_site = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  charge_mode = "95peak_guar"
  size        = 300
  isp         = data.huaweicloud_global_eip_pools.all.geip_pools[0].isp
  name        = "%[1]s"
  type        = data.huaweicloud_global_eip_pools.all.geip_pools[0].allowed_bandwidth_types[0].type

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_global_eip" "test" {
  access_site           = data.huaweicloud_global_eip_pools.all.geip_pools[0].access_site
  geip_pool_name        = data.huaweicloud_global_eip_pools.all.geip_pools[0].name
  internet_bandwidth_id = huaweicloud_global_internet_bandwidth.test.id
  name                  = "%[1]s"

  tags = {
    foo = "bar"
  }
}
`, name)
}

func testDataSourceGlobalEipsByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_global_eips_by_tags" "test" {
  depends_on = [huaweicloud_global_eip.test]

  tags {
    key   = "foo"
    value = "bar"
  }
}

data "huaweicloud_global_eips_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_global_eip.test]

  tags {
    key   = data.huaweicloud_global_eips_by_tags.test.resources[0].tags[0].key
    value = data.huaweicloud_global_eips_by_tags.test.resources[0].tags[0].value
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_global_eips_by_tags.filter_by_tags.resources) > 0
}

data "huaweicloud_global_eips_by_tags" "multiple_tags_filter" {
  depends_on = [huaweicloud_global_eip.test]

  tags {
    key   = "foo"
    value = "bar"
  }

  tags {
    key   = "non-exist-key"
    value = "non-exist-value"
  }
}

output "multiple_tags_filter_is_useful" {
  value = length(data.huaweicloud_global_eips_by_tags.multiple_tags_filter.resources) == 0
}
`, testDataSourceGlobalEipsByTags_base(name))
}
