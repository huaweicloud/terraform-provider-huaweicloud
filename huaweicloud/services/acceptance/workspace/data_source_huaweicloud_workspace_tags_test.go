package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkspaceTagsDataSource_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_tags.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		dcByKeyName = "data.huaweicloud_workspace_tags.test_by_key"
		dcByKey     = acceptance.InitDataSourceCheck(dcByKeyName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkspaceTagsDataSource_base("terraform_test"),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "tags.#"),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dcName, "tags.0.values.#"),
					dcByKey.CheckResourceExists(),
					resource.TestCheckOutput("is_key_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccWorkspaceTagsDataSource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_tags" "test" {}

locals {
  key = "%[1]s"
}

data "huaweicloud_workspace_tags" "filter_by_key" {
  key = "%[1]s"
}

locals {
  tag_filter_result = [
     for v in data.huaweicloud_workspace_tags.filter_by_key.tags[*].key : v == local.key
  ]
}

output "is_key_filter_useful" {
  value = length(local.tag_filter_result) > 0 && alltrue(local.tag_filter_result)
}
`, name)
}
