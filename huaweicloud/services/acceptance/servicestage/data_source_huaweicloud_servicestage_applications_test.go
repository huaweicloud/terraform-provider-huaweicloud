package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplications_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_servicestage_applications.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplications_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_target_application_queried", "true"),
					resource.TestCheckOutput("is_target_application_name_right", "true"),
					resource.TestCheckOutput("is_target_application_description_right", "true"),
					resource.TestCheckOutput("is_target_application_eps_id_right", "true"),
					resource.TestCheckOutput("is_target_application_creator_set", "true"),
					resource.TestCheckOutput("is_target_application_created_at_set", "true"),
					resource.TestCheckOutput("is_target_application_updated_at_set", "true"),
					resource.TestCheckOutput("is_target_application_environments_set", "true"),
				),
			},
		},
	})
}

func testAccDataApplications_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestage_application" "test" {
  name        = "%[2]s"
  description = "Created by terraform test"

  enterprise_project_id = "%[3]s"

  environment {
    id = huaweicloud_servicestage_environment.test.id

    variable {
      name  = "name_1"
      value = "value_1"
    }
    variable {
      name  = "name_2"
      value = "abcdefghijklmnopqrstuvwxyz"
    }
    variable {
      name  = "name_3"
      value = "1234567890"
    }
  }
}

data "huaweicloud_servicestage_applications" "test" {
  depends_on = [huaweicloud_servicestage_application.test]
}

locals {
  app_query_result = [
    for v in data.huaweicloud_servicestage_applications.test.applications : v if v.id == huaweicloud_servicestage_application.test.id
  ]
}

output "is_target_application_queried" {
  value = length(local.app_query_result) == 1
}

output "is_target_application_name_right" {
  value = try(local.app_query_result[0].name == huaweicloud_servicestage_application.test.name, false)
}

output "is_target_application_description_right" {
  value = try(local.app_query_result[0].description == huaweicloud_servicestage_application.test.description, false)
}

output "is_target_application_eps_id_right" {
  value = try(local.app_query_result[0].enterprise_project_id == huaweicloud_servicestage_application.test.enterprise_project_id, false)
}

output "is_target_application_creator_set" {
  value = try(local.app_query_result[0].creator != "", false)
}

output "is_target_application_created_at_set" {
  value = try(local.app_query_result[0].created_at != "", false)
}

output "is_target_application_updated_at_set" {
  value = try(local.app_query_result[0].updated_at != "", false)
}

output "is_target_application_environments_set" {
  value = try(length(local.app_query_result[0].environments[0].variables) == 3, false)
}
`, testAccApplication_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
