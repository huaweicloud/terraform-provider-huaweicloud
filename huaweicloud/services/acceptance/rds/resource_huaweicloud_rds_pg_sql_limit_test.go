package rds

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

func getPgSqlLimitResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit?db_name={db_name}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	dbName := state.Primary.Attributes["db_name"]
	sqlLimitID := state.Primary.Attributes["sql_limit_id"]
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceID)
	getBasePath = strings.ReplaceAll(getBasePath, "{db_name}", dbName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var currentTotal int
	var getPath string
	for {
		getPath = getBasePath + buildGetSqlLimitQueryParams(currentTotal)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving RDS PostgreSQL SQL limit: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}

		sqlLimitObjects := utils.PathSearch("sql_limit_objects", getRespBody, make([]interface{}, 0)).([]interface{})
		sqlLimit := utils.PathSearch(fmt.Sprintf("[?id == '%s']|[0]", sqlLimitID), sqlLimitObjects, nil)
		if sqlLimit != nil {
			return sqlLimit, nil
		}
		currentTotal += len(sqlLimitObjects)
		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		if currentTotal >= int(total) {
			break
		}
	}

	return nil, fmt.Errorf("error retrieving RDS PostgreSQL sql limit by instanceID(%s), account(%s) and id(%s)",
		instanceID, dbName, sqlLimitID)
}

func buildGetSqlLimitQueryParams(offset int) string {
	return fmt.Sprintf("&limit=100&offset=%v", offset)
}

func TestAccPgSqlLimit_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_sql_limit.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgSqlLimitResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgSqlLimit_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "db_name", name),
					resource.TestCheckResourceAttr(rName, "query_id", "100"),
					resource.TestCheckResourceAttr(rName, "query_string", "not know"),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "20"),
					resource.TestCheckResourceAttr(rName, "max_waiting", "5"),
					resource.TestCheckResourceAttr(rName, "search_path", "public"),
					resource.TestCheckResourceAttr(rName, "switch", "open"),
					resource.TestCheckResourceAttrSet(rName, "sql_limit_id"),
					resource.TestCheckResourceAttrSet(rName, "is_effective"),
				),
			},
			{
				Config: testPgSqlLimit_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "max_concurrency", "0"),
					resource.TestCheckResourceAttr(rName, "max_waiting", "0"),
					resource.TestCheckResourceAttr(rName, "switch", "close"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"query_string"},
			},
		},
	})
}

func testPgSqlLimit_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "available" {
  db_type           = "PostgreSQL"
  db_version        = "16"
  instance_mode     = "single"
  group_type        = "general"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vcpus             = 2
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.available.flavors[0].name
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id

  parameters {
    name  = "rds_pg_sql_ccl.enable_ccl"
    value = "on"
  }

  db {
    type    = "PostgreSQL"
    version = "16"
  }

  volume {
    type = "CLOUDSSD"
    size = 100
  }
}

resource "huaweicloud_rds_pg_database" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
}

resource "huaweicloud_rds_pg_plugin" "test" {
  depends_on = [huaweicloud_rds_pg_database.test]

  instance_id   = huaweicloud_rds_instance.test.id
  name          = "rds_pg_sql_ccl"
  database_name = "%[2]s"
}
`, testAccRdsInstance_base(), name)
}

func testPgSqlLimit_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_sql_limit" "test" {
  depends_on = [huaweicloud_rds_pg_plugin.test]

  instance_id     = huaweicloud_rds_instance.test.id
  db_name         = "%[2]s"
  query_id        = "100"
  max_concurrency = 20
  max_waiting     = 5
  search_path     = "public"
  switch          = "open"

  lifecycle {
    ignore_changes = [query_string]
  }
}
`, testPgSqlLimit_base(name), name)
}

func testPgSqlLimit_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_sql_limit" "test" {
  depends_on = [huaweicloud_rds_pg_plugin.test]

  instance_id     = huaweicloud_rds_instance.test.id
  db_name         = "%[2]s"
  query_id        = "100"
  max_concurrency = 0
  max_waiting     = 0
  search_path     = "public"
  switch          = "close"

  lifecycle {
    ignore_changes = [query_string]
  }
}
`, testPgSqlLimit_base(name), name)
}
