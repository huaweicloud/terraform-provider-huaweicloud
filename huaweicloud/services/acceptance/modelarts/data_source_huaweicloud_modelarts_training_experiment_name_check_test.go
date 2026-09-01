package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrainingExperimentNameCheck_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_modelarts_training_experiment_name_check.test"
		dc  = acceptance.InitDataSourceCheck(all)

		duplicateName   = "data.huaweicloud_modelarts_training_experiment_name_check.duplicate"
		dcDuplicateName = acceptance.InitDataSourceCheck(duplicateName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTrainingExperimentNameCheck_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Check with a non-duplicate name.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "is_duplicate", "false"),

					// Check with a duplicate name (the name of the created experiment).
					dcDuplicateName.CheckResourceExists(),
					resource.TestCheckResourceAttr(duplicateName, "is_duplicate", "true"),
				),
			},
		},
	})
}

func testAccDataTrainingExperimentNameCheck_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name        = "%[1]s"
    description = "Created by terraform script"
  }
}

# Check with a non-duplicate name.
data "huaweicloud_modelarts_training_experiment_name_check" "test" {
  experiment_name = "%[1]s-non-duplicate"
  depends_on      = [huaweicloud_modelarts_training_experiment.test]
}

# Check with a duplicate name (the name of the created experiment).
data "huaweicloud_modelarts_training_experiment_name_check" "duplicate" {
  experiment_name = huaweicloud_modelarts_training_experiment.test.metadata[0].name
  depends_on      = [huaweicloud_modelarts_training_experiment.test]
}
`, name)
}
