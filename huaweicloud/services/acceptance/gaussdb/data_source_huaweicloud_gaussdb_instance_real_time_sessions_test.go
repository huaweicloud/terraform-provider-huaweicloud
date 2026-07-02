package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBInstanceRealTimeSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_real_time_sessions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBInstanceRealTimeSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.session_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.pid"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.state"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBInstanceRealTimeSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = "%s"
  node_id      = ""
  component_id = ""
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
