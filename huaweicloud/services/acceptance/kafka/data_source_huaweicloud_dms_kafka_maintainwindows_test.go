package kafka

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataMaintainWindows_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_kafka_maintainwindows.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMaintainWindows_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "maintain_windows.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "maintain_windows.0.default"),
					resource.TestCheckResourceAttrSet(dataSourceName, "maintain_windows.0.begin"),
					resource.TestCheckResourceAttrSet(dataSourceName, "maintain_windows.0.end"),
					resource.TestMatchResourceAttr(dataSourceName, "maintain_windows.0.seq", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataMaintainWindows_basic = `
data "huaweicloud_dms_kafka_maintainwindows" "test" {}
`
