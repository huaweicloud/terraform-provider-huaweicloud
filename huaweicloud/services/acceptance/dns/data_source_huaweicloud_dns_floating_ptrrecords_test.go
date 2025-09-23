package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataFloatingPtrRecords_basic(t *testing.T) {
	var (
		domainName = fmt.Sprintf("acpttest-ptr-%s.com.", acctest.RandString(5))
		rName      = "huaweicloud_dns_ptrrecord.test"

		all                = "data.huaweicloud_dns_floating_ptrrecords.test"
		dcForAllPtrRecords = acceptance.InitDataSourceCheck(all)

		byRecordId      = "data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id"
		dcByRecordId    = acceptance.InitDataSourceCheck(byRecordId)
		byNotRecordId   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_record_id"
		dcByNotRecordId = acceptance.InitDataSourceCheck(byNotRecordId)

		byDomainName           = "data.huaweicloud_dns_floating_ptrrecords.filter_by_domain_name"
		dcByDomainName         = acceptance.InitDataSourceCheck(byDomainName)
		byNotFoundDomainName   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_domain_name"
		dcByNotFoundDomainName = acceptance.InitDataSourceCheck(byNotFoundDomainName)

		byPublicIp      = "data.huaweicloud_dns_floating_ptrrecords.filter_by_public_ip"
		dcByPublicIp    = acceptance.InitDataSourceCheck(byPublicIp)
		byNotPublicIp   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_public_ip"
		dcByNotPublicIp = acceptance.InitDataSourceCheck(byNotPublicIp)

		byEpsId                   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_eps_id"
		dcByEpsId                 = acceptance.InitDataSourceCheck(byEpsId)
		byNotfoundEpsId           = "data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_eps_id"
		dcByNotfoundNotfoundEpsId = acceptance.InitDataSourceCheck(byNotfoundEpsId)

		byTags                   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_tags"
		dcByTags                 = acceptance.InitDataSourceCheck(byTags)
		byNotfoundTags           = "data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_tags"
		dcByNotfoundNotfoundTags = acceptance.InitDataSourceCheck(byNotfoundTags)

		byStatus   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataFloatingPtrRecords_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dcForAllPtrRecords.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "ptrrecords.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by PTR record ID.
					dcByRecordId.CheckResourceExists(),
					resource.TestCheckOutput("is_record_id_filter_useful", "true"),
					dcByNotRecordId.CheckResourceExists(),
					resource.TestCheckOutput("record_id_not_found_validation_pass", "true"),
					// Filter by PTR record name.
					dcByDomainName.CheckResourceExists(),
					resource.TestCheckOutput("is_domain_name_filter_useful", "true"),
					dcByNotFoundDomainName.CheckResourceExists(),
					resource.TestCheckOutput("domain_name_not_found_validation_pass", "true"),
					// Filter by PTR record public IP.
					dcByPublicIp.CheckResourceExists(),
					resource.TestCheckOutput("is_public_ip_filter_useful", "true"),
					dcByNotPublicIp.CheckResourceExists(),
					resource.TestCheckOutput("public_ip_not_found_validation_pass", "true"),
					// Filter by PTR record enterprise project ID.
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcByNotfoundNotfoundEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_not_found_validation_pass", "true"),
					// Filter by PTR record tags.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNotfoundNotfoundTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_not_found_validation_pass", "true"),
					// Filter by PTR record status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrPair(byRecordId, "ptrrecords.0.ttl", rName, "ttl"),
					resource.TestCheckResourceAttrPair(byRecordId, "ptrrecords.0.description", rName, "description"),
				),
			},
		},
	})
}

func testAccDataFloatingPtrRecords_basic(domainName string) string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dns_ptrrecord" "test" {
  name                  = "%[2]s"
  description           = "Created by terraform"
  floatingip_id         = huaweicloud_vpc_eip.test.id
  ttl                   = 300
  enterprise_project_id = "%[3]s"

  tags = {
    foo       = "bar"
    terraform = "ptrrecord"
  }
}

data "huaweicloud_dns_floating_ptrrecords" "test" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
}

# Filter by PTR record ID.
locals {
  record_id = huaweicloud_dns_ptrrecord.test.id
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_record_id" {
  record_id = local.record_id
}

locals {
  record_id_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id.ptrrecords[*].id : v == local.record_id]
}

output "is_record_id_filter_useful" {
  value = length(local.record_id_filter_result) > 0 && alltrue(local.record_id_filter_result)
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_not_found_record_id" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  record_id  = "%[4]s"
}

output "record_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_record_id.ptrrecords) == 0
}

# Filter by PTR record name.
locals {
  domain_name = huaweicloud_dns_ptrrecord.test.name
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_domain_name" {
  depends_on  = [huaweicloud_dns_ptrrecord.test]
  domain_name = local.domain_name
}

locals {
  domain_name_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_domain_name.ptrrecords[*].domain_name :
  v == local.domain_name]
}

output "is_domain_name_filter_useful" {
  value = length(local.domain_name_filter_result) > 0 && alltrue(local.domain_name_filter_result)
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_not_found_domain_name" {
  depends_on  = [huaweicloud_dns_ptrrecord.test]
  domain_name = "not_found_ptrrecord_name"
}

output "domain_name_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_domain_name.ptrrecords) == 0
}

# Filter by PTR record public IP.
locals {
  public_ip = huaweicloud_dns_ptrrecord.test.address
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_public_ip" {
  public_ip = local.public_ip
}

locals {
  public_ip_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_public_ip.ptrrecords[*].public_ip : v == local.public_ip]
}

output "is_public_ip_filter_useful" {
  value = length(local.public_ip_filter_result) > 0 && alltrue(local.public_ip_filter_result)
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_not_found_public_ip" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  public_ip  = "not_found_public_ip"
}

output "public_ip_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_public_ip.ptrrecords) == 0
}

# Filter by PTR record enterprise project ID.
locals {
  enterprise_project_id = huaweicloud_dns_ptrrecord.test.enterprise_project_id
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_eps_id" {
  depends_on = [huaweicloud_dns_ptrrecord.test]

  enterprise_project_id = local.enterprise_project_id
}

locals {
  eps_id_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_eps_id.ptrrecords[*].enterprise_project_id :
  v == local.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_not_found_eps_id" {
  depends_on = [huaweicloud_dns_ptrrecord.test]

  enterprise_project_id = "%[4]s"
}

output "is_eps_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_eps_id.ptrrecords) == 0
}

# Filter by PTR record tags.
locals {
  tags = huaweicloud_dns_ptrrecord.test.tags
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_tags" {
  depends_on = [huaweicloud_dns_ptrrecord.test]

  tags = local.tags
}

locals {
  tags_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_tags.ptrrecords[*].tags : v == local.tags]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_not_found_tags" {
  depends_on = [huaweicloud_dns_ptrrecord.test]

  tags = {
    key = "not_found_value"
  }
}

output "tags_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_not_found_tags.ptrrecords) == 0
}

# Filter by PTR record status.
locals {
  status = try(data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id.ptrrecords[0].status, "")
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_status" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  status     = local.status
}

locals {
  status_filter_result = [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_status.ptrrecords[*].status : v == local.status]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testAccPtrRecord_base(), domainName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randomId)
}
