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

func getResourceGaussDbDrRelationshipFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3.5/{project_id}/disaster-recovery/relations"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetGaussDbDrRelationshipQueryParams(state.Primary.ID)

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

	relationship := utils.PathSearch(fmt.Sprintf("relations[?instance_id=='%s']|[0]", state.Primary.ID), getRespBody, nil)
	if relationship == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	status := utils.PathSearch("status", relationship, "").(string)
	if status == "completed" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func buildGetGaussDbDrRelationshipQueryParams(instanceId string) string {
	return fmt.Sprintf("?instance_id=%s", instanceId)
}

func TestAccGaussDbDrRelationship_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_gaussdb_dr_relationship.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getResourceGaussDbDrRelationshipFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGaussDbDrRelationship_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "instance_id",
						"huaweicloud_gaussdb_instance.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "disaster_type", "stream"),
					resource.TestCheckResourceAttr(rName, "dr_task_name", name),
					resource.TestCheckResourceAttrSet(rName, "synchronization_id"),
					resource.TestCheckResourceAttrSet(rName, "dr_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "disaster_role"),
					resource.TestCheckResourceAttrSet(rName, "created"),
					resource.TestCheckResourceAttrSet(rName, "updated"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.#"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.0.region_code"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.0.project_id"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.0.project_name"),
					resource.TestCheckResourceAttrSet(rName, "master_region_instance_info.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.#"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.0.region_code"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.0.project_id"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.0.project_name"),
					resource.TestCheckResourceAttrSet(rName, "slave_region_instance_info.0.ip_address"),
					resource.TestCheckResourceAttrSet(rName, "instance_name"),
					resource.TestCheckResourceAttrSet(rName, "instance_status"),
					resource.TestCheckResourceAttrSet(rName, "actions.#"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dr_ip",
					"dr_user_name",
					"dr_user_password",
				},
			},
		},
	})
}

func testGaussDbDrRelationship_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_gaussdb_flavors" "test" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[2]s"
}

resource "huaweicloud_gaussdb_instance" "test" {
  count = 2

  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id

  flavor            = data.huaweicloud_gaussdb_flavors.test.flavors[0].spec_code
  name              = "%[2]s_${count.index}"
  password          = "test_1234"
  replica_num       = 3
  availability_zone = join(",", [data.huaweicloud_availability_zones.test.names[0], 
                      data.huaweicloud_availability_zones.test.names[1], 
                      data.huaweicloud_availability_zones.test.names[2]])

  enterprise_project_id = "%[3]s"

  ha {
    mode             = "centralization_standard"
    replication_mode = "sync"
    consistency      = "strong"
    instance_mode    = "basic"
  }

  volume {
    type = "ULTRAHIGH"
    size = 40
  }
}

resource "huaweicloud_gaussdb_dr_configuration_reset" "test" {
  count = 2

  instance_id        = huaweicloud_gaussdb_instance.test[count.index].id
  opposite_data_cidr = huaweicloud_vpc_subnet.test.cidr
}

`, common.TestVpc(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testGaussDbDrRelationship_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  data_ip = [for v in huaweicloud_gaussdb_instance.test[0].nodes : 
             v if split("::", v.component_names)[length(split("::", v.component_names))-1] == "StateLeader:"][0].data_ip
}

resource "huaweicloud_gaussdb_dr_relationship" "test" {
  depends_on = [
    huaweicloud_gaussdb_dr_configuration_reset.test[0],
    huaweicloud_gaussdb_dr_configuration_reset.test[1]
  ]

  instance_id      = huaweicloud_gaussdb_instance.test[1].id
  disaster_type    = "stream"
  dr_ip            = local.data_ip
  dr_user_name     = "root"
  dr_user_password = "test_1234"
  dr_task_name     = "%[2]s"
}
`, testGaussDbDrRelationship_base(name), name)
}
