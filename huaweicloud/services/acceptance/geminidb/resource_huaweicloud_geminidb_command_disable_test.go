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

func getResourceCommandDisableFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("geminidb", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating GeminiDB client: %s", err)
	}

	return geminidb.GetCommandDisableInfo(client, state.Primary.ID, state.Primary.Attributes["disabled_type"])
}

func TestAccResourceCommandDisable_basic(t *testing.T) {
	var (
		rName  = "huaweicloud_geminidb_command_disable.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCommandDisableFunc,
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
				Config: testAccCommandDisable_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "disabled_type", "command"),
					resource.TestCheckResourceAttr(rName, "commands.#", "2"),
				),
			},
			{
				Config: testAccCommandDisable_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "commands.#", "3"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCommandDisableImportStateFunc(rName),
			},
		},
	})
}

func TestAccResourceCommandDisable_keyCommand(t *testing.T) {
	var (
		rName  = "huaweicloud_geminidb_command_disable.test"
		name   = acceptance.RandomAccResourceName()
		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceCommandDisableFunc,
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
				Config: testAccCommandDisable_keyCommand_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_geminidb_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "disabled_type", "key"),
					resource.TestCheckResourceAttr(rName, "keys.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "keys.*", map[string]string{
						"db_id": "0",
						"key":   "name",
					}),
				),
			},
			{
				Config: testAccCommandDisable_keyCommand_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "keys.#", "3"),
					resource.TestCheckTypeSetElemNestedAttrs(rName, "keys.*", map[string]string{
						"db_id": "2",
						"key":   "grade",
					}),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCommandDisableImportStateFunc(rName),
			},
		},
	})
}

func testAccCommandDisable_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_nosql_flavors" "test" {
  vcpus             = 2
  engine            = "redis"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_geminidb_instance" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  password          = "test_1234"
  mode              = "Cluster"

  datastore {
    type           = "redis"
    version        = "5.0"
    storage_engine = "rocksDB"
  }

  flavor {
    num       = "3"
    size      = "16"
    storage   = "ULTRAHIGH"
    spec_code = data.huaweicloud_gaussdb_nosql_flavors.test.flavors[0].name
  }
}
`, common.TestBaseNetwork(name), name)
}

func testAccCommandDisable_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = huaweicloud_geminidb_instance.test.id
  disabled_type = "command"
  commands      = ["keys","hkeys"]
}
`, testAccCommandDisable_base(name))
}

func testAccCommandDisable_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = huaweicloud_geminidb_instance.test.id
  disabled_type = "command"
  commands      = ["hkeys","hvals","hgetall"]
}
`, testAccCommandDisable_base(name))
}

func testAccCommandDisable_keyCommand_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = huaweicloud_geminidb_instance.test.id
  disabled_type = "key"

  keys {
    db_id    = "0"
    key      = "name"
    commands = ["get"]
  }

  keys {
    db_id    = "0"
    key      = "age"
    commands = ["get","set"]
  }

  keys {
    db_id    = "1"
    key      = "address"
    commands = ["lrange"]
  }
}
`, testAccCommandDisable_base(name))
}

func testAccCommandDisable_keyCommand_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_geminidb_command_disable" "test" {
  instance_id   = huaweicloud_geminidb_instance.test.id
  disabled_type = "key"

  keys {
    db_id    = "0"
    key      = "name"
    commands = ["get","set"]
  }

  keys {
    db_id    = "1"
    key      = "age"
    commands = ["set"]
  }

  keys {
    db_id    = "2"
    key      = "grade"
    commands = ["sort"]
  }
}
`, testAccCommandDisable_base(name))
}

func testAccCommandDisableImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, disableType string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		disableType = rs.Primary.Attributes["disabled_type"]

		if instanceId == "" || disableType == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<disabled_type>', but got '%s/%s'",
				instanceId, disableType)
		}

		return fmt.Sprintf("%s/%s", instanceId, disableType), nil
	}
}
