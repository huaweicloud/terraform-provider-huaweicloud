package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCSMSSecretsByTagDataSource_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_csms_secrets_by_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCSMSSecretsByTagDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.auto_rotation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.kms_key_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.secret_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_detail.0.update_time"),
					resource.TestCheckOutput("tags_single_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_double_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_single_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_double_filter_is_useful", "true"),
					resource.TestCheckOutput("count_is_useful", "true"),
				),
			},
		},
	})
}

func testCSMSSecretsByTagDataSource_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name                  = "%s_1"
  secret_text           = "this is a password"
  description           = "csms secret test"
  secret_type           = "COMMON"
  enterprise_project_id = "0"

  tags = {
    test_tag_1 = "test_value_1"
  }
}

resource "huaweicloud_csms_secret" "test2" {
  name                  = "%s_2"
  secret_text           = "this is a password"
  description           = "csms secret test"
  secret_type           = "COMMON"
  enterprise_project_id = "0"

  tags = {
    test_tag_1 = "test_value_2"
  }
}

data "huaweicloud_csms_secrets_by_tags" "test" {
  resource_instances = "resource_instances"
  action             = "filter"
  sequence           = "test_sequence"

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

locals {
  secret1_key   = keys(huaweicloud_csms_secret.test.tags)[0]
  secret1_value = huaweicloud_csms_secret.test.tags[keys(huaweicloud_csms_secret.test.tags)[0]]
  secret2_key   = keys(huaweicloud_csms_secret.test2.tags)[0]
  secret2_value = huaweicloud_csms_secret.test2.tags[keys(huaweicloud_csms_secret.test.tags)[0]]
  secret1_name  = huaweicloud_csms_secret.test.name
  secret2_name  = huaweicloud_csms_secret.test2.name
}

data "huaweicloud_csms_secrets_by_tags" "tags_single_filter" {
  resource_instances = "resource_instances"
  action             = "filter"

  tags {
    key    = local.secret1_key
    values = [local.secret1_value]
  }

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

data "huaweicloud_csms_secrets_by_tags" "tags_double_filter" {
  resource_instances = "resource_instances"
  action             = "filter"

  tags {
    key    = local.secret1_key
    values = [local.secret1_value,local.secret2_value]
  }

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

data "huaweicloud_csms_secrets_by_tags" "matches_single_filter" {
  resource_instances = "resource_instances"
  action             = "filter"

  matches {
    key   = "resource_name"
    value = local.secret1_name
  }

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

data "huaweicloud_csms_secrets_by_tags" "matches_double_filter" {
  resource_instances = "resource_instances"
  action             = "filter"

  matches {
    key   = "resource_name"
    value = local.secret1_name
  }

  matches {
    key   = "resource_name"
    value = local.secret2_name
  }

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

data "huaweicloud_csms_secrets_by_tags" "count" {
  resource_instances = "resource_instances"
  action             = "count"

  depends_on = [huaweicloud_csms_secret.test, huaweicloud_csms_secret.test2]
}

output "tags_single_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_csms_secrets_by_tags.tags_single_filter.resources) == 1,
    data.huaweicloud_csms_secrets_by_tags.tags_single_filter.resources[0].tags[0].key == local.secret1_key,
    data.huaweicloud_csms_secrets_by_tags.tags_single_filter.resources[0].tags[0].value == local.secret1_value
  ])
}

output "tags_double_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources) == 2,
    data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources[0].tags[0].key == local.secret1_key,
    data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources[0].tags[0].value == local.secret1_value,
    data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources[1].tags[0].key == local.secret1_key,
    data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources[1].tags[0].value == local.secret2_value,
  ])
}

output "matches_single_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_csms_secrets_by_tags.matches_single_filter.resources) == 1,
    data.huaweicloud_csms_secrets_by_tags.tags_double_filter.resources[0].resource_name == local.secret1_name,
  ])
}

output "matches_double_filter_is_useful" {
  value = alltrue([
    length(data.huaweicloud_csms_secrets_by_tags.matches_double_filter.resources) == 0,
  ])
}

output "count_is_useful" {
  value = alltrue([
    data.huaweicloud_csms_secrets_by_tags.count.total_count == 2,
  ])
}
`, name, name)
}
