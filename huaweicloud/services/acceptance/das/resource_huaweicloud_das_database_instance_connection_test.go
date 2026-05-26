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

func getDatabaseInstanceConnectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("das", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DAS client: %s", err)
	}
	return das.GetDatabaseInstanceConnectionById(client, state.Primary.ID)
}

func TestAccDatabaseInstanceConnection_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_das_database_instance_connection.test"
		rcCreateUser = acceptance.InitResourceCheck(resourceName, &obj, getDatabaseInstanceConnectionResourceFunc)

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
				Config: testAccDatabaseInstanceConnection_basic_step1(name, password),
				Check: resource.ComposeTestCheckFunc(
					rcCreateUser.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "network_type", "rds"),
					resource.TestCheckResourceAttr(resourceName, "username", name),
					resource.TestCheckResourceAttr(resourceName, "is_save_password", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttr(resourceName, "database_name", name),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "datastore_version"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "conn_share_type"),
				),
			},
			{
				Config: testAccDatabaseInstanceConnection_basic_step2(name, updatePassword),
				Check: resource.ComposeTestCheckFunc(
					rcCreateUser.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttr(resourceName, "engine_type", "mysql"),
					resource.TestCheckResourceAttr(resourceName, "network_type", "rds"),
					resource.TestCheckResourceAttr(resourceName, "username", name+"_test"),
					resource.TestCheckResourceAttr(resourceName, "is_save_password", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script!"),
					resource.TestCheckResourceAttrSet(resourceName, "port"),
					resource.TestCheckResourceAttr(resourceName, "database_name", name+"_test"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "datastore_version"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "conn_share_type"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_ids", "password", "sql_record_flag"},
			},
		},
	})
}

func testAccDatabaseInstanceConnection_basic_step1(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_das_database_instance_connection" "test" {
  depends_on = [
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id      = "%[1]s"
  engine_type      = "mysql"
  network_type     = "rds"
  username         = "%[2]s"
  password         = "%[3]s"
  is_save_password = true
  node_ids         = ["%[1]s"]
  description      = "Created by terraform script"
  database_name    = "%[2]s"
  sql_record_flag  = true
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}

func testAccDatabaseInstanceConnection_basic_step2(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_mysql_account" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  password    = "%[3]s"

  hosts = [
    "%%"
  ]
}

resource "huaweicloud_das_database_instance_connection" "test" {
  depends_on = [
    huaweicloud_rds_mysql_account.test,
  ]

  instance_id      = "%[1]s"
  engine_type      = "mysql"
  network_type     = "rds"
  username         = "%[2]s_test"
  password         = "%[3]s"
  is_save_password = false
  node_ids         = ["%[1]s"]
  description      = "Created by terraform script!"
  database_name    = "%[2]s_test"
  sql_record_flag  = false
}
`, acceptance.HW_RDS_INSTANCE_ID, name, password)
}
