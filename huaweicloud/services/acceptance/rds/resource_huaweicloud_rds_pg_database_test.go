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

func getPgDatabaseResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getPgDatabase: query RDS PostgreSQL database
	var (
		getPgDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?db={db_name}&page=1&limit=100"
		getPgDatabaseProduct = "rds"
	)
	getPgDatabaseClient, err := cfg.NewServiceClient(getPgDatabaseProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getPgDatabasePath := getPgDatabaseClient.Endpoint + getPgDatabaseHttpUrl
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{project_id}", getPgDatabaseClient.ProjectID)
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{instance_id}", instanceId)
	getPgDatabasePath = strings.ReplaceAll(getPgDatabasePath, "{db_name}", dbName)

	getPgDatabaseResp, err := pagination.ListAllItems(
		getPgDatabaseClient,
		"page",
		getPgDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL database: %s", err)
	}

	getPgDatabaseRespJson, err := json.Marshal(getPgDatabaseResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL database: %s", err)
	}
	var getPgDatabaseRespBody interface{}
	err = json.Unmarshal(getPgDatabaseRespJson, &getPgDatabaseRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL database: %s", err)
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getPgDatabaseRespBody, nil)
	if database != nil {
		return database, nil
	}

	return nil, fmt.Errorf("error retrieving RDS PostgreSQL database by instanceID %s and database %s", instanceId,
		dbName)
}

func TestAccPgDatabase_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_database.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgDatabaseResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgDatabase_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "owner", "root"),
					resource.TestCheckResourceAttr(rName, "character_set", "UTF8"),
					resource.TestCheckResourceAttr(rName, "lc_collate", "en_US.UTF-8"),
					resource.TestCheckResourceAttr(rName, "description", "test_description"),
					resource.TestCheckResourceAttrSet(rName, "size"),
				),
			},
			{
				Config: testPgDatabase_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "owner",
						"huaweicloud_rds_pg_account.test", "name"),
					resource.TestCheckResourceAttr(rName, "character_set", "UTF8"),
					resource.TestCheckResourceAttr(rName, "lc_collate", "en_US.UTF-8"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"template", "lc_ctype", "is_revoke_public_privilege"},
			},
		},
	})
}

func testPgDatabase_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = data.huaweicloud_vpc_subnet.test.id
  vpc_id            = data.huaweicloud_vpc.test.id
  time_zone         = "UTC+08:00"

  db {
    type    = "PostgreSQL"
    version = "12"
  }

  volume {
    type = "CLOUDSSD"
    size = 50
  }
}
`, testAccRdsInstance_base(), name)
}

func testPgDatabase_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_database" "test" {
  instance_id                = huaweicloud_rds_instance.test.id
  name                       = "%[2]s"
  owner                      = "root"
  character_set              = "UTF8"
  template                   = "template1"
  lc_collate                 = "en_US.UTF-8"
  lc_ctype                   = "en_US.UTF-8"
  is_revoke_public_privilege = false
  description                = "test_description"
}
`, testPgDatabase_base(name), name)
}

func testPgDatabase_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_account" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "%[2]s"
  password    = "Test@12345678"
}

resource "huaweicloud_rds_pg_database" "test" {
  depends_on = [huaweicloud_rds_pg_account.test]

  instance_id                = huaweicloud_rds_instance.test.id
  name                       = "%[2]s"
  owner                      = huaweicloud_rds_pg_account.test.name
  character_set              = "UTF8"
  template                   = "template1"
  lc_collate                 = "en_US.UTF-8"
  lc_ctype                   = "en_US.UTF-8"
  is_revoke_public_privilege = false
  description                = ""
}
`, testPgDatabase_base(name), name)
}
