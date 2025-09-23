package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponents_basic(t *testing.T) {
	var (
		allWithoutEps = "data.huaweicloud_cae_components.test.0"
		dcWithoutEps  = acceptance.InitDataSourceCheck(allWithoutEps)

		allWithEps = "data.huaweicloud_cae_components.test.1"
		dcWithEps  = acceptance.InitDataSourceCheck(allWithEps)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComponents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dcWithoutEps.CheckResourceExists(),
					resource.TestCheckOutput("is_component_id_set_and_valid_without_eps", "true"),
					resource.TestCheckOutput("is_component_name_set_and_valid_without_eps", "true"),
					resource.TestCheckOutput("is_component_annotations_set_and_valid_without_eps", "true"),
					resource.TestCheckOutput("is_component_spec_set_and_valid_without_eps", "true"),
					resource.TestCheckOutput("is_component_created_at_set_and_valid_without_eps", "true"),
					resource.TestCheckOutput("is_component_updated_at_set_and_valid_without_eps", "true"),
					dcWithEps.CheckResourceExists(),
					resource.TestCheckOutput("is_component_id_set_and_valid_with_eps", "true"),
					resource.TestCheckOutput("is_component_name_set_and_valid_with_eps", "true"),
					resource.TestCheckOutput("is_component_annotations_set_and_valid_with_eps", "true"),
					resource.TestCheckOutput("is_component_spec_set_and_valid_with_eps", "true"),
					resource.TestCheckOutput("is_component_created_at_set_and_valid_with_eps", "true"),
					resource.TestCheckOutput("is_component_updated_at_set_and_valid_with_eps", "true"),
				),
			},
		},
	})
}

func testAccDataSourceComponents_base(name string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

# Query the enterprise project ID of the environment.
data "huaweicloud_cae_environments" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  enterprise_project_id = count.index == 1 ? "%[2]s" : null
}

# The first application belongs to the default enterprise project.
# The second application belongs to the non-default enterprise project.
resource "huaweicloud_cae_application" "test" {
  count = 2

  environment_id        = try(data.huaweicloud_cae_environments.test[count.index].environments[0].id, "NOT_FOUND")
  name                  = format("%[3]s-%%d", count.index)
  enterprise_project_id = try(data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id, null)
}

data "huaweicloud_swr_repositories" "test" {}

locals {
  swr_repositories = [for v in data.huaweicloud_swr_repositories.test.repositories : v if length(v.tags) > 0][0]
}

resource "huaweicloud_cae_component" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  application_id        = huaweicloud_cae_application.test[count.index].id
  enterprise_project_id = try(data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id, null)

  metadata {
    name = format("%[3]s-%%d", count.index)

    annotations = {
      version = "1.0.0"
    }
  }

  spec {
    replica = 1
    runtime = "Docker"

    source {
      type = "image"
      url  = format("%%s:%%s", local.swr_repositories.path, local.swr_repositories.tags[0])
    }

    resource_limit {
      cpu    = "500m"
      memory = "1Gi"
    }
  }
}
`, acceptance.HW_CAE_ENVIRONMENT_IDS, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccDataSourceComponents_basic() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cae_components" "test" {
  depends_on = [
    huaweicloud_cae_component.test
  ]

  count = 2

  environment_id        = local.env_ids[count.index]
  application_id        = huaweicloud_cae_application.test[count.index].id
  enterprise_project_id = try(data.huaweicloud_cae_environments.test[count.index].environments[0].annotations.enterprise_project_id, null)
}

locals {
  component_id_without_eps            = huaweicloud_cae_component.test[0].id
  component_filter_result_without_eps = try([
    for v in data.huaweicloud_cae_components.test[0].components : v if v.id == local.component_id_without_eps
  ][0], null)

  component_id_with_eps            = huaweicloud_cae_component.test[1].id
  component_filter_result_with_eps = try([
    for v in data.huaweicloud_cae_components.test[1].components : v if v.id == local.component_id_with_eps
  ][0], null)
}

output "is_component_id_set_and_valid_without_eps" {
  value = local.component_filter_result_without_eps != null
}

output "is_component_name_set_and_valid_without_eps" {
  value = try(local.component_filter_result_without_eps.name == huaweicloud_cae_component.test[0].metadata[0].name, false)
}

output "is_component_annotations_set_and_valid_without_eps" {
  value = try(alltrue([
    length(local.component_filter_result_without_eps.annotations) > 0,
    local.component_filter_result_without_eps.annotations["version"] == huaweicloud_cae_component.test[0].metadata[0].annotations["version"]
  ]), false)
}

output "is_component_spec_set_and_valid_without_eps" {
  value = try(alltrue([
    length(local.component_filter_result_without_eps.spec) > 0,
    local.component_filter_result_without_eps.spec[0].runtime != "",
    local.component_filter_result_without_eps.spec[0].environment_id != "",
    local.component_filter_result_without_eps.spec[0].replica != 0,
    local.component_filter_result_without_eps.spec[0].available_replica >= 0,
    local.component_filter_result_without_eps.spec[0].source != "",
    local.component_filter_result_without_eps.spec[0].build != "",
    local.component_filter_result_without_eps.spec[0].resource_limit != "",
    local.component_filter_result_without_eps.spec[0].image_url != "",
    local.component_filter_result_without_eps.spec[0].status != "",
  ]), false)
}

output "is_component_created_at_set_and_valid_without_eps" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result_without_eps.created_at)) > 0, false)
}

output "is_component_updated_at_set_and_valid_without_eps" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result_without_eps.updated_at)) > 0, false)
}

output "is_component_id_set_and_valid_with_eps" {
  value = local.component_filter_result_with_eps != null
}

output "is_component_name_set_and_valid_with_eps" {
  value = try(local.component_filter_result_with_eps.name == huaweicloud_cae_component.test[1].metadata[0].name, false)
}

output "is_component_annotations_set_and_valid_with_eps" {
  value = try(alltrue([
    length(local.component_filter_result_with_eps.annotations) > 0,
    local.component_filter_result_with_eps.annotations["version"] == huaweicloud_cae_component.test[1].metadata[0].annotations["version"]
  ]), false)
}

output "is_component_spec_set_and_valid_with_eps" {
  value = try(alltrue([
    length(local.component_filter_result_with_eps.spec) > 0,
    local.component_filter_result_with_eps.spec[0].runtime != "",
    local.component_filter_result_with_eps.spec[0].environment_id != "",
    local.component_filter_result_with_eps.spec[0].replica != 0,
    local.component_filter_result_with_eps.spec[0].available_replica >= 0,
    local.component_filter_result_with_eps.spec[0].source != "",
    local.component_filter_result_with_eps.spec[0].build != "",
    local.component_filter_result_with_eps.spec[0].resource_limit != "",
    local.component_filter_result_with_eps.spec[0].image_url != "",
    local.component_filter_result_with_eps.spec[0].status != "",
  ]), false)
}

output "is_component_created_at_set_and_valid_with_eps" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result_with_eps.created_at)) > 0, false)
}

output "is_component_updated_at_set_and_valid_with_eps" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result_with_eps.updated_at)) > 0, false)
}
`, testAccDataSourceComponents_base(name))
}
