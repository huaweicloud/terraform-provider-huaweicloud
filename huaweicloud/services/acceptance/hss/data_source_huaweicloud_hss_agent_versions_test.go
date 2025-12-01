package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAgentVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_agent_versions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAgentVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_list.0.release_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_list.0.release_note"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version_list.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceAgentVersions_basic() string {
	return `
data "huaweicloud_hss_agent_versions" "test" {
  enterprise_project_id = "0"
  os_type               = "Linux"
}
`
}
