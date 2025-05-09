package lts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHostAccesses_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_lts_host_accesses.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_lts_host_accesses.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNames   = "data.huaweicloud_lts_host_accesses.filter_by_names"
		dcByNames = acceptance.InitDataSourceCheck(byNames)

		notFound   = "data.huaweicloud_lts_host_accesses.filter_by_not_found"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)

		byHostGroupNames   = "data.huaweicloud_lts_host_accesses.filter_by_host_group_names"
		dcByHostGroupNames = acceptance.InitDataSourceCheck(byHostGroupNames)

		byLogGroupNames   = "data.huaweicloud_lts_host_accesses.filter_by_log_group_names"
		dcByLogGroupNames = acceptance.InitDataSourceCheck(byLogGroupNames)

		byLogStreamNames   = "data.huaweicloud_lts_host_accesses.filter_by_log_stream_names"
		dcByLogStreamNames = acceptance.InitDataSourceCheck(byLogStreamNames)

		byTags   = "data.huaweicloud_lts_host_accesses.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHostAccesses_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "accesses.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNames.CheckResourceExists(),
					resource.TestCheckOutput("is_names_filter_useful", "true"),
					dcNotFound.CheckResourceExists(),
					resource.TestCheckOutput("host_access_not_found", "true"),
					dcByHostGroupNames.CheckResourceExists(),
					resource.TestCheckOutput("is_host_group_names_filter_useful", "true"),
					dcByLogGroupNames.CheckResourceExists(),
					resource.TestCheckOutput("is_log_group_names_filter_useful", "true"),
					dcByLogStreamNames.CheckResourceExists(),
					resource.TestCheckOutput("is_log_stream_names_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_windows_log_info_set", "true"),
					resource.TestCheckOutput("is_multi_log_format_set", "true"),
					resource.TestCheckResourceAttrSet(byName, "accesses.0.id"),
					resource.TestCheckResourceAttr(byName, "accesses.0.name", rName),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_group_name", "huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_stream_name", "huaweicloud_lts_stream.test", "stream_name"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.paths.#", "2"),
					resource.TestCheckResourceAttr(byName, "accesses.0.access_config.0.black_paths.#", "2"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.access_config.0.single_log_format.0.mode",
						"huaweicloud_lts_host_access.test", "access_config.0.single_log_format.0.mode"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.access_config.0.single_log_format.0.value",
						"huaweicloud_lts_host_access.test", "access_config.0.single_log_format.0.value"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.access_config.0.custom_key_value",
						"huaweicloud_lts_host_access.test", "access_config.0.custom_key_value"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.access_config.0.system_fields",
						"huaweicloud_lts_host_access.test", "access_config.0.system_fields"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.access_config.0.repeat_collect",
						"huaweicloud_lts_host_access.test", "access_config.0.repeat_collect"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.processor_type",
						"huaweicloud_lts_host_access.test", "processor_type"),
					resource.TestCheckResourceAttr(byName, "accesses.0.processors.#", "2"),
					resource.TestCheckResourceAttrSet(byName, "accesses.0.processors.0.type"),
					resource.TestCheckResourceAttrSet(byName, "accesses.0.processors.0.detail"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.demo_log",
						"huaweicloud_lts_host_access.test", "demo_log"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.demo_fields.0.field_name",
						"huaweicloud_lts_host_access.test", "demo_fields.0.field_name"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.demo_fields.0.field_value",
						"huaweicloud_lts_host_access.test", "demo_fields.0.field_value"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.binary_collect",
						"huaweicloud_lts_host_access.test", "binary_collect"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.encoding_format",
						"huaweicloud_lts_host_access.test", "encoding_format"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.incremental_collect",
						"huaweicloud_lts_host_access.test", "incremental_collect"),
					resource.TestCheckResourceAttrPair(byName, "accesses.0.log_split",
						"huaweicloud_lts_host_access.test", "log_split"),
					resource.TestMatchResourceAttr(byName, "accesses.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataSourceHostAccesses_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_host_access" "with_windows" {
  name           = "%[2]s_with_windows"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]

  access_config {
    paths = ["D:\\data\\log\\*"]

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    multi_log_format {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }
  }

  tags = {
    foo = "bar"
    TF  = "oenwr"
  }
}
`, testHostAccessConfig_basic_step1(name), name)
}

func testAccDataSourceHostAccesses_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_lts_host_accesses" "test" {
  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

# Filter by host access name.
locals {
  host_access_name = huaweicloud_lts_host_access.test.name
}

# Filter by host access name.
data "huaweicloud_lts_host_accesses" "filter_by_name" {
  access_config_name_list = [local.host_access_name]

  depends_on = [huaweicloud_lts_host_access.test]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_lts_host_accesses.filter_by_name.accesses[*].name : v == local.host_access_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by host access names.
locals {
  host_access_names = [huaweicloud_lts_host_access.test.name, huaweicloud_lts_host_access.with_windows.name]
}

data "huaweicloud_lts_host_accesses" "filter_by_names" {
  access_config_name_list = local.host_access_names

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

locals {
  host_assesses       = data.huaweicloud_lts_host_accesses.filter_by_names.accesses
  names_filter_result = [for v in local.host_assesses[*].name : v if contains(local.host_access_names, v)]
}

output "is_names_filter_useful" {
  value = join(",", sort(local.names_filter_result)) == join(",", sort(local.host_access_names))
}

# Host access not found.
data "huaweicloud_lts_host_accesses" "filter_by_not_found" {
  access_config_name_list = local.host_access_names
  host_group_name_list    = ["host_group_not_found"]

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

output "host_access_not_found" {
  value = length(data.huaweicloud_lts_host_accesses.filter_by_not_found.accesses) == 0
}

# Filter by host group names associated with host access.
data "huaweicloud_lts_host_accesses" "filter_by_host_group_names" {
  host_group_name_list = [huaweicloud_lts_host_group.test.name]

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

locals {
  host_group_id = huaweicloud_lts_host_group.test.id
  host_group_names_filter_result = [
    for v in data.huaweicloud_lts_host_accesses.filter_by_host_group_names.accesses[*].host_group_ids : contains(v, local.host_group_id)
  ]
}

output "is_host_group_names_filter_useful" {
  value = length(local.host_group_names_filter_result) > 0 && alltrue(local.host_group_names_filter_result)
}

# Filter by log group names associated with host access.
locals {
  log_group_name = huaweicloud_lts_group.test.group_name
}

data "huaweicloud_lts_host_accesses" "filter_by_log_group_names" {
  log_group_name_list = [local.log_group_name]

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

locals {
  log_group_names_filter_result = [
    for v in data.huaweicloud_lts_host_accesses.filter_by_log_group_names.accesses[*].log_group_name : v == local.log_group_name
  ]
}

output "is_log_group_names_filter_useful" {
  value = length(local.log_group_names_filter_result) > 0 && alltrue(local.log_group_names_filter_result)
}

# Filter by log stream names associated with host access.
locals {
  log_stream_name = huaweicloud_lts_stream.test.stream_name
}

data "huaweicloud_lts_host_accesses" "filter_by_log_stream_names" {
  log_stream_name_list = [local.log_stream_name]

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

locals {
  log_stream_names_filter_result = [
    for v in data.huaweicloud_lts_host_accesses.filter_by_log_stream_names.accesses[*].log_stream_name : v == local.log_stream_name
  ]
}

output "is_log_stream_names_filter_useful" {
  value = length(local.log_stream_names_filter_result) > 0 && alltrue(local.log_stream_names_filter_result)
}

# Filter by host access tags.
locals {
  tags = huaweicloud_lts_host_access.test.tags
}

data "huaweicloud_lts_host_accesses" "filter_by_tags" {
  tags = local.tags

  depends_on = [
    huaweicloud_lts_host_access.test,
    huaweicloud_lts_host_access.with_windows
  ]
}

locals {
  tags_filter_result = [for item in data.huaweicloud_lts_host_accesses.filter_by_tags.accesses[*].tags : true
  if anytrue([for k, v in local.tags : lookup(item, k, "not_found") == v])]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}

# Check attributes.
locals {
  with_windows              = [for v in local.host_assesses : v if v.id == huaweicloud_lts_host_access.with_windows.id]
  windows_log_info          = try(local.with_windows[0].access_config[0].windows_log_info[0], {})
  multi_log_format          = try(local.with_windows[0].access_config[0].multi_log_format[0], {})
  resource_windows_log_info = huaweicloud_lts_host_access.with_windows.access_config[0].windows_log_info[0]
  resource_multi_log_format = huaweicloud_lts_host_access.with_windows.access_config[0].multi_log_format[0]
}

output "is_windows_log_info_set" {
  value = alltrue([local.windows_log_info.time_offset == local.resource_windows_log_info.time_offset,
    local.windows_log_info.time_offset_unit == local.resource_windows_log_info.time_offset_unit,
    length(local.windows_log_info.categorys) == length(local.resource_windows_log_info.categorys),
    length(local.windows_log_info.event_level) == length(local.resource_windows_log_info.event_level)
  ])
}

output "is_multi_log_format_set" {
  value = alltrue([local.multi_log_format.mode == local.resource_multi_log_format.mode,
    local.multi_log_format.value == local.resource_multi_log_format.value
  ])
}
`, testAccDataSourceHostAccesses_base(name), testHostAccessConfig_windows_basic(name))
}
