package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3ComponentRecords_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_servicestagev3_component_records.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			acceptance.TestAccPreCheckImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3ComponentRecords_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_begin_time_set", "true"),
					resource.TestCheckOutput("is_end_time_set", "true"),
					resource.TestCheckOutput("is_description_set", "true"),
					resource.TestCheckOutput("is_instance_id_set", "true"),
					resource.TestCheckOutput("is_version_set", "true"),
					resource.TestCheckOutput("is_current_used_set", "true"),
					resource.TestCheckOutput("is_status_set", "true"),
					resource.TestCheckOutput("is_deploy_type_set", "true"),
					resource.TestCheckOutput("is_jobs_set", "true"),
					resource.TestCheckOutput("is_job_sequence_set", "true"),
					resource.TestCheckOutput("is_job_job_id_set", "true"),
					resource.TestCheckOutput("is_job_info_set", "true"),
					resource.TestCheckOutput("is_job_info_source_url_set", "true"),
					resource.TestCheckOutput("is_job_info_first_batch_weight_set", "true"),
					resource.TestCheckOutput("is_job_info_first_batch_replica_set", "true"),
					resource.TestCheckOutput("is_job_info_replica_set", "true"),
					resource.TestCheckOutput("is_job_info_remaining_batch_set", "true"),
				),
			},
		},
	})
}

func testAccDataV3ComponentRecords_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_cce_clusters" "test" {
  cluster_id = "%[1]s"
}
  
resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[2]s"
  enterprise_project_id = "0"
}
  
resource "huaweicloud_servicestagev3_environment" "test" {
  name                  = "%[2]s"
  vpc_id                = try(data.huaweicloud_cce_clusters.test.clusters[0].vpc_id, "")
  enterprise_project_id = "0"
}

resource "huaweicloud_servicestagev3_environment_associate" "test" {
  environment_id = huaweicloud_servicestagev3_environment.test.id

  resources {
    type = "cce"
    id   = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
  }
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  docker_runtime_stack = try([for o in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks:
    o if o.type == "Docker" && o.deploy_mode == "container"][0], {})
}

resource "huaweicloud_servicestagev3_component" "test" {
  depends_on = [
    huaweicloud_servicestagev3_environment_associate.test
  ]

  application_id = huaweicloud_servicestagev3_application.test.id
  environment_id = huaweicloud_servicestagev3_environment.test.id
  name           = "%[2]s"

  runtime_stack {
    deploy_mode = try(local.docker_runtime_stack.deploy_mode, "container")
    type        = try(local.docker_runtime_stack.type, "Docker")
    name        = try(local.docker_runtime_stack.name, "Docker")
    version     = try(local.docker_runtime_stack.version, null)
  }

  source = jsonencode({
    "auth": "iam",
    "kind": "image",
    "storage": "swr",
    "url": "%[3]s"
  })

  version = "1.0.1"
  replica = 2

  refer_resources {
    id         = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type       = "cce"
    parameters = jsonencode({
      "namespace": "default",
      "type": "VirtualMachine"
    })
  }

  limit_cpu      = 0.5
  limit_memory   = 1
  request_cpu    = 0.5
  request_memory = 1
}

data "huaweicloud_servicestagev3_component_records" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  component_id   = huaweicloud_servicestagev3_component.test.id
}

locals {
  component_record = try(data.huaweicloud_servicestagev3_component_records.test.records[0], null)
}

output "is_begin_time_set" {
  value = lookup(local.component_record, "begin_time", null) != null
}

output "is_end_time_set" {
  value = lookup(local.component_record, "end_time", null) != null
}

output "is_description_set" {
  value = lookup(local.component_record, "description", null) != null
}

output "is_instance_id_set" {
  value = lookup(local.component_record, "instance_id", null) != null
}

output "is_version_set" {
  value = lookup(local.component_record, "version", null) != null
}

output "is_current_used_set" {
  value = lookup(local.component_record, "current_used", null) != null
}

output "is_status_set" {
  value = lookup(local.component_record, "status", null) != null
}

output "is_deploy_type_set" {
  value = lookup(local.component_record, "deploy_type", null) != null
}

output "is_jobs_set" {
  value = length(lookup(local.component_record, "jobs", [])) > 0
}

locals {
  component_record_job = try(lookup(local.component_record, "jobs", [])[0], {})
}

output "is_job_sequence_set" {
  value = lookup(local.component_record_job, "sequence", null) != null
}

output "is_job_job_id_set" {
  value = lookup(local.component_record_job, "job_id", null) != null
}

output "is_job_info_set" {
  value = length(lookup(local.component_record_job, "job_info", [])) > 0
}

locals {
  component_record_job_info = try(lookup(local.component_record_job, "job_info", [])[0], {})
}

output "is_job_info_source_url_set" {
  value = lookup(local.component_record_job_info, "source_url", null) != null
}

output "is_job_info_first_batch_weight_set" {
  value = lookup(local.component_record_job_info, "first_batch_weight", null) != null
}

output "is_job_info_first_batch_replica_set" {
  value = lookup(local.component_record_job_info, "first_batch_replica", null) != null
}

output "is_job_info_replica_set" {
  value = lookup(local.component_record_job_info, "replica", null) != null
}

output "is_job_info_remaining_batch_set" {
  value = lookup(local.component_record_job_info, "remaining_batch", null) != null
}
`, acceptance.HW_CCE_CLUSTER_ID, name, acceptance.HW_BUILD_IMAGE_URL)
}
