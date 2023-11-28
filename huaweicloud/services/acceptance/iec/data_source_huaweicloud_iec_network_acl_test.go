package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkACLDataSource_basic(t *testing.T) {
	var (
		rName          = acceptance.RandomAccResourceName()
		resourceName   = "huaweicloud_iec_network_acl.test"
		dataSourceName = "data.huaweicloud_iec_network_acl.by_name"
		dataSourceById = "data.huaweicloud_iec_network_acl.by_id"
		fwGroup        firewalls.Firewall
		rc             = acceptance.InitResourceCheck(resourceName, &fwGroup, getNetworkACLResourceFunc)
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		dcByID         = acceptance.InitDataSourceCheck(dataSourceById)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dc.CheckResourceExists(),
					dcByID.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttr(dataSourceById, "name", rName),
				),
			},
		},
	})
}

func testAccDataSourceNetworkACL_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "test" {
  name        = "%s"
  description = "IEC network acl for acc test"
}

data "huaweicloud_iec_network_acl" "by_name" {
  name = huaweicloud_iec_network_acl.test.name
}

data "huaweicloud_iec_network_acl" "by_id" {
  id = huaweicloud_iec_network_acl.test.id
}
`, rName)
}
