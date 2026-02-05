package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getMysqlDatabasePrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getMysqlDatabasePrivilege: query RDS Mysql database privilege
	var (
		getMysqlDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getMysqlDatabasePrivilegeProduct = "rds"
	)
	getMysqlDatabasePrivilegeClient, err := cfg.NewServiceClient(getMysqlDatabasePrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePrivilegePath := getMysqlDatabasePrivilegeClient.Endpoint + getMysqlDatabasePrivilegeHttpUrl
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{project_id}",
		getMysqlDatabasePrivilegeClient.ProjectID)
	getMysqlDatabasePrivilegePath = strings.ReplaceAll(getMysqlDatabasePrivilegePath, "{instance_id}", instanceId)

	getMysqlDatabasePrivilegeQueryParams := buildGetMysqlDatabasePrivilegeQueryParams(dbName)
	getMysqlDatabasePrivilegePath += getMysqlDatabasePrivilegeQueryParams

	getMysqlDatabasePrivilegeResp, err := pagination.ListAllItems(
		getMysqlDatabasePrivilegeClient,
		"page",
		getMysqlDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	getMysqlDatabasePrivilegeRespJson, err := json.Marshal(getMysqlDatabasePrivilegeResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}
	var getMysqlDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getMysqlDatabasePrivilegeRespJson, &getMysqlDatabasePrivilegeRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Mysql database privilege: %s", err)
	}

	curJson := utils.PathSearch("users", getMysqlDatabasePrivilegeRespBody, make([]interface{}, 0))
	if len(curJson.([]interface{})) == 0 {
		return nil, fmt.Errorf("error get RDS Mysql database privilege")
	}

	return getMysqlDatabasePrivilegeRespBody, nil
}

func buildGetMysqlDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func TestAccMysqlDatabasePrivilege_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_rds_mysql_database_privilege.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getMysqlDatabasePrivilegeResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"random": {
				Source: "hashicorp/random",
				// The version of the random provider must be greater than 3.3.0 to support the 'numeric' parameter.
				VersionConstraint: "3.3.0",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlDatabasePrivilege_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"huaweicloud_rds_mysql_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
				),
			},
			{
				Config: testAccMysqlDatabasePrivilege_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"huaweicloud_rds_mysql_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "2"),
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

func testAccMysqlDatabasePrivilege_basic_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_vpc" "test" {
  name                  = "%[1]s"
  cidr                  = "192.168.0.0/16"
  enterprise_project_id = "%[2]s"
}

resource "huaweicloud_vpc_subnet" "test" {
  vpc_id     = huaweicloud_vpc.test.id
  name       = "%[1]s"
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_rds_instance" "test" {
  name                  = "%[1]s"
  flavor                = data.huaweicloud_rds_flavors.test.flavors[0].name
  availability_zone     = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)
  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  enterprise_project_id = "%[2]s"

  db {
    type    = "MySQL"
    version = "8.0"
    port    = 3306
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%[1]s"
  character_set = "utf8"
}

resource "random_password" "test" {
  length           = 12
  min_numeric      = 1
  min_upper        = 1
  min_lower        = 1
  special          = true
  min_special      = 1
  override_special = "!#"
}

resource "huaweicloud_rds_mysql_account" "test" {
  count = 3

  instance_id = huaweicloud_rds_instance.test.id
  name        = format("%[1]s_%%d", count.index)
  password    = random_password.test.result
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccMysqlDatabasePrivilege_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  # The behavior of parameter 'name' of the database resource is 'Required', means this parameter does not have
  # 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_mysql_database.test.name
  
  dynamic "users" {
    for_each = slice(huaweicloud_rds_mysql_account.test[*].name, 0, 2)

    content {
      name     = users.value
      readonly = true
    }
  }
}
`, testAccMysqlDatabasePrivilege_basic_base(name))
}

func testAccMysqlDatabasePrivilege_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_mysql_database_privilege" "test" {
  # The behavior of parameter 'name' of the database resource is 'Required', means this parameter does not 
  # have 'Know After Apply' behavior.
  depends_on = [huaweicloud_rds_mysql_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_mysql_database.test.name
  
  dynamic "users" {
    for_each = slice(huaweicloud_rds_mysql_account.test[*].name, 1, 3)

    content {
      name     = users.value
      readonly = true
    }
  }
}
`, testAccMysqlDatabasePrivilege_basic_base(name))
}
