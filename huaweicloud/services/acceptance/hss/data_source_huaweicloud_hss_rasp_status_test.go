package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with premium edition host protection enabled.
			// The host also set the application protection (RASP) at the same time.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRaspStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
				),
			},
		},
	})
}

func testAccDataSourceRaspStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_rasp_status" "test" {
  host_id               = "%[1]s"
  enterprise_project_id = "all_granted_eps"
  app_type              = "java"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
