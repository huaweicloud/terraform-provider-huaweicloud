package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscShowTasks_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dsc_show_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscShowTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "mask_task_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scan_task_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "watermark_embed_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "watermark_extract_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "watermark_task_num"),
				),
			},
		},
	})
}

const testAccDataSourceDscShowTasks_basic = `
data "huaweicloud_dsc_show_tasks" "test" {}
`
