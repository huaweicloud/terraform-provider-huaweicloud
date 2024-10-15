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

func getUserGroupFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	getUserGroupProduct := "vpn"
	client, err := conf.NewServiceClient(getUserGroupProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetUserGroup(client, state.Primary.Attributes["vpn_server_id"], state.Primary.ID)
}

func testUserGroupImportState(name string) resource.ImportStateIdFunc {
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

func TestAccUserGroup_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_user_group.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getUserGroupFunc,
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
				Config: testAccUserGroup_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
					resource.TestCheckResourceAttr(rName, "type", "Custom"),
					resource.TestCheckResourceAttrSet(rName, "vpn_server_id"),
					resource.TestCheckResourceAttrSet(rName, "type"),
					resource.TestCheckResourceAttrSet(rName, "user_number"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccUserGroup_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
					resource.TestCheckResourceAttr(rName, "description", "test update"),
				),
			},
			{
				Config: testAccUserGroup_disassociateAllUsers(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "users.#", "0"),
					resource.TestCheckResourceAttr(rName, "description", "test update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testUserGroupImportState(rName),
			},
		},
	})
}

func testAccUserGroup_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = "%[2]s"
  name          = "%[3]s"
  description   = "test"

  users {
    id = huaweicloud_vpn_user.user[0].id
  }

  users {
    id = huaweicloud_vpn_user.user[1].id
  }
}
`, testAccUserGroup_base(name), acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccUserGroup_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = "%[2]s"
  name          = "%[3]s-update"
  description   = "test update"

  users {
    id = huaweicloud_vpn_user.user[1].id
  }

  users {
    id = huaweicloud_vpn_user.user[2].id
  }
}
`, testAccUserGroup_base(name), acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccUserGroup_disassociateAllUsers(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = "%[2]s"
  name          = "%[3]s-update"
  description   = "test update"
}
`, testAccUserGroup_base(name), acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccUserGroup_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_user" "user" {
  count         = 3
  vpn_server_id = "%[1]s"
  name          = "%[2]s${count.index}"
  password      = "Typhoeus12${count.index}"
  description   = "test"
}
`, acceptance.HW_VPN_P2C_SERVER, name)
}
