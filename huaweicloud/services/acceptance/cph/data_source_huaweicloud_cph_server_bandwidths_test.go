package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphServerBandwidths_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_server_bandwidths.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCphServerBandwidths_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// This data source only supports the broadband list of system-defined networks, and test data cannot be constructed.
				),
			},
		},
	})
}

const testDataSourceCphServerBandwidths_basic = `data "huaweicloud_cph_server_bandwidths" "test" {}`
