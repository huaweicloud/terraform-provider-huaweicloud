package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWatermarkExtractLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_watermark_extract_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWatermarkExtractLogs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "water_mark_log.#"),
				),
			},
		},
	})
}

const testDataSourceWatermarkExtractLogs_basic = `
data "huaweicloud_dsc_watermark_extract_logs" "test" {}
`
