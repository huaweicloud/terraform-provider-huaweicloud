package iec

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSitesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_iec_sites.sites_test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSitesConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "sites.#"),
					resource.TestCheckResourceAttr(dataSourceName, "sites.0.area", "east"),
					resource.TestCheckResourceAttrSet(dataSourceName, "sites.0.lines.#"),
				),
			},
		},
	})
}

func testAccSitesConfig_basic() string {
	return `
data "huaweicloud_iec_sites" "sites_test" {
  area = "east"
}
`
}
