package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV2WorkflowExecutions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_modelartsv2_workflow_executions.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byStatus   = "data.huaweicloud_modelartsv2_workflow_executions.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byLabels   = "data.huaweicloud_modelartsv2_workflow_executions.filter_by_labels"
		dcByLabels = acceptance.InitDataSourceCheck(byLabels)

		bySceneId   = "data.huaweicloud_modelartsv2_workflow_executions.filter_by_scene_id"
		dcBySceneId = acceptance.InitDataSourceCheck(bySceneId)

		byWorkspaceId   = "data.huaweicloud_modelartsv2_workflow_executions.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV2WorkflowExecutions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "executions.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "executions.0.id"),
					resource.TestCheckResourceAttrSet(all, "executions.0.name"),
					resource.TestCheckResourceAttrSet(all, "executions.0.status"),
					resource.TestCheckResourceAttrSet(all, "executions.0.workflow_id"),
					resource.TestCheckResourceAttrSet(all, "executions.0.workflow_name"),
					resource.TestCheckResourceAttrSet(all, "executions.0.created_at"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByLabels.CheckResourceExists(),
					resource.TestCheckOutput("is_labels_filter_useful", "true"),
					dcBySceneId.CheckResourceExists(),
					resource.TestCheckOutput("is_scene_id_filter_useful", "true"),
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV2WorkflowExecution_base(name string) string {
	return fmt.Sprintf(`
variable "workflow_data_requirements" {
  type    = list(object({
    name       = string
    type       = string
    delay      = optional(bool, false)
    used_steps = list(string)
  }))
  default = [
    {
      name       = "input_dataset"
      type       = "dataset"
      delay      = false
      used_steps = ["dataset_step"]
    },
  ]
}

variable "workflow_parameters" {
  type    = list(object({
    name       = string
    type       = string
    default    = optional(string)
    used_steps = list(string)
  }))
  default = [
    {
      name       = "can_execute"
      type       = "bool"
      default    = "true"
      used_steps = ["condition_step"]
    },
    {
      name       = "placeholder_name"
      type       = "str"
      default    = "\"0.8\""
      used_steps = ["dataset_step"]
    },
  ]
}

variable "workflow_policy_scenes" {
  type    = list(object({
    id    = optional(string)
    name  = string
    steps = list(string)
  }))
  default = [
    {
      name  = "scene_a"
      steps = ["condition_step", "dataset_step"]
    },
    {
      name  = "scene_b"
      steps = ["condition_step"]
    },
  ]
}

resource "huaweicloud_modelartsv2_workflow" "test" {
  name        = "%[1]s"
  description = "created by terraform for execution test"
  labels      = ["terraform"]

  steps {
    name  = "condition_step"
    type  = "condition"
    title = "condition_step"

    conditions {
      type  = "=="
      left  = jsonencode("$ref/parameters/can_execute")
      right = jsonencode(true)
    }

    depend_steps    = []
    if_then_steps   = ["dataset_step"]
    else_then_steps = []
  }
  steps {
    name  = "dataset_step"
    type  = "release_dataset"
    title = "dataset release"

    inputs {
      name = "input_name"
      data = jsonencode("$ref/data_requirements/input_dataset")
      type = "dataset"
    }
    outputs {
      name = "output_name"
      type = "dataset"
    }
    properties = jsonencode({
      version_format              = "Default"
      train_evaluate_sample_ratio = "$ref/parameters/placeholder_name"
      clear_hard_property         = true
      remove_sample_usage         = true
      label_task_type             = 0
    })

    depend_steps = []
  }

  dynamic "data_requirements" {
    for_each = var.workflow_data_requirements

    content {
      name       = data_requirements.value.name
      type       = data_requirements.value.type
      delay      = data_requirements.value.delay
      used_steps = data_requirements.value.used_steps
    }
  }

  dynamic "parameters" {
    for_each = var.workflow_parameters

    content {
      name       = parameters.value.name
      type       = parameters.value.type
      default    = parameters.value.default
      used_steps = parameters.value.used_steps
    }
  }

  policy {
    dynamic "scenes" {
      for_each = var.workflow_policy_scenes

      content {
        id    = scenes.value.id
        name  = scenes.value.name
        steps = scenes.value.steps
      }
    }
  }
}

variable "workflow_executions" {
  type    = list(object({
    name        = string
    description = string
    scene_name  = string
    labels      = list(string)
  }))
  default = [
    {
      name        = "terraform-test-execution-scene-a"
      description = "created by terraform with scene_a"
      scene_name  = "scene_a"
      labels      = ["scene-a", "test"]
    },
    {
      name        = "terraform-test-execution-scene-b"
      description = "created by terraform with scene_b"
      scene_name  = "scene_b"
      labels      = ["scene-b", "test"]
    },
  ]
}

resource "huaweicloud_modelartsv2_workflow_execution" "test_scene_a" {
  name        = var.workflow_executions[0].name
  description = var.workflow_executions[0].description
  workflow_id = huaweicloud_modelartsv2_workflow.test.id
  scene_name  = var.workflow_executions[0].scene_name
  labels      = var.workflow_executions[0].labels

  depends_on = [huaweicloud_modelartsv2_workflow.test]

  lifecycle {
    ignore_changes = [
      data_requirements,
      parameters,
    ]
  }
}

resource "huaweicloud_modelartsv2_workflow_execution" "test_scene_b" {
  name        = var.workflow_executions[1].name
  description = var.workflow_executions[1].description
  workflow_id = huaweicloud_modelartsv2_workflow.test.id
  scene_name  = var.workflow_executions[1].scene_name
  labels      = var.workflow_executions[1].labels

  depends_on = [huaweicloud_modelartsv2_workflow_execution.test_scene_a]

  lifecycle {
    ignore_changes = [
      data_requirements,
      parameters,
    ]
  }
}
`, name)
}

func testAccDataV2WorkflowExecutions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Query all workflow executions without any filter.
data "huaweicloud_modelartsv2_workflow_executions" "all" {
  workflow_id = huaweicloud_modelartsv2_workflow.test.id

  depends_on = [
    huaweicloud_modelartsv2_workflow_execution.test_scene_a,
    huaweicloud_modelartsv2_workflow_execution.test_scene_b,
  ]
}

# Filter workflow executions by status.
locals {
  status = try(data.huaweicloud_modelartsv2_workflow_executions.all.executions[0].status, "NOT_FOUND")
}

data "huaweicloud_modelartsv2_workflow_executions" "filter_by_status" {
  workflow_id = huaweicloud_modelartsv2_workflow.test.id
  status      = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflow_executions.filter_by_status.executions : v.status == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 || alltrue(local.status_filter_result)
}

# Filter workflow executions by labels.
locals {
  labels_list = try(data.huaweicloud_modelartsv2_workflow_executions.all.executions[0].labels, [])
  labels      = join(",", local.labels_list)
}

data "huaweicloud_modelartsv2_workflow_executions" "filter_by_labels" {
  workflow_id = huaweicloud_modelartsv2_workflow.test.id
  labels      = local.labels
}

locals {
  labels_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflow_executions.filter_by_labels.executions :
    length(setintersection(v.labels, local.labels_list)) > 0
  ]
}

output "is_labels_filter_useful" {
  value = length(local.labels_filter_result) > 0 || alltrue(local.labels_filter_result)
}

# Filter workflow executions by scene_id.
locals {
  scene_id = try(data.huaweicloud_modelartsv2_workflow_executions.all.executions[0].scene_id, "NOT_FOUND")
}

data "huaweicloud_modelartsv2_workflow_executions" "filter_by_scene_id" {
  workflow_id = huaweicloud_modelartsv2_workflow.test.id
  scene_id    = local.scene_id
}

locals {
  scene_id_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflow_executions.filter_by_scene_id.executions :
    v.scene_id == local.scene_id
  ]
}

output "is_scene_id_filter_useful" {
  value = length(local.scene_id_filter_result) > 0 || alltrue(local.scene_id_filter_result)
}

# Filter workflow executions by workspace_id.
locals {
  workspace_id = try(data.huaweicloud_modelartsv2_workflow_executions.all.executions[0].workspace_id, "NOT_FOUND")
}

data "huaweicloud_modelartsv2_workflow_executions" "filter_by_workspace_id" {
  workflow_id  = huaweicloud_modelartsv2_workflow.test.id
  workspace_id = local.workspace_id
}

locals {
  workspace_id_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflow_executions.filter_by_workspace_id.executions :
    v.workspace_id == local.workspace_id
  ]
}

output "is_workspace_id_filter_useful" {
  value = length(local.workspace_id_filter_result) > 0 || alltrue(local.workspace_id_filter_result)
}
`, testAccDataV2WorkflowExecution_base(name))
}
