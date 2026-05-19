package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getAlgorithmResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetAlgorithmById(client, state.Primary.ID)
}

func TestAccAlgorithm_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_modelarts_algorithm.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getAlgorithmResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsAlgorithm(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAlgorithm_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "metadata.0.workspace_id", "0"),
					resource.TestCheckResourceAttr(rName, "metadata.0.tags.#", "1"),
					resource.TestCheckResourceAttr(rName, "metadata.0.tags.0.key", "auto_search"),
					resource.TestCheckResourceAttr(rName, "job_config.0.code_dir", fmt.Sprintf("/%s/algorithm/", name)),
					resource.TestCheckResourceAttr(rName, "job_config.0.boot_file", fmt.Sprintf("/%s/algorithm/bootfile.py", name)),
					resource.TestCheckResourceAttr(rName, "job_config.0.engine.0.engine_id", acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID),
					resource.TestCheckResourceAttr(rName, "job_config.0.engine.0.engine_name", acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME),
					resource.TestCheckResourceAttr(rName, "job_config.0.engine.0.engine_version", acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters_customization", "true"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.name", "parameter2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.value", "10"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.description", "parameter2 description"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.constraint.0.type", "Float"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.constraint.0.valid_type", "None"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.constraint.0.editable", "false"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.constraint.0.required", "false"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.name", "parameter1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.description", "parameter1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.constraint.0.required", "true"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.constraint.0.editable", "true"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.constraint.0.valid_range.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.constraint.0.valid_range.0", "10"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.1.constraint.0.valid_range.1", "50"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.name", "input1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.description", "Input parameter1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.access_method", "parameter"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.remote_constraints.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.remote_constraints.0.data_type", "obs"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.remote_constraints.1.data_type", "modelarts_dataset"),
					resource.TestCheckResourceAttrSet(rName, "job_config.0.inputs.0.remote_constraints.1.attributes"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.1.name", "a_input2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.1.access_method", "env"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.1.remote_constraints.#", "1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.1.remote_constraints.0.data_type", "obs"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.name", "output1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.description", "Output parameter1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.access_method", "env"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.1.name", "aoutput2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.1.access_method", "parameter"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.#", "3"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.key", "flavor_type"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.operator", "in"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.values.#", "2"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.1.key", "device_distributed_mode"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.1.values.#", "1"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.1.values.0", "singular"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.2.key", "host_distributed_mode"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.2.values.#", "1"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.2.values.0", "multiple"),
					resource.TestCheckResourceAttr(rName, "advanced_config.#", "1"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.#", "1"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.0.name", "search"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.0.mode", "max"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.0.regex", "10.0"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.search_params.#", "2"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.algo_configs.0.name", "bayes_opt_search"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.algo_configs.0.params.#", "2"),
				),
			},
			{
				Config: testAccAlgorithm_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttrPair(rName, "metadata.0.workspace_id", "huaweicloud_modelarts_workspace.test", "id"),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", ""),
					resource.TestCheckResourceAttr(rName, "job_config.0.code_dir", fmt.Sprintf("/%s/algorithm/", name)),
					resource.TestCheckResourceAttr(rName, "job_config.0.boot_file", ""),
					resource.TestCheckResourceAttr(rName, "job_config.0.command", "bash ${MA_JOB_DIR}/algorithm/bootfile.py"),
					resource.TestCheckResourceAttr(rName, "job_config.0.engine.0.image_url", acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters_customization", "false"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.name", "parameter_update"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.description", ""),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.value", "40"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.0.constraint.0.valid_range.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.#", "2"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.name", "input1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.description", ""),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.remote_constraints.#", "1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.0.remote_constraints.0.data_type", "modelarts_dataset"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.1.remote_constraints.#", "0"),
					resource.TestCheckResourceAttrSet(rName, "job_config.0.inputs.0.remote_constraints.0.attributes"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.#", "1"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.name", "output_update"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.description", ""),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.0.access_method", "env"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.#", "3"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.key", "flavor_type"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.values.#", "1"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.0.values.0", "Ascend"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.1.key", "device_distributed_mode"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.1.values.0", "multiple"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.2.key", "host_distributed_mode"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.2.values.0", "singular"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.#", "2"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.1.name", "a_search"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.reward_attrs.1.mode", "min"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.search_params.0.name", "parameter_update"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.algo_configs.0.name", "anneal_search"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.algo_configs.0.params.0.key", "avg_best_idx"),
					resource.TestCheckResourceAttr(rName, "advanced_config.0.auto_search.0.algo_configs.0.params.0.value", "2.0"),
				),
			},
			{
				Config: testAccAlgorithm_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(rName, "metadata.0.tags.#", "0"),
					resource.TestCheckResourceAttr(rName, "job_config.0.inputs.#", "0"),
					resource.TestCheckResourceAttr(rName, "job_config.0.outputs.#", "0"),
					resource.TestCheckResourceAttr(rName, "job_config.0.parameters.#", "0"),
					resource.TestCheckResourceAttr(rName, "resource_requirements.#", "0"),
					resource.TestCheckResourceAttr(rName, "advanced_config.#", "0"),
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

func testAccAlgorithm_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "%[1]s"
  acl    = "private"
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket  = huaweicloud_obs_bucket.test.bucket
  key     = "algorithm/bootfile.py"
  content = <<EOF
#!/usr/bin/env python
import os
def main():
    os.environ['CUDA_VISIBLE_DEVICES'] = '0'
if __name__ == '__main__':
    main()
EOF
  content_type = "text/py"
}

resource "huaweicloud_modelarts_workspace" "test" {
  name = "%[1]s"
}
`, name)
}

func testAccAlgorithm_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name        = "%[2]s"
    description = "Created by terraform script"

    tags {
      key = "auto_search"
    }
  }

  job_config {
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/algorithm/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    engine {
      engine_id      = "%[3]s"
      engine_name    = "%[4]s"
      engine_version = "%[3]s"
    }

    parameters_customization = true

    parameters {
      name        = "parameter2"
      description = "parameter2 description"
      value       = "10"

      constraint {
        type       = "Float"
        valid_type = "None"
      }
    }
    parameters {
      name        = "parameter1"
      description = "parameter1"

      constraint {
        type        = "Float"
        editable    = true
        required    = true
        valid_range = ["10", "50"]
      }
    }

    inputs {
      name          = "input1"
      description   = "Input parameter1"
      access_method = "parameter"

      remote_constraints {
        data_type = "obs"
      }
      remote_constraints {
        data_type  = "modelarts_dataset"
        attributes = jsonencode({
          data_format       = ["CarbonData"]
          data_segmentation = ["false"]
          dataset_type      = ["0"]
        })
      }
    }
    inputs {
      name          = "a_input2"
      access_method = "env"

      remote_constraints {
        data_type = "obs"
      }
    }

    outputs {
      name          = "output1"
      description   = "Output parameter1"
      access_method = "env"
    }
    outputs {
      name          = "aoutput2"
      access_method = "parameter"
    }
  }

  resource_requirements {
    key      = "flavor_type"
    values   = ["Ascend", "CPU"]
    operator = "in"
  }
  resource_requirements {
    key      = "device_distributed_mode"
    values   = ["singular"]
    operator = "in"
  }
  resource_requirements {
    key      = "host_distributed_mode"
    values   = ["multiple"]
    operator = "in"
  }

  advanced_config {
    auto_search {
      reward_attrs {
        name  = "search"
        mode  = "max"
        regex = "10.0"
      }

      search_params {
        name        = "parameter2"
        lower_bound = "10"
        upper_bound = "20"
      }
      search_params {
        name        = "parameter1"
        lower_bound = "20"
        upper_bound = "30"
      }

      algo_configs {
        name = "bayes_opt_search"

        params {
          key   = "num_samples"
          value = "20"
          type  = "Integer"
        }
        params {
          key   = "kind"
          value = "ucb"
          type  = "String"
        }
      }
    }
  }
}
`, testAccAlgorithm_base(name),
		name,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME,
	)
}

func testAccAlgorithm_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name         = "%[2]s"
    workspace_id = huaweicloud_modelarts_workspace.test.id

    tags {
      key = "auto_search"
    }
  }

  job_config {
    code_dir = "/${huaweicloud_obs_bucket.test.bucket}/algorithm/"
    command  = "bash $${MA_JOB_DIR}/algorithm/bootfile.py"

    engine {
      image_url = "%[3]s"
    }

    parameters {
      name  = "parameter_update"
      value = "40"

      constraint {
        type        = "Float"
        valid_range = ["20", "50"]
      }
    }

    inputs {
      name          = "input1"
      access_method = "parameter"

      remote_constraints {
        data_type  = "modelarts_dataset"
        attributes = jsonencode({
          data_format       = ["Default", "CarbonData"]
          data_segmentation = ["no_limit"]
          dataset_type      = ["1"]
        })
      }
    }
    inputs {
      name          = "a_input2"
      access_method = "env"
    }

    outputs {
      name          = "output_update"
      access_method = "env"
    }
  }

  resource_requirements {
    key      = "flavor_type"
    values   = ["Ascend"]
    operator = "in"
  }
  resource_requirements {
    key      = "device_distributed_mode"
    values   = ["multiple"]
    operator = "in"
  }
  resource_requirements {
    key      = "host_distributed_mode"
    values   = ["singular"]
    operator = "in"
  }

  advanced_config {
    auto_search {
      reward_attrs {
        name  = "search"
        mode  = "max"
        regex = "10.0"
      }
      reward_attrs {
        name  = "a_search"
        mode  = "min"
        regex = "50.0"
      }

      search_params {
        name        = "parameter_update"
        lower_bound = "20"
        upper_bound = "30"
      }

      algo_configs {
        name = "anneal_search"

        params {
          key   = "avg_best_idx"
          value = "2.0"
          type  = "Float"
        }
      }
    }
  }
}
`, testAccAlgorithm_base(name), name, acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL)
}

func testAccAlgorithm_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name         = "%[2]s"
    workspace_id = huaweicloud_modelarts_workspace.test.id
  }

  job_config {
    code_dir = "/${huaweicloud_obs_bucket.test.bucket}/algorithm/"
    command  = "bash $${MA_JOB_DIR}/algorithm/bootfile.py"

    engine {
      image_url = "%[3]s"
    }
  }
}
`, testAccAlgorithm_base(name), name, acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL)
}
