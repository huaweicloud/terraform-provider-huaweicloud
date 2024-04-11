package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFloatingPtrrecords_basic(t *testing.T) {
	var (
		domainName   = fmt.Sprintf("acpttest-ptr-%s.com.", acctest.RandString(5))
		byRecordId   = "data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id"
		dcByRecordId = acceptance.InitDataSourceCheck(byRecordId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFloatingPtrrecords_basic(domainName),
				Check: resource.ComposeTestCheckFunc(
					dcByRecordId.CheckResourceExists(),
					resource.TestCheckResourceAttr(byRecordId, "ptrrecords.0.ttl", "300"),
					resource.TestCheckResourceAttr(byRecordId, "ptrrecords.0.description", "Created by terraform"),

					resource.TestCheckOutput("is_record_id_filter_useful", "true"),
					resource.TestCheckOutput("is_public_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_domain_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_domain_name", "true"),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					resource.TestCheckResourceAttr("data.huaweicloud_dns_floating_ptrrecords.filter_by_tags", "ptrrecords.#", "1"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceFloatingPtrrecords_basic(domainName string) string {
	return fmt.Sprintf(`
%s

locals {
  tags = {
    foo       = "bar"
    terraform = "ptrrecord"
  }
}

resource "huaweicloud_dns_ptrrecord" "test" {
  name          = "%s"
  description   = "Created by terraform"
  floatingip_id = huaweicloud_vpc_eip.eip_1.id
  ttl           = 300
  tags          = local.tags
}

locals {
  record_id = huaweicloud_dns_ptrrecord.test.id
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_record_id" {
  record_id = local.record_id
}

output "is_record_id_filter_useful" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id.ptrrecords) > 0 && alltrue(
    [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_record_id.ptrrecords[*].id : v == local.record_id]
  )
}

locals {
  public_ip = huaweicloud_vpc_eip.eip_1.address
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_public_id" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  public_ip  = local.public_ip
}

output "is_public_ip_filter_useful" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_public_id.ptrrecords) > 0 && alltrue(
    [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_public_id.ptrrecords[*].public_ip : v == local.public_ip]
  )
}

locals {
  domain_name = huaweicloud_dns_ptrrecord.test.name
}

data "huaweicloud_dns_floating_ptrrecords" "by_domain_name_filter" {
  depends_on  = [huaweicloud_dns_ptrrecord.test]
  domain_name = local.domain_name
}

output "is_domain_name_filter_useful" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.by_domain_name_filter.ptrrecords) > 0 && alltrue(
    [for v in data.huaweicloud_dns_floating_ptrrecords.by_domain_name_filter.ptrrecords[*].domain_name : v == local.domain_name]
  )
}

data "huaweicloud_dns_floating_ptrrecords" "not_found_domain_name" {
  depends_on  = [huaweicloud_dns_ptrrecord.test]
  domain_name = "not-found.com."
}

output "not_found_domain_name" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.not_found_domain_name.ptrrecords) == 0
}

locals {
  enterprise_project_id = huaweicloud_dns_ptrrecord.test.enterprise_project_id
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_eps_id" {
  enterprise_project_id = local.enterprise_project_id
}

output "is_eps_id_filter_useful" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_eps_id.ptrrecords) > 0 && alltrue(
    [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_eps_id.ptrrecords[*].enterprise_project_id : v == local.enterprise_project_id]
  )
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_tags" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  tags       = local.tags
}

data "huaweicloud_dns_floating_ptrrecords" "filter_by_status" {
  depends_on = [huaweicloud_dns_ptrrecord.test]
  status     = "ACTIVE"
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_dns_floating_ptrrecords.filter_by_status.ptrrecords) > 0 && alltrue(
    [for v in data.huaweicloud_dns_floating_ptrrecords.filter_by_status.ptrrecords[*].status : v == "ACTIVE"]
  )
}
`, testAccDNSPtrRecord_base(), domainName)
}
