package codeartsinspector

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func TestAccDataSourceCodeartsInspectorWebsiteScanTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_codearts_inspector_website_scan_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)
	timer := utils.FormatTimeStampUTC(time.Now().Add(48 * time.Hour).Unix())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCodeartsInspectorWebsiteScanTasks_basic(rName, timer),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.task_name", rName),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.task_type", "normal"),
					resource.TestCheckResourceAttrPair(dataSource, "tasks.0.url",
						"huaweicloud_codearts_inspector_website.test", "website_address"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.timer", timer),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.scan_mode", "deep"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.port_scan", "true"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.weak_pwd_scan", "false"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.cve_check", "false"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.text_check", "false"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.picture_check", "false"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.malicious_code", "false"),
					resource.TestCheckResourceAttr(dataSource, "tasks.0.malicious_link", "false"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.reason"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.pack_num"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.score"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.safe_level"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.high"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.middle"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.low"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.hint"),
				),
			},
		},
	})
}

func testDataSourceCodeartsInspectorWebsiteScanTasks_basic(name, timer string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_codearts_inspector_website_scan_tasks" "test" {
  depends_on = [huaweicloud_codearts_inspector_website_scan.test]
  
  domain_id = huaweicloud_codearts_inspector_website.test.id
}
`, testInspectorWebsiteScan_basic(name, timer))
}
