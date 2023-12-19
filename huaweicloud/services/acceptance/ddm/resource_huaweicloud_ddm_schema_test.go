package ddm

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

func getDdmSchemaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getSchema: Query DDM schema
	var (
		getSchemaHttpUrl = "v1/{project_id}/instances/{instance_id}/databases/{ddm_dbname}"
		getSchemaProduct = "ddm"
	)
	getSchemaClient, err := cfg.NewServiceClient(getSchemaProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DDM client: %s", err)
	}

	parts := strings.SplitN(state.Primary.ID, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, must be <instance_id>/<db_name>")
	}
	instanceID := parts[0]
	schemaName := parts[1]
	getSchemaPath := getSchemaClient.Endpoint + getSchemaHttpUrl
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{project_id}", getSchemaClient.ProjectID)
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{instance_id}", fmt.Sprintf("%v", instanceID))
	getSchemaPath = strings.ReplaceAll(getSchemaPath, "{ddm_dbname}", fmt.Sprintf("%v", schemaName))

	getSchemaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSchemaResp, err := getSchemaClient.Request("GET", getSchemaPath, &getSchemaOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DdmSchema: %s", err)
	}
	getSchemaRespBody, err := utils.FlattenResponse(getSchemaResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("database", getSchemaRespBody, nil), nil
}

func TestAccDdmSchema_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	instanceName := strings.ReplaceAll(name, "_", "-")
	rName := "huaweicloud_ddm_schema.test"
	dbPwd := "test_1234"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDdmSchemaResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDdmSchema_basic(instanceName, name, dbPwd),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "shard_mode", "single"),
					resource.TestCheckResourceAttr(rName, "shard_number", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"instance_id", "data_nodes.0.admin_user", "data_nodes.0.status",
					"data_nodes.0.admin_password", "delete_rds_data"},
			},
		},
	})
}

func testDdmSchema_base(name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_ddm_engines" test {
  version = "3.0.8.5"
}

data "huaweicloud_ddm_flavors" test {
 engine_id = data.huaweicloud_ddm_engines.test.engines[0].id
 cpu_arch  = "X86"
}

resource "huaweicloud_ddm_instance" "test" {
 name              = "%[2]s"
 flavor_id         = data.huaweicloud_ddm_flavors.test.flavors[0].id
 node_num          = 2
 engine_id         = data.huaweicloud_ddm_engines.test.engines[0].id
 vpc_id            = huaweicloud_vpc.test.id
 subnet_id         = huaweicloud_vpc_subnet.test.id
 security_group_id = huaweicloud_networking_secgroup.test.id

 availability_zones = [
   data.huaweicloud_availability_zones.test.names[0]
 ]
}

resource "huaweicloud_rds_instance" "test" {
  name               = "%[2]s"
  flavor             = "rds.mysql.n1.large.4"
  security_group_id  = huaweicloud_networking_secgroup.test.id
  subnet_id          = huaweicloud_vpc_subnet.test.id
  vpc_id             = huaweicloud_vpc.test.id

 availability_zone = [
   data.huaweicloud_availability_zones.test.names[0]
 ]

 db {
   password = "%[3]s"
   type     = "MySQL"
   version  = "5.7"
   port     = 3306
 }

 volume {
   type = "CLOUDSSD"
   size = 40
 }
}
`, common.TestVpc(name), name, dbPwd)
}

func testDdmSchema_basic(instanceName, name, dbPwd string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ddm_schema" "test" {
  instance_id  = huaweicloud_ddm_instance.test.id
  name         = "%[2]s"
  shard_mode   = "single"
  shard_number = "1"

  data_nodes {
    id             = huaweicloud_rds_instance.test.id
    admin_user     = "root"
    admin_password = "%[3]s"
  }

  delete_rds_data = "true"

  lifecycle {
    ignore_changes = [
      data_nodes,
    ]
  }
}
`, testDdmSchema_base(instanceName, dbPwd), name, dbPwd)
}
