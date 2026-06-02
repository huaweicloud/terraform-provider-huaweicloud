package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrainingExperiments_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_modelarts_training_experiments.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byWorkspaceId   = "data.huaweicloud_modelarts_training_experiments.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		bySortByDesc   = "data.huaweicloud_modelarts_training_experiments.filtered_by_sort_by_desc"
		dcBySortByDesc = acceptance.InitDataSourceCheck(bySortByDesc)

		bySortByAsc   = "data.huaweicloud_modelarts_training_experiments.filtered_by_sort_by_asc"
		dcBySortByAsc = acceptance.InitDataSourceCheck(bySortByAsc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTrainingExperiments_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "training_experiments.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),

					// Filter by 'workspace_id' parameter.
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),

					// Filter by 'sort_by' and 'order' parameters.
					dcBySortByDesc.CheckResourceExists(),
					dcBySortByAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_by_and_order_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(byWorkspaceId, "training_experiments.0.metadata.0.id"),
					resource.TestCheckResourceAttrSet(byWorkspaceId, "training_experiments.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(byWorkspaceId, "training_experiments.0.metadata.0.description"),
					resource.TestCheckResourceAttrSet(byWorkspaceId, "training_experiments.0.metadata.0.workspace_id"),
					resource.TestMatchResourceAttr(byWorkspaceId, "training_experiments.0.metadata.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byWorkspaceId, "training_experiments.0.metadata.0.update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataTrainingExperiments_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name      = "%[1]s"
  auth_type = "private"
}

resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name         = "%[1]s"
    description  = "Created by terraform script"
    workspace_id = huaweicloud_modelarts_workspace.test.id
  }
}
`, name)
}

func testAccDataTrainingExperiments_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_modelarts_training_experiments" "test" {
  depends_on = [huaweicloud_modelarts_training_experiment.test]
}

# Filter by 'workspace_id' parameter.
data "huaweicloud_modelarts_training_experiments" "filter_by_workspace_id" {
  workspace_id = huaweicloud_modelarts_workspace.test.id

  depends_on = [huaweicloud_modelarts_training_experiment.test]
}

locals {
  workspace_id_filter_result = [for v in data.huaweicloud_modelarts_training_experiments.filter_by_workspace_id.training_experiments :
    v.metadata[0].workspace_id == huaweicloud_modelarts_workspace.test.id
  ]
}

output "is_workspace_id_filter_useful" {
  value = length(local.workspace_id_filter_result) > 0 && alltrue(local.workspace_id_filter_result)
}

# Filter by 'sort_by' and 'order' parameters.
data "huaweicloud_modelarts_training_experiments" "filtered_by_sort_by_desc" {
  sort_by = "name"

  depends_on = [huaweicloud_modelarts_training_experiment.test]
}

data "huaweicloud_modelarts_training_experiments" "filtered_by_sort_by_asc" {
  sort_by = "name"
  order   = "asc"

  depends_on = [huaweicloud_modelarts_training_experiment.test]
}

locals {
  sort_by_and_order_desc_result = data.huaweicloud_modelarts_training_experiments.filtered_by_sort_by_desc.training_experiments[*].metadata[0].name
  sort_by_and_order_asc_result  = data.huaweicloud_modelarts_training_experiments.filtered_by_sort_by_asc.training_experiments[*].metadata[0].name
}

output "is_sort_by_and_order_useful" {
  value = local.sort_by_and_order_desc_result == reverse(local.sort_by_and_order_asc_result)
}
`, testAccDataTrainingExperiments_base(name))
}
