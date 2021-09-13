package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/keypairs"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeV2Keypair_basic(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_keypair.test"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_basic(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
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

func TestAccComputeV2Keypair_privateKey(t *testing.T) {
	var keypair keypairs.KeyPair
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_keypair.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2KeypairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Keypair_privateKey(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2KeypairExists(resourceName, &keypair),
					resource.TestCheckResourceAttrSet(resourceName, "key_file"),
				),
			},
		},
	})
}

func testAccCheckComputeV2KeypairDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_keypair" {
			continue
		}

		_, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Keypair still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2KeypairExists(n string, kp *keypairs.KeyPair) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := keypairs.Get(computeClient, rs.Primary.ID).Extract()
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

func testAccComputeV2Keypair_basic(rName, keypair string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_keypair" "test" {
  name       = "%s"
  public_key = "%s"
}
`, rName, keypair)
}

func testAccComputeV2Keypair_privateKey(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_keypair" "test" {
  name = "%s"
}
`, rName)
}
