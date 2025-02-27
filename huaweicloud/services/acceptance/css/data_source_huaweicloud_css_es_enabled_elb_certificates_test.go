package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssEsEnabledElbCertificates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_es_enabled_elb_certificates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCssEsEnabledElbCertificates_basic(rName),
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

func testDataSourceCssEsEnabledElbCertificates_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_css_es_enabled_elb_certificates" "test" {
  cluster_id = huaweicloud_css_cluster.test.id

  depends_on = [huaweicloud_css_es_loadbalancer_config.test]
}
`, testAccCssEsLoadbalancerConfig_basic(name))
}
