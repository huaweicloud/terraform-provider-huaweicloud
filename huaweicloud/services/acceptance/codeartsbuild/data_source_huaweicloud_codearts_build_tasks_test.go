package codeartsbuild

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCodeArtsBuildTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_build_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCodeArtsBuildTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.health_score"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_code"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_finished"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.disabled"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.favorite"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_modify"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_delete"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_execute"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_copy"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_forbidden"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.is_view"),

					resource.TestCheckOutput("is_creator_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceCodeArtsBuildTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_build_tasks" "test" {
  depends_on = [huaweicloud_codearts_build_task.test]

  project_id = huaweicloud_codearts_project.test.id
}

// filter by creator_id
data "huaweicloud_codearts_build_tasks" "filter_by_creator_id" {
  project_id = huaweicloud_codearts_project.test.id
  creator_id = data.huaweicloud_codearts_build_tasks.test.tasks[0].creator
}

output "is_creator_id_filter_useful" {
  value = length(data.huaweicloud_codearts_build_tasks.filter_by_creator_id.tasks) > 0
}
`, testBuildTask_basic(name))
}
