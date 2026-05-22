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

func getV2WorkflowResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region     = acceptance.HW_REGION_NAME
		workflowId = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetV2WorkflowById(client, workflowId)
}

func TestAccV2Workflow_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelartsv2_workflow.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getV2WorkflowResourceFunc)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowSubscription(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2Workflow_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "terraform-test-workflow"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "user_name"),
					resource.TestCheckResourceAttrSet(rName, "param_ready"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Check steps
					resource.TestMatchResourceAttr(rName, "steps.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(rName, "steps.0.name", "condition_step"),
					resource.TestCheckResourceAttr(rName, "steps.0.type", "condition"),
					resource.TestCheckResourceAttr(rName, "steps.0.title", "condition_step"),
					resource.TestCheckResourceAttr(rName, "steps.0.conditions.#", "2"),
					resource.TestCheckResourceAttr(rName, "steps.0.if_then_steps.#", "2"),
					resource.TestCheckResourceAttr(rName, "steps.0.else_then_steps.#", "2"),
					resource.TestCheckResourceAttr(rName, "steps.1.name", "dataset_step"),
					resource.TestCheckResourceAttr(rName, "steps.1.type", "release_dataset"),
					resource.TestCheckResourceAttr(rName, "steps.2.name", "train_step"),
					resource.TestCheckResourceAttr(rName, "steps.2.type", "job"),
					resource.TestCheckResourceAttr(rName, "steps.3.name", "other_dataset_step"),
					resource.TestCheckResourceAttr(rName, "steps.4.name", "other_train_step"),
					// Check data_requirements
					resource.TestMatchResourceAttr(rName, "data_requirements.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(rName, "data_requirements.0.name", "input_dataset"),
					resource.TestCheckResourceAttr(rName, "data_requirements.0.type", "dataset"),
					resource.TestCheckResourceAttr(rName, "data_requirements.1.name", "other_input_dataset"),
					// Check parameters
					resource.TestMatchResourceAttr(rName, "parameters.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(rName, "parameters.0.name", "can_execute"),
					resource.TestCheckResourceAttr(rName, "parameters.0.type", "bool"),
					resource.TestCheckResourceAttr(rName, "parameters.0.default", "true"),
					// Check storages
					resource.TestMatchResourceAttr(rName, "storages.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(rName, "storages.0.name", "storage_name"),
					resource.TestCheckResourceAttr(rName, "storages.0.type", "obs"),
					// Check policy
					resource.TestMatchResourceAttr(rName, "policy.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestMatchResourceAttr(rName, "policy.0.scenes.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
				),
			},
			{
				Config: testAccV2Workflow_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "terraform-test-workflow"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "user_name"),
					// Check steps
					resource.TestMatchResourceAttr(rName, "steps.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					resource.TestCheckResourceAttr(rName, "steps.0.name", "condition_step"),
					resource.TestCheckResourceAttr(rName, "steps.0.type", "condition"),
					resource.TestCheckResourceAttr(rName, "steps.0.if_then_steps.#", "1"),
					resource.TestCheckResourceAttr(rName, "steps.0.else_then_steps.#", "1"),
					resource.TestCheckResourceAttr(rName, "steps.1.name", "dataset_step"),
					resource.TestCheckResourceAttr(rName, "steps.1.type", "release_dataset"),
					resource.TestCheckResourceAttr(rName, "steps.2.name", "other_dataset_step"),
					resource.TestCheckResourceAttr(rName, "steps.2.type", "release_dataset"),
					// Check data_requirements
					resource.TestMatchResourceAttr(rName, "data_requirements.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					// Check parameters
					resource.TestMatchResourceAttr(rName, "parameters.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
					// Check storages
					resource.TestMatchResourceAttr(rName, "storages.#", regexp.MustCompile(`^[1-9]([0-9]+)?$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"workspace_id",
					"data",
					"source_workflow_id",
					"gallery_subscription",
					"assets",
					"sub_graphs",
					"extend",
					"policy",
					"with_subscription",
					"subscription_id",
					"exeml_template_id",
					"package",
				},
			},
		},
	})
}

func testAccV2Workflow_basic_step1() string {
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
    {
      name       = "other_input_dataset"
      type       = "dataset"
      delay      = false
      used_steps = ["other_dataset_step"]
    },
  ]
}

variable "shared_algorithm_params" {
  type    = list(object({
    name  = string
    value = string
  }))
  default = [
    { name = "learning_rate_strategy", value = "$ref/parameters/learning_rate_strategy" },
    { name = "batch_size", value = "$ref/parameters/batch_size" },
    { name = "eval_batch_size", value = "$ref/parameters/eval_batch_size" },
    { name = "evaluate_every_n_epochs", value = "$ref/parameters/evaluate_every_n_epochs" },
    { name = "save_model_secs", value = "$ref/parameters/save_model_secs" },
    { name = "save_summary_steps", value = "$ref/parameters/save_summary_steps" },
    { name = "log_every_n_steps", value = "$ref/parameters/log_every_n_steps" },
    { name = "do_data_cleaning", value = "$ref/parameters/do_data_cleaning" },
  ]
}

variable "train_step_specific_params" {
  type    = list(object({
    name  = string
    value = string
  }))
  default = [
    { name = "jpeg_preprocess", value = "$ref/parameters/train_step-jpeg_preprocess" },
    { name = "do_eval_along_train", value = "$ref/parameters/train_step-do_eval_along_train" },
    { name = "use_fp16", value = "$ref/parameters/train_step-use_fp16" },
    { name = "task_type", value = "$ref/parameters/train_step-task_type" },
    { name = "data_format", value = "$ref/parameters/train_step-data_format" },
    { name = "do_model_analysis", value = "$ref/parameters/train_step-do_model_analysis" },
    { name = "model_name", value = "$ref/parameters/train_step-model_name" },
    { name = "best_model", value = "$ref/parameters/train_step-best_model" },
    { name = "xla_complie", value = "$ref/parameters/train_step-xla_complie" },
    { name = "variable_update", value = "$ref/parameters/train_step-variable_update" },
    { name = "do_train", value = "$ref/parameters/train_step-do_train" },
  ]
}

variable "other_train_step_specific_params" {
  type    = list(object({
    name  = string
    value = string
  }))
  default = [
    { name = "jpeg_preprocess", value = "$ref/parameters/other_train_step-jpeg_preprocess" },
    { name = "do_eval_along_train", value = "$ref/parameters/other_train_step-do_eval_along_train" },
    { name = "use_fp16", value = "$ref/parameters/other_train_step-use_fp16" },
    { name = "task_type", value = "$ref/parameters/other_train_step-task_type" },
    { name = "data_format", value = "$ref/parameters/other_train_step-data_format" },
    { name = "do_model_analysis", value = "$ref/parameters/other_train_step-do_model_analysis" },
    { name = "model_name", value = "$ref/parameters/other_train_step-model_name" },
    { name = "best_model", value = "$ref/parameters/other_train_step-best_model" },
    { name = "xla_complie", value = "$ref/parameters/other_train_step-xla_complie" },
    { name = "variable_update", value = "$ref/parameters/other_train_step-variable_update" },
    { name = "do_train", value = "$ref/parameters/other_train_step-do_train" },
  ]
}

variable "workflow_parameters" {
  type    = list(object({
    name       = string
    type       = string
    default    = optional(string)
    format     = optional(string)
    constraint = optional(string)
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
      name       = "train_flavor"
      type       = "json"
      format     = "train_flavor"
      constraint = "{\"device_distributed_mode\":[\"multiple\"], \"flavor_type\":[\"GPU\"]}"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "learning_rate_strategy"
      type       = "str"
      default    = "\"0.002\""
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "batch_size"
      type       = "int"
      default    = "64"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "eval_batch_size"
      type       = "int"
      default    = "64"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "evaluate_every_n_epochs"
      type       = "float"
      default    = "\"1.0\""
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "save_model_secs"
      type       = "int"
      default    = "60"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "save_summary_steps"
      type       = "int"
      default    = "10"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "log_every_n_steps"
      type       = "int"
      default    = "10"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "do_data_cleaning"
      type       = "bool"
      default    = "true"
      used_steps = ["train_step", "other_train_step"]
    },
    {
      name       = "other_placeholder_name"
      type       = "str"
      default    = "\"0.8\""
      used_steps = ["other_dataset_step"]
    },
    {
      name       = "train_step-jpeg_preprocess"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-do_eval_along_train"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-use_fp16"
      type       = "str"
      default    = "\"True\""
      title      = "use_fp16"
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-task_type"
      type       = "str"
      default    = "\"image_classification_v2\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-data_format"
      type       = "str"
      default    = "\"NCHW\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-do_model_analysis"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-model_name"
      type       = "str"
      default    = "\"resnet_v1_50\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-best_model"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-xla_complie"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-variable_update"
      type       = "str"
      default    = "\"horovod\""
      used_steps = ["train_step"]
    },
    {
      name       = "train_step-do_train"
      type       = "str"
      default    = "\"True\""
      used_steps = ["train_step"]
    },
    {
      name       = "other_train_step-jpeg_preprocess"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-do_eval_along_train"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-use_fp16"
      type       = "str"
      default    = "\"True\""
      title      = "use_fp16"
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-task_type"
      type       = "str"
      default    = "\"image_classification_v2\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-data_format"
      type       = "str"
      default    = "\"NCHW\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-do_model_analysis"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-model_name"
      type       = "str"
      default    = "\"resnet_v1_50\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-best_model"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-xla_complie"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-variable_update"
      type       = "str"
      default    = "\"horovod\""
      used_steps = ["other_train_step"]
    },
    {
      name       = "other_train_step-do_train"
      type       = "str"
      default    = "\"True\""
      used_steps = ["other_train_step"]
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
      steps = ["condition_step", "dataset_step", "train_step"]
    },
    {
      name  = "other_data_path"
      steps = ["condition_step", "other_dataset_step", "other_train_step"]
    },
  ]
}

resource "huaweicloud_modelartsv2_workflow" "test" {
  name        = "terraform-test-workflow"
  description = "created by terraform"

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
    if_then_steps   = ["dataset_step", "train_step"]
    else_then_steps = ["other_dataset_step", "other_train_step"]
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
    name  = "train_step"
    type  = "job"
    title = "picture classify"

    inputs {
      name = "data_url"
      data = jsonencode("$ref/consumptions/dataset_step/output_name")
      type = "dataset"
    }
    outputs {
      name   = "train_url"
      type   = "obs"
      config = jsonencode({
        obs_url = "$ref/storages(with_execution_id=True, create_dir=True)/storage_name/directory_path"
      })
    }

    properties = jsonencode({
      kind = "job"

      metadata = {
        name = "workflow-390f1624ac8f4ac8f4dc8a76527cd99fd78f3"
      }
      spec     = {
        resource = {
          node_count = 1
          policy     = "regular"
          flavor     = "$ref/parameters/train_flavor"
        }
      }
      algorithm = {
        subscription_id = "%[1]s"
        item_version_id = "%[2]s"

        parameters = concat(var.shared_algorithm_params, var.train_step_specific_params)
      }
    })

    depend_steps = ["condition_step", "dataset_step"]
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
  steps {
    name  = "other_train_step"
    type  = "job"
    title = "other picture classify"

    inputs {
      name = "data_url"
      data = jsonencode("$ref/consumptions/dataset_step/output_name")
      type = "dataset"
    }
    inputs {
      name = "other_data_url"
      data = jsonencode("$ref/consumptions/other_dataset_step/output_name")
      type = "dataset"
    }
    outputs {
      name   = "train_url"
      type   = "obs"
      config = jsonencode({
        obs_url = "$ref/storages(with_execution_id=True, create_dir=True)/storage_name/directory_path"
      })
    }
    outputs {
      name   = "other_train_url"
      type   = "obs"
      config = jsonencode({
        obs_url = "$ref/storages(with_execution_id=True, create_dir=True)/other_storage_name/directory_path"
      })
    }

    properties = jsonencode({
      kind = "job"

      metadata  = {
        name = "workflow-390f1624ac8f4ac8f4dc8a76527cd99fd78f3"
      }
      spec      = {
        resource = {
          node_count = 1
          policy     = "regular"
          flavor     = "$ref/parameters/train_flavor"
        }
      }
      algorithm = {
        subscription_id = "%[1]s"
        item_version_id = "%[2]s"

        parameters = concat(var.shared_algorithm_params, var.other_train_step_specific_params)
      }
    })

    depend_steps = ["condition_step", "other_dataset_step"]
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
      format     = parameters.value.format
      constraint = parameters.value.constraint
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
`, acceptance.HW_MODELARTS_WORKFLOW_SUBSCRIPTION_ID, acceptance.HW_MODELARTS_WORKFLOW_ITEM_VERSION_ID)
}

func testAccV2Workflow_basic_step2() string {
	return `
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

resource "huaweicloud_modelartsv2_workflow" "test" {
  name        = "terraform-test-workflow"
  description = "created by terraform"

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
}
`
}
