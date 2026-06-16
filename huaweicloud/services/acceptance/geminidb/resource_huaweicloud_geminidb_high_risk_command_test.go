package geminidb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/geminidb"
)

func getResourceHighRiskCommandFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	resourceId := strings.Split(state.Primary.ID, "/")

	return geminidb.GetHighRiskCommandInfo(client, resourceId[0], resourceId[1])
}

func TestAccResourceHighRiskCommand_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_geminidb_high_risk_command.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceHighRiskCommandFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccHighRiskCommand_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "origin_name", "keys"),
					resource.TestCheckResourceAttr(rName, "name", "test"),
				),
			},
			{
				Config: testAccHighRiskCommand_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "origin_name", "keys"),
					resource.TestCheckResourceAttr(rName, "name", "test_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHighRiskCommandImportStateFunc(rName),
			},
		},
	})
}

func testAccHighRiskCommand_base(name string) string {
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

func testAccHighRiskCommand_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_high_risk_command" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  origin_name = "keys"
  name        = "test"
}
`, testAccHighRiskCommand_base(name))
}

func testAccHighRiskCommand_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_high_risk_command" "test" {
  instance_id = huaweicloud_geminidb_instance.test.id
  origin_name = "keys"
  name        = "test_update"
}
`, testAccHighRiskCommand_base(name))
}

func testAccHighRiskCommandImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var originName, instanceId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		originName = rs.Primary.Attributes["origin_name"]

		if originName == "" || instanceId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<origin_name>', but got '%s/%s'",
				instanceId, originName)
		}

		return fmt.Sprintf("%s/%s", instanceId, originName), nil
	}
}
