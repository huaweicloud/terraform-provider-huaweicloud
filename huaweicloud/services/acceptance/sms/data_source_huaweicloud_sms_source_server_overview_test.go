package sms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsSourceServerOverview_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_source_server_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsSourceServerOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "waiting"),
					resource.TestCheckResourceAttrSet(dataSource, "replicate"),
					resource.TestCheckResourceAttrSet(dataSource, "syncing"),
					resource.TestCheckResourceAttrSet(dataSource, "stopped"),
					resource.TestCheckResourceAttrSet(dataSource, "deleting"),
					resource.TestCheckResourceAttrSet(dataSource, "cutovering"),
					resource.TestCheckResourceAttrSet(dataSource, "unavailable"),
					resource.TestCheckResourceAttrSet(dataSource, "stopping"),
					resource.TestCheckResourceAttrSet(dataSource, "skipping"),
					resource.TestCheckResourceAttrSet(dataSource, "finished"),
					resource.TestCheckResourceAttrSet(dataSource, "initialize"),
					resource.TestCheckResourceAttrSet(dataSource, "error"),
					resource.TestCheckResourceAttrSet(dataSource, "cloning"),
					resource.TestCheckResourceAttrSet(dataSource, "unconfigured"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsSourceServerOverview_basic() string {
	return `
data "huaweicloud_sms_source_server_overview" "test" {}
`
}
