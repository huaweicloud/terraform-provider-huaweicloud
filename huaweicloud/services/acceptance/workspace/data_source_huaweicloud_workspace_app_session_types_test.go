package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppSessionTypes_basic(t *testing.T) {
	rName := "data.huaweicloud_workspace_app_session_types.all"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppSessionTypes_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_spec_code_set", "true"),
					resource.TestCheckOutput("is_session_type_set", "true"),
					resource.TestCheckOutput("is_resource_type_set", "true"),
					resource.TestCheckOutput("is_cloud_service_type_set", "true"),
				),
			},
		},
	})
}

const testAccDataAppSessionTypes_basic string = `
data "huaweicloud_workspace_app_session_types" "all" {}

locals {
  session_types      = data.huaweicloud_workspace_app_session_types.all.session_types
  first_session_type = try(local.session_types[0], {})
}

output "is_resource_spec_code_set" {
  value = length(local.session_types) != 0 ? try(local.first_session_type.resource_spec_code != "", false) : true 
}

output "is_session_type_set" {
  value = length(local.session_types) != 0 ? try(local.first_session_type.session_type != "", false) : true 
}

output "is_resource_type_set" {
  value = length(local.session_types) != 0 ? try(local.first_session_type.resource_type != "", false) : true 
}

output "is_cloud_service_type_set" {
  value = length(local.session_types) != 0 ? try(local.first_session_type.cloud_service_type != "", false) : true 
}
`
