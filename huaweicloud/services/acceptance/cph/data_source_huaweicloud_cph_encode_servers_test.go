package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphEncodeServers_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_encode_servers.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCphEncodeServers_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// no test data
				),
			},
		},
	})
}

const testDataSourceCphEncodeServers_basic = `data "huaweicloud_cph_encode_servers" "test" {}`
