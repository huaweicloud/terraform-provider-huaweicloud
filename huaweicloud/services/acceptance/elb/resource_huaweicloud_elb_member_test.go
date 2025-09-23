package elb

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

func getELBMemberResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}/members/{member_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", state.Primary.Attributes["pool_id"])
	getPath = strings.ReplaceAll(getPath, "{member_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccElbV3Member_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	updateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_member.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBMemberResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3MemberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "weight", "20"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "ipv4_subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "address",
						"huaweicloud_compute_instance.test", "access_ip_v4"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8000"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_version"),
					resource.TestCheckResourceAttrSet(resourceName, "member_type"),
					resource.TestCheckResourceAttrSet(resourceName, "operating_status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbV3MemberConfig_update(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "weight", "40"),
					resource.TestCheckResourceAttrPair(resourceName, "subnet_id",
						"huaweicloud_vpc_subnet.test", "ipv4_subnet_id"),
					resource.TestCheckResourceAttrPair(resourceName, "address",
						"huaweicloud_compute_instance.test", "access_ip_v4"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "9000"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testElbMemberResourceImportState(resourceName),
			},
		},
	})
}

func testAccElbV3MemberConfig_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  type        = "instance"
  vpc_id      = huaweicloud_vpc.test.id
  description = "test pool description"

  minimum_healthy_member_count = 1

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccElbV3MemberConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 8000
  name          = "%[2]s"
  weight        = 20
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, testAccElbV3MemberConfig_base(rName), rName)
}

func testAccElbV3MemberConfig_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 9000
  name          = "%[2]s"
  weight        = 40
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, testAccElbV3MemberConfig_base(rName), rName)
}

func testElbMemberResourceImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		poolId := rs.Primary.Attributes["pool_id"]
		return fmt.Sprintf("%s/%s", poolId, rs.Primary.ID), nil
	}
}
