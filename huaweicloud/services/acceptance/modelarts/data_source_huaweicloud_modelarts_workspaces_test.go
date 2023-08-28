package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceWorkspaces_basic(t *testing.T) {
	rName := "data.huaweicloud_modelarts_workspaces.name_filter"
	dc := acceptance.InitDataSourceCheck(rName)
	name := acceptance.RandomAccResourceName()
	name2 := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWorkspaces_basic(name, name2),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.id"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.name"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.auth_type"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.description"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.owner"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "workspaces.0.status"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceWorkspaces_basic(name, name2 string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%[1]s"
  description = "%[1]s"
}

resource "huaweicloud_modelarts_workspace" "test2" {
  name        = "%[2]s"
  description = "%[2]s"
}

data "huaweicloud_modelarts_workspaces" "name_filter" {
  name = "%[1]s"
  depends_on = [huaweicloud_modelarts_workspace.test, huaweicloud_modelarts_workspace.test2]
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_workspaces.name_filter.workspaces) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_workspaces.name_filter.workspaces[*].name : v == "%[1]s"]
  )  
}

data "huaweicloud_modelarts_workspaces" "enterprise_project_id_filter" {
  enterprise_project_id = "0"
  depends_on = [huaweicloud_modelarts_workspace.test, huaweicloud_modelarts_workspace.test2]
}
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_modelarts_workspaces.enterprise_project_id_filter.workspaces) > 0 && alltrue(
    [for v in data.huaweicloud_modelarts_workspaces.enterprise_project_id_filter.workspaces[*].enterprise_project_id : v == "0"]
  )
}
`, name, name2)
}
