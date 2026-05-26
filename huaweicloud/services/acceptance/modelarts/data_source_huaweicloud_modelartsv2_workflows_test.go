package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV2Workflows_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelartsv2_workflows.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_modelartsv2_workflows.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byDescription   = "data.huaweicloud_modelartsv2_workflows.filter_by_description"
		dcByDescription = acceptance.InitDataSourceCheck(byDescription)

		byLabels   = "data.huaweicloud_modelartsv2_workflows.filter_by_labels"
		dcByLabels = acceptance.InitDataSourceCheck(byLabels)

		bySearchType   = "data.huaweicloud_modelartsv2_workflows.filter_by_search_type"
		dcBySearchType = acceptance.InitDataSourceCheck(bySearchType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV2Workflows_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "workflows.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "workflows.0.id"),
					resource.TestCheckResourceAttrSet(all, "workflows.0.name"),
					resource.TestCheckResourceAttrSet(all, "workflows.0.description"),
					resource.TestMatchResourceAttr(all, "workflows.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "workflows.0.user_name"),
					resource.TestCheckResourceAttrSet(all, "workflows.0.param_ready"),
					resource.TestMatchResourceAttr(all, "workflows.0.last_modified_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "workflows.0.parameters.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttrSet(all, "workflows.0.parameters.0.name"),
					resource.TestCheckResourceAttrSet(all, "workflows.0.parameters.0.type"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					dcByLabels.CheckResourceExists(),
					resource.TestCheckOutput("is_labels_filter_useful", "true"),
					dcBySearchType.CheckResourceExists(),
					resource.TestCheckOutput("is_search_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataV2Workflow_base = `
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
    {
      name       = "other_input_dataset"
      type       = "dataset"
      delay      = false
      used_steps = ["other_dataset_step"]
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
      name       = "other_route_path"
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
    {
      name       = "other_placeholder_name"
      type       = "str"
      default    = "\"0.8\""
      used_steps = ["other_dataset_step"]
    },
  ]
}

variable "workflow_storages" {
  type    = list(object({
    name = string
    type = string
  }))
  default = [
    {
      name = "storage_name"
      type = "obs"
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
      name  = "data_path"
      steps = ["condition_step", "dataset_step"]
    },
    {
      name  = "other_data_path"
      steps = ["condition_step", "other_dataset_step"]
    },
  ]
}

resource "huaweicloud_modelartsv2_workflow" "test" {
  count = 2

  name        = format("terraform-test-workflow-%d", count.index)
  description = format("created by terraform %d", count.index)
  labels      = [format("terraform-%d", count.index)]

  steps {
    name  = "condition_step"
    type  = "condition"
    title = "condition_step"

    conditions {
      type  = "=="
      left  = jsonencode("$ref/parameters/can_execute")
      right = jsonencode(true)
    }
    conditions {
      type  = "=="
      left  = jsonencode("$ref/parameters/other_route_path")
      right = jsonencode(true)
    }

    depend_steps    = []
    if_then_steps   = ["dataset_step"]
    else_then_steps = ["other_dataset_step"]
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
    properties   = jsonencode({
      version_format              = "Default"
      train_evaluate_sample_ratio = "$ref/parameters/placeholder_name"
      clear_hard_property         = true
      remove_sample_usage         = true
      label_task_type             = 0
    })

    depend_steps = []
  }
  steps {
    name  = "other_dataset_step"
    type  = "release_dataset"
    title = "other dataset release"

    inputs {
      name = "other_input_name"
      data = jsonencode("$ref/data_requirements/other_input_dataset")
      type = "dataset"
    }
    outputs {
      name = "output_name"
      type = "dataset"
    }

    properties = jsonencode({
      version_format              = "Default"
      train_evaluate_sample_ratio = "$ref/parameters/other_placeholder_name"
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

  dynamic "storages" {
    for_each = var.workflow_storages

    content {
      name = storages.value.name
      type = storages.value.type
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
`

func testAccDataV2Workflows_basic() string {
	return fmt.Sprintf(`
%[1]s

# Query all workflows without any filter.
data "huaweicloud_modelartsv2_workflows" "all" {
  depends_on = [huaweicloud_modelartsv2_workflow.test]
}

# Filter workflows by name.
locals {
  name = huaweicloud_modelartsv2_workflow.test[0].name
}

data "huaweicloud_modelartsv2_workflows" "filter_by_name" {
  name = local.name

  depends_on = [huaweicloud_modelartsv2_workflow.test]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflows.filter_by_name.workflows : v.name == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter workflows by description.
locals {
  description = huaweicloud_modelartsv2_workflow.test[0].description
}

data "huaweicloud_modelartsv2_workflows" "filter_by_description" {
  description = local.description

  depends_on = [huaweicloud_modelartsv2_workflow.test]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflows.filter_by_description.workflows : v.description == local.description
  ]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}

# Filter workflows by labels.
locals {
  labels = huaweicloud_modelartsv2_workflow.test[0].labels
}

data "huaweicloud_modelartsv2_workflows" "filter_by_labels" {
  labels = local.labels

  depends_on = [huaweicloud_modelartsv2_workflow.test]
}

locals {
  labels_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflows.filter_by_labels.workflows :
    length(setintersection(v.labels, local.labels)) > 0
  ]
}

output "is_labels_filter_useful" {
  value = length(local.labels_filter_result) > 0 && alltrue(local.labels_filter_result)
}

# Filter workflows by search_type.
locals {
  search_type = "equal"
}

data "huaweicloud_modelartsv2_workflows" "filter_by_search_type" {
  name        = local.name
  search_type = local.search_type

  depends_on = [huaweicloud_modelartsv2_workflow.test]
}

locals {
  search_type_filter_result = [
    for v in data.huaweicloud_modelartsv2_workflows.filter_by_search_type.workflows : local.name == v.name
  ]
}

output "is_search_type_filter_useful" {
  value = length(local.search_type_filter_result) == 1 && alltrue(local.search_type_filter_result)
}
`, testAccDataV2Workflow_base)
}
