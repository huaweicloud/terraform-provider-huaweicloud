package esw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEswConnectionVportBind_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEswConnectionVportBind_basic(rName),
			},
		},
	})
}

func testAccEswConnectionVportBind_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_subnet_private_ip" "test" {
  subnet_id    = huaweicloud_vpc_subnet.test[1].id
  device_owner = "neutron:VIP_PORT"
}

resource "huaweicloud_esw_connection_vport_bind" "test" {
  connection_id = huaweicloud_esw_connection.test.id
  port_id       = huaweicloud_vpc_subnet_private_ip.test.id
}
`, testAccEswConnection_basic(rName), rName)
}
