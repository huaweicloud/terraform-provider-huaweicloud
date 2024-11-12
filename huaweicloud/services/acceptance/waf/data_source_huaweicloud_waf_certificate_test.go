package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourceCertificate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		datasourceName = "data.huaweicloud_waf_certificate.test"
		dc             = acceptance.InitDataSourceCheck(datasourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceWafCertificate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttrSet(datasourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "created_at"),
				),
			},
		},
	})
}

func testAccDatasourceWafCertificate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_certificate" "test" {
  name                  = huaweicloud_waf_certificate.test.name
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_certificate.test
  ]
}
`, testAccWafCertificate_basic(name, generateCertificateBody()), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
