package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePrivateCas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ccm_private_cas.test"
		rName      = acceptance.RandomAccResourceNameWithDash()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_ccm_private_cas.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byStatus   = "data.huaweicloud_ccm_private_cas.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byType   = "data.huaweicloud_ccm_private_cas.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		bySortAsc   = "data.huaweicloud_ccm_private_cas.filter_by_sort_asc"
		dcBySortAsc = acceptance.InitDataSourceCheck(bySortAsc)

		bySortDesc   = "data.huaweicloud_ccm_private_cas.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourcePrivateCas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cas.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.distinguished_name.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.crl_configuration.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.expired_at"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.path_length"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.issuer_id"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.issuer_name"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.key_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.signature_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.gen_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.free_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "cas.0.charging_mode"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),

					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourcePrivateCas_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_private_cas" "test" {
  depends_on = [
    huaweicloud_ccm_private_ca.test_root,
    huaweicloud_ccm_private_ca.test_subordinate,
  ]
}

# Fuzzy search based on name, so only verify the length of the result.
locals {
  name = data.huaweicloud_ccm_private_cas.test.cas[0].distinguished_name[0].common_name
}

data "huaweicloud_ccm_private_cas" "filter_by_name" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_ccm_private_cas.filter_by_name.cas) > 0
}

# Search results by status.
locals {
  status = data.huaweicloud_ccm_private_cas.test.cas[0].status
}

data "huaweicloud_ccm_private_cas" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ccm_private_cas.filter_by_status.cas[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Search results by type.
locals {
  type = data.huaweicloud_ccm_private_cas.test.cas[0].type
}

data "huaweicloud_ccm_private_cas" "filter_by_type" {
  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_ccm_private_cas.filter_by_type.cas[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

# Search results by sort.
data "huaweicloud_ccm_private_cas" "filter_by_sort_asc" {
  depends_on = [
    huaweicloud_ccm_private_ca.test_root,
    huaweicloud_ccm_private_ca.test_subordinate,
  ]
  sort_key = "create_time"
  sort_dir = "ASC"
}

data "huaweicloud_ccm_private_cas" "filter_by_sort_desc" {
  depends_on = [
    huaweicloud_ccm_private_ca.test_root,
    huaweicloud_ccm_private_ca.test_subordinate,
  ]
  sort_key = "create_time"
  sort_dir = "DESC"
}

locals {
  cas_length = length(data.huaweicloud_ccm_private_cas.filter_by_sort_desc.cas)
  asc_first_id = data.huaweicloud_ccm_private_cas.filter_by_sort_asc.cas[0].id
  desc_last_id = data.huaweicloud_ccm_private_cas.filter_by_sort_desc.cas[local.cas_length - 1].id
}

output "sort_filter_is_useful" {
  value = local.asc_first_id == local.desc_last_id
}
`, tesPrivateCA_postpaid_subordinate(name))
}
