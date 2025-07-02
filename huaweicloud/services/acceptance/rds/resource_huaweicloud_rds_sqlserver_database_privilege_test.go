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

func getSQLServerDatabasePrivilegeResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSQLServerDatabasePrivilege: query RDS SQLServer database privilege
	var (
		getSQLServerDatabasePrivilegeHttpUrl = "v3/{project_id}/instances/{instance_id}/database/db_user"
		getSQLServerDatabasePrivilegeProduct = "rds"
	)
	getSQLServerDatabasePrivilegeClient, err := cfg.NewServiceClient(getSQLServerDatabasePrivilegeProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<db_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getSQLServerDatabasePrivilegePath := getSQLServerDatabasePrivilegeClient.Endpoint + getSQLServerDatabasePrivilegeHttpUrl
	getSQLServerDatabasePrivilegePath = strings.ReplaceAll(getSQLServerDatabasePrivilegePath, "{project_id}",
		getSQLServerDatabasePrivilegeClient.ProjectID)
	getSQLServerDatabasePrivilegePath = strings.ReplaceAll(getSQLServerDatabasePrivilegePath, "{instance_id}", instanceId)

	getSQLServerDatabasePrivilegeQueryParams := buildGetSQLServerDatabasePrivilegeQueryParams(dbName)
	getSQLServerDatabasePrivilegePath += getSQLServerDatabasePrivilegeQueryParams

	getSQLServerDatabasePrivilegeResp, err := pagination.ListAllItems(
		getSQLServerDatabasePrivilegeClient,
		"page",
		getSQLServerDatabasePrivilegePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL Server database privilege: %s", err)
	}

	getSQLServerDatabasePrivilegeRespJson, err := json.Marshal(getSQLServerDatabasePrivilegeResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL Server database privilege: %s", err)
	}
	var getSQLServerDatabasePrivilegeRespBody interface{}
	err = json.Unmarshal(getSQLServerDatabasePrivilegeRespJson, &getSQLServerDatabasePrivilegeRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQL Server database privilege: %s", err)
	}

	curJson := utils.PathSearch("users[?name != 'rdsuser']", getSQLServerDatabasePrivilegeRespBody,
		make([]interface{}, 0))
	if len(curJson.([]interface{})) == 0 {
		return nil, fmt.Errorf("error get RDS SQL Server database privilege")
	}

	return getSQLServerDatabasePrivilegeRespBody, nil
}

func buildGetSQLServerDatabasePrivilegeQueryParams(dbName string) string {
	return fmt.Sprintf("?db-name=%s&page=1&limit=100", dbName)
}

func TestAccSQLServerDatabasePrivilege_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_sqlserver_database_privilege.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLServerDatabasePrivilegeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLServerDatabasePrivilege_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"huaweicloud_rds_sqlserver_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "users.0.name",
						"huaweicloud_rds_sqlserver_account.account_1", "name"),
				),
			},
			{
				Config: testSQLServerDatabasePrivilege_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"huaweicloud_rds_sqlserver_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "users.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "users.0.name",
						"huaweicloud_rds_sqlserver_account.account_2", "name"),
					resource.TestCheckResourceAttr(rName, "users.0.readonly", "true"),
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

func testAccRdsSQLServerDatabasePrivilege_base(name string) string {
	return fmt.Sprintf(`
%[1]s
resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mssql.spec.se.s6.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  charging_mode     = "postPaid"
  db {
    password = "Huangwei!120521"
    type     = "SQLServer"
    version  = "2022_SE"
  }
  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}
resource "huaweicloud_rds_sqlserver_database" "test" {
  depends_on  = [huaweicloud_rds_instance.test]
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[3]s"
}
`, testAccRdsInstance_base(), name, name)
}

func testSQLServerDatabasePrivilege_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_sqlserver_account" "account_1" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_1"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_sqlserver_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_sqlserver_database.test.name

  users {
    name = huaweicloud_rds_sqlserver_account.account_1.name
  }
}
`, testAccRdsSQLServerDatabasePrivilege_base(name), name)
}

func testSQLServerDatabasePrivilege_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_sqlserver_account" "account_1" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_1"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_sqlserver_account" "account_2" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s_2"
  password    = "Terraform145@"
}

resource "huaweicloud_rds_sqlserver_database_privilege" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  db_name     = huaweicloud_rds_sqlserver_database.test.name

  users {
    name     = huaweicloud_rds_sqlserver_account.account_2.name
    readonly = true
  }
}
`, testSQLServerDatabase_basic(name), name)
}
