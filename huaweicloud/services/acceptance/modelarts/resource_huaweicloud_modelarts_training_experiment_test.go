package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getTrainingExperimentResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetTrainingExperimentById(client, state.Primary.ID)
}

func TestAccTrainingExperiment_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_modelarts_training_experiment.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getTrainingExperimentResourceFunc)

		rNameWithWorkspaceId = "huaweicloud_modelarts_training_experiment.test_with_workspace_id"
		rcWithWorkspaceId    = acceptance.InitResourceCheck(rNameWithWorkspaceId, &obj, getTrainingExperimentResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcWithWorkspaceId.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTrainingExperiment_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttrSet(rName, "metadata.0.workspace_id"),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", "acceptance test create"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcWithWorkspaceId.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rNameWithWorkspaceId, "metadata.0.workspace_id",
						"huaweicloud_modelarts_workspace.test", "id"),
					resource.TestCheckResourceAttr(rNameWithWorkspaceId, "metadata.0.description", ""),
				),
			},
			{
				Config: testAccTrainingExperiment_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", updateName),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", ""),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					rcWithWorkspaceId.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameWithWorkspaceId, "metadata.0.name", updateName),
					resource.TestCheckResourceAttrPair(rNameWithWorkspaceId, "metadata.0.workspace_id",
						"huaweicloud_modelarts_workspace.test", "id"),
					resource.TestCheckResourceAttr(rNameWithWorkspaceId, "metadata.0.description", "Updated by terraform script"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTrainingExperiment_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name      = "%[1]s"
  auth_type = "private"
}
`, name)
}

func testAccTrainingExperiment_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name        = "%[2]s"
    description = "acceptance test create"
  }
}

resource "huaweicloud_modelarts_training_experiment" "test_with_workspace_id" {
  metadata {
    name         = "%[2]s"
    workspace_id = huaweicloud_modelarts_workspace.test.id
  }
}
`, testAccTrainingExperiment_base(name), name)
}

func testAccTrainingExperiment_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name = "%[2]s"
  }
}

resource "huaweicloud_modelarts_training_experiment" "test_with_workspace_id" {
  metadata {
    name         = "%[2]s"
    workspace_id = huaweicloud_modelarts_workspace.test.id
    description  = "Updated by terraform script"
  }
}
`, testAccTrainingExperiment_base(name), updateName)
}
