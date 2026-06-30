package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getV2WorkflowExecutionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetV2WorkflowExecutionById(client, state.Primary.Attributes["workflow_id"], state.Primary.ID)
}

func TestAccV2WorkflowExecution_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelartsv2_workflow_execution.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV2WorkflowExecutionResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccV2WorkflowExecution_basic_step1(name),
				ExpectError: regexp.MustCompile(`error creating workflow execution`),
			},
			{
				Config: testAccV2WorkflowExecution_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "test workflow execution"),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccV2WorkflowExecution_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestMatchResourceAttr(rName, "parameters.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.name"),
					resource.TestCheckResourceAttrSet(rName, "parameters.0.type"),
					resource.TestMatchResourceAttr(rName, "data_requirements.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(rName, "data_requirements.0.name"),
					resource.TestCheckResourceAttrSet(rName, "data_requirements.0.type"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccV2WorkflowExecutionImportStateIDFunc(rName),
				ImportStateVerifyIgnore: []string{
					"name",
					"description",
					"workspace_id",
					"workflow_id",
					"workflow_name",
					"scene_id",
					"scene_name",
					"policies",
					"duration",
				},
			},
		},
	})
}

func testAccV2WorkflowExecution_basic_step1(name string) string {
	randomUUID, _ := uuid.NewRandom()

	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = "%[1]s"
  description = "test workflow execution"
  workflow_id = "%[2]s"
}
`, name, randomUUID.String())
}

func testAccV2WorkflowExecution_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = "%[1]s"
  description = "test workflow execution"
  workflow_id = "%[2]s"

  lifecycle {
    ignore_changes = [
      data_requirements,
      parameters,
    ]
  }
}
`, name, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccV2WorkflowExecution_basic_step3(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = "%[1]s"
  description = "test workflow execution"
  workflow_id = "%[2]s"

  parameters {
    name        = "param1"
    type        = "string"
    description = "test parameter"
    value       = jsonencode("test-value")
  }

  data_requirements {
    name       = "data1"
    type       = "obs"
    value      = jsonencode({"path" = "obs://bucket/path"})
    used_steps = ["step1"]
  }

  lifecycle {
    ignore_changes = [
      data_requirements,
      parameters,
    ]
  }
}
`, name, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccV2WorkflowExecutionImportStateIDFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		workflowID := rs.Primary.Attributes["workflow_id"]
		if workflowID == "" {
			return "", fmt.Errorf("attribute (workflow_id) of resource (%s) not found", name)
		}
		return fmt.Sprintf("%s/%s", workflowID, rs.Primary.ID), nil
	}
}
