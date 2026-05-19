package das

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/das"
)

func getDatabaseUserResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}
	return das.GetDatabaseUserById(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccDatabaseUser_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_das_database_user.test"
		rcCreateUser = acceptance.InitResourceCheck(resourceName, &obj, getDatabaseUserResourceFunc)

		name           = acceptance.RandomAccResourceName()
		password       = acceptance.RandomPassword()
		updatePassword = acceptance.RandomPassword()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcCreateUser.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseUser_basic_step1(name, password),
				Check: resource.ComposeTestCheckFunc(
					rcCreateUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				Config: testAccDatabaseUser_basic_step2(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rcCreateUser.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateIdFunc:       testAccDatabaseUserImportStateFunc(resourceName),
			},
		},
	})
}

func testAccDatabaseUser_basic_step1(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = "%[1]s"
  name          = "%[2]s"
  character_set = "utf8"
}

resource "huaweicloud_rds_mysql_account" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
  ]
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id = "%[1]s"
  db_name     = huaweicloud_rds_mysql_database.test.name

  users {
    name     = huaweicloud_rds_mysql_account.test.name
    readonly = false
  }
}

resource "huaweicloud_das_database_user" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database_privilege.test,
  ]

  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}

func testAccDatabaseUser_basic_step2(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = "%[1]s"
  name          = "%[2]s"
  character_set = "utf8"
}

resource "huaweicloud_rds_mysql_account" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
  ]

  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  depends_on = [
    huaweicloud_rds_mysql_database.test,
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id = "%[1]s"
  db_name     = huaweicloud_rds_mysql_database.test.name

  users {
    name     = huaweicloud_rds_mysql_account.test.name
    readonly = false
  }
}

resource "huaweicloud_das_database_user" "test" {
  depends_on = [
	huaweicloud_rds_mysql_database_privilege.test,
  ]

  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}

func testAccDatabaseUserImportStateFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		dbUserId := rs.Primary.ID
		if instanceId == "" || dbUserId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<instance_id>/<id>', but got '%s/%s'",
				instanceId, dbUserId)
		}
		return fmt.Sprintf("%s/%s", instanceId, dbUserId), nil
	}
}
