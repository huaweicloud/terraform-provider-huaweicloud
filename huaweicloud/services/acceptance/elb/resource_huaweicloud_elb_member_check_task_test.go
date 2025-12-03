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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBMemberCheckTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/members/check/jobs/{job_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func TestAccElbMemberCheckTask_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_member_check_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBMemberCheckTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccElbMemberCheckTaskConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "member_id",
						"huaweicloud_elb_member.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "subject", "securityGroup"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "result.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_result"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.severity"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.subject"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.reason_template"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.reason_params.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "check_item_total_num"),
					resource.TestCheckResourceAttrSet(resourceName, "check_item_finished_num"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func TestAccElbMemberCheckTask_all(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_member_check_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getELBMemberCheckTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccElbMemberCheckTaskConfig_all(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "member_id",
						"huaweicloud_elb_member.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "subject", "all"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "result.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_result"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.severity"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.subject"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.reason_template"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_items.0.reason_params.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.config.0.check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_result"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.severity"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.subject"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.reason_template"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_items.0.reason_params.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.acl.0.check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_result"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.severity"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.subject"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.reason_template"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_items.0.reason_params.#"),
					resource.TestCheckResourceAttrSet(resourceName, "result.0.security_group.0.check_status"),
					resource.TestCheckResourceAttrSet(resourceName, "check_item_total_num"),
					resource.TestCheckResourceAttrSet(resourceName, "check_item_finished_num"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func testAccElbMemberCheckTaskConfig_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

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
  name               = "%[1]s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name               = "%[1]s"
  vpc_id             = data.huaweicloud_vpc.test.id
  ipv4_subnet_id     = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
  waf_failure_action = "discard"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    data.huaweicloud_vpc_subnet.test.id
  ]

  lifecycle {
    ignore_changes = [
      l4_flavor_id, l7_flavor_id
    ]
  }
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[1]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  forward_eip                 = true
  idle_timeout                = 62
  request_timeout             = 63
  response_timeout            = 64
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%[1]s"
  protocol    = "HTTP"
  lb_method   = "LEAST_CONNECTIONS"
  listener_id = huaweicloud_elb_listener.test.id
}

resource "huaweicloud_elb_member" "test" {
  address       = huaweicloud_compute_instance.test.access_ip_v4
  protocol_port = 8000
  name          = "%[1]s"
  weight        = 20
  pool_id       = huaweicloud_elb_pool.test.id
  subnet_id     = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_elb_monitor" "test" {
  pool_id          = huaweicloud_elb_pool.test.id
  name             = "%[1]s"
  protocol         = "HTTP"
  interval         = 30
  timeout          = 20
  max_retries      = 8
  max_retries_down = 8
  url_path         = "/bb"
  domain_name      = "www.bb.com"
  port             = 8888
  status_code      = "200,301,404-500,504"
  http_method      = "HEAD"
  enabled          = true
}
`, rName)
}

func testAccElbMemberCheckTaskConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_member_check_task" "test" {
  depends_on = [huaweicloud_elb_monitor.test]

  member_id   = huaweicloud_elb_member.test.id
  listener_id = huaweicloud_elb_listener.test.id
  subject     = "securityGroup"
}
`, testAccElbMemberCheckTaskConfig_base(rName))
}

func testAccElbMemberCheckTaskConfig_all(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_member_check_task" "test" {
  member_id   = huaweicloud_elb_member.test.id
  listener_id = huaweicloud_elb_listener.test.id
  subject     = "all"
}
`, testAccElbMemberCheckTaskConfig_base(rName))
}
