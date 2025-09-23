package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dbss"
)

func getAddEcsDatabaseFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dbss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DBSS client: %s", err)
	}
	return dbss.GetDatabases(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccAddEcsDatabase_basic(t *testing.T) {
	var (
		addEcsDatabase interface{}
		rName          = "huaweicloud_dbss_ecs_database.test"
		name           = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&addEcsDatabase,
		getAddEcsDatabaseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddEcsDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_dbss_instance.test", "instance_id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "MYSQL"),
					resource.TestCheckResourceAttr(rName, "version", "8"),
					resource.TestCheckResourceAttr(rName, "ip", "192.168.0.88"),
					resource.TestCheckResourceAttr(rName, "port", "3306"),
					resource.TestCheckResourceAttr(rName, "os", "LINUX64"),
					resource.TestCheckResourceAttr(rName, "charset", "UTF8"),
					resource.TestCheckResourceAttr(rName, "instance_name", name),
					resource.TestCheckResourceAttr(rName, "status", "ON"),
					resource.TestCheckResourceAttrSet(rName, "db_classification"),
				),
			},
			{
				Config: testAccAddEcsDatabase_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "OFF"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"lts_audit_switch",
				},
				ImportStateIdFunc: testAccAddEcsDatabaseImportState(rName),
			},
		},
	})
}

func testAccAddEcsDatabase_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dbss_ecs_database" "test" {
  instance_id   = huaweicloud_dbss_instance.test.instance_id
  name          = "%[2]s"
  type          = "MYSQL"
  version       = "8"
  ip            = "192.168.0.88"
  port          = "3306"
  os            = "LINUX64"
  charset       = "UTF8"
  instance_name = "%[2]s"
  status        = "ON"
}
`, testInstance_basic(name), name)
}

func testAccAddEcsDatabase_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dbss_ecs_database" "test" {
  instance_id   = huaweicloud_dbss_instance.test.instance_id
  name          = "%[2]s"
  type          = "MYSQL"
  version       = "8"
  ip            = "192.168.0.88"
  port          = "3306"
  os            = "LINUX64"
  charset       = "UTF8"
  instance_name = "%[2]s"
  status        = "OFF"
}
`, testInstance_basic(name), name)
}

func testAccAddEcsDatabaseImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, databaseId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		databaseId = rs.Primary.ID
		if instanceId == "" || databaseId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, databaseId)
		}
		return fmt.Sprintf("%s/%s", instanceId, databaseId), nil
	}
}
