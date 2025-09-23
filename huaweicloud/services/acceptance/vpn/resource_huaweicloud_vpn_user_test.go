package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
)

func getUserFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	getUserProduct := "vpn"
	client, err := conf.NewServiceClient(getUserProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetUser(client, state.Primary.Attributes["vpn_server_id"], state.Primary.ID)
}

func testUserImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["vpn_server_id"] == "" {
			return "", fmt.Errorf("attribute (vpn_server_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (id) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["vpn_server_id"] + "/" + rs.Primary.ID, nil
	}
}

func TestAccUser_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_user.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testUser_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "password", "Typhoeus123"),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttrSet(rName, "vpn_server_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccUser_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "password", "Gaia1234"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testUserImportState(rName),
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}

func testUser_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_user" "test" {
  vpn_server_id = "%[1]s"
  name          = "%[2]s"
  password      = "Typhoeus123"
  description   = "test"
}
`, acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccUser_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_user" "test" {
  vpn_server_id = "%[1]s"
  name          = "%[2]s"
  password      = "Gaia1234"
  description   = ""
}
`, acceptance.HW_VPN_P2C_SERVER, name)
}
