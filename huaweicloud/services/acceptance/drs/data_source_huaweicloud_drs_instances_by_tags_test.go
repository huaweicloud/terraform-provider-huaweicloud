package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceInstancesByTags_basic(t *testing.T) {
	var (
		datasource         = "data.huaweicloud_drs_instances_by_tags.test"
		filterDatasource   = "data.huaweicloud_drs_instances_by_tags.filter_by_name"
		nonExistDatasource = "data.huaweicloud_drs_instances_by_tags.non_exist"
		dc                 = acceptance.InitDataSourceCheck(datasource)
		dcFilter           = acceptance.InitDataSourceCheck(filterDatasource)
		dcNonExist         = acceptance.InitDataSourceCheck(nonExistDatasource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Note: Please ensure that test data (DRS instances) is prepared in the test environment.
			acceptance.TestAccPreCheckDRSEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceInstancesByTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcFilter.CheckResourceExists(),
					dcNonExist.CheckResourceExists(),
					resource.TestCheckOutput("filter_result_matches", "true"),
					resource.TestCheckOutput("non_exist_result_is_empty", "true"),
					resource.TestCheckResourceAttrSet(datasource, "resources.#"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(datasource, "resources.0.resource_name"),
				),
			},
		},
	})
}

const testAccDatasourceInstancesByTags_basic = `
data "huaweicloud_drs_instances_by_tags" "test" {
  resource_type = "sync"
}

data "huaweicloud_drs_instances_by_tags" "filter_by_name" {
  resource_type = "sync"
  
  matches {
    key   = "resource_name"
    value = data.huaweicloud_drs_instances_by_tags.test.resources[0].resource_name
  }
}

data "huaweicloud_drs_instances_by_tags" "non_exist" {
  resource_type = "sync"
  
  tags {
    key    = "non-exist-tag-key"
    values = ["non-exist-tag-value"]
  }
}

output "filter_result_matches" {
  value = length(data.huaweicloud_drs_instances_by_tags.filter_by_name.resources) >= 1
}

output "non_exist_result_is_empty" {
  value = length(data.huaweicloud_drs_instances_by_tags.non_exist.resources) == 0
}
`
