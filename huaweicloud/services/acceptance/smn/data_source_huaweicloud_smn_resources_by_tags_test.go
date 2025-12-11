package smn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmnResourcesByTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_smn_resources_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSmnResourcesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.detail_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.0.value"),

					resource.TestCheckResourceAttrSet("data.huaweicloud_smn_resources_by_tags.total_count_test",
						"total_count"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSmnResourcesByTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "topic_1" {
  name = "%[1]s_1"
}

resource "huaweicloud_smn_topic" "topic_2" {
  name = "%[1]s_2"

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_smn_topic" "topic_3" {
  name = "%[1]s_3"

  tags = {
    terraform = "test"
    key       = "value"
  }
}
`, name)
}

func testDataSourceSmnResourcesByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_smn_resources_by_tags" "test" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}

data "huaweicloud_smn_resources_by_tags" "tags_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }
}
output "tags_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.tags_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "tags_any_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  tags_any {
    key    = "foo"
    values = ["bar"]
  }
}
output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.tags_any_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "not_tags_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}
output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.not_tags_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "not_tags_any_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  not_tags {
    key    = "no_foo"
    values = ["no_bar"]
  }
}
output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.not_tags_any_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "matches_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "filter"

  matches {
    key   = "resource_name"
    value = "%[2]s"
  }
}
output "matches_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.matches_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "without_any_tag_filter" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type   = "smn_topic"
  action          = "filter"
  without_any_tag = "true"
}
output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_smn_resources_by_tags.without_any_tag_filter.resources) > 0
}

data "huaweicloud_smn_resources_by_tags" "total_count_test" {
  depends_on = [
    huaweicloud_smn_topic.topic_1,
    huaweicloud_smn_topic.topic_2,
    huaweicloud_smn_topic.topic_3
  ]

  resource_type = "smn_topic"
  action        = "count"
}
`, testDataSourceSmnResourcesByTags_base(name), name)
}
