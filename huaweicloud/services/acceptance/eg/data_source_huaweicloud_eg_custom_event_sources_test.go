package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/eg/v1/source/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCustomEventSources_basic(t *testing.T) {
	var (
		baseRes      = "huaweicloud_eg_custom_event_source.test"
		byName       = "data.huaweicloud_eg_custom_event_sources.filter_by_name"
		nameNotFound = "data.huaweicloud_eg_custom_event_sources.name_not_found"

		obj            custom.Source
		rc             = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventSourceFunc)
		dcByName       = acceptance.InitDataSourceCheck(byName)
		dcNameNotFound = acceptance.InitDataSourceCheck(nameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventSources_basic(),
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

func testAccDataCustomEventSources_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = "%[1]s"
  name       = "%[2]s"
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}

func testAccDataCustomEventSources_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_sources" "filter_by_name" {
  // The behavior of parameter 'name' of the resource is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  name = huaweicloud_eg_custom_event_source.test.name
}

data "huaweicloud_eg_custom_event_sources" "name_not_found" {
  // Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  // needs to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  name = "resource_not_found"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_sources.filter_by_name.sources[*].id :
                   v == huaweicloud_eg_custom_event_source.test.id]
}

output "is_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_sources.name_not_found.sources) == 0
}
`, testAccDataCustomEventSources_base())
}

func TestAccDataCustomEventSources_filterById(t *testing.T) {
	var (
		baseRes    = "huaweicloud_eg_custom_event_source.test"
		byId       = "data.huaweicloud_eg_custom_event_sources.filter_by_id"
		idNotFound = "data.huaweicloud_eg_custom_event_sources.id_not_found"

		obj          custom.Source
		rc           = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventSourceFunc)
		dcById       = acceptance.InitDataSourceCheck(byId)
		dcIdNotFound = acceptance.InitDataSourceCheck(idNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventSources_filterById(),
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

func testAccDataCustomEventSources_filterById() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_sources" "filter_by_id" {
  source_id = huaweicloud_eg_custom_event_source.test.id
}

data "huaweicloud_eg_custom_event_sources" "id_not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  source_id = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_sources.filter_by_id.sources[*].id :
                   v == huaweicloud_eg_custom_event_source.test.id]
}

output "is_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "id_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_sources.id_not_found.sources) == 0
}
`, testAccDataCustomEventSources_base(), randUUID)
}

func TestAccDataCustomEventSources_filterByChannelId(t *testing.T) {
	var (
		baseRes           = "huaweicloud_eg_custom_event_source.test"
		byChannelId       = "data.huaweicloud_eg_custom_event_sources.filter_by_channel_id"
		channelIdNotFound = "data.huaweicloud_eg_custom_event_sources.channel_id_not_found"

		obj                 custom.Source
		rc                  = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventSourceFunc)
		dcByChannelId       = acceptance.InitDataSourceCheck(byChannelId)
		dcChannelIdNotFound = acceptance.InitDataSourceCheck(channelIdNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventSources_filterByChannelId(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByChannelId.CheckResourceExists(),
					resource.TestCheckOutput("is_channel_id_filter_useful", "true"),
					dcChannelIdNotFound.CheckResourceExists(),
					resource.TestCheckOutput("channel_id_not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventSources_filterByChannelId() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_sources" "filter_by_channel_id" {
  // The behavior of parameter 'channel_id' of the resource is 'Required', means this parameter does not have 'Know After Apply' behavior.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  channel_id = "%[2]s"
}

data "huaweicloud_eg_custom_event_sources" "channel_id_not_found" {
  // Since a random ID is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  channel_id = "%[3]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_sources.filter_by_channel_id.sources[*].channel_id : v == "%[2]s"]
}

output "is_channel_id_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "channel_id_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_sources.channel_id_not_found.sources) == 0
}
`, testAccDataCustomEventSources_base(), acceptance.HW_EG_CHANNEL_ID, randUUID)
}

func TestAccDataCustomEventSources_filterByFuzzyName(t *testing.T) {
	var (
		baseRes           = "huaweicloud_eg_custom_event_source.test"
		byFuzzyName       = "data.huaweicloud_eg_custom_event_sources.filter_by_fuzzy_name"
		fuzzyNameNotFound = "data.huaweicloud_eg_custom_event_sources.fuzzy_name_not_found"

		obj                 custom.Source
		rc                  = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventSourceFunc)
		dcByFuzzyName       = acceptance.InitDataSourceCheck(byFuzzyName)
		dcFuzzyNameNotFound = acceptance.InitDataSourceCheck(fuzzyNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventSources_filterByFuzzyName(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcByFuzzyName.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_name_filter_useful", "true"),
					dcFuzzyNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("fuzzy_name_not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDataCustomEventSources_filterByFuzzyName() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_eg_custom_event_sources" "filter_by_fuzzy_name" {
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  fuzzy_name = huaweicloud_eg_custom_event_source.test.name
}

data "huaweicloud_eg_custom_event_sources" "fuzzy_name_not_found" {
  // Since a random name is used, there is no dependency relationship with resource attachment, and the dependency needs
  // to be manually set.
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  fuzzy_name = "%[2]s"
}

locals {
  filter_result = [for v in data.huaweicloud_eg_custom_event_sources.filter_by_fuzzy_name.sources[*].name : 
                   length(regexall("${huaweicloud_eg_custom_event_source.test.name}", v)) > 0]
}

output "is_fuzzy_name_filter_useful" {
  value = alltrue(local.filter_result) && length(local.filter_result) > 0
}

output "fuzzy_name_not_found_validation_pass" {
  value = length(data.huaweicloud_eg_custom_event_sources.fuzzy_name_not_found.sources) == 0
}
`, testAccDataCustomEventSources_base(), randUUID)
}

func TestAccDataCustomEventSources_filterBySort(t *testing.T) {
	var (
		baseRes = "huaweicloud_eg_custom_event_source.test"
		bySort  = "data.huaweicloud_eg_custom_event_sources.filter_by_sort"

		obj      custom.Source
		rc       = acceptance.InitResourceCheck(baseRes, &obj, getCustomEventSourceFunc)
		dcBySort = acceptance.InitDataSourceCheck(bySort)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomEventSources_filterBySort(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					dcBySort.CheckResourceExists(),
					resource.TestCheckResourceAttr(bySort, "sources.#", "1"),
				),
			},
		},
	})
}

func testAccDataCustomEventSources_filterBySort() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = "%[1]s"
  name       = "%[2]s"
}

data "huaweicloud_eg_custom_event_sources" "filter_by_sort" {
  depends_on = [
    huaweicloud_eg_custom_event_source.test,
  ]

  sort = "created_time:asc"
}

output "is_sort_filter_useful" {
  value = length(data.huaweicloud_eg_custom_event_sources.filter_by_sort.sources) > 0
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}
