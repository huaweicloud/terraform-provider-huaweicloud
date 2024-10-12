package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnUsers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_users.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnUsers_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceVpnUsers_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_users" "test" {
  vpn_server_id = "%[2]s"

  depends_on = [huaweicloud_vpn_user.test]
}
`, testUser_basic(name), acceptance.HW_VPN_P2C_SERVER)
}
