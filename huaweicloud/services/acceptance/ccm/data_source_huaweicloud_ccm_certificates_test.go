package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCertificates_basic(t *testing.T) {
	var (
		datasource = "data.huaweicloud_ccm_certificates.test"
		dc         = acceptance.InitDataSourceCheck(datasource)

		byStatus   = "data.huaweicloud_ccm_certificates.test"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byEpsID   = "data.huaweicloud_ccm_certificates.test"
		dcByEpsID = acceptance.InitDataSourceCheck(byEpsID)

		byDeploySupport   = "data.huaweicloud_ccm_certificates.test"
		dcByDeploySupport = acceptance.InitDataSourceCheck(byDeploySupport)

		notFound   = "data.huaweicloud_ccm_certificates.test"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMCertificateName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.domain"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.expire_time"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.status"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.domain_count"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.wildcard_count"),
					resource.TestCheckResourceAttrSet(datasource, "certificates.0.enterprise_project_id"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEpsID.CheckResourceExists(),
					resource.TestCheckOutput("epsID_filter_is_useful", "true"),

					dcByDeploySupport.CheckResourceExists(),
					resource.TestCheckOutput("deploy_support_filter_is_useful", "true"),

					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("not_found_filter_is_empty", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCertificates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ccm_certificates" "test" {
  name = "%s"
}

# Search results by status.
locals {
  status = data.huaweicloud_ccm_certificates.test.certificates[0].status
}

data "huaweicloud_ccm_certificates" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ccm_certificates.filter_by_status.certificates[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Search results by enterprise project ID.
locals {
  epsID = data.huaweicloud_ccm_certificates.test.certificates[0].enterprise_project_id
}

data "huaweicloud_ccm_certificates" "filter_by_epsID" {
  enterprise_project_id = local.epsID
}

locals {
  epsID_filter_result = [
    for v in data.huaweicloud_ccm_certificates.filter_by_epsID.certificates[*].enterprise_project_id : v == local.epsID
  ]
}

output "epsID_filter_is_useful" {
  value = alltrue(local.epsID_filter_result) && length(local.epsID_filter_result) > 0
}

# Search results by deploy support.
locals {
  deploy_support = data.huaweicloud_ccm_certificates.test.certificates[0].deploy_support
}

data "huaweicloud_ccm_certificates" "filter_by_deploy_support" {
  deploy_support = local.deploy_support
}

locals {
  deploy_support_filter_result = [
    for v in data.huaweicloud_ccm_certificates.filter_by_deploy_support.certificates[*].deploy_support : v == local.deploy_support
  ]
}

output "deploy_support_filter_is_useful" {
  value = alltrue(local.deploy_support_filter_result) && length(local.deploy_support_filter_result) > 0
}

# Not found by name.
data "huaweicloud_ccm_certificates" "not_found" {
  name = "not_found"
}

output "not_found_filter_is_empty" {
  value = length(data.huaweicloud_ccm_certificates.not_found.certificates) == 0
}
`, acceptance.HW_CCM_CERTIFICATE_NAME)
}
