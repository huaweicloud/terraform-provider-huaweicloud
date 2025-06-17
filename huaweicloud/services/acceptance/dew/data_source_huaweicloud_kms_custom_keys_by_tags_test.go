package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsCustomKeysByTagsDataSource_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		tagKeyName     = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_kms_custom_keys_by_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testKmsCustomKeysByTagsDataSource_basic(name, tagKeyName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.key_alias"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.key_description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.creation_date"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.key_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.default_key_flag"),
					resource.TestCheckOutput("single_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("multi_value_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("multi_tag_filter_is_useful", "true"),
					resource.TestCheckOutput("resource_name_match_filter_is_useful", "true"),
					resource.TestCheckOutput("tag_and_match_filter_is_useful", "true"),
					resource.TestCheckOutput("count_action_is_useful", "true"),
					resource.TestCheckOutput("empty_count_action_is_useful", "true"),
				),
			},
		},
	})
}

func testKmsCustomKeysByTagsDataSource_basic(name string, tagKeyName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test1" {
  key_alias       = "%[1]s_1"
  pending_days    = "7"
  key_description = "test description 1"
  
  tags = {
    %[2]s_1 = "test_value_1"
    %[2]s_2 = "test_value_2"
  }
}

resource "huaweicloud_kms_key" "test2" {
  key_alias       = "%[1]s_2"
  pending_days    = "7"
  key_description = "test description 2"
  
  tags = {
    %[2]s_1 = "test_value_1"
    %[2]s_3 = "test_value_3"
  }
}

resource "huaweicloud_kms_key" "test3" {
  key_alias       = "%[1]s_3"
  pending_days    = "7"
  key_description = "test description 3"
  
  tags = {
    %[2]s_4 = "test_value_4"
  }
}

data "huaweicloud_kms_custom_keys_by_tags" "test" {
  action = "filter"

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

# Test single tag filter
data "huaweicloud_kms_custom_keys_by_tags" "single_tag_filter" {
  action = "filter"

  tags {
    key    = "%[2]s_1"
    values = ["test_value_1"]
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "single_tag_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_kms_custom_keys_by_tags.single_tag_filter.resources) == 2,
  ])
}

# Test multiple values for a single tag
data "huaweicloud_kms_custom_keys_by_tags" "multi_value_tag_filter" {
  action = "filter"

  tags {
    key    = "%[2]s_1"
    values = ["test_value_1", "non_existent_value"]
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "multi_value_tag_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_kms_custom_keys_by_tags.multi_value_tag_filter.resources) == 2,
  ])
}

# Test multiple tags (intersection)
data "huaweicloud_kms_custom_keys_by_tags" "multi_tag_filter" {
  action = "filter"

  tags {
    key    = "%[2]s_1"
    values = ["test_value_1"]
  }
  
  tags {
    key    = "%[2]s_2"
    values = ["test_value_2"]
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "multi_tag_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_kms_custom_keys_by_tags.multi_tag_filter.resources) == 1,
  ])
}

# Test resource_name match filter
data "huaweicloud_kms_custom_keys_by_tags" "resource_name_match_filter" {
  action = "filter"

  matches {
    key   = "resource_name"
    value = "%[1]s_1"
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "resource_name_match_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_kms_custom_keys_by_tags.resource_name_match_filter.resources) == 1,
  ])
}

# Test tag and match filter together (intersection)
data "huaweicloud_kms_custom_keys_by_tags" "tag_and_match_filter" {
  action = "filter"

  tags {
    key    = "%[2]s_1"
    values = ["test_value_1"]
  }

  matches {
    key   = "resource_name"
    value = "%[1]s_1"
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "tag_and_match_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_kms_custom_keys_by_tags.tag_and_match_filter.resources) == 1,
  ])
}

# Test count action
data "huaweicloud_kms_custom_keys_by_tags" "count_action" {
  action = "count"

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "count_action_is_useful" {
  value = alltrue([
    data.huaweicloud_kms_custom_keys_by_tags.count_action.total_count >= 3,
    length(data.huaweicloud_kms_custom_keys_by_tags.count_action.resources) == 0,
  ])
}

# Test count action with empty result
data "huaweicloud_kms_custom_keys_by_tags" "empty_count_action" {
  action = "count"
  
  tags {
    key    = "non_existent_tag"
    values = ["non_existent_value"]
  }

  depends_on = [
    huaweicloud_kms_key.test1,
    huaweicloud_kms_key.test2,
    huaweicloud_kms_key.test3
  ]
}

output "empty_count_action_is_useful" {
  value = alltrue([
    data.huaweicloud_kms_custom_keys_by_tags.empty_count_action.total_count == 0,
  ])
}
`, name, tagKeyName)
}
