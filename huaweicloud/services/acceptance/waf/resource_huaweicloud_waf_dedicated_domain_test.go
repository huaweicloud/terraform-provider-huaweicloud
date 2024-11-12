package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getWafDedicateDomainResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.WafDedicatedV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF dedicated client: %s", err)
	}

	epsID := state.Primary.Attributes["enterprise_project_id"]
	return domains.GetWithEpsID(client, state.Primary.ID, epsID)
}

// Before running the test case, please ensure that there is at least one WAF dedicated instance in the current region.
func TestAccDedicateDomain_basic(t *testing.T) {
	var (
		obj             interface{}
		certificateBody = generateCertificateBody()

		randName     = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_waf_dedicated_domain.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getWafDedicateDomainResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomain_basic(randName, certificateBody),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "false"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.1"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_1"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "1"),
					resource.TestCheckResourceAttr(resourceName, "website_name", randName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.http_return_code", "404"),
					resource.TestCheckResourceAttr(resourceName, "custom_page.0.block_page_type", "application/json"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key1", "$time_local"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$tenant_id"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "false"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "50"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "200"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag"),
					resource.TestCheckResourceAttrSet(resourceName, "custom_page.0.page_content"),
					resource.TestCheckResourceAttrSet(resourceName, "server.0.vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_id"),
					resource.TestCheckResourceAttrSet(resourceName, "certificate_name"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "tls"),
					resource.TestCheckResourceAttrSet(resourceName, "cipher"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.template_name"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_3ds"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.pci_dss"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.error_threshold"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.error_percentage"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.initial_downtime"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.pending_url_request_threshold"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.0.duration"),
				),
			},
			{
				Config: testAccWafDedicatedDomain_update1(randName, certificateBody),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "cipher", "cipher_2"),
					resource.TestCheckResourceAttr(resourceName, "pci_3ds", "true"),
					resource.TestCheckResourceAttr(resourceName, "pci_dss", "true"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "website_name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "redirect_url", "${http_host}/error.html"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8443"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.15"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key2", "$request_length"),
					resource.TestCheckResourceAttr(resourceName, "forward_header_map.key3", "$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "1000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "87.5"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "200"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "5"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "7000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "10000"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "100"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "1000"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "1000"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "ip_tag_another"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag_update"),
				),
			},
			{
				Config: testAccWafDedicatedDomain_update2(randName, certificateBody),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "99"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "2147483647"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "false"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "180"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "3600"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.0", "ip_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.ip_tags.1", "ip_tag_another"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.session_tag", "session_tag_update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_mark.0.user_tag", "user_tag_update"),
				),
			},
			{
				Config: testAccWafDedicatedDomain_update3(randName, certificateBody),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.error_percentage", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.initial_downtime", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.multiplier_for_consecutive_breakdowns", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.pending_url_request_threshold", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.duration", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_protection.0.status", "true"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.connection_timeout", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.read_timeout", "0"),
					resource.TestCheckResourceAttr(resourceName, "timeout_settings.0.write_timeout", "0"),
				),
			},
			{
				Config: testAccWafDedicatedDomain_policy(randName, certificateBody),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "domain", fmt.Sprintf("www.%s.com", randName)),
					resource.TestCheckResourceAttr(resourceName, "proxy", "true"),
					resource.TestCheckResourceAttr(resourceName, "tls", "TLS v1.2"),
					resource.TestCheckResourceAttr(resourceName, "protect_status", "0"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.0.client_protocol", "HTTPS"),
					resource.TestCheckResourceAttr(resourceName, "server.0.server_protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "server.0.type", "ipv4"),
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.0.14"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"keep_policy"},
				ImportStateIdFunc:       testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafDedicatedDomain_base(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

# Using this datasource to get the vpc_id to which the dedicated instance belongs.
data "huaweicloud_waf_dedicated_instances" "test" {
  enterprise_project_id = "%[2]s"
}
`, testAccWafCertificate_basic(name, certificateBody), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedDomain_basic(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain                = "www.%[2]s.com"
  certificate_id        = huaweicloud_waf_certificate.test.id
  keep_policy           = false
  proxy                 = false
  tls                   = "TLS v1.1"
  cipher                = "cipher_1"
  protect_status        = 1
  website_name          = "%[2]s"
  description           = "test description"
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  custom_page {
    http_return_code = "404"
    block_page_type  = "application/json"
    page_content     = <<EOF
{
  "event_id": "$${waf_event_id}",
  "error_msg": "error message"
}
EOF
  }

  forward_header_map = {
    "key1" = "$time_local"
    "key2" = "$tenant_id"
  }

  timeout_settings {
    connection_timeout = 50
    read_timeout       = 200
    write_timeout      = 200
  }

  traffic_mark {
    ip_tags     = ["ip_tag", "$remote_addr"]
    session_tag = "session_tag"
    user_tag    = "user_tag"
  }
}
`, testAccWafDedicatedDomain_base(name, certificateBody), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedDomain_update1(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain                = "www.%[2]s.com"
  certificate_id        = huaweicloud_waf_certificate.test.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  protect_status        = 0
  website_name          = "%[2]s_update"
  description           = "test description update"
  redirect_url          = "$${http_host}/error.html"
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 1000
    error_percentage                      = 87.5
    initial_downtime                      = 200
    multiplier_for_consecutive_breakdowns = 5
    pending_url_request_threshold         = 7000
    duration                              = 10000
    status                                = true
  }

  timeout_settings {
    connection_timeout = 100
    read_timeout       = 1000
    write_timeout      = 1000
  }

  traffic_mark {
    ip_tags     = ["ip_tag_update", "ip_tag_another"]
    session_tag = "session_tag_update"
    user_tag    = "user_tag_update"
  }
}
`, testAccWafDedicatedDomain_base(name, certificateBody), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedDomain_update2(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain                = "www.%[2]s.com"
  certificate_id        = huaweicloud_waf_certificate.test.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  protect_status        = 0
  website_name          = "%[2]s_update"
  description           = "test description update"
  redirect_url          = "$${http_host}/error.html"
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 2147483647
    error_percentage                      = 99
    initial_downtime                      = 2147483647
    multiplier_for_consecutive_breakdowns = 2147483647
    pending_url_request_threshold         = 2147483647
    duration                              = 2147483647
    status                                = false
  }

  timeout_settings {
    connection_timeout = 180
    read_timeout       = 3600
    write_timeout      = 3600
  }
}
`, testAccWafDedicatedDomain_base(name, certificateBody), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedDomain_update3(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain                = "www.%[2]s.com"
  certificate_id        = huaweicloud_waf_certificate.test.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  protect_status        = 0
  website_name          = "%[2]s_update"
  description           = "test description update"
  redirect_url          = "$${http_host}/error.html"
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.15"
    port            = 8443
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  forward_header_map = {
    "key2" = "$request_length"
    "key3" = "$remote_addr"
  }

  connection_protection {
    error_threshold                       = 0
    error_percentage                      = 0
    initial_downtime                      = 0
    multiplier_for_consecutive_breakdowns = 0
    pending_url_request_threshold         = 0
    duration                              = 0
    status                                = true
  }

  timeout_settings {
    connection_timeout = 0
    read_timeout       = 0
    write_timeout      = 0
  }
}
`, testAccWafDedicatedDomain_base(name, certificateBody), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafDedicatedDomain_policy(name, certificateBody string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "%[3]s"
}

resource "huaweicloud_waf_dedicated_domain" "test" {
  domain         = "www.%[2]s.com"
  certificate_id = huaweicloud_waf_certificate.test.id
  policy_id      = huaweicloud_waf_policy.test.id
  keep_policy    = true
  proxy          = true
  tls            = "TLS v1.2"
  protect_status = 0
  enterprise_project_id = "%[3]s"

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }
}
`, testAccWafDedicatedDomain_base(name, certificateBody), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
