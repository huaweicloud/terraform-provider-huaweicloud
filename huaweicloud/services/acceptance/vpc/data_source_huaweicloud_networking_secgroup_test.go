package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNetworkingSecGroupV3DataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_networking_secgroup.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func TestAccNetworkingSecGroupV3DataSource_byID(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_networking_secgroup.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupV3DataSource_secGroupID(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

func testAccNetworkingSecGroupV3DataSource_group(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s"
  description = "My neutron security group"
}
`, rName)
}

func testAccNetworkingSecGroupV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  name = huaweicloud_networking_secgroup.test.name
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}

func testAccNetworkingSecGroupV3DataSource_secGroupID(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_networking_secgroup" "test" {
  secgroup_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupV3DataSource_group(rName))
}
