package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcGlobalConnectionBandwidthTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_global_connection_bandwidth_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcGlobalConnectionBandwidthTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceCcGlobalConnectionBandwidthTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name        = "%[1]s"
  type        = "Region"  
  bordercross = false
  charge_mode = "bwd"
  size        = 5
  description = "test"
  sla_level   = "Ag"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testDataSourceCcGlobalConnectionBandwidthTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_global_connection_bandwidth_tags" "test" {
  depends_on = [huaweicloud_cc_global_connection_bandwidth.test]
}
`, testDataSourceCcGlobalConnectionBandwidthTags_base(name))
}
