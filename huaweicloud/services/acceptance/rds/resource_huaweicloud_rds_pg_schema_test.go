package rds

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getPgSchemaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schema/detail?db_name={db_name}&page=1&limit=100"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])
	getPath = strings.ReplaceAll(getPath, "{db_name}", state.Primary.Attributes["db_name"])

	getResp, err := pagination.ListAllItems(
		client,
		"page",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL schema: %s", err)
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL schema: %s", err)
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL schema: %s", err)
	}

	schemaName := state.Primary.Attributes["schema_name"]
	schemaBody := utils.PathSearch(fmt.Sprintf("database_schemas[?schema_name=='%s']|[0]", schemaName), getRespBody, nil)
	if schemaBody != nil {
		return schemaBody, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func TestAccPgSchema_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_schema.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgSchemaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPgSchema_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "db_name", name),
					resource.TestCheckResourceAttr(rName, "schema_name", name),
					resource.TestCheckResourceAttr(rName, "owner", "root"),
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

func testPgSchema_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  description       = "test_description"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = data.huaweicloud_networking_secgroup.test.id
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
`, testAccRdsInstance_base(), name)
}

func testPgSchema_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_schema" "test" {
  depends_on = [huaweicloud_rds_pg_database.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
  schema_name = "%[2]s"
  owner       = "root"
}
`, testPgSchema_base(name), name)
}
