package eip

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

func getVpcEipv3AssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/eip/publicips/{publicip_id}"
		product = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, err
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{publicip_id}", state.Primary.ID)

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

	associateInstanceId := utils.PathSearch("publicip.associate_instance_id", getRespBody, "").(string)
	if associateInstanceId == "" {
		return nil, golangsdk.ErrDefault404{}
	}
	return getRespBody, nil
}

func TestAccVpcEipv3Associate_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_eipv3_associate.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getVpcEipv3AssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcEipv3Associate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "publicip_id",
						"huaweicloud_vpc_eip.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "associate_instance_type", "ELB"),
					resource.TestCheckResourceAttrPair(resourceName, "associate_instance_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVpcEipv3Associate_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                = "%[2]s"
  vpc_id              = huaweicloud_vpc.test.id
  ipv4_subnet_id      = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  waf_failure_action  = "discard"
  autoscaling_enabled = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]
}
`, common.TestVpc(rName), rName)
}

func testAccVpcEipv3Associate_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_eipv3_associate" "test" {
  publicip_id             = huaweicloud_vpc_eip.test.id
  associate_instance_type = "ELB"
  associate_instance_id   = huaweicloud_elb_loadbalancer.test.id
}
`, testAccVpcEipv3Associate_base(rName))
}
