package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBRtsOverview_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_rts_overview.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBRtsOverview_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "active_num"),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_num"),
					resource.TestCheckResourceAttrSet(dataSource, "lock_num"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBRtsOverview_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_rts_overview" "test" {
  instance_id = "%s"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
