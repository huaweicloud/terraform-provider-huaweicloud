package drs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSupportLinks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_drs_support_links.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSupportLinks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "support_links.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("job_type_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_links.0.engine_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_links.0.net_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_links.0.task_modes.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_links.0.job_direction"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_links.0.cluster_mode"),
				),
			},
		},
	})
}

func testAccDataSourceSupportLinks_basic() string {
	return `
data "huaweicloud_drs_support_links" "test" {
  job_type = "migration"
}

output "job_type_filter_is_useful" {
  value = length(data.huaweicloud_drs_support_links.test.support_links) > 0
}
`
}
