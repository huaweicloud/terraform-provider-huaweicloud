package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV3ComponentAction_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// At least one of JAR package must be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV3ComponentAction_basic_step1(),
			},
		},
	})
}

func testAccV3ComponentAction_basic_base() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_cce_clusters" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_servicestagev3_runtime_stacks" "test" {}

locals {
  docker_runtime_stack = try([for o in data.huaweicloud_servicestagev3_runtime_stacks.test.runtime_stacks:
    o if o.type == "Docker" && o.deploy_mode == "container"][0], {})
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
    id   = try(data.huaweicloud_cce_clusters.test.clusters[0].id, "")
    type = "cce"
  }
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
    name        = try(local.docker_runtime_stack.name, "Docker")
    type        = try(local.docker_runtime_stack.type, "Docker")
    version     = try(local.docker_runtime_stack.version, "1.0")
  }

  source = jsonencode({
    kind    = "image"
    storage = "swr"
    url     = try(element(split(",", "%[3]s"), 0), "")
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
`, acceptance.HW_CCE_CLUSTER_ID, name, acceptance.HW_SERVICESTAGE_JAR_PKG_STORAGE_URLS)
}

func testAccV3ComponentAction_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component_action" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  component_id   = huaweicloud_servicestagev3_component.test.id
  action         = "sync_workload"
}
`, testAccV3ComponentAction_basic_base())
}
