package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkspaceDesktopTags_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_desktop_tags.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkspaceDesktopTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_valid", "true"),
					resource.TestCheckOutput("is_key_valid", "true"),
					resource.TestCheckOutput("is_value_valid", "true"),
				),
			},
			{
				Config:      testAccDataSourceWorkspaceDesktopTags_NotExists(),
				ExpectError: regexp.MustCompile(`The resource does not exist or is a resource of another project. DesktopId is not exist.`),
			},
		},
	})
}

func testAccDataSourceWorkspaceDesktopTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags" "test" {
  desktop_id = "%[1]s"
}

locals {
  tag_result = try(data.huaweicloud_workspace_desktop_tags.test.tags[0], {})
}

output "is_tags_valid" {
  value = length(data.huaweicloud_workspace_desktop_tags.test.tags) > 0
}

output "is_key_valid" {
  value = try(local.tag_result.key != null && local.tag_result.key != "", false)
}

output "is_value_valid" {
  value = try(local.tag_result.value != null, false)
}
`, acceptance.HW_WORKSPACE_DESKTOP_ID)
}

func testAccDataSourceWorkspaceDesktopTags_NotExists() string {
	random_id, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_workspace_desktop_tags" "not_exist" {
  desktop_id = "%[1]s"
}
`, random_id)
}
