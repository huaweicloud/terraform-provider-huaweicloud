package taurusdb

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccTaurusDBInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTaurusDBInstanceDataSourceID("data.huaweicloud_taurusdb_instance.test"),
				),
			},
		},
	})
}

func testAccCheckTaurusDBInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find TaurusDB instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("the TaurusDB instance data source ID not set")
		}

		return nil
	}
}

func testAccTaurusDBInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {
}

resource "huaweicloud_taurusdb_instance" "test" {
  name                  = "%s"
  password              = "Test@12345678"
  flavor                = "gaussdb.mysql.4xlarge.x86.4"
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "0"
  master_availability_zone = data.huaweicloud_availability_zones.test.names[0]
  availability_zone_mode   = "multi"
}

data "huaweicloud_taurusdb_instance" "test" {
  name = huaweicloud_taurusdb_instance.test.name
  depends_on = [
    huaweicloud_taurusdb_instance.test,
  ]
}
`, common.TestBaseNetwork(rName), rName)
}
