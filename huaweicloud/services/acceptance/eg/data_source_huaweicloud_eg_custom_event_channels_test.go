package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCustomEventChannels_basic(t *testing.T) {
	var (
		baseRes      = "huaweicloud_eg_custom_event_channel.test"
		byName       = "data.huaweicloud_eg_custom_event_channels.filter_by_name"
		nameNotFound = "data.huaweicloud_eg_custom_event_channels.name_not_found"

		obj            interface{}
		rc             = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventChannelFunc)
		dcByName       = acceptance.InitDataSourceCheck(byName)
		dcNameNotFound = acceptance.InitDataSourceCheck(nameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventChannels_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventChannels_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}
`, name)
}

func testAccDataCustomEventChannels_basic() string {
	return fmt.Sprintf(`
%[1]s
data "huaweicloud_eg_custom_event_channels" "filter_by_name" {
  // The behavior of parameter 'name' of the resource is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  name = huaweicloud_eg_custom_event_channel.test.name
}

data "huaweicloud_eg_custom_event_channels" "name_not_found" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  name = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_channels.filter_by_name.channels[*].id :
                   v == huaweicloud_eg_custom_event_channel.test.id]
}

output "is_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_channels.name_not_found.channels) == 0
}
`, testAccDataCustomEventChannels_base())
}

func TestAccDataCustomEventChannels_filterById(t *testing.T) {
	var (
		baseRes    = "huaweicloud_eg_custom_event_channel.test"
		byId       = "data.huaweicloud_eg_custom_event_channels.filter_by_id"
		idNotFound = "data.huaweicloud_eg_custom_event_channels.id_not_found"

		obj          interface{}
		rc           = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventChannelFunc)
		dcById       = acceptance.InitDataSourceCheck(byId)
		dcIdNotFound = acceptance.InitDataSourceCheck(idNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventChannels_filterById(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					dcIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("id_not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventChannels_filterById() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_channels" "filter_by_id" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id
}

data "huaweicloud_eg_custom_event_channels" "id_not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  channel_id = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_channels.filter_by_id.channels[*].id :
                   v == huaweicloud_eg_custom_event_channel.test.id]
}

output "is_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "id_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_channels.id_not_found.channels) == 0
}
`, testAccDataCustomEventChannels_base(), randUUID)
}

func TestAccDataCustomEventChannels_filterByEpsId(t *testing.T) {
	var (
		baseRes           = "huaweicloud_eg_custom_event_channel.test"
		byChannelId       = "data.huaweicloud_eg_custom_event_channels.filter_by_eps_id"
		channelIdNotFound = "data.huaweicloud_eg_custom_event_channels.eps_id_not_found"

		obj             interface{}
		rc              = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventChannelFunc)
		dcByEpsId       = acceptance.InitDataSourceCheck(byChannelId)
		dcEpsIdNotFound = acceptance.InitDataSourceCheck(channelIdNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventChannels_filterByEpsId(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					dcEpsIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("eps_id_not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventChannels_base_withEpsId() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataCustomEventChannels_filterByEpsId() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_channels" "filter_by_eps_id" {
  // The behavior of parameter 'channel_id' of the resource is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  enterprise_project_id = "%[2]s"
}

data "huaweicloud_eg_custom_event_channels" "eps_id_not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  enterprise_project_id = "%[3]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_channels.filter_by_eps_id.channels[*].enterprise_project_id : v == "%[2]s"]
}

output "is_eps_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "eps_id_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_channels.eps_id_not_found.channels) == 0
}
`, testAccDataCustomEventChannels_base_withEpsId(), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, randUUID)
}

func TestAccDataCustomEventChannels_fuzzyName(t *testing.T) {
	var (
		baseRes = "huaweicloud_eg_custom_event_channel.test"
		byFuzzy = "data.huaweicloud_eg_custom_event_channels.filter_by_fuzzy_name"

		obj     interface{}
		rc      = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventChannelFunc)
		dcFuzzy = acceptance.InitDataSourceCheck(byFuzzy)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventChannels_fuzzyName(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventChannels_fuzzyName() string {
	name := acceptance.RandomAccResourceName()
	fuzzyName := name[:len(name)-3]

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

data "huaweicloud_eg_custom_event_channels" "filter_by_fuzzy_name" {
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  fuzzy_name = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_channels.filter_by_fuzzy_name.channels[*].name :
                   can(regex(".*%[2]s.*", v))]
}

output "is_fuzzy_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}
`, name, fuzzyName)
}

func TestAccDataCustomEventChannels_sort(t *testing.T) {
	var (
		baseRes = "huaweicloud_eg_custom_event_channel.test"
		bySort  = "data.huaweicloud_eg_custom_event_channels.sort_by_name"

		obj    interface{}
		rc     = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventChannelFunc)
		dcSort = acceptance.InitDataSourceCheck(bySort)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventChannels_sort(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcSort.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_useful", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventChannels_sort() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

data "huaweicloud_eg_custom_event_channels" "sort_by_name" {
  depends_on = [
    huaweicloud_eg_custom_event_channel.test,
  ]

  sort = "created_time:DESC"
}

output "is_sort_useful" {
  value = length(data.huaweicloud_eg_custom_event_channels.sort_by_name.channels) > 0
}
`, name)
}
