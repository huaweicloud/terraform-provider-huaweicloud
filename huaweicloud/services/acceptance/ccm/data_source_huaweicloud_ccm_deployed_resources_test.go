package ccm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcmDeployeds_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ccm_deployed_resources.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCCMSSLCertificateId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCcmDeployeds_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.certificate_id"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.deployed_resources.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.deployed_resources.0.resource_num"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.deployed_resources.0.resource_location"),
					resource.TestCheckResourceAttrSet(dataSource, "results.0.deployed_resources.0.region_resources.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCcmDeployeds_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_ccm_deployed_resources" "test" {
  service_names   = ["WAF"]
  certificate_ids = ["%s"]
}
`, acceptance.HW_CCM_SSL_CERTIFICATE_ID)
}
