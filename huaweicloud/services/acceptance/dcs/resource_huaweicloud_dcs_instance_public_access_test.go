package dcs

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

func getDcsInstancePublicAccessResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating DCS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	publicInfo := utils.PathSearch("publicip_info", getRespBody, nil)
	if publicInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccDcsInstancePublicAccess_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dcs_instance_public_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDcsInstancePublicAccessResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDcsInstancePublicAccess_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "elb_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "eip_id"),
					resource.TestCheckResourceAttrSet(rName, "eip_address"),
					resource.TestCheckResourceAttrSet(rName, "elb_listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "elb_listeners.0.id"),
					resource.TestCheckResourceAttrSet(rName, "elb_listeners.0.port"),
					resource.TestCheckResourceAttrSet(rName, "elb_listeners.0.name"),
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

func testDcsInstancePublicAccess_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  cache_mode       = "ha"
  capacity         = 1
  cpu_architecture = "x86_64"
}

resource "huaweicloud_dcs_instance" "test" {
  name               = "%[1]s"
  engine_version     = "5.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[1]s"
  cross_vpc_backend = true
  vpc_id            = data.huaweicloud_vpc.test.id
  ipv4_subnet_id    = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]

  backend_subnets = [
    data.huaweicloud_vpc_subnet.test.id
  ]
}

resource "huaweicloud_vpc_eipv3_associate" "test" {
  publicip_id             = huaweicloud_vpc_eip.test.id
  associate_instance_type = "ELB"
  associate_instance_id   = huaweicloud_elb_loadbalancer.test.id
}
`, name)
}

func testDcsInstancePublicAccess_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_instance_public_access" "test" {
  depends_on = [huaweicloud_vpc_eipv3_associate.test]

  instance_id = huaweicloud_dcs_instance.test.id
  elb_id      = huaweicloud_elb_loadbalancer.test.id
}
`, testDcsInstancePublicAccess_base(name))
}
