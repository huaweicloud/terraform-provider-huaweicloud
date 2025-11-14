package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceIngressAssociatedDomains_basic(t *testing.T) {
	var (
		httpPort       = "data.huaweicloud_apig_instance_ingress_associated_domains.filter_by_http_port"
		dcForHTTPPort  = acceptance.InitDataSourceCheck(httpPort)
		httpsPort      = "data.huaweicloud_apig_instance_ingress_associated_domains.filter_by_https_port"
		dcForHTTPSPort = acceptance.InitDataSourceCheck(httpsPort)

		notFound   = "data.huaweicloud_apig_instance_ingress_associated_domains.filter_by_not_found_port"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckDnsZoneNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataInstanceIngressAssociatedDomains_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcForHTTPPort.CheckResourceExists(),
					resource.TestCheckOutput("http_port_filter_is_useful", "true"),
					resource.TestMatchResourceAttr(httpPort, "domains.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(httpPort, "domains.0.name"),
					resource.TestCheckResourceAttrSet(httpPort, "domains.0.group_id"),
					resource.TestCheckResourceAttrSet(httpPort, "domains.0.group_name"),
					dcForHTTPSPort.CheckResourceExists(),
					resource.TestCheckOutput("https_port_filter_is_useful", "true"),
					resource.TestMatchResourceAttr(httpsPort, "domains.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(httpsPort, "domains.0.name"),
					resource.TestCheckResourceAttrSet(httpsPort, "domains.0.group_id"),
					resource.TestCheckResourceAttrSet(httpsPort, "domains.0.group_name"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(notFound, "domains.#", "0"),
				),
			},
		},
	})
}

func testAccDataInstanceIngressAssociatedDomains_basic(name string) string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_apig_instances" "test" {
  instance_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_apig_instances.test.instances[0].id
  domain_name = try(split(",", "%[2]s")[0], "NOT_FOUND")
}

data "huaweicloud_dns_zones" "test" {
  zone_type   = "public"
  name        = local.domain_name
  search_mode = "equal"
}

resource "huaweicloud_apig_group" "test" {
  instance_id = local.instance_id
  name        = "%[3]s"
}

resource "huaweicloud_apig_instance_ingress_port" "http" {
  instance_id = local.instance_id
  protocol    = "HTTP"
  port        = 8080
}

resource "huaweicloud_apig_instance_ingress_port" "https" {
  instance_id = local.instance_id
  protocol    = "HTTPS"
  port        = 8443
}

resource "huaweicloud_apig_group_domain_associate" "test" {
  depends_on = [
    huaweicloud_apig_instance_ingress_port.http,
    huaweicloud_apig_instance_ingress_port.https
  ]

  instance_id = local.instance_id
  group_id    = huaweicloud_apig_group.test.id
  url_domain  = try(data.huaweicloud_dns_zones.test.zones[0].name, "NOT_FOUND")

  min_ssl_version           = "TLSv1.1"
  ingress_http_port         = huaweicloud_apig_instance_ingress_port.http.port
  ingress_https_port        = huaweicloud_apig_instance_ingress_port.https.port
  is_http_redirect_to_https = false
}

data "huaweicloud_apig_instance_ingress_associated_domains" "filter_by_http_port" {
  depends_on = [
    huaweicloud_apig_group_domain_associate.test
  ]

  instance_id     = local.instance_id
  ingress_port_id = huaweicloud_apig_instance_ingress_port.http.id
}

locals {
  http_port_filter_result = [
    for v in data.huaweicloud_apig_instance_ingress_associated_domains.filter_by_http_port.domains : v.name == local.domain_name
  ]
}

output "http_port_filter_is_useful" {
  value = length(local.http_port_filter_result) > 0 && alltrue(local.http_port_filter_result)
}

data "huaweicloud_apig_instance_ingress_associated_domains" "filter_by_https_port" {
  depends_on = [
    huaweicloud_apig_group_domain_associate.test
  ]

  instance_id     = local.instance_id
  ingress_port_id = huaweicloud_apig_instance_ingress_port.https.id
}

locals {
  https_port_filter_result = [
    for v in data.huaweicloud_apig_instance_ingress_associated_domains.filter_by_https_port.domains : v.name == local.domain_name
  ]
}

output "https_port_filter_is_useful" {
  value = length(local.https_port_filter_result) > 0 && alltrue(local.https_port_filter_result)
}

data "huaweicloud_apig_instance_ingress_associated_domains" "filter_by_not_found_port" {
  depends_on = [
    huaweicloud_apig_group_domain_associate.test
  ]

  instance_id     = local.instance_id
  ingress_port_id = "%[4]s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, acceptance.HW_DNS_ZONE_NAMES, name, randUUID)
}
