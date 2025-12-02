package deh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDehInstancesByTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_deh_instances_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDehInstancesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.%"),
					resource.TestCheckResourceAttrSet("data.huaweicloud_deh_instances_by_tags.total_count_test",
						"total_count"),

					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDehInstancesByTags_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_deh_instance" "test" {
  availability_zone = "cn-south-1e"
  name              = "%[1]s"
  host_type         = "c3"
  auto_placement    = "on"

  metadata = {
    "ha_enabled" = "true"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "true"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testDataSourceDehInstancesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_deh_instances_by_tags" "test" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"
}

data "huaweicloud_deh_instances_by_tags" "matches_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances_by_tags.matches_filter.resources) > 0
}

data "huaweicloud_deh_instances_by_tags" "tags_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances_by_tags.tags_filter.resources) > 0
}

data "huaweicloud_deh_instances_by_tags" "not_tags_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances_by_tags.not_tags_filter.resources) > 0
}

data "huaweicloud_deh_instances_by_tags" "tags_any_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"

  tags_any {
    key    = "foo"
    values = ["bar"]
  }
}

output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances_by_tags.tags_any_filter.resources) > 0
}

data "huaweicloud_deh_instances_by_tags" "not_tags_any_filter" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}

output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_deh_instances_by_tags.not_tags_any_filter.resources) > 0
}

data "huaweicloud_deh_instances_by_tags" "total_count_test" {
  depends_on = [huaweicloud_deh_instance.test]

  action = "count"
}
`, testDehInstancesByTags_basic(name), name)
}
