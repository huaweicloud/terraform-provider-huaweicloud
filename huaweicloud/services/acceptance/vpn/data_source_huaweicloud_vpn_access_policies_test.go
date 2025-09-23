package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnAccessPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_access_policies.test"
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
				Config: testDataSourceVpnAccessPolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.user_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.user_group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "access_policies.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceVpnAccessPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpn_access_policies" "test" {
  vpn_server_id = "%[2]s"

  depends_on = [huaweicloud_vpn_access_policy.test]
}
`, testAccAccessPolicy_basic(name), acceptance.HW_VPN_P2C_SERVER)
}
