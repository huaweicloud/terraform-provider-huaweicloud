package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getL7ListenerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB v2 Client: %s", err)
	}
	resp, err := listeners.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("unable to find the listener (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_listener.listener_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getL7ListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2ListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "connection_limit", "-1"),
				),
			},
			{
				Config: testAccLBV2ListenerConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccLBV2Listener_https(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_listener.listener_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&listener,
		getL7ListenerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2ListenerConfig_https(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "false"),
				),
			},
			{
				Config: testAccLBV2ListenerConfig_https_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "http2_enable", "true"),
				),
			},
		},
	})
}

func testAccLBV2ListenerConfig_base(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  description   = "created by acceptance test"
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}
`, common.TestVpc(rName), rName)
}

func testAccLBV2ListenerConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_listener" "listener_1" {
  name             = "%s"
  description      = "created by acceptance test"
  protocol         = "HTTP"
  protocol_port    = 8080
  loadbalancer_id  = huaweicloud_lb_loadbalancer.loadbalancer_1.id

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccLBV2ListenerConfig_base(rName), rName)
}

func testAccLBV2ListenerConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id

  tags = {
    foo   = "bar"
    owner = "terraform_update"
  }
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName), rNameUpdate)
}

func testAccLBV2ListenerConfig_https(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_lb_listener" "listener_1" {
  name                      = "%s"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 443
  loadbalancer_id           = huaweicloud_lb_loadbalancer.loadbalancer_1.id
  default_tls_container_ref = huaweicloud_lb_certificate.certificate_1.id
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName), rName)
}

func testAccLBV2ListenerConfig_https_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_lb_listener" "listener_1" {
  name                      = "%s"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 443
  loadbalancer_id           = huaweicloud_lb_loadbalancer.loadbalancer_1.id
  default_tls_container_ref = huaweicloud_lb_certificate.certificate_1.id
  http2_enable              = true
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName), rNameUpdate)
}
