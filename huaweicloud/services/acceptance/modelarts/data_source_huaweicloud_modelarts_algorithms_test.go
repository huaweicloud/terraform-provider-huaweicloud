package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAlgorithms_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all   = "data.huaweicloud_modelarts_algorithms.all"
		dcAll = acceptance.InitDataSourceCheck(all)

		byWorkspaceId   = "data.huaweicloud_modelarts_algorithms.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		bySearches   = "data.huaweicloud_modelarts_algorithms.filter_by_searches"
		dcBySearches = acceptance.InitDataSourceCheck(bySearches)

		bySortByAndOrderDesc   = "data.huaweicloud_modelarts_algorithms.filter_by_sort_by_and_order_desc"
		dcBySortByAndOrderDesc = acceptance.InitDataSourceCheck(bySortByAndOrderDesc)

		bySortByAndOrderAsc   = "data.huaweicloud_modelarts_algorithms.filter_by_sort_by_and_order_asc"
		dcBySortByAndOrderAsc = acceptance.InitDataSourceCheck(bySortByAndOrderAsc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsAlgorithm(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAlgorithms_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dcAll.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "algorithms.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'workspace_id' parameter.
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					// Filter by 'searches' (fuzzy name matching) parameter.
					dcBySearches.CheckResourceExists(),
					resource.TestCheckOutput("is_searches_filter_useful", "true"),
					// Filter by 'sort_by' and 'order' parameters.
					dcBySortByAndOrderDesc.CheckResourceExists(),
					dcBySortByAndOrderAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_by_and_order_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.id"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.description"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.workspace_id"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.create_time"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.source"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.metadata.0.is_valid"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.job_config.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.code_dir"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.boot_file"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.parameters_customization"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.engine.0.engine_id"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.engine.0.engine_name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.engine.0.engine_version"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.job_config.0.inputs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.inputs.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.inputs.0.description"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.job_config.0.outputs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.outputs.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.outputs.0.description"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.job_config.0.parameters.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.parameters.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.job_config.0.parameters.0.description"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.resource_requirements.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.resource_requirements.0.key"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.resource_requirements.0.operator"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.resource_requirements.0.values.#"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.resource_requirements.1.key"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.resource_requirements.1.values.#"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.advanced_config.0.auto_search.0.reward_attrs.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.reward_attrs.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.reward_attrs.0.mode"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.reward_attrs.0.regex"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.advanced_config.0.auto_search.0.search_params.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.search_params.0.name"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.search_params.0.lower_bound"),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.search_params.0.upper_bound"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.advanced_config.0.auto_search.0.algo_configs.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(bySearches, "algorithms.0.advanced_config.0.auto_search.0.algo_configs.0.name"),
					resource.TestMatchResourceAttr(bySearches, "algorithms.0.advanced_config.0.auto_search.0.algo_configs.0.params.#",
						regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataAlgorithms_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%[1]s/bootfile.py"
  content      = <<EOF
#!/usr/bin/env python
import os
print(os.getcwd())
EOF
  content_type = "text/py"
}

resource "huaweicloud_modelarts_workspace" "test" {
  name = "%[1]s"
}

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name         = "%[1]s"
    description  = "Created by terraform script"
    workspace_id = huaweicloud_modelarts_workspace.test.id

    tags {
      key = "auto_search"
    }
  }

  job_config {
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    engine {
      engine_id      = "%[2]s"
      engine_version = "%[2]s"
      engine_name    = "%[3]s"
    }

    parameters_customization = true

    inputs {
      name          = "%[1]s"
      description   = "Input parameter1"
      access_method = "parameter"
      remote_constraints {
        data_type = "obs"
      }

      remote_constraints {
        data_type = "modelarts_dataset"
        attributes = jsonencode({
          data_format       = ["CarbonData"]
          data_segmentation = ["false"]
          dataset_type      = ["0"]
        })
      }
    }

    outputs {
      name          = "%[1]s"
      description   = "Output parameter1"
      access_method = "env"
    }

    parameters {
      name        = "%[1]s"
      description = "parameter1"

      constraint {
        type        = "Float"
        editable    = true
        required    = true
        valid_range = ["10", "50"]
      }
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
        name        = "%[1]s"
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
`, name, acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID, acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME)
}

func testAccDataAlgorithms_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_modelarts_algorithms" "all" {
  depends_on = [huaweicloud_modelarts_algorithm.test]
}

# Filter by 'workspace_id' parameter.
locals {
  workspace_id = huaweicloud_modelarts_workspace.test.id
}
data "huaweicloud_modelarts_algorithms" "filter_by_workspace_id" {
  workspace_id = local.workspace_id

  depends_on = [huaweicloud_modelarts_algorithm.test]
}

locals {
  workspace_id_filter_result = [for v in data.huaweicloud_modelarts_algorithms.filter_by_workspace_id.algorithms :
  v.metadata[0].workspace_id == local.workspace_id]
}

output "is_workspace_id_filter_useful" {
  value = length(local.workspace_id_filter_result) > 0 && alltrue(local.workspace_id_filter_result)
}

# Filter by 'searches' (fuzzy name matching) parameter.
locals {
  name = huaweicloud_modelarts_algorithm.test.metadata[0].name
}

data "huaweicloud_modelarts_algorithms" "filter_by_searches" {
  searches = "name:${local.name}"

  depends_on = [huaweicloud_modelarts_algorithm.test]
}

locals {
  searches_filter_result = [for v in data.huaweicloud_modelarts_algorithms.filter_by_searches.algorithms :
  v.metadata[0].name == local.name]
}

output "is_searches_filter_useful" {
  value = length(local.searches_filter_result) > 0 && alltrue(local.searches_filter_result)
}

# Filter by 'sort_by' and 'order' parameters.
data "huaweicloud_modelarts_algorithms" "filter_by_sort_by_and_order_desc" {
  sort_by = "name"

  depends_on = [huaweicloud_modelarts_algorithm.test]
}

data "huaweicloud_modelarts_algorithms" "filter_by_sort_by_and_order_asc" {
  sort_by = "name"
  order   = "asc"

  depends_on = [huaweicloud_modelarts_algorithm.test]
}

locals {
  sort_by_and_order_desc_result = data.huaweicloud_modelarts_algorithms.filter_by_sort_by_and_order_desc.algorithms[*].metadata[0].name
  sort_by_and_order_asc_result  = data.huaweicloud_modelarts_algorithms.filter_by_sort_by_and_order_asc.algorithms[*].metadata[0].name
}

output "is_sort_by_and_order_useful" {
  value = local.sort_by_and_order_desc_result == reverse(local.sort_by_and_order_asc_result)
}
`, testAccDataAlgorithms_base(name))
}
