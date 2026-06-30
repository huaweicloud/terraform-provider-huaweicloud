package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureDirectories_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_directories.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureDirectories_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture directories"),
			},
			{
				Config: testAccDataArchitectureDirectories_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Query all directories
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "directories.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_all_directories_queried", "true"),
					resource.TestCheckResourceAttrSet(all, "directories.0.id"),
					resource.TestCheckResourceAttrSet(all, "directories.0.name"),
					resource.TestCheckResourceAttr(all, "directories.0.type", "STANDARD_ELEMENT"),
					resource.TestCheckResourceAttrSet(all, "directories.0.qualified_name"),
					resource.TestCheckResourceAttrSet(all, "directories.0.created_by"),
					resource.TestCheckResourceAttrSet(all, "directories.0.updated_by"),
					resource.TestMatchResourceAttr(all, "directories.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "directories.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataArchitectureDirectories_nonExistentWorkspace() string {
	randUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_directories" "test" {
  workspace_id = "%[1]s"
  type         = "STANDARD_ELEMENT"
}
`, randUUID.String())
}

func testAccDataArchitectureDirectories_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_directory" "test" {
  count = 2

  workspace_id = "%[1]s"
  name         = format("%[2]s_standard_%%d", count.index)
  type         = "STANDARD_ELEMENT"
}

data "huaweicloud_dataarts_architecture_directories" "all" {
  depends_on = [huaweicloud_dataarts_architecture_directory.test]

  workspace_id = "%[1]s"
  type         = "STANDARD_ELEMENT"
}

locals {
  all_directories_queried = [
    for v in huaweicloud_dataarts_architecture_directory.test[*].id :
      contains(data.huaweicloud_dataarts_architecture_directories.all.directories[*].id, v)
  ]
}

output "is_all_directories_queried" {
  value = length(local.all_directories_queried) > 0 && alltrue(local.all_directories_queried)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}
