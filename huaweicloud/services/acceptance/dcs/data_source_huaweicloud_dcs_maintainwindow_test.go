package dcs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsMaintainWindowDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_maintainwindow.test"
	seqFilterDataSource := "data.huaweicloud_dcs_maintainwindow.seq_filter"
	beginFilterDataSource := "data.huaweicloud_dcs_maintainwindow.begin_filter"
	endFilterDataSource := "data.huaweicloud_dcs_maintainwindow.end_filter"
	defaultFilterDataSource := "data.huaweicloud_dcs_maintainwindow.default_filter"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsMaintainWindowDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "seq"),
					resource.TestCheckResourceAttrSet(dataSource, "begin"),
					resource.TestCheckResourceAttrSet(dataSource, "end"),
					resource.TestCheckResourceAttrSet(dataSource, "default"),
					resource.TestCheckResourceAttrSet(seqFilterDataSource, "seq"),
					resource.TestCheckResourceAttrSet(seqFilterDataSource, "begin"),
					resource.TestCheckResourceAttrSet(seqFilterDataSource, "end"),
					resource.TestCheckResourceAttrSet(seqFilterDataSource, "default"),
					resource.TestCheckResourceAttrSet(beginFilterDataSource, "seq"),
					resource.TestCheckResourceAttrSet(beginFilterDataSource, "begin"),
					resource.TestCheckResourceAttrSet(beginFilterDataSource, "end"),
					resource.TestCheckResourceAttrSet(beginFilterDataSource, "default"),
					resource.TestCheckResourceAttrSet(endFilterDataSource, "seq"),
					resource.TestCheckResourceAttrSet(endFilterDataSource, "begin"),
					resource.TestCheckResourceAttrSet(endFilterDataSource, "end"),
					resource.TestCheckResourceAttrSet(endFilterDataSource, "default"),
					resource.TestCheckResourceAttrSet(defaultFilterDataSource, "seq"),
					resource.TestCheckResourceAttrSet(defaultFilterDataSource, "begin"),
					resource.TestCheckResourceAttrSet(defaultFilterDataSource, "end"),
					resource.TestCheckResourceAttrSet(defaultFilterDataSource, "default"),
				),
			},
		},
	})
}

func testAccDcsMaintainWindowDataSource_basic() string {
	return `
data "huaweicloud_dcs_maintainwindow" "test" {}

data "huaweicloud_dcs_maintainwindow" "seq_filter" {
  seq = 1
}

data "huaweicloud_dcs_maintainwindow" "begin_filter" {
  begin = 10
}

data "huaweicloud_dcs_maintainwindow" "end_filter" {
  end = 15
}

data "huaweicloud_dcs_maintainwindow" "default_filter" {
  default = true
}
`
}
