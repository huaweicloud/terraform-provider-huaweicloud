package gaussdb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGaussDbReadReplicaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/readonly-nodes"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	readReplica := utils.PathSearch(fmt.Sprintf("nodes[?id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if readReplica == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return readReplica, nil
}

func TestAccGaussDbReadReplica_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_gaussdb_read_replica.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGaussDbReadReplicaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDBReadReplica_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_GAUSSDB_INSTANCE_ID),
					resource.TestCheckResourceAttrPair(rName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(rName, "flavor_ref", "gaussdb.bs.s6.2xlarge.x864.ha"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "private_ip"),
					resource.TestCheckResourceAttrSet(rName, "data_ip"),
					resource.TestCheckResourceAttrSet(rName, "component_names"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testGaussDbReadReplicaImportState(rName),
				ImportStateVerifyIgnore: []string{"configuration_id"},
			},
		},
	})
}

func testGaussDBReadReplica_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_parameter_templates" "test" {}

locals {
  configurations = data.huaweicloud_gaussdb_parameter_templates.test.configurations
}

resource "huaweicloud_gaussdb_read_replica" "test" {
  instance_id       = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  flavor_ref        = "gaussdb.bs.s6.2xlarge.x864.ha"
  configuration_id  = [for v in local.configurations : v if v.node_type == "ha:readonly"][0].id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}

func testGaussDbReadReplicaImportState(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceId := rs.Primary.Attributes["instance_id"]
		if instanceId == "" {
			return "", fmt.Errorf("resource (%s) attributes are missing", resourceName)
		}

		return fmt.Sprintf("%s/%s", instanceId, rs.Primary.ID), nil
	}
}
