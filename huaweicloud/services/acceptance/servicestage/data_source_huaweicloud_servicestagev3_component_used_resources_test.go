package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3ComponentUsedResources_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_servicestagev3_component_used_resources.test"
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
				Config: testAccDataV3ComponentUsedResources_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "applications.0.id"),
					resource.TestCheckResourceAttrSet(all, "applications.0.label"),
					resource.TestMatchResourceAttr(all, "enterprise_projects.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestMatchResourceAttr(all, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckResourceAttrSet(all, "environments.0.id"),
					resource.TestCheckResourceAttrSet(all, "environments.0.label"),
				),
			},
		},
	})
}

func testAccDataV3ComponentUsedResources_basic(name string) string {
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

data "huaweicloud_servicestagev3_component_used_resources" "test" {
  depends_on = [
    huaweicloud_servicestagev3_component.test,
  ]
}
`, acceptance.HW_CCE_CLUSTER_ID, name, acceptance.HW_BUILD_IMAGE_URL)
}
