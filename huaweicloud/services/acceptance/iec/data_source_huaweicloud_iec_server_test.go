package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccServerDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_iec_server.server_1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServer_basic(rName),
			},
			{
				Config: testAccServerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "server-"+rName),
					resource.TestCheckResourceAttr(dataSourceName, "image_name", "Ubuntu 16.04 server 64bit"),
					resource.TestCheckResourceAttr(dataSourceName, "nics.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "volume_attached.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "volume_attached.0.boot_index", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "volume_attached.0.type", "GPSSD"),
					resource.TestCheckResourceAttr(dataSourceName, "volume_attached.0.size", "40"),
					resource.TestCheckResourceAttr(dataSourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(dataSourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "coverage_sites.0.site_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "coverage_sites.0.site_info"),
				),
			},
		},
	})
}

func testAccServerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_server" "server_1" {
  name = huaweicloud_iec_server.server_test.name
}
`, testAccServer_basic(rName))
}
