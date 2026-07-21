package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscWatermarkEmbedLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_watermark_embed_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscWatermarkEmbedLogs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.dest_url"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.doc_path"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.file_exist"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.file_url"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.finish_num"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_embed_logs.0.total_file_num"),
				),
			},
		},
	})
}

const testDataSourceDscWatermarkEmbedLogs_basic = `
data "huaweicloud_dsc_watermark_embed_logs" "test" {}
`
