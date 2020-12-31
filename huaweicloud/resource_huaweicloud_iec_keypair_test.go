package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/keypairs"
)

func TestAccIECKeypairResource_basic(t *testing.T) {
	var keypair common.KeyPair
	resourceName := "huaweicloud_iec_keypair.kp_1"
	rName := fmt.Sprintf("KeyPair-%s", acctest.RandString(4))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIECKeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIECKeypair_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIECKeypairExists(resourceName, &keypair),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "public_key"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIECKeypairDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_keypair" {
			continue
		}

		_, err := keypairs.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Keypair still exists")
		}
	}

	return nil
}

func testAccCheckIECKeypairExists(n string, kp *common.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
		}

		found, err := keypairs.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmt.Errorf("Keypair not found")
		}

		*kp = *found

		return nil
	}
}

func testAccIECKeypair_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_keypair" "kp_1" {
  name = "%s"
}
`, rName)
}
