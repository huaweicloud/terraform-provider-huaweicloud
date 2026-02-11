package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmCsrPrivateKey_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_csr_private_key.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmCsrPrivateKey_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "private_key"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcmCsrPrivateKey_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ccm_csr_private_key" "test" {
  csr_id = huaweicloud_ccm_csr.test.id
}
`, testCCMCsr_basic(name))
}
