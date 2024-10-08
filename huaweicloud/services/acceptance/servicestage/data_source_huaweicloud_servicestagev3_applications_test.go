package servicestage

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3Applications_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_servicestagev3_applications.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3Applications_basic_step1(name),
			},
			{
				// Update the application name and make sure the attribute 'updated_at' not empty.
				Config: testAccDataV3Applications_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`[1-9]\d*`)),
					resource.TestCheckOutput("is_application_id_set", "true"),
					resource.TestCheckOutput("is_application_name_set", "true"),
					resource.TestCheckOutput("is_application_description_set", "true"),
					resource.TestCheckOutput("is_application_eps_id_set", "true"),
					resource.TestCheckOutput("is_application_creator_set", "true"),
					resource.TestCheckOutput("is_application_created_at_set", "true"),
					resource.TestCheckOutput("is_application_updated_at_set", "true"),
				),
			},
		},
	})
}

func testAccDataV3Applications_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_servicestagev3_application" "test" {
  name                  = "%[1]s"
  description           = "Created by terraform test"
  enterprise_project_id = "%[2]s"

  tags = {
    foo   = "bar"
    owner = "terraform"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataV3Applications_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_servicestagev3_applications" "test" {
  depends_on = [
    huaweicloud_servicestagev3_application.test
  ]
}

locals {
  application_id            = huaweicloud_servicestagev3_application.test.id
  application_filter_result = try([
    for v in data.huaweicloud_servicestagev3_applications.test.applications : v if v.id == local.application_id
  ][0], null)
}

output "is_application_id_set" {
  value = local.application_filter_result != null
}

output "is_application_name_set" {
  value = try(local.application_filter_result.name == huaweicloud_servicestagev3_application.test.name, false)
}

output "is_application_description_set" {
  value = try(local.application_filter_result.description == huaweicloud_servicestagev3_application.test.description, false)
}

output "is_application_eps_id_set" {
  value = try(local.application_filter_result.enterprise_project_id == huaweicloud_servicestagev3_application.test.enterprise_project_id, false)
}

output "is_application_creator_set" {
  value = try(local.application_filter_result.creator != "", false)
}

output "is_application_created_at_set" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.application_filter_result.created_at)) > 0, false)
}

output "is_application_updated_at_set" {
  value = try(length(regexall("^\\d{4}\\-\\d{2}\\-\\d{2}T\\d{2}:\\d{2}:\\d{2}(?:Z|[+-]\\d{2}:\\d{2})$",
    local.application_filter_result.updated_at)) > 0, false)
}
`, testAccDataV3Applications_basic_step1(name))
}
