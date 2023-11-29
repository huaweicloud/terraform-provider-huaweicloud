package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/keypairs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getKeypairResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := cfg.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IEC client: %s", err)
	}
	return keypairs.Get(c, state.Primary.ID).Extract()
}

func TestAccKeypairResource_basic(t *testing.T) {
	var (
		keypair      common.KeyPair
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_iec_keypair.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&keypair,
		getKeypairResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccKeypair_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
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

func testAccKeypair_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_keypair" "test" {
  name = "%s"
}
`, name)
}
