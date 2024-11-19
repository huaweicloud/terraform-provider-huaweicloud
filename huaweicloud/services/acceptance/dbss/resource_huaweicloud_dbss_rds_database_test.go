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

func getAddRdsDatabaseFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("dbss", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DBSS client: %s", err)
	}
	return dbss.GetDatabaseList(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccAddRdsDatabase_basic(t *testing.T) {
	var (
		addRdsDatabase interface{}
		rName          = "huaweicloud_dbss_rds_database.test"
		name           = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&addRdsDatabase,
		getAddRdsDatabaseFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAddRdsDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id", "huaweicloud_dbss_instance.test", "instance_id"),
					resource.TestCheckResourceAttrPair(rName, "rds_id", "huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "type", "MYSQL"),
					resource.TestCheckResourceAttr(rName, "status", "ON"),
					resource.TestCheckResourceAttrSet(rName, "db_id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "charset"),
					resource.TestCheckResourceAttrSet(rName, "ip"),
					resource.TestCheckResourceAttrSet(rName, "port"),
					resource.TestCheckResourceAttrSet(rName, "os"),
					resource.TestCheckResourceAttrSet(rName, "instance_name"),
					resource.TestCheckResourceAttrSet(rName, "db_classification"),
				),
			},
			{
				Config: testAccAddRdsDatabase_update(name),
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
				ImportStateIdFunc: testAccAddRdsDatabaseImportState(rName),
			},
		},
	})
}

// Before test, you need to set the default security group rule, enable `3306` port
func testAccAddRdsDatabase_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  db {
    type    = "MySQL"
    version = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 100
  }
}
`, testInstance_basic(name), name)
}

func testAccAddRdsDatabase_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dbss_rds_database" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
  rds_id      = huaweicloud_rds_instance.test.id
  type        = "MYSQL"
  status      = "ON"
}
`, testAccAddRdsDatabase_base(name))
}

func testAccAddRdsDatabase_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dbss_rds_database" "test" {
  instance_id = huaweicloud_dbss_instance.test.instance_id
  rds_id      = huaweicloud_rds_instance.test.id
  type        = "MYSQL"
  status      = "OFF"
}
`, testAccAddRdsDatabase_base(name))
}

func testAccAddRdsDatabaseImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var instanceId, rdsId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", rName)
		}

		instanceId = rs.Primary.Attributes["instance_id"]
		rdsId = rs.Primary.ID
		if instanceId == "" || rdsId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, rdsId)
		}
		return fmt.Sprintf("%s/%s", instanceId, rdsId), nil
	}
}
