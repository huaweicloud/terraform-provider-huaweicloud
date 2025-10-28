package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDomains_basic(t *testing.T) {
	var (
		rName      = "data.huaweicloud_cdn_domains.test"
		dc         = acceptance.InitDataSourceCheck(rName)
		domainName = generateDomainName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDomains_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "domains.0.id"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.name"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.type"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.domain_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.cname"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.https_status"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.service_area"),
					resource.TestCheckResourceAttrSet(rName, "domains.0.enterprise_project_id"),

					resource.TestCheckOutput("domain_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("domain_status_filter_is_useful", "true"),
					resource.TestCheckOutput("service_area_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDomains_basic(domainName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cdn_domains" "test" {
  depends_on = [huaweicloud_cdn_domain.test]
}

data "huaweicloud_cdn_domains" "domain_id_filter" {
  domain_id = huaweicloud_cdn_domain.test.id
}
locals {
  domain_id = huaweicloud_cdn_domain.test.id
}
output "domain_id_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.domain_id_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.domain_id_filter.domains[*].id : v == local.domain_id]
  )
}

data "huaweicloud_cdn_domains" "name_filter" {
  name = data.huaweicloud_cdn_domains.test.domains.0.name
}
locals {
  name = data.huaweicloud_cdn_domains.test.domains.0.name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.name_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.name_filter.domains[*].name : v == local.name]
  )
}

data "huaweicloud_cdn_domains" "type_filter" {
  type = data.huaweicloud_cdn_domains.test.domains.0.type
}
locals {
  type = data.huaweicloud_cdn_domains.test.domains.0.type
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.type_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.type_filter.domains[*].type : v == local.type]
  )
}

data "huaweicloud_cdn_domains" "service_area_filter" {
  service_area = data.huaweicloud_cdn_domains.test.domains.0.service_area
}
locals {
  service_area = data.huaweicloud_cdn_domains.test.domains.0.service_area
}
output "service_area_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.service_area_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.service_area_filter.domains[*].service_area 
	: v == local.service_area]
  )
}

data "huaweicloud_cdn_domains" "domain_status_filter" {
  domain_status = data.huaweicloud_cdn_domains.test.domains.0.domain_status
}
locals {
  domain_status = data.huaweicloud_cdn_domains.test.domains.0.domain_status
}
output "domain_status_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.domain_status_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.domain_status_filter.domains[*].domain_status 
	: v == local.domain_status]
  )
}

data "huaweicloud_cdn_domains" "enterprise_project_id_filter" {
  enterprise_project_id = data.huaweicloud_cdn_domains.test.domains.0.enterprise_project_id
}
locals {
  enterprise_project_id = data.huaweicloud_cdn_domains.test.domains.0.enterprise_project_id
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_cdn_domains.enterprise_project_id_filter.domains) > 0 && alltrue(
    [for v in data.huaweicloud_cdn_domains.enterprise_project_id_filter.domains[*].enterprise_project_id 
	: v == local.enterprise_project_id]
  )
}
`, testAccDomain_basic_step1(domainName))
}
