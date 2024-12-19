package live

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRecordCallbacks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_live_record_callbacks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byDomainName   = "data.huaweicloud_live_record_callbacks.filter_by_domain_name"
		dcByDomainName = acceptance.InitDataSourceCheck(byDomainName)

		byAppName   = "data.huaweicloud_live_record_callbacks.filter_by_app_name"
		dcByAppName = acceptance.InitDataSourceCheck(byAppName)

		allByAppName   = "data.huaweicloud_live_record_callbacks.filter_all_by_app_name"
		dcAllByAppName = acceptance.InitDataSourceCheck(allByAppName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLiveIngestDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRecordCallbacks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.sign_type"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "callbacks.0.url"),

					dcByDomainName.CheckResourceExists(),
					resource.TestCheckOutput("domain_name_filter_is_useful", "true"),

					dcByAppName.CheckResourceExists(),
					resource.TestCheckOutput("app_name_filter_is_useful", "true"),

					dcAllByAppName.CheckResourceExists(),
					resource.TestCheckOutput("all_app_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRecordCallbacks_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_live_record_callbacks" "test" {
  depends_on = [huaweicloud_live_record_callback.test]
}

# Filter by domain_name
locals {
  domain_name = data.huaweicloud_live_record_callbacks.test.callbacks[0].domain_name
}

data "huaweicloud_live_record_callbacks" "filter_by_domain_name" {
  depends_on = [huaweicloud_live_record_callback.test]

  domain_name = local.domain_name
}

locals {
  domain_name_filter_result = [
    for v in data.huaweicloud_live_record_callbacks.filter_by_domain_name.callbacks[*].domain_name : v == local.domain_name
  ]
}

output "domain_name_filter_is_useful" {
  value = alltrue(local.domain_name_filter_result) && length(local.domain_name_filter_result) > 0
}

# Filter by app_name
locals {
  app_name = data.huaweicloud_live_record_callbacks.test.callbacks[0].app_name
}

data "huaweicloud_live_record_callbacks" "filter_by_app_name" {
  depends_on = [huaweicloud_live_record_callback.test]

  app_name = local.app_name
}

locals {
  app_name_filter_result = [
    for v in data.huaweicloud_live_record_callbacks.filter_by_app_name.callbacks[*].app_name : v == local.app_name || v == "*"
  ]
}

output "app_name_filter_is_useful" {
  value = alltrue(local.app_name_filter_result) && length(local.app_name_filter_result) > 0
}

# Filter all by app_name (*)
data "huaweicloud_live_record_callbacks" "filter_all_by_app_name" {
  depends_on = [huaweicloud_live_record_callback.test]

  app_name = "*"
}

output "all_app_name_filter_is_useful" {
  value = length(data.huaweicloud_live_record_callbacks.filter_all_by_app_name.callbacks) > 0
}
`, testCallBack_basic())
}
