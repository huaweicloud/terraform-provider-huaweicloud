package taurusdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDBDatabaseResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGaussDBDatabase: Query the GaussDB MySQL database
	var (
		getGaussDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		getGaussDBDatabaseProduct = "gaussdb"
	)
	getGaussDBDatabaseClient, err := cfg.NewServiceClient(getGaussDBDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	databaseName := parts[1]

	getGaussDBDatabaseBasePath := getGaussDBDatabaseClient.Endpoint + getGaussDBDatabaseHttpUrl
	getGaussDBDatabaseBasePath = strings.ReplaceAll(getGaussDBDatabaseBasePath, "{project_id}",
		getGaussDBDatabaseClient.ProjectID)
	getGaussDBDatabaseBasePath = strings.ReplaceAll(getGaussDBDatabaseBasePath, "{instance_id}", instanceID)

	getGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	var currentTotal int
	var database interface{}
	getGaussDBDatabasePath := getGaussDBDatabaseBasePath + buildGaussDBMysqlQueryParams(currentTotal)

	for {
		getGaussDBDatabaseResp, err := getGaussDBDatabaseClient.Request("GET", getGaussDBDatabasePath,
			&getGaussDBDatabaseOpt)

		if err != nil {
			return nil, fmt.Errorf("error retrieving GaussDB MySQL database: %s", err)
		}

		getGaussDBDatabaseRespBody, err := utils.FlattenResponse(getGaussDBDatabaseResp)
		if err != nil {
			return nil, err
		}
		db, pageNum := flattenGetGaussDBDatabaseResponseBodyGetDatabase(getGaussDBDatabaseRespBody, databaseName)
		if db != nil {
			database = db
			break
		}
		total := utils.PathSearch("total_count", getGaussDBDatabaseRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getGaussDBDatabasePath = getGaussDBDatabaseBasePath + buildGaussDBMysqlQueryParams(currentTotal)
	}
	if database == nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL database")
	}
	return database, nil
}

func flattenGetGaussDBDatabaseResponseBodyGetDatabase(resp interface{}, databaseName string) (interface{}, int) {
	if resp == nil {
		return nil, 0
	}
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		if databaseName == name {
			return v, len(curArray)
		}
	}
	return nil, len(curArray)
}

func buildGaussDBMysqlQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func TestAccGaussDBDatabase_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_mysql_database.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDBDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "character_set", "gbk"),
					resource.TestCheckResourceAttr(rName, "description", "gaussdb mysql database"),
				),
			},
			{
				Config: testGaussDBDatabase_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_mysql_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "character_set", "gbk"),
					resource.TestCheckResourceAttr(rName, "description", "gaussdb mysql database update"),
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

func testGaussDBDatabase_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_mysql_database" "test" {
  instance_id   = huaweicloud_gaussdb_mysql_instance.test.id
  name          = "%s"
  character_set = "gbk"
  description   = "gaussdb mysql database"
}
`, testAccGaussDBInstanceConfig_basic(name), name)
}

func testGaussDBDatabase_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_gaussdb_mysql_database" "test" {
  instance_id   = huaweicloud_gaussdb_mysql_instance.test.id
  name          = "%s"
  character_set = "gbk"
  description   = "gaussdb mysql database update"
}
`, testAccGaussDBInstanceConfig_basic(name), name)
}
