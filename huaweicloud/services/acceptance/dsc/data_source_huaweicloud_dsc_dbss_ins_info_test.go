package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDbssInsInfo_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_dbss_ins_info.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscDbssInsInfo_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "dbss_instance_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dbss_rds_database_list.#"),
				),
			},
		},
	})
}

const testDataSourceDscDbssInsInfo_basic = `
data "huaweicloud_dsc_dbss_ins_info" "test" {}
`
