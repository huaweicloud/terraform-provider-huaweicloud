package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKmsKeys_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_kms_keys.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byKeyState   = "data.huaweicloud_kms_keys.filter_by_key_state"
		dcByKeyState = acceptance.InitDataSourceCheck(byKeyState)

		bykeyAlgorithm   = "data.huaweicloud_kms_keys.filter_by_key_algorithm"
		dcBykeyAlgorithm = acceptance.InitDataSourceCheck(bykeyAlgorithm)

		byEps   = "data.huaweicloud_kms_keys.filter_by_enterprise_project_id"
		dcByEps = acceptance.InitDataSourceCheck(byEps)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsKeys_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keys.#"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_alias"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_algorithm"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_usage"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_type"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_state"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.created_at"),

					dcByKeyState.CheckResourceExists(),
					resource.TestCheckOutput("key_state_filter_is_useful", "true"),

					dcBykeyAlgorithm.CheckResourceExists(),
					resource.TestCheckOutput("key_algorithm_filter_is_useful", "true"),

					dcByEps.CheckResourceExists(),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceKmsKeys_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_kms_keys" "test" {
  depends_on = [
    huaweicloud_kms_key.test
  ]
}

locals {
  key_state = data.huaweicloud_kms_keys.test.keys[0].key_state
}

data "huaweicloud_kms_keys" "filter_by_key_state" {
  key_state = local.key_state
}

locals {
  key_state_filter_result = [
    for v in data.huaweicloud_kms_keys.filter_by_key_state.keys[*].key_state : v == local.key_state
  ]
}

output "key_state_filter_is_useful" {
  value = alltrue(local.key_state_filter_result) && length(local.key_state_filter_result) > 0
}

locals {
  key_algorithm = data.huaweicloud_kms_keys.test.keys[0].key_algorithm
}

data "huaweicloud_kms_keys" "filter_by_key_algorithm" {
  key_algorithm = local.key_algorithm
}

locals {
  key_algorithm_filter_result = [ 
    for v in data.huaweicloud_kms_keys.filter_by_key_algorithm.keys[*].key_algorithm : v == local.key_algorithm
  ]
}

output "key_algorithm_filter_is_useful" {
  value = alltrue(local.key_algorithm_filter_result) && length(local.key_algorithm_filter_result) > 0
}

locals {
  enterprise_project_id = data.huaweicloud_kms_keys.test.keys[0].enterprise_project_id
}

data "huaweicloud_kms_keys" "filter_by_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id
}

locals {
  enterprise_project_id_filter_result = [
    for v in data.huaweicloud_kms_keys.filter_by_enterprise_project_id.keys[*].enterprise_project_id : 
    v == local.enterprise_project_id
  ]
}

output "enterprise_project_id_filter_is_useful" {
  value = alltrue(local.enterprise_project_id_filter_result) && length(local.enterprise_project_id_filter_result) > 0
}
`, testAccKmsKey_basic(name))
}
