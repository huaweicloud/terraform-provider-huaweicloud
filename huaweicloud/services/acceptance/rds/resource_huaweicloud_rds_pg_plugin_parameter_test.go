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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	defaultValues = map[string][]string{
		"shared_preload_libraries": {"passwordcheck.so", "pg_stat_statements", "pg_sql_history", "pgaudit"},
	}
)

func getPgPluginParameterResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/parameter/{name}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	instanceID := state.Primary.Attributes["instance_id"]
	name := state.Primary.Attributes["name"]
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{name}", name)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RDS PostgreSQL plugin parameter values: %s", err)
	}
	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	respValueString := utils.PathSearch("value", respBody, nil)
	if respValueString == nil {
		return nil, fmt.Errorf("error found RDS PostgreSQL plugin parameter values")
	}
	respValues := strings.Split(respValueString.(string), ",")

	defaults := defaultValues[name]
	defaultsMap := make(map[string]bool)
	for _, value := range defaults {
		defaultsMap[value] = true
	}
	values := make([]string, 0)
	for _, value := range respValues {
		if !defaultsMap[value] {
			values = append(values, value)
		}
	}
	if len(values) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return respBody, nil
}

func TestAccPgPluginParameter_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_rds_pg_plugin_parameter.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getPgPluginParameterResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPgPluginParameter_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_rds_instance.test", "id"),
					resource.TestCheckResourceAttr(rName, "name", "shared_preload_libraries"),
					resource.TestCheckResourceAttrPair(rName, "values.0",
						"data.huaweicloud_rds_pg_plugin_parameter_value_range.test", "values.0"),
					resource.TestCheckResourceAttrSet(rName, "restart_required"),
					resource.TestCheckResourceAttrSet(rName, "default_values.#"),
				),
			},
			{
				Config: testPgPluginParameter_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "values.0",
						"data.huaweicloud_rds_pg_plugin_parameter_value_range.test", "values.1"),
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

func testPgPluginParameter_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = "rds.pg.n1.large.2"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  db {
    type    = "PostgreSQL"
    version = "16"
    port    = 8634
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }
}

data "huaweicloud_rds_pg_plugin_parameter_value_range" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "shared_preload_libraries"
}
`, common.TestBaseNetwork(name), name)
}

func testPgPluginParameter_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_plugin_parameter" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "shared_preload_libraries"
  values      = [data.huaweicloud_rds_pg_plugin_parameter_value_range.test.values[0]]
}
`, testPgPluginParameter_base(name))
}

func testPgPluginParameter_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_pg_plugin_parameter" "test" {
  instance_id = huaweicloud_rds_instance.test.id
  name        = "shared_preload_libraries"
  values      = [data.huaweicloud_rds_pg_plugin_parameter_value_range.test.values[1]]
}
`, testPgPluginParameter_base(name))
}
