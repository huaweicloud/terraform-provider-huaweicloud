package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnUserGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_user_groups.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnUserGroups_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "user_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "user_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "user_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "user_groups.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "user_groups.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceVpnUserGroups_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_user_groups" "test" {
  vpn_server_id = "%[2]s"

  depends_on = [huaweicloud_vpn_user_group.test]
}
`, testDataSourceVpnUserGroups_base(), acceptance.HW_VPN_P2C_SERVER)
}

func testDataSourceVpnUserGroups_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = "%[1]s"
  name          = "%[2]s"
  description   = "test"
}
`, acceptance.HW_VPN_P2C_SERVER, name)
}
