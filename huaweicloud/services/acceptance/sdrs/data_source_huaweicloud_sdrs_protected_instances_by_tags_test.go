package sdrs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSdrsProtectedInstancesByTags_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_sdrs_protected_instances_by_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSdrsProtectedInstancesByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.source_server"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.target_server"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.server_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.metadata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.attachment.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_detail.0.priority_station"),

					resource.TestCheckOutput("action_count_is_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_exist_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_non_exist_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_any_exist_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_any_all_non_exist_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_not_all_non_exist_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_not_any_all_non_exist_useful", "true"),
					resource.TestCheckOutput("is_matches_filter_exist_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSdrsProtectedInstancesByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sdrs_protected_instances_by_tags" "test" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"
}

# Query count
data "huaweicloud_sdrs_protected_instances_by_tags" "total_count" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "count"
}

output "action_count_is_useful" {
  value = data.huaweicloud_sdrs_protected_instances_by_tags.total_count.total_count > 0
}

# Filter by tags exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_tags_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  tags {
    key    = "foo"
    values = ["bar"]
  }

  tags {
    key    = "key"
    values = ["value"]
  }
}

output "is_tags_filter_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_tags_exist.resources) > 0
}

# Filter by tags non-exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_tags_non_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  tags {
    key    = "foo-non-exist"
    values = ["bar"]
  }

  tags {
    key    = "key"
    values = ["value"]
  }
}

output "is_tags_filter_non_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_tags_non_exist.resources) == 0
}

# Filter by tags_any one exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_tags_any_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  tags_any {
    key    = "foo-non-exist"
    values = ["bar"]
  }

  tags_any {
    key    = "key"
    values = ["value"]
  }
}

output "is_tags_filter_any_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_tags_any_exist.resources) > 0
}

# Filter by tags_any all non exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_tags_any_all_non_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter	"

  tags_any {
    key    = "foo-non-exist"
    values = ["bar"]
  }

  tags_any {
    key    = "key-non-exist"
    values = ["value"]
  }
}

output "is_tags_filter_any_all_non_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_tags_any_all_non_exist.resources) == 0
}

# Filter by not_tags all non exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_not_tags_all_non_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  not_tags {
    key    = "foo-non-exist"
    values = ["bar"]
  }

  not_tags {
    key    = "key-non-exist"
    values = ["value"]
  }
}

output "is_tags_filter_not_all_non_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_not_tags_all_non_exist.resources) > 0
}

# Filter by not_tags_any all non exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_not_tags_any_all_non_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  not_tags_any {
    key    = "foo-non-exist"
    values = ["bar"]
  }

  not_tags_any {
    key    = "key-non-exist"
    values = ["value"]
  }
}

output "is_tags_filter_not_any_all_non_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_not_tags_any_all_non_exist.resources) > 0
}

# Filter by matches exist
data "huaweicloud_sdrs_protected_instances_by_tags" "filter_by_matches_exist" {
  depends_on = [huaweicloud_sdrs_protected_instance.test]

  action = "filter"

  matches {
    key   = "resource_name"
    value = huaweicloud_sdrs_protected_instance.test.name
  }
}

output "is_matches_filter_exist_useful" {
  value = length(data.huaweicloud_sdrs_protected_instances_by_tags.filter_by_matches_exist.resources) > 0
}
`, testProtectedInstance_basic(name))
}
