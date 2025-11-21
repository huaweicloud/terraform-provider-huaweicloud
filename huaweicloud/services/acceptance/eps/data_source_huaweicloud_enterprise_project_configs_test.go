package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectConfigs_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_configs.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectConfigs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("huaweicloud_enterprise_project_configs", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectConfigs_basic = `
data "huaweicloud_enterprise_project_configs" "test" {}

output "huaweicloud_enterprise_project_configs" {
  value = data.huaweicloud_enterprise_project_configs.test.support_item.delete_ep_support == true
}
`
