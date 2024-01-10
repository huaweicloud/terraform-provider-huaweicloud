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

func getMysqlDatabaseResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getMysqlDatabase: query RDS Mysql database
	var (
		getMysqlDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		getMysqlDatabaseProduct = "rds"
	)
	getMysqlDatabaseClient, err := cfg.NewServiceClient(getMysqlDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getMysqlDatabasePath := getMysqlDatabaseClient.Endpoint + getMysqlDatabaseHttpUrl
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{project_id}", getMysqlDatabaseClient.ProjectID)
	getMysqlDatabasePath = strings.ReplaceAll(getMysqlDatabasePath, "{instance_id}", instanceId)

	getMysqlDatabaseResp, err := pagination.ListAllItems(
		getMysqlDatabaseClient,
		"page",
		getMysqlDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, fmt.Errorf("error retrieving MysqlDatabase")
	}

	getMysqlDatabaseRespJson, err := json.Marshal(getMysqlDatabaseResp)
	if err != nil {
		return nil, err
	}
	var getMysqlDatabaseRespBody interface{}
	err = json.Unmarshal(getMysqlDatabaseRespJson, &getMysqlDatabaseRespBody)
	if err != nil {
		return nil, err
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getMysqlDatabaseRespBody, nil)
	if database != nil {
		return database, nil
	}

	return nil, fmt.Errorf("error get RDS Mysql database by instanceID %s and database %s", instanceId, dbName)
}

func TestAccMysqlDatabase_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_mysql_database.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getMysqlDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testMysqlDatabase_basic(name, "test database"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "character_set", "utf8"),
					resource.TestCheckResourceAttr(rName, "description", "test database"),
				),
			},
			{
				Config: testMysqlDatabase_basic(name, "test database update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "character_set", "utf8"),
					resource.TestCheckResourceAttr(rName, "description", "test database update"),
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

func testMysqlDatabase_basic(name, description string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_mysql_database" "test" {
  instance_id   = huaweicloud_rds_instance.test.id
  name          = "%s"
  character_set = "utf8"
  description   = "%s"
}
`, testAccRdsInstance_mysql_step1(name), name, description)
}
