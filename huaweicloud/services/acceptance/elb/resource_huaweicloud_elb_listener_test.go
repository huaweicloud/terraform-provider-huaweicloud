package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getELBListenerResourceFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.ElbV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ELB client: %s", err)
	}
	return listeners.Get(client, state.Primary.ID).Extract()
}

func TestAccElbV3Listener_basic(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_listener.test"

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
				Config: testAccElbV3ListenerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
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

func TestAccElbV3Listener_with_port_ranges(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_listener.test"

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
				Config: testAccElbV3ListenerConfig_with_port_ranges(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "true"),
					resource.TestCheckResourceAttr(resourceName, "port_ranges.0.start_port", "8000"),
					resource.TestCheckResourceAttr(resourceName, "port_ranges.0.end_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
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

func TestAccElbV3Listener_with_default_pool(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_listener.test"

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
				Config: testAccElbV3ListenerConfig_with_default_pool(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_with_default_pool_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
		},
	})
}

func TestAccElbV3Listener_Https(t *testing.T) {
	var listener listeners.Listener
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_listener.test"

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
				Config: testAccElbV3Listener_https(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_proto", "true"),
					resource.TestCheckResourceAttr(resourceName, "real_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_elb_id", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_member_retry", "true"),
					resource.TestCheckResourceAttr(resourceName, "transparent_client_ip_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "sni_match_algo", "longest_suffix"),
					resource.TestCheckResourceAttrSet(resourceName, "security_policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccElbV3Listener_https_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_proto", "false"),
					resource.TestCheckResourceAttr(resourceName, "real_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_elb_id", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_member_retry", "false"),
					resource.TestCheckResourceAttr(resourceName, "transparent_client_ip_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "sni_match_algo", "wildcard"),
					resource.TestCheckResourceAttrSet(resourceName, "security_policy_id"),
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

func testAccElbV3ListenerConfig_basic(rName string) string {
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
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
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

func testAccElbV3ListenerConfig_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true

  idle_timeout     = 62
  request_timeout  = 63
  response_timeout = 64

  forward_eip          = true
  forward_port         = true
  forward_request_port = true
  forward_host         = false

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, rNameUpdate, rNameUpdate)
}

func testAccElbV3ListenerConfig_with_port_ranges(rName string) string {
	return fmt.Sprintf(`

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%[1]s"
  description     = "test description"
  protocol        = "UDP"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id

  port_ranges {
    start_port = 8000
    end_port   = 8080
  }

}
`, rName)
}

func testAccElbV3ListenerConfig_with_default_pool(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  force_delete = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_pool" "test" {
  name      = "%[1]s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "instance"
  vpc_id    = data.huaweicloud_vpc.test.id
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[1]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  default_pool_id             = huaweicloud_elb_pool.test.id
  advanced_forwarding_enabled = false

  idle_timeout     = 62
  request_timeout  = 63
  response_timeout = 64

  force_delete = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3ListenerConfig_with_default_pool_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[1]s"
  cross_vpc_backend = true
  ipv4_subnet_id    = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  force_delete = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_pool" "test" {
  name      = "%[1]s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "instance"
  vpc_id    = data.huaweicloud_vpc.test.id
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[1]s"
  description                 = "test description"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  default_pool_id             = huaweicloud_elb_pool.test.id
  advanced_forwarding_enabled = true

  idle_timeout     = 62
  request_timeout  = 63
  response_timeout = 64

  forward_eip          = true
  forward_port         = true
  forward_request_port = true
  forward_host         = false

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  force_delete = true

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, rNameUpdate)
}

func testAccElbV3Listener_https(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener" "test" {
  name               = "%[2]s"
  protocol           = "HTTPS"
  protocol_port      = 8080
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  server_certificate = huaweicloud_elb_certificate.test.id
  security_policy_id = huaweicloud_elb_security_policy.test[0].id

  forward_eip          = true
  forward_port         = true
  forward_request_port = true
  forward_host         = true
  forward_proto        = true
  real_ip              = true
  forward_elb_id       = true  
  enable_member_retry  = true
  sni_match_algo       = "longest_suffix"
}
`, testAccElbV3Listener_https_base(rName), rName)
}

func testAccElbV3Listener_https_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener" "test" {
  name               = "%[2]s"
  protocol           = "HTTPS"
  protocol_port      = 8080
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  server_certificate = huaweicloud_elb_certificate.test.id
  security_policy_id = huaweicloud_elb_security_policy.test[1].id

  forward_eip          = false
  forward_port         = false
  forward_request_port = false
  forward_host         = false
  forward_proto        = false
  real_ip              = false
  forward_elb_id       = false
  enable_member_retry  = false
  sni_match_algo       = "wildcard"  

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"  
}
`, testAccElbV3Listener_https_base(rName), rName)
}

func testAccElbV3Listener_https_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_security_policy" "test" {
  protocols = [
    "TLSv1.1",
    "TLSv1.2"
  ]

  ciphers = [
    "ECDHE-ECDSA-AES128-SHA",
    "ECDHE-RSA-AES256-SHA"
  ]

  name = "%[2]s-${count.index}"

  count = 2
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[2]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbV3CertificateConfig_basic(rName), rName)
}
