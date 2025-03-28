package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/loadbalancers"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getElbLoadBalancerCopyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving ELB LoadBalancer copy: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccElbLoadBalancerCopy_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_loadbalancer_copy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getElbLoadBalancerCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbLoadBalancerCopyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "availability_zone.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv4_subnet_id",
						"huaweicloud_vpc_subnet.test.0", "ipv4_subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_address", "192.168.0.216"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_network_id",
						"huaweicloud_vpc_subnet.test.0", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6_address"),
					resource.TestCheckResourceAttrPair(resourceName, "backend_subnets.0",
						"huaweicloud_vpc_subnet.test.0", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "l4_flavor_id",
						"data.huaweicloud_elb_flavors.l4_flavors", "flavors.0.id"),
					resource.TestCheckResourceAttrPair(resourceName, "l7_flavor_id",
						"data.huaweicloud_elb_flavors.l7_flavors", "flavors.0.id"),
					resource.TestCheckResourceAttr(resourceName, "description", "test elb description"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_bandwidth_id",
						"huaweicloud_vpc_bandwidth.test.0", "id"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "deletion_protection_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "waf_failure_action", "discard"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "elb_virsubnet_type"),
					resource.TestCheckResourceAttrSet(resourceName, "operating_status"),
					resource.TestCheckResourceAttrSet(resourceName, "public_border_group"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbLoadBalancerCopyConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "availability_zone.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone.0",
						"data.huaweicloud_availability_zones.test", "names.1"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv4_subnet_id",
						"huaweicloud_vpc_subnet.test.1", "ipv4_subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "ipv4_address", "192.168.1.118"),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_network_id",
						"huaweicloud_vpc_subnet.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "backend_subnets.#", "2"),
					resource.TestCheckResourceAttrPair(resourceName, "l4_flavor_id",
						"data.huaweicloud_elb_flavors.l4_flavors", "flavors.1.id"),
					resource.TestCheckResourceAttrPair(resourceName, "l7_flavor_id",
						"data.huaweicloud_elb_flavors.l7_flavors", "flavors.1.id"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "ipv6_bandwidth_id",
						"huaweicloud_vpc_bandwidth.test.1", "id"),
					resource.TestCheckResourceAttr(resourceName, "cross_vpc_backend", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "deletion_protection_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "waf_failure_action", "forward"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner_update", "terraform_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"loadbalancer_id",
					"ipv6_bandwidth_id",
					"deletion_protection_enable",
					"reuse_pool",
					"period_unit",
					"period",
					"auto_renew",
				},
			},
		},
	})
}
func testAccElbLoadBalancerCopyConfig_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_elb_flavors" "l4_flavors" {
  type = "L4"
}

data "huaweicloud_elb_flavors" "l7_flavors" {
  type = "L7"
}

resource "huaweicloud_vpc_bandwidth" "test" {
  count = 2

  name = "%[2]s_${count.index}"
  size = 5
}

resource "huaweicloud_vpc" "test" {
  name = "%[2]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 2

  name        = "%[2]s_${count.index}"
  vpc_id      = huaweicloud_vpc.test.id
  cidr        = "192.168.${count.index}.0/24"
  gateway_ip  = "192.168.${count.index}.1"
  ipv6_enable = true
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name               = "%[2]s"
  vpc_id             = huaweicloud_vpc.test.id
  ipv4_subnet_id     = huaweicloud_vpc_subnet.test[0].ipv4_subnet_id
  waf_failure_action = "discard"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test[0].id
  ]
}
`, common.TestSecGroup(rName), rName)
}

func testAccElbLoadBalancerCopyConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_loadbalancer_copy" "test" {
  loadbalancer_id   = huaweicloud_elb_loadbalancer.test.id
  name              = "%[2]s"
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test[0].ipv4_subnet_id
  ipv4_address      = "192.168.0.216"
  ipv6_network_id   = huaweicloud_vpc_subnet.test[0].id
  backend_subnets   = [huaweicloud_vpc_subnet.test[0].id]
  l4_flavor_id      = data.huaweicloud_elb_flavors.l4_flavors.flavors[0].id
  l7_flavor_id      = data.huaweicloud_elb_flavors.l7_flavors.flavors[0].id
  reuse_pool        = true


  description                = "test elb description"
  ipv6_bandwidth_id          = huaweicloud_vpc_bandwidth.test[0].id
  cross_vpc_backend          = "false"
  protection_status          = "consoleProtection"
  protection_reason          = "test protection reason"
  deletion_protection_enable = "true"
  waf_failure_action         = "discard"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbLoadBalancerCopyConfig_base(rName), rName)
}

func testAccElbLoadBalancerCopyConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_loadbalancer_copy" "test" {
  loadbalancer_id   = huaweicloud_elb_loadbalancer.test.id
  name              = "%[2]s"
  availability_zone = [data.huaweicloud_availability_zones.test.names[1]]
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test[1].ipv4_subnet_id
  ipv4_address      = "192.168.1.118"
  ipv6_network_id   = huaweicloud_vpc_subnet.test[1].id
  backend_subnets   = [huaweicloud_vpc_subnet.test[0].id, huaweicloud_vpc_subnet.test[1].id]
  l4_flavor_id      = data.huaweicloud_elb_flavors.l4_flavors.flavors[1].id
  l7_flavor_id      = data.huaweicloud_elb_flavors.l7_flavors.flavors[1].id
  reuse_pool        = true

  description                = ""
  ipv6_bandwidth_id          = huaweicloud_vpc_bandwidth.test[1].id
  cross_vpc_backend          = "true"
  protection_status          = "nonProtection"
  deletion_protection_enable = "false"
  waf_failure_action         = "forward"

  tags = {
    key_update   = "value_update"
    owner_update = "terraform_update"
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "false"
}
`, testAccElbLoadBalancerCopyConfig_base(rName), rNameUpdate)
}
