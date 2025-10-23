package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCpcsAppAccessKeys_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cpcs_app_access_keys.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCpcsAppAccessKeys_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.access_key_id"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.access_key"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.is_downloaded"),
					resource.TestCheckResourceAttrSet(dataSource, "access_keys.0.is_imported"),

					resource.TestCheckOutput("is_key_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCpcsAppAccessKeys_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cpcs_app_access_keys" "test" {
  depends_on = [huaweicloud_cpcs_app_access_key.test]
  app_id     = huaweicloud_cpcs_app.test.id
}

locals {
  key_name = data.huaweicloud_cpcs_app_access_keys.test.access_keys.0.key_name
}

# Filter by key_name
data "huaweicloud_cpcs_app_access_keys" "key_name_filter" {
  app_id   = huaweicloud_cpcs_app.test.id
  key_name = local.key_name
}

locals {
  key_name_filter_result = [
    for v in data.huaweicloud_cpcs_app_access_keys.key_name_filter.access_keys[*].key_name : v == local.key_name
  ]
}

output "is_key_name_filter_useful" {
  value = length(local.key_name_filter_result) > 0 && alltrue(local.key_name_filter_result)
}
`, testCpcsAppAccessKey_basic(name))
}
