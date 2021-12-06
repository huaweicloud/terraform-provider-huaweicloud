package dms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDmsMaintainWindowV1DataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dms_maintainwindow_v1.maintainwindow1"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDmsMaintainWindowV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "seq", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "begin", "22"),
				),
			},
		},
	})
}

var testAccDmsMaintainWindowV1DataSource_basic = fmt.Sprintf(`
data "huaweicloud_dms_maintainwindow_v1" "maintainwindow1" {
  seq = 1
}
`)
