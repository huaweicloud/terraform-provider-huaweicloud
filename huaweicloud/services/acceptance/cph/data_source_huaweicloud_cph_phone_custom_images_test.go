package cph

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphPhoneCustomImages_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_phone_custom_images.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCphPhoneCustomImages_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					// no test data
				),
			},
		},
	})
}

const testDataSourceCphPhoneCustomImages_basic = `data "huaweicloud_cph_phone_custom_images" "test" {}`
