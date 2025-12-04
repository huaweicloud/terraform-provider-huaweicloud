package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAttCkStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_event_att_ck_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceAttCkStatistics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.att_ck"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.num"),
				),
			},
		},
	})
}

const testAccDataSourceAttCkStatistics_basic = `
data "huaweicloud_hss_event_att_ck_statistics" "test" {
  category              = "host"
  enterprise_project_id = "0"
}
`
