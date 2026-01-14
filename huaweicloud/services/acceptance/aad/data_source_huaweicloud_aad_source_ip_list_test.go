package antiddos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAadSourceIpList_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_source_ip_list.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAadSourceIpList_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ips.0.data_center"),
					resource.TestCheckResourceAttrSet(dataSource, "ips.0.isp"),
					resource.TestCheckResourceAttrSet(dataSource, "ips.0.ip.#"),
				),
			},
		},
	})
}

const testDataSourceDataSourceAadSourceIpList_basic = `
data "huaweicloud_aad_source_ip_list" "test" {
}
`
