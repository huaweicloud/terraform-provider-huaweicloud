package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourceLockedStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_resource_locked_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourceLockedStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "locked_status"),
				),
			},
		},
	})
}

const testDataSourceResourceLockedStatus_base string = `
resource "huaweicloud_hss_quota" "test" {
  version               = "hss.version.premium"
  period_unit           = "month"
  period                = 1
  auto_renew            = "true"
  enterprise_project_id = "0"
}
`

func testDataSourceResourceLockedStatus_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_resource_locked_status" "test" {
  resource_id           = huaweicloud_hss_quota.test.id
  enterprise_project_id = "0"
}
`, testDataSourceResourceLockedStatus_base)
}
