package eps

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataEnterpriseProjectResourcesMapping_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_enterprise_project_resources_mapping.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEnterpriseProjectResourcesMapping_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_mapping.%"),
				),
			},
		},
	})
}

const testAccDataEnterpriseProjectResourcesMapping_basic = `
data "huaweicloud_enterprise_project_resources_mapping" "test" {}
`
