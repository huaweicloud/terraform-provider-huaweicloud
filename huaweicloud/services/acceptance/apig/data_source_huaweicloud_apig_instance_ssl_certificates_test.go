package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceSSLCertificates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_instance_ssl_certificates.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byInstanceId   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byName   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byCommonName   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_common_name"
		dcByCommonName = acceptance.InitDataSourceCheck(byCommonName)

		bySignatureAlgorithm   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_signature_algorithm"
		dcBySignatureAlgorithm = acceptance.InitDataSourceCheck(bySignatureAlgorithm)

		byType   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byAlgorithmType   = "data.huaweicloud_apig_instance_ssl_certificates.filter_by_algorithm_type"
		dcByAlgorithmType = acceptance.InitDataSourceCheck(byAlgorithmType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceSSLCertificates_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "certificates.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.common_name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.algorithm_type"),
					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("certificate_name_filter_is_useful", "true"),
					dcByCommonName.CheckResourceExists(),
					resource.TestCheckOutput("certificate_common_name_filter_is_useful", "true"),
					dcBySignatureAlgorithm.CheckResourceExists(),
					resource.TestCheckOutput("certificate_signature_algorithm_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("certificate_type_filter_is_useful", "true"),
					dcByAlgorithmType.CheckResourceExists(),
					resource.TestCheckOutput("certificate_algorithm_type_filter_is_useful", "true"),
				),
			},
			{
				Config:      testAccDataSourceInstanceSSLCertificates_step2(),
				ExpectError: regexp.MustCompile("The instance does not exist"),
			},
		},
	})
}

func testAccDataSourceInstanceSSLCertificates_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instance_ssl_certificates" "test" {
  instance_id = "%[1]s"
}

# Filter by Instance ID
locals {
  instance_id = "%[1]s"
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_instance_id" {
  instance_id = local.instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_instance_id.certificates[*].instance_id : v == local.instance_id || v == "common"
  ]
}

output "instance_id_filter_is_useful" {
  value = length(local.instance_id_filter_result) > 0 && alltrue(local.instance_id_filter_result)
}

# Filter by Name
locals {
  name = data.huaweicloud_apig_instance_ssl_certificates.test.certificates[0].name
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_name" {
  instance_id = local.instance_id
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_name.certificates[*].name : v == local.name
  ]
}

output "certificate_name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by Common Name
locals {
  common_name = data.huaweicloud_apig_instance_ssl_certificates.test.certificates[0].common_name
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_common_name" {
  instance_id = local.instance_id
  common_name = local.common_name
}

locals {
  common_name_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_common_name.certificates[*].common_name : v == local.common_name
  ]
}

output "certificate_common_name_filter_is_useful" {
  value = length(local.common_name_filter_result) > 0 && alltrue(local.common_name_filter_result)
}

# Filter by Signature Algorithm
locals {
  signature_algorithm = data.huaweicloud_apig_instance_ssl_certificates.test.certificates[0].signature_algorithm
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_signature_algorithm" {
  instance_id = local.instance_id
  signature_algorithm = local.signature_algorithm
}

locals {
  signature_algorithm_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_signature_algorithm.certificates[*].signature_algorithm : v == local.signature_algorithm
  ]
}

output "certificate_signature_algorithm_filter_is_useful" {
  value = length(local.signature_algorithm_filter_result) > 0 && alltrue(local.signature_algorithm_filter_result)
}

# Filter by Type
locals {
  type = data.huaweicloud_apig_instance_ssl_certificates.test.certificates[0].type
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_type" {
  instance_id = local.instance_id
  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_type.certificates[*].type : v == local.type
  ]
}

output "certificate_type_filter_is_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by Algorithm Type
locals {
  algorithm_type = data.huaweicloud_apig_instance_ssl_certificates.test.certificates[0].algorithm_type
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_algorithm_type" {
  instance_id = local.instance_id
  algorithm_type = local.algorithm_type
}

locals {
  algorithm_type_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_algorithm_type[*].algorithm_type : v == local.algorithm_type
  ]
}

output "certificate_algorithm_type_filter_is_useful" {
  value = length(local.algorithm_type_result) > 0 && alltrue(local.algorithm_type_result)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}

func testAccDataSourceInstanceSSLCertificates_step2() string {
	randomUUID, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
# Filter by Invalid Instance ID
locals {
  instance_id = "%[1]s"
}

data "huaweicloud_apig_instance_ssl_certificates" "filter_by_invalid_instance_id" {
  instance_id = local.instance_id
}

locals {
  invalid_instance_id_filter_result = [
    for v in data.huaweicloud_apig_instance_ssl_certificates.filter_by_invalid_instance_id.certificates[*].instance_id : v == local.instance_id || v == "common"
  ]
}

output "invalid_instance_id_filter_is_useful" {
  value = length(local.invalid_instance_id_filter_result) == 0
}
`, randomUUID)
}
