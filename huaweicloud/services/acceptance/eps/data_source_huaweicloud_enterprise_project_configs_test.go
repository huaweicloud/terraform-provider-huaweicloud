package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectConfigs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_enterprise_project_configs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectConfigs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_item_attribute.0.delete_ep_support_attribute"),
					resource.TestCheckResourceAttrSet(dataSourceName, "support_item.delete_ep_support"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectConfigs_basic = `
data "huaweicloud_enterprise_project_configs" "test" {}
`
