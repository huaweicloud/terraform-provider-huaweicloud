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

func getSQLServerDatabaseResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSQLServerDatabase: query RDS SQLServer database
	var (
		getSQLServerDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		getSQLServerDatabaseProduct = "rds"
	)
	getSQLServerDatabaseClient, err := cfg.NewServiceClient(getSQLServerDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database name from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getSQLServerDatabasePath := getSQLServerDatabaseClient.Endpoint + getSQLServerDatabaseHttpUrl
	getSQLServerDatabasePath = strings.ReplaceAll(getSQLServerDatabasePath, "{project_id}",
		getSQLServerDatabaseClient.ProjectID)
	getSQLServerDatabasePath = strings.ReplaceAll(getSQLServerDatabasePath, "{instance_id}", instanceId)

	getSQLServerDatabaseResp, err := pagination.ListAllItems(
		getSQLServerDatabaseClient,
		"page",
		getSQLServerDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer database: %s", err)
	}

	getSQLServerDatabaseRespJson, err := json.Marshal(getSQLServerDatabaseResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer database: %s", err)
	}
	var getSQLServerDatabaseRespBody interface{}
	err = json.Unmarshal(getSQLServerDatabaseRespJson, &getSQLServerDatabaseRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS SQLServer database: %s", err)
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getSQLServerDatabaseRespBody, nil)
	if database != nil {
		return database, nil
	}

	return nil, fmt.Errorf("error get RDS SQLServer database by instanceID %s and name %s", instanceId, dbName)
}

func TestAccSQLServerDatabase_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_sqlserver_database.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSQLServerDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testSQLServerDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "character_set"),
					resource.TestCheckResourceAttrSet(rName, "state"),
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

func testSQLServerDatabase_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.mssql.spec.se.s3.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  db {
    password = "Terraform145@"
    type     = "SQLServer"
    version  = "2022_SE"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_rds_sqlserver_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
}
`, testAccRdsInstance_base(), name)
}
