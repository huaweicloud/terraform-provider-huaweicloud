package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccGeminiDBInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBInstancesDataSourceID("data.huaweicloud_gaussdb_cassandra_instances.test"),
					resource.TestCheckResourceAttr("data.huaweicloud_gaussdb_cassandra_instances.test", "instances.#", "1"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB cassandra instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB cassandra instances data source ID not set ")
		}

		return nil
	}
}

func testAccGeminiDBInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_gaussdb_cassandra_instance" "test" {
  name        = "%s"
  password    = "Test@12345678"
  flavor      = "geminidb.cassandra.xlarge.4"
  volume_size = 100
  vpc_id      = huaweicloud_vpc.test.id
  subnet_id   = huaweicloud_vpc_subnet.test.id
  node_num    = 4

  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  backup_strategy {
    start_time = "03:00-04:00"
    keep_days  = 14
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_gaussdb_cassandra_instances" "test" {
  name = huaweicloud_gaussdb_cassandra_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_cassandra_instance.test,
  ]
}
`, common.TestBaseNetwork(rName), rName)
}
