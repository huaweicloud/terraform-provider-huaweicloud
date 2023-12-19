package swr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizations_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_swr_organizations.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizations_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "organizations.0.name"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful_not_found", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizations_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_swr_organizations" "test" {}

data "huaweicloud_swr_organizations" "filter_by_name" {
  name = data.huaweicloud_swr_organizations.test.organizations[0].name
}
output "is_name_filter_useful" {
  value = length(data.huaweicloud_swr_organizations.filter_by_name.organizations) == 1
}

data "huaweicloud_swr_organizations" "filter_by_name_not_found" {
  name = "%s"
}
output "is_name_filter_useful_not_found" {
  value = length(data.huaweicloud_swr_organizations.filter_by_name_not_found.organizations) == 0
}
`, rName)
}
