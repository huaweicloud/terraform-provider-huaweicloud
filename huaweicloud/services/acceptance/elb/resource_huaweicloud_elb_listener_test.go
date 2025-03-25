package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v3/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "max_connection", "1000"),
					resource.TestCheckResourceAttr(resourceName, "cps", "100"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "forward_eip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_request_port", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_host", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "advanced_forwarding_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "max_connection", "2000"),
					resource.TestCheckResourceAttr(resourceName, "cps", "200"),
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

func TestAccElbV3Listener_with_protocol_https(t *testing.T) {
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
				Config: testAccElbV3ListenerConfig_protocol_https(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_elb", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_proto", "false"),
					resource.TestCheckResourceAttr(resourceName, "real_ip", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_certificate", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_cipher", "false"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_protocol", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_member_retry", "false"),
					resource.TestCheckResourceAttr(resourceName, "ssl_early_data_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "sni_match_algo", "wildcard"),
					resource.TestCheckResourceAttr(resourceName, "access_policy", "white"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_group",
						"huaweicloud_elb_ipgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_group_enable", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "quic_listener_id",
						"huaweicloud_elb_listener.quic", "id"),
					resource.TestCheckResourceAttr(resourceName, "enable_quic_upgrade", "true"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_protocol_https_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "forward_elb", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_proto", "true"),
					resource.TestCheckResourceAttr(resourceName, "real_ip", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_certificate", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_cipher", "true"),
					resource.TestCheckResourceAttr(resourceName, "forward_tls_protocol", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_member_retry", "true"),
					resource.TestCheckResourceAttr(resourceName, "sni_match_algo", "longest_suffix"),
					resource.TestCheckResourceAttr(resourceName, "ssl_early_data_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "enable_quic_upgrade", "false"),
					resource.TestCheckResourceAttr(resourceName, "access_policy", "black"),
					resource.TestCheckResourceAttrPair(resourceName, "ip_group",
						"huaweicloud_elb_ipgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ip_group_enable", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "quic_listener_id",
						"huaweicloud_elb_listener.quic", "id"),
					resource.TestCheckResourceAttr(resourceName, "enable_quic_upgrade", "false"),
					resource.TestCheckResourceAttrPair(resourceName, "security_policy_id",
						"huaweicloud_elb_security_policy.test", "id"),
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

func TestAccElbV3Listener_with_protocol_tls(t *testing.T) {
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
				Config: testAccElbV3ListenerConfig_protocol_tls(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TLS"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_protocol_tls_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "proxy_protocol_enable", "true"),
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

func TestAccElbV3Listener_with_ip_protocol(t *testing.T) {
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3ListenerConfig_with_ip_protocol(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "IP"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
			{
				Config: testAccElbV3ListenerConfig_with_ip_protocol_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "IP"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
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
  max_connection              = 1000
  cps                         = 100

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
  description                 = "test description update"
  protocol                    = "HTTP"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true
  max_connection              = 2000
  cps                         = 200

  idle_timeout     = 63
  request_timeout  = 64
  response_timeout = 65

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
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id

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
  ipv6_network_id   = data.huaweicloud_vpc_subnet.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  force_delete = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_pool" "test_update" {
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
  default_pool_id             = huaweicloud_elb_pool.test_update.id
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

func testAccElbV3ListenerConfig_loadbalancer_basic(rName string) string {
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
}`, rName)
}

func testAccElbV3ListenerConfig_protocol_quic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_listener" "quic" {
  name               = "%s-quic"
  protocol           = "QUIC"
  protocol_port      = 80
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  server_certificate = huaweicloud_elb_certificate.server.id
}`, rName)
}

func testAccSecurityPolicies_tlsv13(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_security_policy" "test" {
  protocols = [
    "TLSv1.3"
  ]
  ciphers = [
    "TLS_AES_128_GCM_SHA256"
  ]
  name = "%s"
}
`, name)
}

func testAccElbV3ListenerConfig_protocol_https_base(rName string) string {
	lb := testAccElbV3ListenerConfig_loadbalancer_basic(rName)
	serverCertificate := testAccElbV3CertificateConfig_basic(rName)
	quicListener := testAccElbV3ListenerConfig_protocol_quic(rName)
	ipGroup := testAccElbV3IpGroupConfig_basic(rName)
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

%[4]s
`, lb, serverCertificate, quicListener, ipGroup)
}

func testAccElbV3ListenerConfig_protocol_https(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener" "test" {
  name                    = "%[2]s"
  protocol                = "HTTPS"
  protocol_port           = 8080
  loadbalancer_id         = huaweicloud_elb_loadbalancer.test.id
  server_certificate      = huaweicloud_elb_certificate.server.id
  forward_elb             = false
  forward_proto           = false
  real_ip                 = false
  forward_tls_certificate = false
  forward_tls_cipher      = false
  forward_tls_protocol    = false
  enable_member_retry     = false
  ssl_early_data_enable   = true
  gzip_enable             = false
  http2_enable            = true
  sni_match_algo          = "wildcard"
  quic_listener_id        = huaweicloud_elb_listener.quic.id
  enable_quic_upgrade     = "true"
  access_policy           = "white"
  ip_group                = huaweicloud_elb_ipgroup.test.id
  ip_group_enable         = "false"
  tls_ciphers_policy      = "tls-1-0-with-1-3"
}
`, testAccElbV3ListenerConfig_protocol_https_base(rName), rName)
}

func testAccElbV3ListenerConfig_protocol_https_update(rName string) string {
	clientCertificate := testAccElbV3CertificateConfig_client(rName)
	securityPolicy := testAccSecurityPolicies_tlsv13(rName)
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_elb_listener" "test" {
  name                    = "%[4]s"
  protocol                = "HTTPS"
  protocol_port           = 8080
  loadbalancer_id         = huaweicloud_elb_loadbalancer.test.id
  forward_elb             = true
  forward_proto           = true
  real_ip                 = true
  forward_tls_certificate = true
  forward_tls_cipher      = true
  forward_tls_protocol    = true
  security_policy_id      = huaweicloud_elb_security_policy.test.id
  ssl_early_data_enable   = false
  gzip_enable             = true
  enable_member_retry     = true
  http2_enable            = false
  sni_match_algo          = "longest_suffix"
  quic_listener_id        = huaweicloud_elb_listener.quic.id
  enable_quic_upgrade     = "false"
  access_policy           = "black"
  ip_group                = huaweicloud_elb_ipgroup.test.id
  ip_group_enable         = "true"
  tls_ciphers_policy      = "tls-1-2-fs-with-1-3"
  ca_certificate          = huaweicloud_elb_certificate.client.id
  server_certificate      = huaweicloud_elb_certificate.server.id
}
`, testAccElbV3ListenerConfig_protocol_https_base(rName), clientCertificate, securityPolicy, rName)
}

func testAccElbV3ListenerConfig_protocol_tls(rName string) string {
	lb := testAccElbV3ListenerConfig_loadbalancer_basic(rName)
	certificate := testAccElbV3CertificateConfig_basic(rName)
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_elb_listener" "test" {
  name               = "%s"
  protocol           = "TLS"
  protocol_port      = 80
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  server_certificate = huaweicloud_elb_certificate.server.id
  proxy_protocol_enable = true
}`, lb, certificate, rName)
}

func testAccElbV3ListenerConfig_protocol_tls_update(rName string) string {
	lb := testAccElbV3ListenerConfig_loadbalancer_basic(rName)
	certificate := testAccElbV3CertificateConfig_basic(rName)
	return fmt.Sprintf(`
%s

%s

resource "huaweicloud_elb_listener" "test" {
  name                  = "%s"
  protocol              = "TLS"
  protocol_port         = 80
  loadbalancer_id       = huaweicloud_elb_loadbalancer.test.id
  server_certificate    = huaweicloud_elb_certificate.server.id
  proxy_protocol_enable = true
}`, lb, certificate, rName)
}

func testAccElbV3ListenerConfig_with_ip_protocol_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  loadbalancer_type = "gateway"
  description       = "test gateway description"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3ListenerConfig_with_ip_protocol(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[2]s"
  description                 = "test description"
  protocol                    = "IP"
  protocol_port               = 0
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccElbV3ListenerConfig_with_ip_protocol_base(rName), rName)
}

func testAccElbV3ListenerConfig_with_ip_protocol_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[2]s"
  description                 = "test description update"
  protocol                    = "IP"
  protocol_port               = 0
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  advanced_forwarding_enabled = true

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }
}
`, testAccElbV3ListenerConfig_with_ip_protocol_base(rName), rNameUpdate)
}
