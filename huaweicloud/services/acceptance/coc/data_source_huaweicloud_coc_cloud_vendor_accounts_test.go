package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocCloudVendorAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_cloud_vendor_accounts.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCocCloudVendorAccounts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.vendor"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.ak"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.failure_msg"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sync_date"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckOutput("vendor_filter_is_useful", "true"),
					resource.TestCheckOutput("account_id_filter_is_useful", "true"),
					resource.TestCheckOutput("account_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCocCloudVendorAccounts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_coc_cloud_vendor_accounts" "test" {
  depends_on = [huaweicloud_coc_cloud_vendor_account.test]
}

locals {
  vendor = huaweicloud_coc_cloud_vendor_account.test.vendor
}
data "huaweicloud_coc_cloud_vendor_accounts" "vendor_filter" {
  depends_on = [huaweicloud_coc_cloud_vendor_account.test]

  vendor = huaweicloud_coc_cloud_vendor_account.test.vendor
}
output "vendor_filter_is_useful" {
  value = length(data.huaweicloud_coc_cloud_vendor_accounts.vendor_filter.data) > 0 && alltrue(
  [for v in data.huaweicloud_coc_cloud_vendor_accounts.vendor_filter.data[*].vendor : v == local.vendor]
  )
}

locals {
  account_id = huaweicloud_coc_cloud_vendor_account.test.account_id
}
data "huaweicloud_coc_cloud_vendor_accounts" "account_id_filter" {
  account_id = huaweicloud_coc_cloud_vendor_account.test.account_id
}
output "account_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_cloud_vendor_accounts.account_id_filter.data) > 0 && alltrue(
  [for v in data.huaweicloud_coc_cloud_vendor_accounts.account_id_filter.data[*].account_id : v == local.account_id]
  )
}

locals {
  account_name = huaweicloud_coc_cloud_vendor_account.test.account_name
}
data "huaweicloud_coc_cloud_vendor_accounts" "account_name_filter" {
  depends_on = [huaweicloud_coc_cloud_vendor_account.test]

  account_name = huaweicloud_coc_cloud_vendor_account.test.account_name
}
output "account_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_cloud_vendor_accounts.account_name_filter.data) > 0 && alltrue(
  [for v in data.huaweicloud_coc_cloud_vendor_accounts.account_name_filter.data[*].account_name : v == local.account_name]
  )
}
`, testAccCloudVendorAccount_basic(name))
}
