package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceHighRiskCommands_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_high_risk_commands.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHighRiskCommands_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "commands.#"),
					resource.TestCheckResourceAttrSet(dataSource, "commands.0.origin_name"),
					resource.TestCheckResourceAttrSet(dataSource, "commands.0.name"),
				),
			},
		},
	})
}

func testAccDatasourceHighRiskCommands_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data"huaweicloud_availability_zones" "test" {}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test_357159"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "2"
    size      = "4"
    storage   = "ULTRAHIGH"
    spec_code = "geminidb.redis.medium.2"
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccDataSourceHighRiskCommands_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_geminidb_high_risk_commands" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
}
`, testAccDatasourceHighRiskCommands_base(name))
}
