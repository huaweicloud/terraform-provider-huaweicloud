package taurusdb

import (
	"errors"
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

func getTaurusDBDatabaseResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getTaurusDBDatabase: Query the TaurusDB database
	var (
		getTaurusDBDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/databases"
		getTaurusDBDatabaseProduct = "gaussdb"
	)
	getTaurusDBDatabaseClient, err := cfg.NewServiceClient(getTaurusDBDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating TaurusDB Client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid id format, must be <instance_id>/<name>")
	}
	instanceID := parts[0]
	databaseName := parts[1]

	getTaurusDBDatabaseBasePath := getTaurusDBDatabaseClient.Endpoint + getTaurusDBDatabaseHttpUrl
	getTaurusDBDatabaseBasePath = strings.ReplaceAll(getTaurusDBDatabaseBasePath, "{project_id}",
		getTaurusDBDatabaseClient.ProjectID)
	getTaurusDBDatabaseBasePath = strings.ReplaceAll(getTaurusDBDatabaseBasePath, "{instance_id}", instanceID)

	getTaurusDBDatabaseOpt := golangsdk.RequestOpts{
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
	getTaurusDBDatabasePath := getTaurusDBDatabaseBasePath + buildTaurusDBQueryParams(currentTotal)

	for {
		getTaurusDBDatabaseResp, err := getTaurusDBDatabaseClient.Request("GET", getTaurusDBDatabasePath,
			&getTaurusDBDatabaseOpt)

		if err != nil {
			return nil, fmt.Errorf("error retrieving TaurusDB database: %s", err)
		}

		getTaurusDBDatabaseRespBody, err := utils.FlattenResponse(getTaurusDBDatabaseResp)
		if err != nil {
			return nil, err
		}
		db, pageNum := flattenGetTaurusDBDatabaseResponseBodyGetDatabase(getTaurusDBDatabaseRespBody, databaseName)
		if db != nil {
			database = db
			break
		}
		total := utils.PathSearch("total_count", getTaurusDBDatabaseRespBody, float64(0)).(float64)
		currentTotal += pageNum
		if currentTotal == int(total) {
			break
		}
		getTaurusDBDatabasePath = getTaurusDBDatabaseBasePath + buildTaurusDBQueryParams(currentTotal)
	}
	if database == nil {
		return nil, errors.New("error retrieving TaurusDB database")
	}
	return database, nil
}

func flattenGetTaurusDBDatabaseResponseBodyGetDatabase(resp interface{}, databaseName string) (interface{}, int) {
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

func buildTaurusDBQueryParams(offset int) string {
	return fmt.Sprintf("?limit=100&offset=%v", offset)
}

func TestAccTaurusDBDatabase_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_taurusdb_database.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getTaurusDBDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testTaurusDBDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "character_set", "gbk"),
					resource.TestCheckResourceAttr(rName, "description", "gaussdb mysql database"),
				),
			},
			{
				Config: testTaurusDBDatabase_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_taurusdb_instance.test", "id"),
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

func testTaurusDBDatabase_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_database" "test" {
  instance_id   = huaweicloud_taurusdb_instance.test.id
  name          = "%s"
  character_set = "gbk"
  description   = "gaussdb mysql database"
}
`, testAccTaurusDBInstanceConfig_basic(name), name)
}

func testTaurusDBDatabase_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_taurusdb_database" "test" {
  instance_id   = huaweicloud_taurusdb_instance.test.id
  name          = "%s"
  character_set = "gbk"
  description   = "gaussdb mysql database update"
}
`, testAccTaurusDBInstanceConfig_basic(name), name)
}
