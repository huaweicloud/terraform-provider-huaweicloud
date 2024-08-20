package eg

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEventChannels_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_eg_event_channels.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byProviderType   = "data.huaweicloud_eg_event_channels.filter_by_provider_type"
		dcByProviderType = acceptance.InitDataSourceCheck(byProviderType)

		byName   = "data.huaweicloud_eg_event_channels.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_eg_event_channels.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byEpsId   = "data.huaweicloud_eg_event_channels.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEventChannels_basic_step1(name),
			},
			{
				// Make sure the attribute 'updated_at' has been configured.
				Config: testAccDataEventChannels_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "channels.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_query_result_contains_at_least_two_types", "true"),
					dcByProviderType.CheckResourceExists(),
					resource.TestCheckOutput("is_provider_type_filter_useful", "true"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byName, "channels.0.id", "huaweicloud_eg_custom_event_channel.test", "id"),
					resource.TestCheckResourceAttr(byName, "channels.0.name", name),
					resource.TestCheckResourceAttr(byName, "channels.0.description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(byName, "channels.0.cross_account_ids.#", "2"),
					resource.TestCheckResourceAttr(byName, "channels.0.enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(byName, "channels.0.provider_type", "CUSTOM"),
					resource.TestCheckResourceAttrSet(byName, "channels.0.created_at"),
					resource.TestCheckResourceAttrSet(byName, "channels.0.updated_at"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataEventChannels_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = "%[1]s"
  description           = "Created by terraform script"
  enterprise_project_id = "%[2]s"
  cross_account_ids     = ["account1", "account2"]
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataEventChannels_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = "%[1]s"
  description           = "Updated by terraform script"
  enterprise_project_id = "%[2]s"
  cross_account_ids     = ["account1", "account2"]
}

data "huaweicloud_eg_event_channels" "test" {
  depends_on = [huaweicloud_eg_custom_event_channel.test]
}

// At least one of official event channel exist, e.g. official channel names Default.
output "is_query_result_contains_at_least_two_types" {
  value = length(distinct(data.huaweicloud_eg_event_channels.test.channels[*].provider_type)) >= 2
}

# Filter by provider type
data "huaweicloud_eg_event_channels" "filter_by_provider_type" {
  depends_on = [huaweicloud_eg_custom_event_channel.test]

  provider_type = "CUSTOM"
}

locals {
  provider_type_filter_result = [
    for v in data.huaweicloud_eg_event_channels.filter_by_provider_type.channels[*].provider_type : v == "CUSTOM"
  ]
}

output "is_provider_type_filter_useful" {
  value = length(local.provider_type_filter_result) > 0 && alltrue(local.provider_type_filter_result)
}

# Filter by name
data "huaweicloud_eg_event_channels" "filter_by_name" {
  depends_on = [huaweicloud_eg_custom_event_channel.test]

  name = "%[1]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_eg_event_channels.filter_by_name.channels[*].name : v == "%[1]s"
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not found name
data "huaweicloud_eg_event_channels" "filter_by_not_found_name" {
  name = "not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_eg_event_channels.filter_by_not_found_name.channels) < 1
}

# Filter by enterprise project ID
data "huaweicloud_eg_event_channels" "filter_by_eps_id" {
  depends_on = [huaweicloud_eg_custom_event_channel.test]

  enterprise_project_id = "%[2]s"
}

locals {
  eps_id_filter_result = [
    for v in data.huaweicloud_eg_event_channels.filter_by_name.channels[*].enterprise_project_id : v == "%[2]s"
  ]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
