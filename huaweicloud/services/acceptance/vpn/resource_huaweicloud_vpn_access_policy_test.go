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

func getAccessPolicyFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	getAccessPolicyProduct := "vpn"
	client, err := conf.NewServiceClient(getAccessPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPN client: %s", err)
	}

	return vpn.GetAccessPolicy(client, state.Primary.Attributes["vpn_server_id"], state.Primary.ID)
}

func testAccessPolicyImportState(name string) resource.ImportStateIdFunc {
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

func TestAccAccessPolicy_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_access_policy.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getAccessPolicyFunc,
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
				Config: testAccAccessPolicy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "user_group_id", "huaweicloud_vpn_user_group.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "dest_ip_cidrs.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "vpn_server_id"),
					resource.TestCheckResourceAttr(rName, "description", "test"),
					resource.TestCheckResourceAttrSet(rName, "user_group_name"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccAccessPolicy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttrPair(rName, "user_group_id", "huaweicloud_vpn_user_group.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "dest_ip_cidrs.#", "1"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccessPolicyImportState(rName),
			},
		},
	})
}

func testAccAccessPolicy_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_access_policy" "test" {
  vpn_server_id = "%[2]s"
  name          = "%[3]s"
  user_group_id = huaweicloud_vpn_user_group.test[0].id
  dest_ip_cidrs = ["192.168.0.0/16", "192.168.34.0/24"]
  description   = "test"
}
`, testAccessPolicy_base(name), acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccAccessPolicy_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpn_access_policy" "test" {
  vpn_server_id = "%[2]s"
  name          = "%[3]s-update"
  user_group_id = huaweicloud_vpn_user_group.test[1].id
  dest_ip_cidrs = ["192.168.0.0/30"]
  description   = ""
}
`, testAccessPolicy_base(name), acceptance.HW_VPN_P2C_SERVER, name)
}

func testAccessPolicy_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpn_user_group" "test" {
  count         = 2
  vpn_server_id = "%[1]s"
  name          = "%[2]s${count.index}"
  description   = "test${count.index} "
}
`, acceptance.HW_VPN_P2C_SERVER, name)
}
