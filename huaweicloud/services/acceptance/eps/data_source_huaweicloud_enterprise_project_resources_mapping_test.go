package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectResourcesMapping_basic(t *testing.T) {
	all := "data.huaweicloud_enterprise_project_resources_mapping.test"
	dc := acceptance.InitDataSourceCheck(all)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectResourcesMapping_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "resource_mapping.vpcs", "vpc"),
					resource.TestCheckOutput("huaweicloud_enterprise_project_resources_mapping", "true"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectResourcesMapping_basic = `
data "huaweicloud_enterprise_project_resources_mapping" "test" {}

output "huaweicloud_enterprise_project_resources_mapping" {
  value = length(data.huaweicloud_enterprise_project_resources_mapping.test.resource_mapping) > 0
}
`
