package dsc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, ensure that there is at least two database asset under the HW_DSC_TYPE_ID,
// and that the metadata scan has been completed.
func TestAccDataSourceCatalogInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dsc_catalog_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byInstanceName   = "data.huaweicloud_dsc_catalog_instances.filter_by_instance_name"
		dcByInstanceName = acceptance.InitDataSourceCheck(byInstanceName)

		byAddress   = "data.huaweicloud_dsc_catalog_instances.filter_by_address"
		dcByAddress = acceptance.InitDataSourceCheck(byAddress)

		byUser   = "data.huaweicloud_dsc_catalog_instances.filter_by_user"
		dcByUser = acceptance.InitDataSourceCheck(byUser)

		byColIdAndSortDesc   = "data.huaweicloud_dsc_catalog_instances.filter_by_col_id_and_sort_desc"
		dcByColIdAndSortDesc = acceptance.InitDataSourceCheck(byColIdAndSortDesc)

		byColIdAndSortAsc   = "data.huaweicloud_dsc_catalog_instances.filter_by_col_id_and_sort_asc"
		dcByColIdAndSortAsc = acceptance.InitDataSourceCheck(byColIdAndSortAsc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscTypeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCatalogInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestMatchResourceAttr(all, "instances.0.db_infos.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.db_infos.0.db_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.db_infos.0.db_name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.db_infos.0.asset_id"),
					resource.TestMatchResourceAttr(all, "instances.0.db_infos.0.latest_scan_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by 'instance_name' parameter.
					dcByInstanceName.CheckResourceExists(),
					resource.TestCheckOutput("is_instance_name_filter_useful", "true"),
					// Filter by 'address' parameter.
					dcByAddress.CheckResourceExists(),
					resource.TestCheckOutput("is_address_filter_useful", "true"),
					// Filter by 'user' parameter.
					dcByUser.CheckResourceExists(),
					resource.TestCheckOutput("is_user_filter_useful", "true"),
					// Filter by 'col_id' and 'sort' parameters.
					dcByColIdAndSortDesc.CheckResourceExists(),
					dcByColIdAndSortAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_col_id_and_sort_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCatalogInstances_basic() string {
	return fmt.Sprintf(`
locals {
  type_id = "%[1]s"
}

data "huaweicloud_dsc_catalog_instances" "test" {
  type_id = local.type_id
}

# Filter by 'instance_name' parameter.
locals {
  instance_name = try(data.huaweicloud_dsc_catalog_instances.test.instances[0].instance_name, "NOT FOUND")
}

data "huaweicloud_dsc_catalog_instances" "filter_by_instance_name" {
  type_id       = local.type_id
  instance_name = local.instance_name
}

locals {
  instance_name_filter_result = [for v in data.huaweicloud_dsc_catalog_instances.filter_by_instance_name.instances[*].instance_name :
  strcontains(v, local.instance_name)]
}

output "is_instance_name_filter_useful" {
  value = length(local.instance_name_filter_result) > 0 && alltrue(local.instance_name_filter_result)
}

# Filter by 'address' parameter.
locals {
  address = try(data.huaweicloud_dsc_catalog_instances.test.instances[0].address, "NOT FOUND")
}

data "huaweicloud_dsc_catalog_instances" "filter_by_address" {
  type_id = local.type_id
  address = local.address
}

locals {
  address_filter_result = [for v in data.huaweicloud_dsc_catalog_instances.filter_by_address.instances[*].address : strcontains(v, local.address)]
}

output "is_address_filter_useful" {
  value = length(local.address_filter_result) > 0 && alltrue(local.address_filter_result)
}

# Filter by 'user' parameter.
locals {
  user = try(data.huaweicloud_dsc_catalog_instances.test.instances[0].user, "NOT FOUND")
}

data "huaweicloud_dsc_catalog_instances" "filter_by_user" {
  type_id = local.type_id
  user    = local.user
}

locals {
  user_filter_result = [for v in data.huaweicloud_dsc_catalog_instances.filter_by_user.instances[*].user : strcontains(v, local.user)]
}

output "is_user_filter_useful" {
  value = length(local.user_filter_result) > 0 && alltrue(local.user_filter_result)
}

# Filter by 'col_id' and 'sort' parameters.
data "huaweicloud_dsc_catalog_instances" "filter_by_col_id_and_sort_desc" {
  type_id = local.type_id
  col_id  = "address"
  sort    = "desc"
}

data "huaweicloud_dsc_catalog_instances" "filter_by_col_id_and_sort_asc" {
  type_id = local.type_id
  col_id  = "address"
  sort    = "asc"
}

locals {
  col_id_and_sort_desc_result = data.huaweicloud_dsc_catalog_instances.filter_by_col_id_and_sort_desc.instances[*].address
  col_id_and_sort_asc_result  = data.huaweicloud_dsc_catalog_instances.filter_by_col_id_and_sort_asc.instances[*].address
}

output "is_col_id_and_sort_useful" {
  value = length(local.col_id_and_sort_desc_result) > 1 && local.col_id_and_sort_desc_result == reverse(local.col_id_and_sort_asc_result)
}
`, acceptance.HW_DSC_TYPE_ID)
}
