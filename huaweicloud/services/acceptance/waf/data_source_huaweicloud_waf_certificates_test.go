package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafCertificates_basic(t *testing.T) {
	var (
		dataSourceName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_waf_certificates.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_waf_certificates.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byHost   = "data.huaweicloud_waf_certificates.host_filter"
		dcByHost = acceptance.InitDataSourceCheck(byHost)

		byExpirationStatus   = "data.huaweicloud_waf_certificates.expiration_status_filter"
		dcByExpirationStatus = acceptance.InitDataSourceCheck(byExpirationStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafCertificates_basic(dataSourceName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.created_at"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByHost.CheckResourceExists(),
					resource.TestCheckOutput("host_filter_is_useful", "true"),

					dcByExpirationStatus.CheckResourceExists(),
					resource.TestCheckOutput("expiration_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceWafCertificates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_certificates" "test" {
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_certificate.test
  ]
}

# Filter by name
locals {
  name = data.huaweicloud_waf_certificates.test.certificates.0.name
}

data "huaweicloud_waf_certificates" "name_filter" {
  enterprise_project_id = "%[2]s"
  name                  = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_certificates.name_filter.certificates[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by host
data "huaweicloud_waf_certificates" "host_filter" {
  enterprise_project_id = "%[2]s"
  host                  = true
}

output "host_filter_is_useful" {
  value = length(data.huaweicloud_waf_certificates.host_filter.certificates.0.bind_host) > 0
}

# Filter by expiration_status
locals {
  expiration_status = data.huaweicloud_waf_certificates.test.certificates.0.expiration_status
}

data "huaweicloud_waf_certificates" "expiration_status_filter" {
  enterprise_project_id = "%[2]s"
  expiration_status     = local.expiration_status
}

locals {
  expiration_status_filter_result = [
    for v in data.huaweicloud_waf_certificates.expiration_status_filter.certificates[*].expiration_status : v == local.expiration_status
  ]
}

output "expiration_status_filter_is_useful" {
  value = length(local.expiration_status_filter_result) > 0 && alltrue(local.expiration_status_filter_result)
}
`, testAccWafCertificate_basic(name, generateCertificateBody()), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
