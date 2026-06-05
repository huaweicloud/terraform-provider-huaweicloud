package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/elb/v3/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccElbListenerCopy_basic(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_listener_copy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getELBListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbListenerCopy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", "62"),
					resource.TestCheckResourceAttr(resourceName, "request_timeout", "63"),
					resource.TestCheckResourceAttr(resourceName, "response_timeout", "64"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbListenerCopy_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "idle_timeout", "63"),
					resource.TestCheckResourceAttr(resourceName, "request_timeout", "64"),
					resource.TestCheckResourceAttr(resourceName, "response_timeout", "65"),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"listener_id",
					"reuse_pool",
				},
			},
		},
	})
}

func testAccElbListenerCopy_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      l4_flavor_id, l7_flavor_id
    ]
  }
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8000
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = false

  idle_timeout     = 62
  request_timeout  = 63
  response_timeout = 64

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, rName)
}

func testAccElbListenerCopy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener_copy" "test" {
  listener_id     = huaweicloud_elb_listener.test.id
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  name            = "%[2]s"
  protocol_port   = 8080
  description     = "test description"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbListenerCopy_base(rName), rName)
}

func testAccElbListenerCopy_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener_copy" "test" {
  listener_id     = huaweicloud_elb_listener.test.id
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  name            = "%[2]s"
  description     = "test description update"
  protocol_port   = 8080

  idle_timeout     = 63
  request_timeout  = 64
  response_timeout = 65

  forward_eip          = true
  forward_port         = true
  forward_request_port = true
  forward_host         = false

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, testAccElbListenerCopy_base(rName), rNameUpdate, rNameUpdate)
}
