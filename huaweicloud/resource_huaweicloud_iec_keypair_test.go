package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/keypairs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
	config := testAccProvider.Meta().(*config.Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_keypair" {
			continue
		}

		_, err := keypairs.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Keypair still exists")
		}
	}

	return nil
}

func testAccCheckIECKeypairExists(n string, kp *common.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud IEC client: %s", err)
		}

		found, err := keypairs.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Name != rs.Primary.ID {
			return fmtp.Errorf("Keypair not found")
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
