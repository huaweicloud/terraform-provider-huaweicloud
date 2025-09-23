package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	domains "github.com/chnsz/golangsdk/openstack/waf_hw/v1/premium_domains"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
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
					resource.TestCheckResourceAttr(resourceName, "server.0.address", "119.8.2.14"),
					resource.TestCheckResourceAttr(resourceName, "server.0.port", "8083"),
					resource.TestCheckResourceAttr(resourceName, "server.1.address", "119.8.0.14"),
					resource.TestCheckResourceAttr(resourceName, "server.1.port", "8080"),
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
					resource.TestCheckResourceAttr(resourceName, "custom_page.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "server.#", "1"),
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
resource "huaweicloud_waf_certificate" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"

  certificate = <<EOT
%[3]s
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCvmuH5ViGtGOle
vJ8vOoN3Ak4pp3SescdAfQa/r4cOz/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc
8OrcaIsjns92XITVDpFW0ThGyjhTZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK
73OUYQY2E6l44U9G8Id763Bnws9NRn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu
9zu8c8i2+8qLjEsonx5PrwzNlYP3JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd
5F4Rdxl+SAIY+6mr7qu1dAlcVMLSQcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iK
aedPo2GfAgMBAAECggEAeMAvDS3uAEI2Dx/y8h3xUn9yUfBFH+6tTanrXxoK6+OT
Kj96O64qL4l3eQRflkdJiGx74FFomglCtDXxudflfXvxurkJ2hunUySQ5xScwLQt
mB6w8kP6a8IqD+bVdbn32ohk6u5dU0JZ+ErJlklVZRAGJAoCYox5DXwrEh6CP+bJ
pItgjv71tEEnX5sScQwV7FMRbjsPzXoJp8vCQjlUdetM1fk9rs3R2WSeFbPgLLtC
xY0+8Hexy0q6BLmyPZvFCaVIAzAHCYeCyzPK3xcm4odbrBmRL/amOg24CCny065N
MU9RFhEjQsY1RaK7dgkvjsntUZvU+aDcL8o6djOTuQKBgQDlDN/j2ntpGCtbTWH0
cVTW13Ze7U7iE3BfDO3m4VYP3Xi/v5FI8nHlmLrcl30H1dPKvMTec0dCBOqD1wzF
KiqHy8ELowO2CbXMYJpjuPzXH40/AE3eOJVTJM8mOeuFdeFgYCd/9cB7o5jfTA5Y
4zj8EmcRzsH1rNSnvo7/O9q6+wKBgQDERDSvP8RScEbzDKuN6uhzj1K2CAEnY6//
rDA1so18UhAie9NcAvlKa46jQTOcYD77g5h0WSlNt9ZbK9Plq9CY9psI0KNqN3Fl
YVKOKdD5m6Rifmg+lt8KLc/WocQ10DXpPTXzzuRlN/TaMDdN2pedEre/0AAMs8Ia
MIUnu4oyrQKBgQC6b6BNdqi9Ak9IIdR5g0XrGbXfzolGu0vcEkoSg5fpkfuXF/bJ
yY2rtIVkyGmc1w9tFfmol2yI8Ddy2LgsRAYaQl7/edCre3vev0LrqMck0ynE/hpj
purkojF6i+qI10p7h8ie/wmNmbv1BZMoBst7Yf9DH2gA8IynfRQn7DA9wQKBgGaU
M2kJDgX8UsjDbYKuLTIAzb0AMAIzUxBxIX1fRh2dEnvDdjOYBk1EK/fdoyjvENwJ
6ouc8j6BgBKEtKpMg6j+8wbHbTGdqrHPDQPqjSN4mpEz+i4EUqySRxep0tBBc3vl
FybHko3okhvbqXwSbL2Ww90HzI7XAPMJOv8KQO+9AoGBAJxxftNWvypBXGkPCdH2
f3ikvT2Vef9QZjqkvtipCecAkjM6ReLshVsdqFSv/ZmsVUeNKoTHvX2GnhweJM44
x7N2mFK4skBzVtMVbjAHVjG78UitVu+FrzqGreaJXHaduhgUH2iFWfw09joOotAM
X7ioLbTeWGBqFM+C80PkdBNp
-----END PRIVATE KEY-----
EOT
}

resource "huaweicloud_waf_certificate" "test-update" {
  name                  = "%[1]s_update"
  enterprise_project_id = "%[2]s"

  certificate = <<EOT
%[4]s
EOT

  private_key = <<EOT
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCvmuH5ViGtGOle
vJ8vOoN3Ak4pp3SescdAfQa/r4cOz/bmBqBcZJTX9HODhiQzdemyLLs9aOkQXYIc
8OrcaIsjns92XITVDpFW0ThGyjhTZdELj9LsbIcVzNPPclTcebZBlzAyX0oLqpHK
73OUYQY2E6l44U9G8Id763Bnws9NRn3cg0qufrlUgdim/pYZ8ubjvlDJ9eEIhcsu
9zu8c8i2+8qLjEsonx5PrwzNlYP3JqAmZ2dcbQeSPfv5U6ZceKEZfegK+Cxv4rFd
5F4Rdxl+SAIY+6mr7qu1dAlcVMLSQcLlJLRWQ5NmqL9xju7Fbj2VZt+L6nb512iK
aedPo2GfAgMBAAECggEAeMAvDS3uAEI2Dx/y8h3xUn9yUfBFH+6tTanrXxoK6+OT
Kj96O64qL4l3eQRflkdJiGx74FFomglCtDXxudflfXvxurkJ2hunUySQ5xScwLQt
mB6w8kP6a8IqD+bVdbn32ohk6u5dU0JZ+ErJlklVZRAGJAoCYox5DXwrEh6CP+bJ
pItgjv71tEEnX5sScQwV7FMRbjsPzXoJp8vCQjlUdetM1fk9rs3R2WSeFbPgLLtC
xY0+8Hexy0q6BLmyPZvFCaVIAzAHCYeCyzPK3xcm4odbrBmRL/amOg24CCny065N
MU9RFhEjQsY1RaK7dgkvjsntUZvU+aDcL8o6djOTuQKBgQDlDN/j2ntpGCtbTWH0
cVTW13Ze7U7iE3BfDO3m4VYP3Xi/v5FI8nHlmLrcl30H1dPKvMTec0dCBOqD1wzF
KiqHy8ELowO2CbXMYJpjuPzXH40/AE3eOJVTJM8mOeuFdeFgYCd/9cB7o5jfTA5Y
4zj8EmcRzsH1rNSnvo7/O9q6+wKBgQDERDSvP8RScEbzDKuN6uhzj1K2CAEnY6//
rDA1so18UhAie9NcAvlKa46jQTOcYD77g5h0WSlNt9ZbK9Plq9CY9psI0KNqN3Fl
YVKOKdD5m6Rifmg+lt8KLc/WocQ10DXpPTXzzuRlN/TaMDdN2pedEre/0AAMs8Ia
MIUnu4oyrQKBgQC6b6BNdqi9Ak9IIdR5g0XrGbXfzolGu0vcEkoSg5fpkfuXF/bJ
yY2rtIVkyGmc1w9tFfmol2yI8Ddy2LgsRAYaQl7/edCre3vev0LrqMck0ynE/hpj
purkojF6i+qI10p7h8ie/wmNmbv1BZMoBst7Yf9DH2gA8IynfRQn7DA9wQKBgGaU
M2kJDgX8UsjDbYKuLTIAzb0AMAIzUxBxIX1fRh2dEnvDdjOYBk1EK/fdoyjvENwJ
6ouc8j6BgBKEtKpMg6j+8wbHbTGdqrHPDQPqjSN4mpEz+i4EUqySRxep0tBBc3vl
FybHko3okhvbqXwSbL2Ww90HzI7XAPMJOv8KQO+9AoGBAJxxftNWvypBXGkPCdH2
f3ikvT2Vef9QZjqkvtipCecAkjM6ReLshVsdqFSv/ZmsVUeNKoTHvX2GnhweJM44
x7N2mFK4skBzVtMVbjAHVjG78UitVu+FrzqGreaJXHaduhgUH2iFWfw09joOotAM
X7ioLbTeWGBqFM+C80PkdBNp
-----END PRIVATE KEY-----
EOT

  lifecycle {
    ignore_changes = [
      certificate,
    ]
  }
}

resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
}

# Using this datasource to get the vpc_id to which the dedicated instance belongs.
data "huaweicloud_waf_dedicated_instances" "test" {
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, certificateBody, generateCertificateBody())
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
  policy_id             = huaweicloud_waf_policy.test.id
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
    address         = "119.8.2.14"
    port            = 8083
    type            = "ipv4"
    vpc_id          = data.huaweicloud_waf_dedicated_instances.test.instances[0].vpc_id
  }

  server {
    client_protocol = "HTTPS"
    server_protocol = "HTTP"
    address         = "119.8.0.14"
    port            = 8080
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
  certificate_id        = huaweicloud_waf_certificate.test-update.id
  policy_id             = huaweicloud_waf_policy.test.id
  keep_policy           = false
  proxy                 = true
  tls                   = "TLS v1.2"
  cipher                = "cipher_2"
  pci_3ds               = true
  pci_dss               = true
  protect_status        = 0
  website_name          = "%[2]s_update"
  description           = "test description update"
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
  policy_id             = huaweicloud_waf_policy.test.id
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
    port            = 8080
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

func TestAccDedicateDomain_cloud_elb_domain(t *testing.T) {
	var (
		obj interface{}

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
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafDedicatedDomain_cloud_elb_domain_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id", "huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id", "huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "domain", "119.8.5.14"),
					resource.TestCheckResourceAttr(resourceName, "mode", "elb-shared"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "test-web-tag"),
					resource.TestCheckResourceAttrSet(resourceName, "access_status"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.%"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.%"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "proxy"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_identifier.%"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.#"),
					resource.TestCheckResourceAttrSet(resourceName, "timeout_settings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_mark.#"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "extend.%"),
					resource.TestCheckResourceAttrSet(resourceName, "flag.#"),
					resource.TestCheckResourceAttrSet(resourceName, "timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "block_page.#"),
				),
			},
			{
				Config: testAccWafDedicatedDomain_cloud_elb_domain_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id", "huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id", "huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "policy_id", "huaweicloud_waf_policy.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "domain", "119.8.5.14"),
					resource.TestCheckResourceAttr(resourceName, "mode", "elb-shared"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "website_name", "test-web-tag-update"),
					resource.TestCheckResourceAttr(resourceName, "traffic_identifier.ip_tags", "ip_tag,$remote_addr"),
					resource.TestCheckResourceAttr(resourceName, "traffic_identifier.session_tag", "session_tag"),
					resource.TestCheckResourceAttr(resourceName, "traffic_identifier.user_tag", "user_tag"),
					resource.TestCheckResourceAttrSet(resourceName, "access_status"),
					resource.TestCheckResourceAttrSet(resourceName, "alarm_page.%"),
					resource.TestCheckResourceAttrSet(resourceName, "compliance_certification.%"),
					resource.TestCheckResourceAttrSet(resourceName, "protect_status"),
					resource.TestCheckResourceAttrSet(resourceName, "protocol"),
					resource.TestCheckResourceAttrSet(resourceName, "proxy"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_protection.#"),
					resource.TestCheckResourceAttrSet(resourceName, "timeout_settings.#"),
					resource.TestCheckResourceAttrSet(resourceName, "traffic_mark.#"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "extend.%"),
					resource.TestCheckResourceAttrSet(resourceName, "flag.#"),
					resource.TestCheckResourceAttrSet(resourceName, "timestamp"),
					resource.TestCheckResourceAttrSet(resourceName, "block_page.#"),
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

func testAccWafDedicatedDomain_cloud_elb_domain_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[2]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name                = "%[2]s"
  vpc_id              = huaweicloud_vpc.test.id
  ipv4_subnet_id      = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  ipv4_eip_id         = huaweicloud_vpc_eip.test.id
  waf_failure_action  = "forward"
  autoscaling_enabled = true

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  backend_subnets = [
    huaweicloud_vpc_subnet.test.id
  ]

  protection_status = "nonProtection"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "test" {
  name                        = "%[2]s"
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
`, common.TestVpc(name), name)
}

func testAccWafDedicatedDomain_cloud_elb_domain_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_dedicated_domain" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  listener_id     = huaweicloud_elb_listener.test.id
  domain          = "119.8.5.14"
  mode            = "elb-shared"
  description     = "test description"
  protocol_port   = 8080
  website_name    = "test-web-tag"

  lifecycle {
    ignore_changes = [ proxy ]
  }
}
`, testAccWafDedicatedDomain_cloud_elb_domain_base(name))
}

func testAccWafDedicatedDomain_cloud_elb_domain_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_policy" "test" {
  name = "%[2]s"
}

resource "huaweicloud_waf_dedicated_domain" "test" {
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  listener_id     = huaweicloud_elb_listener.test.id
  domain          = "119.8.5.14"
  mode            = "elb-shared"
  description     = "test description update"
  protocol_port   = 8080
  website_name    = "test-web-tag-update"
  policy_id       = huaweicloud_waf_policy.test.id

  traffic_mark {
    ip_tags     = ["ip_tag", "$remote_addr"]
    session_tag = "session_tag"
    user_tag    = "user_tag"
  }

  lifecycle {
    ignore_changes = [ proxy ]
  }
}
`, testAccWafDedicatedDomain_cloud_elb_domain_base(name), name)
}
