package cbh

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceAdminUrl_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cbh_instance_admin_url.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Ensure the CBH instance ID is set before running the test
			acceptance.TestAccPreCheckCbhInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceInstanceAdminUrl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "admin_url"),
				),
			},
		},
	})
}

func testDataSourceInstanceAdminUrl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cbh_instance_admin_url" "test" {
  server_id = "%s"
}
`, acceptance.HW_CBH_INSTANCE_ID)
}
