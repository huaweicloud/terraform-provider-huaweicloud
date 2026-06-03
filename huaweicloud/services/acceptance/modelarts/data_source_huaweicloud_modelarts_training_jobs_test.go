package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTrainingJobs_basic(t *testing.T) {
	var (
		name                 = acceptance.RandomAccResourceNameWithDash()
		rName                = "huaweicloud_modelarts_training_job.test"
		rNameWithCustomImage = "huaweicloud_modelarts_training_job.with_custom_image"

		all = "data.huaweicloud_modelarts_training_jobs.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byWorkspaceId   = "data.huaweicloud_modelarts_training_jobs.filter_by_workspace_id"
		dcByWorkspaceId = acceptance.InitDataSourceCheck(byWorkspaceId)

		byTrainType   = "data.huaweicloud_modelarts_training_jobs.filter_by_train_type"
		dcByTrainType = acceptance.InitDataSourceCheck(byTrainType)

		byFilters   = "data.huaweicloud_modelarts_training_jobs.filter_by_filters"
		dcByFilters = acceptance.InitDataSourceCheck(byFilters)

		bySortByAndOrderDesc   = "data.huaweicloud_modelarts_training_jobs.filter_by_sort_by_and_order_desc"
		dcBySortByAndOrderDesc = acceptance.InitDataSourceCheck(bySortByAndOrderDesc)
		bySortByAndOrderAsc    = "data.huaweicloud_modelarts_training_jobs.filter_by_sort_by_and_order_asc"
		dcBySortByAndOrderAsc  = acceptance.InitDataSourceCheck(bySortByAndOrderAsc)

		byIdWithCustomImage   = "data.huaweicloud_modelarts_training_jobs.filter_by_id_with_custom_image"
		dcByIdWithCustomImage = acceptance.InitDataSourceCheck(byIdWithCustomImage)
		byIdWithAlgorithm     = "data.huaweicloud_modelarts_training_jobs.filter_by_id_with_algorithm"
		dcByIdWithAlgorithm   = acceptance.InitDataSourceCheck(byIdWithAlgorithm)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsAlgorithm(t)
			acceptance.TestAccPreCheckModelArtsTrainingJobPublicResourcePoolFlavorID(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTrainingJobs_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "jobs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'workspace_id' parameter.
					dcByWorkspaceId.CheckResourceExists(),
					resource.TestCheckOutput("is_workspace_id_filter_useful", "true"),
					// Filter by 'train_type' parameter.
					dcByTrainType.CheckResourceExists(),
					resource.TestCheckOutput("is_train_type_filter_useful", "true"),
					// Filter by 'filters' parameter.
					dcByFilters.CheckResourceExists(),
					resource.TestCheckOutput("is_filters_filter_useful", "true"),
					// Filter by 'sort_by' and 'order' parameters.
					dcBySortByAndOrderDesc.CheckResourceExists(),
					dcBySortByAndOrderAsc.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_by_and_order_useful", "true"),
					// Check Attributes.
					resource.TestCheckResourceAttr(byFilters, "jobs.#", "1"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.metadata.0.id", rName, "id"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.kind", "job"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.metadata.0.name", rName, "metadata.0.name"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.metadata.0.description", "Created by Terraform script"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.metadata.0.workspace_id", rName, "metadata.0.workspace_id"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.metadata.0.annotations", rName, "metadata.0.annotations"),
					resource.TestCheckResourceAttrSet(byFilters, "jobs.0.metadata.0.user_name"),
					resource.TestCheckResourceAttrSet(byFilters, "jobs.0.status.0.phase"),
					resource.TestMatchResourceAttr(byFilters, "jobs.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.code_dir", rName, "algorithm.0.code_dir"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.boot_file", rName, "algorithm.0.boot_file"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.autosearch_config_path", rName,
						"algorithm.0.autosearch_config_path"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.autosearch_framework_path", rName,
						"algorithm.0.autosearch_framework_path"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.local_code_dir", rName, "algorithm.0.local_code_dir"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.working_dir", rName, "algorithm.0.working_dir"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.engine.0.engine_id", rName,
						"algorithm.0.engine.0.engine_id"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.engine.0.engine_name", rName,
						"algorithm.0.engine.0.engine_name"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.engine.0.engine_version", rName,
						"algorithm.0.engine.0.engine_version"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.inputs.#", "1"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.inputs.0.name", "data_url"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.inputs.0.description", "Training data input"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.inputs.0.local_dir", rName,
						"algorithm.0.inputs.0.local_dir"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.inputs.0.access_method", "parameter"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.inputs.0.remote.0.obs.0.obs_url", rName,
						"algorithm.0.inputs.0.remote.0.obs.0.obs_url"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.outputs.#", "1"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.outputs.0.name", "train_url"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.outputs.0.description", "Model output"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.outputs.0.local_dir", rName,
						"algorithm.0.outputs.0.local_dir"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.outputs.0.access_method", "parameter"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.algorithm.0.outputs.0.remote.0.obs.0.obs_url", rName,
						"algorithm.0.outputs.0.remote.0.obs.0.obs_url"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.parameters.0.name", "workers"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.parameters.0.value", "4"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.parameters.0.constraint.0.type", "Float"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.parameters.0.constraint.0.valid_type", "None"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.algorithm.0.environments.RUN_TYPE", "custom_job"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.spec.0.resource.0.flavor_id", rName,
						"spec.0.resource.0.flavor_id"),
					resource.TestCheckResourceAttrPair(byFilters, "jobs.0.spec.0.resource.0.node_count", rName,
						"spec.0.resource.0.node_count"),
					resource.TestCheckResourceAttrSet(byFilters, "jobs.0.spec.0.log_export_path.0.obs_url"),
					resource.TestCheckResourceAttr(byFilters, "jobs.0.spec.0.custom_metrics.#", "1"),
					resource.TestCheckResourceAttrSet(byFilters, "jobs.0.spec.0.custom_metrics.0.http_get.0.path"),
					resource.TestCheckResourceAttrSet(byFilters, "jobs.0.spec.0.custom_metrics.0.http_get.0.port"),
					dcByIdWithCustomImage.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.algorithm.0.command",
						rNameWithCustomImage, "algorithm.0.command"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.algorithm.0.engine.0.image_url",
						rNameWithCustomImage, "algorithm.0.engine.0.image_url"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.spec.0.resource.0.pool_id",
						rNameWithCustomImage, "spec.0.resource.0.pool_id"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage,
						"jobs.0.spec.0.resource.0.main_container_customized_flavor.0.cpu_core_num",
						rNameWithCustomImage, "spec.0.resource.0.main_container_customized_flavor.0.cpu_core_num"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.spec.0.resource.0.main_container_customized_flavor.0.mem_size",
						rNameWithCustomImage, "spec.0.resource.0.main_container_customized_flavor.0.mem_size"),
					resource.TestCheckResourceAttr(byIdWithCustomImage, "jobs.0.spec.0.runtime_type", "debug"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.spec.0.volumes.0.nfs.0.nfs_server_path",
						rNameWithCustomImage, "spec.0.volumes.0.nfs.0.nfs_server_path"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.spec.0.volumes.0.nfs.0.local_path",
						rNameWithCustomImage, "spec.0.volumes.0.nfs.0.local_path"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.spec.0.volumes.0.nfs.0.read_only",
						rNameWithCustomImage, "spec.0.volumes.0.nfs.0.read_only"),
					resource.TestCheckResourceAttr(byIdWithCustomImage, "jobs.0.spec.0.schedule_policy.0.priority", "1"),
					resource.TestCheckResourceAttr(byIdWithCustomImage, "jobs.0.spec.0.schedule_policy.0.preemptible", "true"),
					resource.TestCheckResourceAttr(byIdWithCustomImage,
						"jobs.0.spec.0.schedule_policy.0.required_affinity.0.node_affinity.0.node_selector_terms.0.match_expressions.0.key",
						"kubernetes.io/hostname"),
					resource.TestCheckResourceAttr(byIdWithCustomImage,
						"jobs.0.spec.0.schedule_policy.0.required_affinity.0.node_affinity.0.node_selector_terms.0.match_expressions.0.operator",
						"In"),
					resource.TestCheckResourceAttrPair(byIdWithCustomImage, "jobs.0.endpoints.0.ssh.0.key_pair_names.0",
						"huaweicloud_kps_keypair.test", "id"),
					dcByIdWithAlgorithm.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(byIdWithAlgorithm, "jobs.0.metadata.0.id",
						"huaweicloud_modelarts_training_job.with_algorithm", "id"),
					resource.TestCheckResourceAttrPair(byIdWithAlgorithm, "jobs.0.algorithm.0.id",
						"huaweicloud_modelarts_algorithm.test", "id"),
					resource.TestCheckResourceAttrPair(byIdWithAlgorithm, "jobs.0.algorithm.0.name",
						"huaweicloud_modelarts_algorithm.test", "metadata.0.name"),
					resource.TestCheckResourceAttr(byIdWithAlgorithm,
						"jobs.0.spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.weight", "100"),
					resource.TestCheckResourceAttr(byIdWithAlgorithm, "jobs.0.spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.#", "1"),
					resource.TestCheckResourceAttr(byIdWithAlgorithm,
						"jobs.0.spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.preference.0.match_expressions.0.key",
						"kubernetes.io/hostname"),
					resource.TestCheckResourceAttr(byIdWithAlgorithm,
						"jobs.0.spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.preference.0.match_expressions.0.operator", "In"),
					resource.TestCheckResourceAttrSet(byIdWithAlgorithm,
						"jobs.0.spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.preference.0.match_expressions.0.values.0"),
				),
			},
		},
	})
}

func testAccDataTrainingJobs_base(name string) string {
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

# Preset image algorithm job with public resource pool.
resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name         = "%[1]s"
    description  = "Created by Terraform script"
    workspace_id = "0"

    annotations = {
      "performance_diagnosis_enabled" = "true"
    }
  }

  algorithm {
    code_dir                  = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
    boot_file                 = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"
    autosearch_config_path    = "/config/autosearch.yaml"
    autosearch_framework_path = "/config/framework/"
    local_code_dir            = "/home/ma-user/modelarts/user-job-dir"
    working_dir               = "/home/ma-user/modelarts/user-job-dir"

    engine {
      engine_id      = "%[2]s"
      engine_version = "%[2]s"
      engine_name    = "%[3]s"
    }

    inputs {
      name          = "data_url"
      description   = "Training data input"
      local_dir     = "/home/ma-user/modelarts/user-job-dir/data_url_0"
      access_method = "parameter"

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
        }
      }
    }

    outputs {
      name          = "train_url"
      description   = "Model output"
      local_dir     = "/home/ma-user/modelarts/user-job-dir/train_url_0"
      access_method = "parameter"

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
        }
      }
    }

    parameters {
      name  = "workers"
      value = "4"

      constraint {
        type       = "Float"
        valid_type = "None"
      }
    }

    environments = {
      RUN_TYPE = "custom_job"
    }
  }

  spec {
    resource {
      flavor_id  = "%[4]s"
      node_count = 1
    }

    log_export_path {
      obs_url = "/${huaweicloud_obs_bucket.test.bucket}/logs/"
    }
    
    custom_metrics {
      http_get {
        path = "/raw_text"
        port = 10001
      }
    }
  }

  depends_on = [huaweicloud_obs_bucket_object.test]
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  resource_pool_id = try(split(",", "%[5]s")[0], "")
}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = local.resource_pool_id
}

locals {
  private_ip = try(data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[0].status[0].private_ip, "")
}

# Custom image algorithm job with dedicated resource pool.
resource "huaweicloud_modelarts_training_job" "with_custom_image" {
  kind = "job"

  metadata {
    name        = "%[1]s_debug_job"
    annotations = {
      "jupyter-lab/enable" = "true"
    }
  }

  algorithm {
    code_dir = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
    command  = "python bootfile.py --epochs 10"

    engine {
      image_url = "%[6]s"
    }
  }

  spec {
    resource {
      pool_id    = local.resource_pool_id
      node_count = 1

      main_container_customized_flavor {
        cpu_core_num = 2
        mem_size     = 10
      }
    }

    runtime_type = "debug"

    volumes {
      nfs {
        nfs_server_path = "nas.local:/share/work"
        local_path      = "/mnt/nfs/"
        read_only       = false
      }
    }

    schedule_policy {
      priority    = 1
      preemptible = true

      required_affinity {
        node_affinity {
          node_selector_terms {
            match_expressions {
              key      = "kubernetes.io/hostname"
              operator = "In"
              values   = try(local.private_ip, "") != "" ? [local.private_ip] : []
            }
          }
        }
      }
    }
  }

  endpoints {
    ssh {
      key_pair_names = [huaweicloud_kps_keypair.test.id]
    }
  }
}

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name = "%[1]s"
  }

  job_config {
    code_dir = "/${huaweicloud_obs_bucket.test.bucket}/%[1]s/"
    command  = "python bootfile.py --epochs 10"

    engine {
      image_url = "%[6]s"
    }
  }

  depends_on = [huaweicloud_obs_bucket_object.test]
}

# Existing algorithm job with dedicated resource pool.
resource "huaweicloud_modelarts_training_job" "with_algorithm" {
  kind = "job"

  metadata {
    name = "%[1]s_with_dedicated_resource_pool"
  }

  algorithm {
    id = huaweicloud_modelarts_algorithm.test.id
  }

  spec {
    resource {
      pool_id    = local.resource_pool_id
      node_count = 1
    }

    schedule_policy {
      priority    = 1
      preemptible = true

      preferred_affinity {
        node_affinity {
          weight = 100
          preference {
            match_expressions {
              key      = "kubernetes.io/hostname"
              operator = "In"
              values   = local.private_ip != "" ? [local.private_ip] : []
            }
          }
        }
      }
    }
  }
}
`, name,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME,
		acceptance.HW_MODELARTS_TRAINING_JOB_PUBLIC_RESOURCE_POOL_FLAVOR_ID,
		acceptance.HW_MODELARTS_RESOURCE_POOL_IDS,
		acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL,
	)
}

func testAccDataTrainingJobs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_modelarts_training_jobs" "test" {
  depends_on = [huaweicloud_modelarts_training_job.test]
}

# Filter by 'workspace_id' parameter.
locals {
  workspace_id = huaweicloud_modelarts_training_job.test.metadata[0].workspace_id
}

data "huaweicloud_modelarts_training_jobs" "filter_by_workspace_id" {
  workspace_id = local.workspace_id

  depends_on = [huaweicloud_modelarts_training_job.test]
}

locals {
  workspace_id_filter_result = [for v in data.huaweicloud_modelarts_training_jobs.filter_by_workspace_id.jobs[*].metadata[0] :
  v.workspace_id == local.workspace_id]
}

output "is_workspace_id_filter_useful" {
  value = length(local.workspace_id_filter_result) > 0 && alltrue(local.workspace_id_filter_result)
}

# Filter by 'train_type' parameter.
data "huaweicloud_modelarts_training_jobs" "filter_by_train_type" {
  train_type = "job"

  depends_on = [huaweicloud_modelarts_training_job.test]
}

output "is_train_type_filter_useful" {
  value = length(data.huaweicloud_modelarts_training_jobs.filter_by_train_type.jobs) > 0
}

# Filter by 'filters' parameter.
locals {
  job_id = huaweicloud_modelarts_training_job.test.id
}

data "huaweicloud_modelarts_training_jobs" "filter_by_filters" {
  filters {
    key      = "id"
    operator = "in"
    value    = [local.job_id]
  }
}

locals {
  filters_filter_result = [for v in data.huaweicloud_modelarts_training_jobs.filter_by_filters.jobs[*].metadata[0].id :
  v == local.job_id]
}

output "is_filters_filter_useful" {
  value = length(local.filters_filter_result) > 0 && alltrue(local.filters_filter_result)
}

# Filter by 'sort_by' and 'order' parameters.
data "huaweicloud_modelarts_training_jobs" "filter_by_sort_by_and_order_desc" {
  depends_on = [huaweicloud_modelarts_training_job.test]
}

data "huaweicloud_modelarts_training_jobs" "filter_by_sort_by_and_order_asc" {
  sort_by = "create_time"
  order   = "asc"

  depends_on = [huaweicloud_modelarts_training_job.test]
}

locals {
  sort_by_and_order_desc_result = data.huaweicloud_modelarts_training_jobs.filter_by_sort_by_and_order_desc.jobs[*].create_time
  sort_by_and_order_asc_result  = data.huaweicloud_modelarts_training_jobs.filter_by_sort_by_and_order_asc.jobs[*].create_time
}

output "is_sort_by_and_order_useful" {
  value = local.sort_by_and_order_desc_result == reverse(local.sort_by_and_order_asc_result)
}

data "huaweicloud_modelarts_training_jobs" "filter_by_id_with_custom_image" {
  filters {
    key      = "id"
    operator = "in"
    value    = [huaweicloud_modelarts_training_job.with_custom_image.id]
  }
}

data "huaweicloud_modelarts_training_jobs" "filter_by_id_with_algorithm" {
  filters {
    key      = "id"
    operator = "in"
    value    = [huaweicloud_modelarts_training_job.with_algorithm.id]
  }
}
`, testAccDataTrainingJobs_base(name))
}
