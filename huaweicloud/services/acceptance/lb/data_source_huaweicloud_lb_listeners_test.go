package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceListeners_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_lb_listeners.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceListeners_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "listeners.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.name"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protocol"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protocol_port"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.client_ca_tls_container_ref"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.description"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.connection_limit"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.http2_enable"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.insert_headers.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.insert_headers.0.x_forwarded_elb_ip"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.insert_headers.0.x_forwarded_host"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.default_tls_container_ref"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.sni_container_refs.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.protection_status"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.tls_ciphers_policy"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.updated_at"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.loadbalancers.#"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.loadbalancers.0.id"),
					resource.TestCheckResourceAttrSet(rName, "listeners.0.tags.%"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),
					resource.TestCheckOutput("protocol_port_filter_is_useful", "true"),
					resource.TestCheckOutput("client_ca_tls_container_ref_filter_is_useful", "true"),
					resource.TestCheckOutput("default_tls_container_ref_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
					resource.TestCheckOutput("http2_enable_filter_is_useful", "true"),
					resource.TestCheckOutput("listener_id_filter_is_useful", "true"),
					resource.TestCheckOutput("loadbalancer_id_filter_is_useful", "true"),
					resource.TestCheckOutput("tls_ciphers_policy_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceListeners_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_lb_listener" "test" {
  name                        = "%[4]s"
  description                 = "test description"
  protocol                    = "TERMINATED_HTTPS"
  protocol_port               = 443
  loadbalancer_id             = huaweicloud_lb_loadbalancer.loadbalancer_1.id
  default_tls_container_ref   = huaweicloud_lb_certificate.certificate_1.id
  client_ca_tls_container_ref = huaweicloud_lb_certificate.certificate_client.id
  sni_container_refs          = [huaweicloud_lb_certificate.certificate_1.id]
  tls_ciphers_policy          = "tls-1-1"
  protection_status           = "consoleProtection"
  protection_reason           = "test protection reason"

  insert_headers {
    x_forwarded_elb_ip = true
    x_forwarded_host   = true
  }

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, testAccLBV2ListenerConfig_base(rName), testAccLBV2CertificateConfig_basic(rName),
		testAccLBV2CertificateConfig_client(rName), rName)
}

func testAccDatasourceListeners_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_lb_listeners" "test" {
  depends_on = [huaweicloud_lb_listener.test]
}

data "huaweicloud_lb_listeners" "name_filter" {
  name = huaweicloud_lb_listener.test.name

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  name = huaweicloud_lb_listener.test.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.name_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.name_filter.listeners[*].name : v == local.name]
  )  
}

data "huaweicloud_lb_listeners" "protocol_filter" {
  protocol = huaweicloud_lb_listener.test.protocol

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  protocol = huaweicloud_lb_listener.test.protocol
}
output "protocol_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.protocol_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.protocol_filter.listeners[*].protocol : v == local.protocol]
  )  
}

data "huaweicloud_lb_listeners" "protocol_port_filter" {
  protocol_port = huaweicloud_lb_listener.test.protocol_port

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  protocol_port = huaweicloud_lb_listener.test.protocol_port
}
output "protocol_port_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.protocol_port_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.protocol_port_filter.listeners[*].protocol_port : v == local.protocol_port]
  )  
}

data "huaweicloud_lb_listeners" "client_ca_tls_container_ref_filter" {
  client_ca_tls_container_ref = huaweicloud_lb_listener.test.client_ca_tls_container_ref

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  client_ca_tls_container_ref = huaweicloud_lb_listener.test.client_ca_tls_container_ref
}
output "client_ca_tls_container_ref_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.client_ca_tls_container_ref_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.client_ca_tls_container_ref_filter.listeners[*].client_ca_tls_container_ref :
  v == local.client_ca_tls_container_ref]
  )  
}

data "huaweicloud_lb_listeners" "default_tls_container_ref_filter" {
  default_tls_container_ref = huaweicloud_lb_listener.test.default_tls_container_ref

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  default_tls_container_ref = huaweicloud_lb_listener.test.default_tls_container_ref
}
output "default_tls_container_ref_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.default_tls_container_ref_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.default_tls_container_ref_filter.listeners[*].default_tls_container_ref :
  v == local.default_tls_container_ref]
  )  
}

data "huaweicloud_lb_listeners" "description_filter" {
  description = huaweicloud_lb_listener.test.description

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  description = huaweicloud_lb_listener.test.description
}
output "description_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.description_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.description_filter.listeners[*].description : v == local.description]
  )  
}

data "huaweicloud_lb_listeners" "http2_enable_filter" {
  http2_enable = huaweicloud_lb_listener.test.http2_enable

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  http2_enable = huaweicloud_lb_listener.test.http2_enable
}
output "http2_enable_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.http2_enable_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.http2_enable_filter.listeners[*].http2_enable : v == local.http2_enable]
  )  
}

data "huaweicloud_lb_listeners" "listener_id_filter" {
  listener_id = huaweicloud_lb_listener.test.id

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  listener_id = huaweicloud_lb_listener.test.id
}
output "listener_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.listener_id_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.listener_id_filter.listeners[*].id : v == local.listener_id]
  )  
}

data "huaweicloud_lb_listeners" "loadbalancer_id_filter" {
  loadbalancer_id = huaweicloud_lb_listener.test.loadbalancer_id

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  loadbalancer_id = huaweicloud_lb_listener.test.loadbalancer_id
}
output "loadbalancer_id_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.loadbalancer_id_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.loadbalancer_id_filter.listeners[*].loadbalancers[0].id : v == local.loadbalancer_id]
  )  
}

data "huaweicloud_lb_listeners" "tls_ciphers_policy_filter" {
  tls_ciphers_policy = huaweicloud_lb_listener.test.tls_ciphers_policy

  depends_on = [huaweicloud_lb_listener.test]
}
locals {
  tls_ciphers_policy = huaweicloud_lb_listener.test.tls_ciphers_policy
}
output "tls_ciphers_policy_filter_is_useful" {
  value = length(data.huaweicloud_lb_listeners.tls_ciphers_policy_filter.listeners) > 0 && alltrue(
  [for v in data.huaweicloud_lb_listeners.tls_ciphers_policy_filter.listeners[*].tls_ciphers_policy : v == local.tls_ciphers_policy]
  )  
}
`, testAccDatasourceListeners_base(rName))
}
