package gaussdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getOpenGaussSchemaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/schemas"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := state.Primary.Attributes["instance_id"]
	getBasePath := client.Endpoint + httpUrl
	getBasePath = strings.ReplaceAll(getBasePath, "{project_id}", client.ProjectID)
	getBasePath = strings.ReplaceAll(getBasePath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var offset int
	var dbSchema interface{}
	dbName := state.Primary.Attributes["db_name"]
	name := state.Primary.Attributes["name"]

	for {
		getPath := getBasePath + buildOpenGaussSchemaQueryParams(dbName, offset)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, err
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, err
		}
		dbSchema = utils.PathSearch(fmt.Sprintf("database_schemas[?schema_name=='%s']|[0]", name), getRespBody, nil)
		if dbSchema != nil {
			break
		}
		totalCount := utils.PathSearch("total_count", getRespBody, float64(0)).(float64)
		if int(totalCount) <= (offset+1)*100 {
			break
		}
		offset++
	}
	if dbSchema == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return dbSchema, nil
}

func buildOpenGaussSchemaQueryParams(dbName string, offset int) string {
	return fmt.Sprintf("?db_name=%s&limit=100&offset=%v", dbName, offset)
}

func TestAccOpenGaussSchema_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_opengauss_schema.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getOpenGaussSchemaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testOpenGaussSchema_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_opengauss_instance.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "db_name",
						"huaweicloud_gaussdb_opengauss_database.test", "name"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "owner", "root"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccOpenGaussSchemaImportStateFunc(rName),
			},
		},
	})
}

func testOpenGaussSchema_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_opengauss_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "ingress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "in_v4_tcp_opengauss_egress" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  ethertype         = "IPv4"
  direction         = "egress"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_gaussdb_opengauss_instance" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss,
    huaweicloud_networking_secgroup_rule.in_v4_tcp_opengauss_egress
  ]

  vpc_id                = huaweicloud_vpc.test.id
  subnet_id             = huaweicloud_vpc_subnet.test.id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  flavor                = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
  name                  = "%[2]s"
  password              = "Huangwei!120521"
  enterprise_project_id = "%[3]s"

  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "eventual"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_gaussdb_opengauss_database" "test" {
  instance_id   = huaweicloud_gaussdb_opengauss_instance.test.id
  name          = "%[2]s"
  character_set = "UTF8"
  owner         = "root"
  template      = "template0"
  lc_collate    = "C"
  lc_ctype      = "C"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testOpenGaussSchema_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_gaussdb_opengauss_schema" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_database.test]

  instance_id = huaweicloud_gaussdb_opengauss_instance.test.id
  db_name     = huaweicloud_gaussdb_opengauss_database.test.name
  name        = "%[2]s"
  owner       = "root"
}
`, testOpenGaussSchema_base(rName), rName)
}

func testAccOpenGaussSchemaImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("attribute (instance_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["db_name"] == "" {
			return "", fmt.Errorf("attribute (db_name) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["name"] == "" {
			return "", fmt.Errorf("attribute (name) of Resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.Attributes["db_name"],
			rs.Primary.Attributes["name"]), nil
	}
}
