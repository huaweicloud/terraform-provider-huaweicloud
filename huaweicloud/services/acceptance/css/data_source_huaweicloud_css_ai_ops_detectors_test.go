package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAiOpsDetectors_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_css_ai_ops_detectors.test"
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
				Config: testAccDataSourceAiOpsDetectors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "full_detection.#"),
					resource.TestCheckResourceAttrSet(dataSource, "full_detection.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_detection.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "full_detection.0.desc"),
					resource.TestCheckResourceAttrSet(dataSource, "unavailability_detection.#"),
					resource.TestCheckResourceAttrSet(dataSource, "unavailability_detection.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "unavailability_detection.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "unavailability_detection.0.desc"),
				),
			},
		},
	})
}

func testAccDataSourceAiOpsDetectors_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_css_ai_ops_detectors" "test" {
  cluster_id = "%[1]s"
}
`, acceptance.HW_CSS_CLUSTER_ID)
}
