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
