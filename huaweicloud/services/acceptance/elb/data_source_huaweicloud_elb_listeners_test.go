package elb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	rName := "data.huaweicloud_elb_listeners.test"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.name"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.loadbalancer_id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.description"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.default_pool_id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.http2_enable"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_elb"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_eip"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_port"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_request_port"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_host"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_proto"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_tls_certificate"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_tls_cipher"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.forward_tls_protocol"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.real_ip"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.ipgroup.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.ipgroup.0.enable_ipgroup"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.ipgroup.0.ipgroup_id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.ipgroup.0.type"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.port_ranges.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.sni_certificate.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.server_certificate"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.idle_timeout"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.request_timeout"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.response_timeout"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.loadbalancer_id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.advanced_forwarding_enabled"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protection_status"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.proxy_protocol_enable"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.quic_config.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.sni_match_algo"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.ssl_early_data_enable"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.max_connection"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.cps"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.enable_member_retry"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.gzip_enable"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.tags.%"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.updated_at"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("request_timeout_filter_is_useful", "true"),
					resource.TestCheckOutput("server_certificate_filter_is_useful", "true"),
					resource.TestCheckOutput("enable_member_retry_filter_is_useful", "true"),
					resource.TestCheckOutput("advanced_forwarding_enabled_filter_is_useful", "true"),
					resource.TestCheckOutput("http2_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("idle_timeout_filter_is_useful", "true"),
					resource.TestCheckOutput("response_timeout_filter_is_useful", "true"),
					resource.TestCheckOutput("protection_status_filter_is_useful", "true"),
					resource.TestCheckOutput("proxy_protocol_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("ssl_early_data_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbListenerConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name           = "%[3]s"
  vpc_id         = huaweicloud_vpc.test.id
  ipv4_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  
  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]
}

resource "huaweicloud_elb_pool" "test" {
  name      = "%[3]s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "instance"
  vpc_id    = huaweicloud_vpc.test.id
}

resource "huaweicloud_elb_ipgroup" "test"{
  name                  = "%[3]s"
  enterprise_project_id = "0"

  ip_list {
    ip          = "192.168.10.10"
    description = "ECS01"
  }
}

resource "huaweicloud_elb_listener" "quic" {
  name               = "%[3]s-quic"
  protocol           = "QUIC"
  protocol_port      = 80
  loadbalancer_id    = huaweicloud_elb_loadbalancer.test.id
  server_certificate = huaweicloud_elb_certificate.server.id
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[3]s"
  description                 = "test description"
  protocol                    = "HTTPS"
  protocol_port               = 8080
  loadbalancer_id             = huaweicloud_elb_loadbalancer.test.id
  server_certificate          = huaweicloud_elb_certificate.server.id
  forward_elb                 = false
  forward_proto               = false
  real_ip                     = false
  forward_tls_certificate     = false
  forward_tls_cipher          = false
  forward_tls_protocol        = false
  enable_member_retry         = false
  ssl_early_data_enable       = true
  gzip_enable                 = false
  http2_enable                = true
  sni_match_algo              = "wildcard"
  quic_listener_id            = huaweicloud_elb_listener.quic.id
  enable_quic_upgrade         = "false"
  access_policy               = "white"
  ip_group                    = huaweicloud_elb_ipgroup.test.id
  ip_group_enable             = "false"
  tls_ciphers_policy          = "tls-1-0-with-1-3"
  default_pool_id             = huaweicloud_elb_pool.test.id
  advanced_forwarding_enabled = false
  force_delete                = true

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, common.TestVpc(name), testAccElbV3CertificateConfig_basic(name), name)
}

func testAccDatasourceListeners_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_listeners" "test" {
  depends_on = [huaweicloud_elb_listener.test]
}

data "huaweicloud_elb_listeners" "name_filter" {
  name       = "%[2]s"
  depends_on = [huaweicloud_elb_listener.test]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.name_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.name_filter.listeners[*].name :v == "%[2]s"]
  )  
}

data "huaweicloud_elb_listeners" "description_filter" {
  description = huaweicloud_elb_listener.test.description
  depends_on  = [huaweicloud_elb_listener.test]
}
locals {
  description = huaweicloud_elb_listener.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.description_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.description_filter.listeners[*].description :v == local.description]
  )  
}

data "huaweicloud_elb_listeners" "request_timeout_filter" {
  request_timeout = huaweicloud_elb_listener.test.request_timeout

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  request_timeout = huaweicloud_elb_listener.test.request_timeout
}
output "request_timeout_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.request_timeout_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.request_timeout_filter.listeners[*].request_timeout :v == local.request_timeout]
  )  
}

data "huaweicloud_elb_listeners" "server_certificate_filter" {
  server_certificate = huaweicloud_elb_listener.test.server_certificate

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  server_certificate = huaweicloud_elb_listener.test.server_certificate
}
output "server_certificate_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.server_certificate_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.server_certificate_filter.listeners[*].server_certificate :
  v == local.server_certificate]
  )  
}

data "huaweicloud_elb_listeners" "enable_member_retry_filter" {
  enable_member_retry = huaweicloud_elb_listener.test.enable_member_retry

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  enable_member_retry = huaweicloud_elb_listener.test.enable_member_retry
}
output "enable_member_retry_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.enable_member_retry_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.enable_member_retry_filter.listeners[*].enable_member_retry :
  v == local.enable_member_retry]
  )  
}

data "huaweicloud_elb_listeners" "advanced_forwarding_enabled_filter" {
  advanced_forwarding_enabled = huaweicloud_elb_listener.test.advanced_forwarding_enabled

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  advanced_forwarding_enabled = huaweicloud_elb_listener.test.advanced_forwarding_enabled
}
output "advanced_forwarding_enabled_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.advanced_forwarding_enabled_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.advanced_forwarding_enabled_filter.listeners[*].advanced_forwarding_enabled :
  v == local.advanced_forwarding_enabled]
  )  
}

data "huaweicloud_elb_listeners" "http2_enable_filter" {
  http2_enable = huaweicloud_elb_listener.test.http2_enable

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  http2_enable = huaweicloud_elb_listener.test.http2_enable
}
output "http2_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.http2_enable_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.http2_enable_filter.listeners[*].http2_enable : v == local.http2_enable]
  )  
}

data "huaweicloud_elb_listeners" "idle_timeout_filter" {
  idle_timeout = huaweicloud_elb_listener.test.idle_timeout

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  idle_timeout = huaweicloud_elb_listener.test.idle_timeout
}
output "idle_timeout_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.idle_timeout_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.idle_timeout_filter.listeners[*].idle_timeout : v == local.idle_timeout]
  )  
}

data "huaweicloud_elb_listeners" "response_timeout_filter" {
  response_timeout = huaweicloud_elb_listener.test.response_timeout

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  response_timeout = huaweicloud_elb_listener.test.response_timeout
}
output "response_timeout_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.response_timeout_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.response_timeout_filter.listeners[*].response_timeout : v == local.response_timeout]
  )  
}

data "huaweicloud_elb_listeners" "protection_status_filter" {
  protection_status = huaweicloud_elb_listener.test.protection_status

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  protection_status = huaweicloud_elb_listener.test.protection_status
}
output "protection_status_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.protection_status_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.protection_status_filter.listeners[*].protection_status : v == local.protection_status]
  )  
}

data "huaweicloud_elb_listeners" "proxy_protocol_enable_filter" {
  proxy_protocol_enable = huaweicloud_elb_listener.test.proxy_protocol_enable

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  proxy_protocol_enable = huaweicloud_elb_listener.test.proxy_protocol_enable
}
output "proxy_protocol_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.proxy_protocol_enable_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.proxy_protocol_enable_filter.listeners[*].proxy_protocol_enable :
  v == local.proxy_protocol_enable]
  )  
}

data "huaweicloud_elb_listeners" "ssl_early_data_enable_filter" {
  ssl_early_data_enable = huaweicloud_elb_listener.test.ssl_early_data_enable

  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  ssl_early_data_enable = huaweicloud_elb_listener.test.ssl_early_data_enable
}
output "ssl_early_data_enable_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.ssl_early_data_enable_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.ssl_early_data_enable_filter.listeners[*].ssl_early_data_enable :
  v == local.ssl_early_data_enable]
  )  
}

data "huaweicloud_elb_listeners" "protocol_filter" {
  protocol   = huaweicloud_elb_listener.test.protocol
  depends_on = [huaweicloud_elb_listener.test]
}
locals {
  protocol = huaweicloud_elb_listener.test.protocol
}
output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.protocol_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.protocol_filter.listeners[*].protocol :v == local.protocol]
  )  
}

data "huaweicloud_elb_listeners" "protocol_port_filter" {
  protocol_port = huaweicloud_elb_listener.test.protocol_port
  depends_on    = [huaweicloud_elb_listener.test]
}
locals {
  protocol_port = huaweicloud_elb_listener.test.protocol_port
}
output "protocol_port_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.protocol_port_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.protocol_port_filter.listeners[*].protocol_port :v == local.protocol_port]
  )  
}

data "huaweicloud_elb_listeners" "loadbalancer_id_filter" {
   depends_on      = [huaweicloud_elb_listener.test]
   loadbalancer_id = huaweicloud_elb_listener.test.loadbalancer_id
}
locals {
  loadbalancer_id = huaweicloud_elb_listener.test.loadbalancer_id
}
output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.loadbalancer_id_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.loadbalancer_id_filter.listeners[*].loadbalancer_id :v == local.loadbalancer_id]
  )  
}

data "huaweicloud_elb_listeners" "listener_id_filter" {
   depends_on  = [huaweicloud_elb_listener.test]
   listener_id = huaweicloud_elb_listener.test.id
}
locals {
   listener_id = huaweicloud_elb_listener.test.id
}
output "listener_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_listeners.listener_id_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_elb_listeners.listener_id_filter.listeners[*].id :v == local.listener_id]
  )  
}
`, testAccElbListenerConfig_basic(name), name)
}
