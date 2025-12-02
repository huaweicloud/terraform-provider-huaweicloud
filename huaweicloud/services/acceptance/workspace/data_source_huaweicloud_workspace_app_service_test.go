package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppService_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_workspace_app_service.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppService_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_open_with_ad_set_and_valid", "true"),
					resource.TestCheckResourceAttrSet(dcName, "service_status"),
					resource.TestCheckResourceAttrSet(dcName, "tenant_domain_id"),
					resource.TestCheckResourceAttrSet(dcName, "tenant_domain_name"),
					resource.TestMatchResourceAttr(dcName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

const testAccDataAppService_basic string = `
data "huaweicloud_workspace_app_service" "test" {}

output "is_open_with_ad_set_and_valid" {
  value = data.huaweicloud_workspace_app_service.test.open_with_ad != null
}
`
