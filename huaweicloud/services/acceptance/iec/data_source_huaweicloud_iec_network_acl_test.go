package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECNetworkACLDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIECNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_network_acl.by_name", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_iec_network_acl.by_id", "name", rName),
				),
			},
		},
	})
}

func testAccDataSourceIECNetworkACL_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "test" {
  name        = "%s"
  description = "IEC network acl for acc test"
}

data "huaweicloud_iec_network_acl" "by_name" {
  name = huaweicloud_iec_network_acl.test.name
}

data "huaweicloud_iec_network_acl" "by_id" {
  id = huaweicloud_iec_network_acl.test.id
}
`, rName)
}
