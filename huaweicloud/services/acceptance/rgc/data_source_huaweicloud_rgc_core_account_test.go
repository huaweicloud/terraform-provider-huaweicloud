package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCoreAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_core_account.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCCoreAccountType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCoreAccounts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "core_resource_mappings"),
				),
			},
		},
	})
}

func testDataSourceCoreAccounts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_core_account" "test" {
  account_type = "%[1]s"
}
`, acceptance.HW_RGC_CORE_ACCOUNT_TYPE)
}
