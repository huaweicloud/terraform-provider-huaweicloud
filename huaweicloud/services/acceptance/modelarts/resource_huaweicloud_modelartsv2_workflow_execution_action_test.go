package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkflowExecutionAction_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		rName          = "huaweicloud_modelartsv2_workflow_execution_action.test"
		rNameWithParam = "huaweicloud_modelartsv2_workflow_execution_action.test_with_param"
	)

	// Avoid CheckDestroy because this resource is a one-time action resource.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccWorkflowExecutionAction_basic_step1,
				ExpectError: regexp.MustCompile(`(?i)error creating workflow execution action:.+`),
			},
			{
				Config: testAccWorkflowExecutionAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "workflow_id"),
					resource.TestCheckResourceAttrSet(rName, "execution_id"),
					resource.TestCheckResourceAttr(rName, "action_name", "stop"),
				),
			},
			{
				Config: testAccWorkflowExecutionAction_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rNameWithParam, "workflow_id"),
					resource.TestCheckResourceAttrSet(rNameWithParam, "execution_id"),
				),
			},
		},
	})
}

func testAccWorkflowExecutionAction_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow" "test" {
  name        = "%[1]s"
  description = "created by terraform"

  steps {
    name  = "simple_step"
    type  = "job"
    title = "simple job step"

    inputs {
      name = "input_data"
      data = jsonencode("$ref/data_requirements/input_dataset")
      type = "dataset"
    }
    outputs {
      name = "output_data"
      type = "obs"
      config = jsonencode({
        obs_url = "$ref/storages/storage_name/output"
      })
    }

    properties = jsonencode({
      kind = "job"
      metadata = {
        name = "simple-job"
      }
      spec = {
        resource = {
          node_count = 1
          policy     = "regular"
          flavor     = "modelarts.vm.cpu.2u"
        }
      }
    })

    depend_steps = []
  }

  data_requirements {
    name       = "input_dataset"
    type       = "dataset"
    delay      = false
    used_steps = ["simple_step"]
  }

  parameters {
    name       = "test_param"
    type       = "str"
    default    = "\"test_value\""
    used_steps = ["simple_step"]
  }

  storages {
    name = "storage_name"
    type = "obs"
  }
}

resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = "%[1]s-execution"
  description = "test workflow execution"
  workflow_id = huaweicloud_modelartsv2_workflow.test.id

  data_requirements {
    name       = "input_dataset"
    type       = "dataset"
    value      = jsonencode({"dataset_id" = "test-dataset-id"})
    used_steps = ["simple_step"]
  }

  parameters {
    name        = "test_param"
    type        = "string"
    description = "test parameter"
    value       = jsonencode("execution-test-value")
  }

  lifecycle {
    ignore_changes = [
      data_requirements,
      parameters,
    ]
  }
}
`, name)
}

const testAccWorkflowExecutionAction_basic_step1 = `
resource "huaweicloud_modelartsv2_workflow_execution_action" "invalid_workflow" {
  workflow_id  = "invalid-workflow-id"
  execution_id = "invalid-execution-id"
  action_name  = "stop"
}
`

func testAccWorkflowExecutionAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_execution_action" "test" {
  workflow_id  = huaweicloud_modelartsv2_workflow.test.id
  execution_id = huaweicloud_modelartsv2_workflow_execution.test.id
  action_name  = "stop"
}
`, testAccWorkflowExecutionAction_base(name))
}

func testAccWorkflowExecutionAction_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelartsv2_workflow_execution_action" "test_with_param" {
  workflow_id  = huaweicloud_modelartsv2_workflow.test.id
  execution_id = huaweicloud_modelartsv2_workflow_execution.test.id
  action_name  = "rerun"

  data_requirements {
    name       = "data_input"
    type       = "OBS"
    value      = jsonencode({"obs_url": "obs://bucket/path"})
    used_steps = ["step1"]
    delay      = false
  }

  parameters {
    name        = "param1"
    type        = "String"
    description = "Test parameter"
    value       = jsonencode("test_value")
    used_steps  = ["step1"]
  }

  policies {
    rerun_steps = ["step1"]
  }
}
`, testAccWorkflowExecutionAction_base(name))
}
