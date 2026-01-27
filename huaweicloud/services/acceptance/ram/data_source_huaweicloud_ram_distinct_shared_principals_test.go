package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDistinctSharedPrincipals_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_distinct_shared_principals.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDistinctSharedPrincipals_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "distinct_shared_principals.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "distinct_shared_principals.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceDistinctSharedPrincipals_basic() string {
	return `
data "huaweicloud_ram_distinct_shared_principals" "test" {
  resource_owner = "other-accounts"
}
`
}
