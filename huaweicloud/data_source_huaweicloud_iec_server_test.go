package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECServerDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_iec_server.server_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecServerDestory,
		Steps: []resource.TestStep{
			{
				Config: testAccIecServer_basic(rName),
			},
			{
				Config: testAccIECServerDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "server-"+rName),
					resource.TestCheckResourceAttr(resourceName, "image_name", "Ubuntu 16.04 server 64bit"),
					resource.TestCheckResourceAttr(resourceName, "nics.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "volume_attached.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "volume_attached.0.boot_index", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_attached.0.type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "volume_attached.0.size", "40"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "public_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "coverage_sites.0.site_id"),
					resource.TestCheckResourceAttrSet(resourceName, "coverage_sites.0.site_info"),
				),
			},
		},
	})
}

func testAccIECServerDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iec_server" "server_1" {
  name = huaweicloud_iec_server.server_test.name
}
`, testAccIecServer_basic(rName))
}
