package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/geminidb"
)

func getResourceMemoryMappingFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	queryParams := map[string]string{
		"id": state.Primary.ID,
	}

	return geminidb.GetMemoryMappingInfo(client, queryParams)
}

func TestAccResourceMemoryMapping_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_geminidb_memory_mapping.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceMemoryMappingFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMemoryMapping_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "source_instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "target_instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "source_instance_name"),
					resource.TestCheckResourceAttrSet(rName, "target_instance_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created"),
					resource.TestCheckResourceAttrSet(rName, "updated"),
					resource.TestCheckResourceAttrSet(rName, "rule_count"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemoryMapping_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data"huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  group_type    = "dedicated"
  instance_mode = "ha"
}

resource "huaweicloud_rds_instance" "test" {
  name                = "%[2]s"
  ha_replication_mode = "async"
  flavor              = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id   = huaweicloud_networking_secgroup.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  vpc_id              = huaweicloud_vpc.test.id
  availability_zone   = [
    data.huaweicloud_availability_zones.test.names[0],
    data.huaweicloud_availability_zones.test.names[1],
  ]

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus  = 1
  engine = "redis"
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].availability_zones[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "Test_357159"
  mode              = "Replication"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "2"
    size      = "8"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccMemoryMapping_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_memory_mapping" "test" {
  source_instance_id = huaweicloud_rds_instance.test.id
  target_instance_id = huaweicloud_geminidb_instance.test.id
}
`, testAccMemoryMapping_base(name))
}
