package secmaster

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccDataSourceSecmasterPlaybookMonitors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_playbook_monitors.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecmasterPlaybookMonitors_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceSecmasterPlaybookMonitors_basic(name string) string {
	nowStamp := time.Now().Unix()
	startTime := utils.FormatTimeStampRFC3339(nowStamp-24*60*60, false, "2006-01-02T15:04:05.000Z+0800")
	endTime := utils.FormatTimeStampRFC3339(nowStamp, false, "2006-01-02T15:04:05.000Z+0800")
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_secmaster_playbook_monitors" "test" {
  workspace_id       = "%[2]s"
  playbook_id        = huaweicloud_secmaster_playbook.test.id
  start_time         = "%[3]s"
  end_time           = "%[4]s"
  version_query_type = "ALL"

  depends_on = [huaweicloud_secmaster_playbook_enable.test]
}
`, testPlaybookEnable_basic(name), acceptance.HW_SECMASTER_WORKSPACE_ID, startTime, endTime)
}
