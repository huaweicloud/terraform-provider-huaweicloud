package modelarts

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getTrainingJobResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetTrainingJobById(client, state.Primary.ID)
}

func TestAccTrainingJob_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_modelarts_training_job.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getTrainingJobResourceFunc)

		rNameCustomAlgorithm = "huaweicloud_modelarts_training_job.custom_algorithm"
		rcCustomAlgorithm    = acceptance.InitResourceCheck(rNameCustomAlgorithm, &obj, getTrainingJobResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsAlgorithmEngine(t)
			acceptance.TestAccPreCheckModelArtsTrainingJobPublicResourcePoolFlavorID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcCustomAlgorithm.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTrainingJob_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "kind", "job"),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", "Created by Terraform script"),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.%", "3"),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.fault-tolerance/job-retry-num", "3"),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.fault-tolerance/hang-retry", "true"),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.ssh-plugin/ssh-key-file-path", "/home/ma-user/.ssh/"),
					resource.TestCheckResourceAttrPair(rName, "metadata.0.training_experiment_reference.0.id",
						"huaweicloud_modelarts_training_experiment.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "algorithm.0.id", "huaweicloud_modelarts_algorithm.test", "id"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.code_dir",
						fmt.Sprintf("/%s/training_job/", name)),
					resource.TestCheckResourceAttr(rName, "algorithm.0.boot_file",
						fmt.Sprintf("/%s/training_job/bootfile.py", name)),
					resource.TestCheckResourceAttr(rName, "algorithm.0.environments.TRAIN_URL",
						"/home/ma-user/modelarts/user-job-dir/envs"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.inputs.#", "1"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.inputs.0.name", "epochs"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.outputs.#", "1"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.outputs.0.name", "train_url"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.parameters.#", "1"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.parameters.0.name", "workers"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.parameters.0.value", "4"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.parameters.0.constraint.0.type", "Float"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.parameters.0.constraint.0.valid_type", "None"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.flavor_id",
						acceptance.HW_MODELARTS_TRAINING_JOB_PUBLIC_RESOURCE_POOL_FLAVOR_ID),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.node_count", "1"),
					resource.TestCheckResourceAttr(rName, "spec.0.runtime_type", "production"),
					resource.TestCheckResourceAttr(rName, "spec.0.log_export_path.0.obs_url",
						fmt.Sprintf("/%s/logs/", name)),
					resource.TestCheckResourceAttr(rName, "spec.0.log_export_path.0.host_path", "/var/log/training"),
					resource.TestCheckResourceAttr(rName, "spec.0.log_export_config.0.version", "v1"),
					resource.TestCheckResourceAttr(rName, "spec.0.log_export_config.0.rotation_enabled", "true"),
					resource.TestCheckResourceAttr(rName, "spec.0.auto_stop.0.time_unit", "HOURS"),
					resource.TestCheckResourceAttr(rName, "spec.0.auto_stop.0.duration", "2"),
					resource.TestCheckResourceAttrPair(rName, "spec.0.notification.0.topic_urn",
						"huaweicloud_smn_topic.test", "id"),
					resource.TestCheckResourceAttr(rName, "spec.0.notification.0.events.#", "2"),
					resource.TestCheckResourceAttr(rName, "spec.0.custom_metrics.0.http_get.0.path", "/raw_text"),
					resource.TestCheckResourceAttr(rName, "spec.0.custom_metrics.0.http_get.0.port", "10001"),
					resource.TestCheckResourceAttr(rName, "tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "status"),

					rcCustomAlgorithm.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "kind", "job"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "metadata.0.description", "Created by Terraform script"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.code_dir",
						fmt.Sprintf("/%s/training_job/", name)),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.boot_file",
						fmt.Sprintf("/%s/training_job/bootfile.py", name)),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.autosearch_config_path", "/config/autosearch.yaml"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.autosearch_framework_path", "/framework/"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.local_code_dir", "/home/ma-user/modelarts/user-job-dir"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.working_dir", "/home/ma-user/modelarts/user-job-dir"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.engine.0.engine_id",
						acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.engine.0.engine_name",
						acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "algorithm.0.engine.0.engine_version",
						acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "spec.0.resource.0.flavor_id",
						acceptance.HW_MODELARTS_TRAINING_JOB_PUBLIC_RESOURCE_POOL_FLAVOR_ID),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "spec.0.resource.0.node_count", "1"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "tags.%", "1"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "tags.foo", "bar"),
				),
			},
			{
				Config: testAccTrainingJob_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "metadata.0.description", ""),
					resource.TestCheckResourceAttr(rName, "tags.%", "1"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),

					rcCustomAlgorithm.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "metadata.0.description", "Updated by Terraform script"),
					resource.TestCheckResourceAttr(rNameCustomAlgorithm, "tags.%", "0"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.annotations", "algorithm.0.inputs"},
			},
			{
				ResourceName:            rNameCustomAlgorithm,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.annotations", "algorithm.0.inputs"},
			},
		},
	})
}

func testAccTrainingJob_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "training_job/bootfile.py"
  content      = <<EOF
#!/usr/bin/env python
import os
print(os.getcwd())
EOF
  content_type = "text/py"
}

resource "huaweicloud_smn_topic" "test" {
  name = "%[1]s"
}

resource "huaweicloud_modelarts_algorithm" "test" {
  metadata {
    name = "%[1]s"
  }

  job_config {
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    engine {
      engine_id      = "%[2]s"
      engine_version = "%[2]s"
      engine_name    = "%[3]s"
    }

    parameters_customization = true

    parameters {
      name  = "workers"
      value = "4"

      constraint {
        type       = "Float"
        valid_type = "None"
      }
    }

    inputs {
      name          = "epochs"
      access_method = "parameter"
      description   = "Training epochs"

      remote_constraints {
        data_type = "obs"
      }
    }

    outputs {
      name          = "train_url"
      access_method = "env"
      description   = "Training URL"
    }
  }

  advanced_config {
    auto_search {
      reward_attrs {
        name  = "search"
        mode  = "max"
        regex = "10.0"
      }

      search_params {
        name        = "workers"
        lower_bound = "10"
        upper_bound = "20"
      }

      algo_configs {
        name = "tpe_search"

        params {
          key   = "seed"
          value = "1"
          type  = "Integer"
        }
        params {
          key   = "gamma"
          value = "0.25"
          type  = "Float"
        }
        params {
          key   = "n_initial_points"
          value = "20"
          type  = "Integer"
        }
        params {
          key   = "num_samples"
          value = "20"
          type  = "Integer"
        }
      }
    }
  }
}

resource "huaweicloud_modelarts_training_experiment" "test" {
  metadata {
    name = "%[1]s"
  }
}
`, name, acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID, acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME)
}

func testAccTrainingJob_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# Existing algorithm + public resource pool
resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name        = "%[2]s"
    description = "Created by Terraform script"
    annotations = {
      "fault-tolerance/job-retry-num" = "3"
      "fault-tolerance/hang-retry"    = "true"
      "ssh-plugin/ssh-key-file-path"  = "/home/ma-user/.ssh/"
    }

    training_experiment_reference {
      id = huaweicloud_modelarts_training_experiment.test.id
    }
  }

  algorithm {
    id        = huaweicloud_modelarts_algorithm.test.id
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    environments = {
      "TRAIN_URL" = "/home/ma-user/modelarts/user-job-dir/envs"
    }

    inputs {
      name          = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.name
      description   = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.description
      local_dir     = "/home/ma-user/modelarts/inputs/${huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.name}_0"
      access_method = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.access_method

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
        }
      }
    }

    outputs {
      name          = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.name
      description   = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.description
      local_dir     = "/home/ma-user/modelarts/outputs/${huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.name}_0"
      access_method = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.access_method

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
        }
      }
    }

    dynamic "parameters" {
      for_each = huaweicloud_modelarts_algorithm.test.job_config.0.parameters

      content {
        name  = parameters.value.name
        value = parameters.value.value

        dynamic "constraint" {
          for_each = length(parameters.value.constraint) > 0 ? parameters.value.constraint : []

          content {
            type       = constraint.value.type
            valid_type = constraint.value.valid_type
          }
        }
      }
    }
  }

  spec {
    resource {
      flavor_id  = "%[3]s"
      node_count = 1
    }

    log_export_path {
      obs_url   = "/${huaweicloud_obs_bucket.test.bucket}/logs/"
      host_path = "/var/log/training"
    }

    log_export_config {
      version          = "v1"
      rotation_enabled = true
    }

    auto_stop {
      time_unit = "HOURS"
      duration  = 2
    }

    notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = ["JobStarted", "JobCompleted"]
    }

    custom_metrics {
      http_get {
        path = "/raw_text"
        port = 10001
      }
    }
  }

  tags = {
    foo = "bar"
  }
}

# Custom algorithm + public resource pool
resource "huaweicloud_modelarts_training_job" "custom_algorithm" {
  kind = "job"

  metadata {
    name        = "%[2]s_custom"
    description = "Created by Terraform script"
  }

  algorithm {
    code_dir                  = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
    boot_file                 = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"
    autosearch_config_path    = "/config/autosearch.yaml"
    autosearch_framework_path = "/framework/"
    local_code_dir            = "/home/ma-user/modelarts/user-job-dir"
    working_dir               = "/home/ma-user/modelarts/user-job-dir"

    engine {
      engine_id      = "%[4]s"
      engine_version = "%[4]s"
      engine_name    = "%[5]s"
    }
  }

  spec {
    resource {
      flavor_id  = "%[3]s"
      node_count = 1
    }
  }

  tags = {
    foo = "bar"
  }
}
`, testAccTrainingJob_base(name), name,
		acceptance.HW_MODELARTS_TRAINING_JOB_PUBLIC_RESOURCE_POOL_FLAVOR_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME)
}

func testAccTrainingJob_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name        = "%[2]s"
    annotations = {
      "fault-tolerance/job-retry-num" = "3"
      "fault-tolerance/hang-retry"    = "true"
      "ssh-plugin/ssh-key-file-path"  = "/home/ma-user/.ssh/"
    }

    training_experiment_reference {
      id = huaweicloud_modelarts_training_experiment.test.id
    }
  }

  algorithm {
    id        = huaweicloud_modelarts_algorithm.test.id
    code_dir  = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
    boot_file = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"

    environments = {
      "TRAIN_URL" = "/home/ma-user/modelarts/user-job-dir/envs"
    }

    inputs {
      name          = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.name
      description   = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.description
      local_dir     = "/home/ma-user/modelarts/inputs/${huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.name}_0"
      access_method = huaweicloud_modelarts_algorithm.test.job_config.0.inputs.0.access_method

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
        }
      }
    }

    outputs {
      name          = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.name
      description   = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.description
      local_dir     = "/home/ma-user/modelarts/outputs/${huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.name}_0"
      access_method = huaweicloud_modelarts_algorithm.test.job_config.0.outputs.0.access_method

      remote {
        obs {
          obs_url = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
        }
      }
    }

    dynamic "parameters" {
      for_each = huaweicloud_modelarts_algorithm.test.job_config.0.parameters

      content {
        name  = parameters.value.name
        value = parameters.value.value

        dynamic "constraint" {
          for_each = length(parameters.value.constraint) > 0 ? parameters.value.constraint : []

          content {
            type       = constraint.value.type
            valid_type = constraint.value.valid_type
          }
        }
      }
    }
  }

  spec {
    resource {
      flavor_id  = "%[3]s"
      node_count = 1
    }

    log_export_path {
      obs_url   = "/${huaweicloud_obs_bucket.test.bucket}/logs/"
      host_path = "/var/log/training"
    }

    log_export_config {
      version          = "v1"
      rotation_enabled = true
    }

    auto_stop {
      time_unit = "HOURS"
      duration  = 2
    }

    notification {
      topic_urn = huaweicloud_smn_topic.test.id
      events    = ["JobStarted", "JobCompleted"]
    }

    custom_metrics {
      http_get {
        path = "/raw_text"
        port = 10001
      }
    }
  }

  tags = {
    owner = "terraform"
  }
}

resource "huaweicloud_modelarts_training_job" "custom_algorithm" {
  kind = "job"

  metadata {
    name        = "%[2]s_custom"
    description = "Updated by Terraform script"
  }

  algorithm {
    code_dir                  = "/${huaweicloud_obs_bucket.test.bucket}/training_job/"
    boot_file                 = "/${huaweicloud_obs_bucket.test.bucket}/${huaweicloud_obs_bucket_object.test.key}"
    autosearch_config_path    = "/config/autosearch.yaml"
    autosearch_framework_path = "/framework/"
    local_code_dir            = "/home/ma-user/modelarts/user-job-dir"
    working_dir               = "/home/ma-user/modelarts/user-job-dir"

    engine {
      engine_id      = "%[4]s"
      engine_version = "%[4]s"
      engine_name    = "%[5]s"
    }
  }

  spec {
    resource {
      flavor_id  = "%[3]s"
      node_count = 1
    }
  }

  tags = {}
}
`, testAccTrainingJob_base(name), name,
		acceptance.HW_MODELARTS_TRAINING_JOB_PUBLIC_RESOURCE_POOL_FLAVOR_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_ID,
		acceptance.HW_MODELARTS_ALGORITHM_ENGINE_NAME)
}

func TestAccTrainingJob_debug_job(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		obj   interface{}
		rName = "huaweicloud_modelarts_training_job.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getTrainingJobResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsAlgorithmImageUrl(t)
			acceptance.TestAccPreCheckModelArtsResourcePoolIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccTrainingJob_debug_job_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "kind", "job"),
					resource.TestCheckResourceAttr(rName, "metadata.0.name", name),
					resource.TestCheckResourceAttr(rName, "metadata.0.annotations.%", "1"),
					resource.TestCheckResourceAttrSet(rName, "metadata.0.workspace_id"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.code_dir", fmt.Sprintf("/%s/%s/", name, name)),
					resource.TestCheckResourceAttr(rName, "algorithm.0.command", "python train.py --epochs 10"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.autosearch_config_path", "/config/autosearch.yaml"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.autosearch_framework_path", "/framework/"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.local_code_dir", "/home/ma-user/modelarts/user-job-dir"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.working_dir", "/home/ma-user/modelarts/user-job-dir"),
					resource.TestCheckResourceAttr(rName, "algorithm.0.engine.0.image_url",
						acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.pool_id",
						strings.Split(acceptance.HW_MODELARTS_RESOURCE_POOL_IDS, ",")[0]),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.node_count", "1"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.main_container_customized_flavor.0.cpu_core_num", "2"),
					resource.TestCheckResourceAttr(rName, "spec.0.resource.0.main_container_customized_flavor.0.mem_size", "10"),
					resource.TestCheckResourceAttr(rName, "spec.0.runtime_type", "debug"),
					resource.TestCheckResourceAttr(rName, "spec.0.volumes.0.nfs.0.nfs_server_path", "nas.local:/share/work"),
					resource.TestCheckResourceAttr(rName, "spec.0.volumes.0.nfs.0.local_path", "/mnt/nfs/"),
					resource.TestCheckResourceAttr(rName, "spec.0.volumes.0.nfs.0.read_only", "false"),
					resource.TestCheckResourceAttr(rName, "spec.0.schedule_policy.0.priority", "1"),
					resource.TestCheckResourceAttr(rName, "spec.0.schedule_policy.0.preemptible", "true"),
					resource.TestCheckResourceAttr(rName,
						"spec.0.schedule_policy.0.required_affinity.0.node_affinity.0.node_selector_terms.0.match_expressions.0.key",
						"kubernetes.io/hostname"),
					resource.TestCheckResourceAttr(rName,
						"spec.0.schedule_policy.0.required_affinity.0.node_affinity.0.node_selector_terms.0.match_expressions.0.operator", "In"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.ssh.0.key_pair_names.0",
						"huaweicloud_kps_keypair.test", "id"),
				),
			},
			{
				Config: testAccTrainingJob_debug_job_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "spec.0.schedule_policy.0.priority", "1"),
					resource.TestCheckResourceAttr(rName, "spec.0.schedule_policy.0.preemptible", "true"),
					resource.TestCheckResourceAttr(rName, "spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.weight", "100"),
					resource.TestCheckResourceAttr(rName,
						"spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.preference.0.match_expressions.0.key",
						"kubernetes.io/hostname"),
					resource.TestCheckResourceAttr(rName,
						"spec.0.schedule_policy.0.preferred_affinity.0.node_affinity.0.preference.0.match_expressions.0.operator", "In"),
					resource.TestCheckResourceAttrPair(rName, "endpoints.0.ssh.0.key_pair_names.0",
						"huaweicloud_kps_keypair.test", "id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.annotations", "enable_force_new"},
			},
		},
	})
}

func testAccTrainingJob_debug_job_base(name string) string {
	return fmt.Sprintf(`
locals {
  resource_pool_id = try(split(",", "%[1]s")[0], "")
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%[2]s/bootfile.py"
  content      = <<EOF
#!/usr/bin/env python
import os
print(os.getcwd())
EOF
  content_type = "text/py"
}


resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

data "huaweicloud_modelartsv2_resource_pools" "test" {}

locals {
  resource_pool = try([for v in data.huaweicloud_modelartsv2_resource_pools.test.resource_pools :
  v if v.metadata[0].name == local.resource_pool_id][0], {})
}

data "huaweicloud_modelartsv2_resource_pool_nodes" "test" {
  resource_pool_name = local.resource_pool_id
}

locals {
  private_ip = data.huaweicloud_modelartsv2_resource_pool_nodes.test.nodes[0].status[0].private_ip
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_IDS, name)
}

func testAccTrainingJob_debug_job_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name         = "%[2]s"
    workspace_id = try(local.resource_pool.metadata[0].labels["os.modelarts/workspace.id"], "")
    annotations  = {
      "jupyter-lab/enable" = "true"
    }
  }

  algorithm {
    code_dir                  = "/${huaweicloud_obs_bucket.test.bucket}/%[2]s/"
    command                   = "python train.py --epochs 10"
    autosearch_config_path    = "/config/autosearch.yaml"
    autosearch_framework_path = "/framework/"
    local_code_dir            = "/home/ma-user/modelarts/user-job-dir"
    working_dir               = "/home/ma-user/modelarts/user-job-dir"

    engine {
      image_url = "%[3]s"
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
`, testAccTrainingJob_debug_job_base(name),
		name,
		acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL)
}

// Only for testing the 'spec.schedule_policy.preferred_affinity' parameter.
func testAccTrainingJob_debug_job_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_training_job" "test" {
  kind = "job"

  metadata {
    name         = "%[2]s"
    workspace_id = try(local.resource_pool.metadata[0].labels["os.modelarts/workspace.id"], "")
    annotations  = {
      "jupyter-lab/enable" = "true"
    }
  }

  algorithm {
    code_dir                  = "/${huaweicloud_obs_bucket.test.bucket}/%[2]s/"
    command                   = "python train.py --epochs 10"
    autosearch_config_path    = "/config/autosearch.yaml"
    autosearch_framework_path = "/framework/"
    local_code_dir            = "/home/ma-user/modelarts/user-job-dir"
    working_dir               = "/home/ma-user/modelarts/user-job-dir"

    engine {
      image_url = "%[3]s"
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

      preferred_affinity {
        node_affinity {
          weight = 100
          preference {
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

  enable_force_new = "true"
}
`, testAccTrainingJob_debug_job_base(name),
		name,
		acceptance.HW_MODELARTS_ALGORITHM_IMAGE_URL)
}
