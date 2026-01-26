package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectServices_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_enterprise_project_services.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.service"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.service_i18n_display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.resource_type_i18n_display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "services.0.resource_types.0.regions.#"),

					resource.TestCheckOutput("huaweicloud_enterprise_project_services_locale", "true"),
					resource.TestCheckOutput("huaweicloud_enterprise_project_services_service", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectServices_basic = `
data "huaweicloud_enterprise_project_services" "test" {}

output "huaweicloud_enterprise_project_services" {
  value = length(data.huaweicloud_enterprise_project_services.test.services) > 0
}

data "huaweicloud_enterprise_project_services" "test_locale" {
  locale  = "en-us"
}

output "huaweicloud_enterprise_project_services_locale" {
  value = length(data.huaweicloud_enterprise_project_services.test_locale.services) > 0
}

data "huaweicloud_enterprise_project_services" "test_service" {
  service  = "vpc"
}

output "huaweicloud_enterprise_project_services_service" {
  value = length(data.huaweicloud_enterprise_project_services.test_service.services) > 0
}
`
