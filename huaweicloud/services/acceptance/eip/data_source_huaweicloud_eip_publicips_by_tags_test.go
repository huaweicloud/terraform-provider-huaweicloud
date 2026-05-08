package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePublicipsByTags_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_eip_publicips_by_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Need to ensure that the set EIP resource have created labels.
			acceptance.TestAccPrecheckEipIDAndIP(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePublicipsByTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),

					resource.TestCheckOutput("results_is_not_empty", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePublicipsByTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_eip_publicips_by_tags" "test" {
  action = "filter"
}

data "huaweicloud_eip_publicips_by_tags" "filter_by_count" {
  action = "count"
}

data "huaweicloud_eip_publicips_by_tags" "filter_by_matches" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = "%[1]s"
  }
}

data "huaweicloud_eip_publicips_by_tags" "filter_by_tags" {
  action = "filter"

  tags {
    key    = data.huaweicloud_eip_publicips_by_tags.test.resources.0.tags.0.key
    values = [data.huaweicloud_eip_publicips_by_tags.test.resources.0.tags.0.value]
  }
}

output "results_is_not_empty" {
  value = data.huaweicloud_eip_publicips_by_tags.filter_by_count.total_count > 0
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_eip_publicips_by_tags.filter_by_matches.resources) > 0
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_eip_publicips_by_tags.filter_by_tags.resources) > 0
}
`, acceptance.HW_EIP_ADDRESS)
}
