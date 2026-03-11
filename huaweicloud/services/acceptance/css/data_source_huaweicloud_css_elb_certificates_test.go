package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssElbCertificates_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_elb_certificates.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCSSClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssElbCertificates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "certificates.0.type"),
				),
			},
		},
	})
}

func testAccDataSourceCssElbCertificates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_elb_certificates" "test" {
  cluster_id = "%[1]s"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
