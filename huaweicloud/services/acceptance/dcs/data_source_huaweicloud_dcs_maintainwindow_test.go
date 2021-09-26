package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsMaintainWindowDataSource_basic(t *testing.T) {
	sourceName := "data.huaweicloud_dcs_maintainwindow.maintainwindow1"
	dc := acceptance.InitDataSourceCheck(sourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsMaintainWindowDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(sourceName, "seq", "1"),
					resource.TestCheckResourceAttr(sourceName, "begin", "22"),
				),
			},
		},
	})
}

var testAccDcsMaintainWindowDataSource_basic = fmt.Sprintf(`
data "huaweicloud_dcs_maintainwindow" "maintainwindow1" {
seq = 1
default = false
}
`)
