package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppService_basic(t *testing.T) {
	dataSource := "data.huaweicloud_workspace_app_service.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppService_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_open_with_ad_set_and_valid", "true"),
					resource.TestCheckResourceAttrSet(dataSource, "service_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tenant_domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tenant_domain_name"),
					resource.TestMatchResourceAttr(dataSource, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

const testAccDataSourceAppService_basic = `
data "huaweicloud_workspace_app_service" "test" {}

output "is_open_with_ad_set_and_valid" {
  value = data.huaweicloud_workspace_app_service.test.open_with_ad != null
}
`
