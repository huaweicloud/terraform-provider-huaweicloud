package identitycenter

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIdentityCenterApplicationTemplates_basic(t *testing.T) {
	rName := "data.huaweicloud_identitycenter_application_templates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIdentityCenterApplicationTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "application_templates.0.application_type", ""),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.description"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.sso_protocol"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.template_id"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.template_version"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.response_config"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.response_schema_config"),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.security_config.#"),
					resource.TestCheckResourceAttr(rName, "application_templates.0.security_config.0.ttl", ""),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.service_provider_config.#"),
					resource.TestCheckResourceAttr(rName, "application_templates.0.service_provider_config.0.audience", ""),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.service_provider_config.0.require_request_signature"),
					resource.TestCheckResourceAttr(rName, "application_templates.0.service_provider_config.0.start_url", ""),
					resource.TestCheckResourceAttrSet(rName, "application_templates.0.service_provider_config.0.consumers.#")),
			},
		},
	})
}

const testAccDatasourceIdentityCenterApplicationTemplates_basic = `
data "huaweicloud_identitycenter_catalog_applications" "test"{}
 
data "huaweicloud_identitycenter_application_templates" "test"{
  application_id = data.huaweicloud_identitycenter_catalog_applications.test.applications[0].application_id
}
`
