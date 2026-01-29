package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceVpcepResourcesByTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_vpcep_resources_by_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpcepServicesByTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					resource.TestCheckResourceAttrSet(dataSourceName, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources.0.tags.#"),

					resource.TestCheckOutput("count_filter_is_useful", "true"),
					resource.TestCheckOutput("matches_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_filter_is_useful", "true"),
					resource.TestCheckOutput("tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("not_tags_any_filter_is_useful", "true"),
					resource.TestCheckOutput("without_any_tag_filter_is_useful", "true"),
				),
			},
		},
	},
	)
}

func testDataSourceVpcepServicesByTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpcep_public_services" "test" {
  service_name = "dns"
}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id = data.huaweicloud_vpcep_public_services.test.services.0.id
  vpc_id     = huaweicloud_vpc.test.id
  network_id = huaweicloud_vpc_subnet.test.id

  tags = {
    owner = "terraform"
    foo   = "bar"
    test  = "acc"
    key   = "value" 
  }
}
`, common.TestVpc(name))
}

func testDataSourceVpcepServicesByTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_resources_by_tags" "test" {
  depends_on    = [huaweicloud_vpcep_endpoint.test]
  resource_type = "endpoint"
  action        = "filter"
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_count" {
  depends_on    = [huaweicloud_vpcep_endpoint.test]
  resource_type = "endpoint"
  action        = "count"
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_matches" {
  resource_type = "endpoint"
  action        = "filter"

  matches {
    key   = "resource_name"
    value = data.huaweicloud_vpcep_resources_by_tags.test.resources.0.resource_name
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_tags" {
  resource_type = "endpoint"
  action        = "filter"

  tags {
    key    = data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.0.key
    values = [data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.0.value]
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_not_tags" {
  resource_type = "endpoint"
  action        = "filter"

  not_tags {
    key    = data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.1.key
    values = [data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.1.value]
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_tags_any" {
  resource_type = "endpoint"
  action        = "filter"

  tags_any {
    key    = data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.2.key
    values = [data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.2.value]
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_not_tags_any" {
  resource_type = "endpoint"
  action        = "filter"

  not_tags_any {
    key    = data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.3.key
    values = [data.huaweicloud_vpcep_resources_by_tags.test.resources.0.tags.3.value]
  }
}

data "huaweicloud_vpcep_resources_by_tags" "filter_by_without_any_tag" {
  resource_type   = "endpoint"
  action          = "filter"
  without_any_tag = true
  depends_on      = [huaweicloud_vpcep_endpoint.test]
}

output "count_filter_is_useful" {
  value = data.huaweicloud_vpcep_resources_by_tags.filter_by_count.total_count > 0
}

output "matches_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_matches.resources) == 1
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_tags.resources) > 0
}

output "not_tags_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_not_tags.resources) == 0
}

output "tags_any_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_tags_any.resources) > 0
}

output "not_tags_any_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_not_tags_any.resources) == 0
}

output "without_any_tag_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_resources_by_tags.filter_by_without_any_tag.resources) == 0
}
`, testDataSourceVpcepServicesByTags_base(name))
}
