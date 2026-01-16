package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstancesByTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dds_instances_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstancesByTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.tags.#"),

					resource.TestCheckOutput("results_is_not_empty", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceInstancesByTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_instances_by_tags" "test" {
  action = "filter"
}

data "huaweicloud_dds_instances_by_tags" "filter_by_count" {
  action = "count"
}

data "huaweicloud_dds_instances_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "instance_id"
    value = "%s"
  }
}

data "huaweicloud_dds_instances_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = data.huaweicloud_dds_instances_by_tags.test.instances.0.tags.0.key
    values = [data.huaweicloud_dds_instances_by_tags.test.instances.0.tags.0.value]
  }
}


output "results_is_not_empty" {
  value = data.huaweicloud_dds_instances_by_tags.filter_by_count.total_count > 0
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_dds_instances_by_tags.filter_by_matches.instances) > 0
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_dds_instances_by_tags.filter_by_tags.instances) > 0
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
