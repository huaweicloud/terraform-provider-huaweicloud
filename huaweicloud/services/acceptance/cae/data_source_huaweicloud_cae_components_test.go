package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceComponents_basic(t *testing.T) {
	all := "data.huaweicloud_cae_components.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCaeApplication(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceComponents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_component_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_name_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_annotations_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_spec_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_created_at_set_and_valid", "true"),
					resource.TestCheckOutput("is_component_updated_at_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataSourceComponents_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_repositories" "test" {}

locals {
  swr_repositories = [for v in data.huaweicloud_swr_repositories.test.repositories : v if length(v.tags) > 0][0]
}

resource "huaweicloud_cae_component" "test" {
  environment_id = "%[1]s"
  application_id = "%[2]s"

  metadata {
    name = "%[3]s"

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
`, acceptance.HW_CAE_ENVIRONMENT_ID, acceptance.HW_CAE_APPLICATION_ID, name)
}

func testAccDataSourceComponents_basic() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cae_components" "test" {
  depends_on = [
    huaweicloud_cae_component.test
  ]

  environment_id = "%[2]s"
  application_id = "%[3]s"
}

locals {
  component_id            = huaweicloud_cae_component.test.id
  component_filter_result = try([
    for v in data.huaweicloud_cae_components.test.components : v if v.id == local.component_id
  ][0], null)
}

output "is_component_id_set_and_valid" {
  value = local.component_filter_result != null
}

output "is_component_name_set_and_valid" {
  value = try(local.component_filter_result.name == huaweicloud_cae_component.test.metadata[0].name, false)
}

output "is_component_annotations_set_and_valid" {
  value = try(alltrue([
    length(local.component_filter_result.annotations) > 0,
    local.component_filter_result.annotations["version"] == huaweicloud_cae_component.test.metadata[0].annotations["version"]
  ]), false)
}

output "is_component_spec_set_and_valid" {
  value = try(alltrue([
    length(local.component_filter_result.spec) > 0,
    local.component_filter_result.spec[0].runtime != "",
    local.component_filter_result.spec[0].environment_id != "",
    local.component_filter_result.spec[0].replica != 0,
    local.component_filter_result.spec[0].available_replica >= 0,
    local.component_filter_result.spec[0].source != "",
    local.component_filter_result.spec[0].build != "",
    local.component_filter_result.spec[0].resource_limit != "",
    local.component_filter_result.spec[0].image_url != "",
    local.component_filter_result.spec[0].status != "",
  ]), false)
}

output "is_component_created_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result.created_at)) > 0, false)
}

output "is_component_updated_at_set_and_valid" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.component_filter_result.updated_at)) > 0, false)
}
`, testAccDataSourceComponents_base(name),
		acceptance.HW_CAE_ENVIRONMENT_ID,
		acceptance.HW_CAE_APPLICATION_ID)
}
