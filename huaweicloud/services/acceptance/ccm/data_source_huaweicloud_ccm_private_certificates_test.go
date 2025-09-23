package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateCertificates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ccm_private_certificates.test"
		rName      = acceptance.RandomAccResourceNameWithDash()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_ccm_private_certificates.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_ccm_private_certificates.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		bySort   = "data.huaweicloud_ccm_private_certificates.filter_by_sort"
		dcBySort = acceptance.InitDataSourceCheck(bySort)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePrivateCertificates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.distinguished_name.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.issuer_id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.key_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.gen_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.issuer_name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.start_at"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.expired_at"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.enterprise_project_id"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcBySort.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePrivateCertificates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_certificates" "test" {
  depends_on = [
    huaweicloud_ccm_private_certificate.test
  ]
}

# Fuzzy search based on name, so only the length of the result is verified.
locals {
  name = data.huaweicloud_ccm_private_certificates.test.certificates[0].distinguished_name[0].common_name
}

data "huaweicloud_ccm_private_certificates" "filter_by_name" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_ccm_private_certificates.filter_by_name.certificates) > 0
}

# Search results by status.
locals {
  status = data.huaweicloud_ccm_private_certificates.test.certificates[0].status
}

data "huaweicloud_ccm_private_certificates" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ccm_private_certificates.filter_by_status.certificates[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Search results by sort_key and sort_dir and only test the length of the result.
data "huaweicloud_ccm_private_certificates" "filter_by_sort" {
  depends_on = [
    huaweicloud_ccm_private_certificate.test
  ]

  sort_key = "create_time"
  sort_dir = "ASC"
}

output "sort_filter_is_useful" {
  value = length(data.huaweicloud_ccm_private_certificates.filter_by_sort.certificates) > 0
}
`, testPrivateCertificate_basic(name))
}
