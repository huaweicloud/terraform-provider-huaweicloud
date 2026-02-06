package ccm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmPrivateCertificateTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_private_certificate_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmPrivateCertificateTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
				),
			},
		},
	})
}

const testDataSourceDataSourceCcmPrivateCertificateTags_basic = `
data "huaweicloud_ccm_private_certificate_tags" "test" {
}
`
