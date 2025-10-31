package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBlueprint_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_blueprint.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCAccountID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBlueprint_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_product_id"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_product_name"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_product_version"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_status"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_product_impl_type"),
					resource.TestCheckResourceAttrSet(dataSource, "blueprint_product_impl_detail"),
				),
			},
		},
	})
}

func testAccDataSourceBlueprint_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_blueprint" "test" {
  managed_account_id = "%[1]s"
}`, acceptance.HW_RGC_ACCOUNT_ID)
}
