package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafAllDomains_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_all_domains.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byHostName   = "data.huaweicloud_waf_all_domains.filter_by_hostname"
		dcByHostName = acceptance.InitDataSourceCheck(byHostName)

		byProtectStatus   = "data.huaweicloud_waf_all_domains.filter_by_protect_status"
		dcByProtectStatus = acceptance.InitDataSourceCheck(byProtectStatus)

		byWafType   = "data.huaweicloud_waf_all_domains.filter_by_waf_type"
		dcByWafType = acceptance.InitDataSourceCheck(byWafType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Prepare a WAF domain before test
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWafAllDomains_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.policyid"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.hostname"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.proxy"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.waf_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.flag.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.flag.0.pci_3ds"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.flag.0.ipv6"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.server.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.server.0.front_protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.server.0.back_protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.server.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.server.0.type"),

					dcByHostName.CheckResourceExists(),
					resource.TestCheckOutput("hostname_filter_is_useful", "true"),

					dcByProtectStatus.CheckResourceExists(),
					resource.TestCheckOutput("protect_status_filter_is_useful", "true"),

					dcByWafType.CheckResourceExists(),
					resource.TestCheckOutput("waf_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceWafAllDomains_basic = `

data "huaweicloud_waf_all_domains" "test" {}

locals {
  hostname = data.huaweicloud_waf_all_domains.test.items[0].hostname
}

data "huaweicloud_waf_all_domains" "filter_by_hostname" {
  hostname = local.hostname
}

locals {
  hostname_filter_result = [
    for v in data.huaweicloud_waf_all_domains.filter_by_hostname.items[*].hostname : v == local.hostname
  ]
}

output "hostname_filter_is_useful" {
  value = alltrue(local.hostname_filter_result) && length(local.hostname_filter_result) > 0
}

locals {
  protect_status = data.huaweicloud_waf_all_domains.test.items[0].protect_status
}

data "huaweicloud_waf_all_domains" "filter_by_protect_status" {
  protect_status = local.protect_status
}

locals {
  protect_status_filter_result = [ 
    for v in data.huaweicloud_waf_all_domains.filter_by_protect_status.items[*].protect_status : v == local.protect_status
  ]
}

output "protect_status_filter_is_useful" {
  value = alltrue(local.protect_status_filter_result) && length(local.protect_status_filter_result) > 0
}

locals {
  waf_type = data.huaweicloud_waf_all_domains.test.items[0].waf_type
}

data "huaweicloud_waf_all_domains" "filter_by_waf_type" {
  waf_type = local.waf_type
}

locals {
  waf_type_filter_result = [
    for v in data.huaweicloud_waf_all_domains.filter_by_waf_type.items[*].waf_type : v == local.waf_type
  ]
}

output "waf_type_filter_is_useful" {
  value = alltrue(local.waf_type_filter_result) && length(local.waf_type_filter_result) > 0
}
`
