package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIECKeypairDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("KeyPair-%s", acctest.RandString(4))
	resourceName := "data.huaweicloud_iec_keypair.by_name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIECKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceIECKeypair_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
		},
	})
}

func testAccDataSourceIECKeypair_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_keypair" "kp_1" {
  name = "%s"
}

data "huaweicloud_iec_keypair" "by_name" {
  name = huaweicloud_iec_keypair.kp_1.name
}
`, rName)
}
