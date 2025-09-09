package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageTags_basic(rName),
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

func testDataSourceCcBandwidthPackageTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cc_bandwidth_package" "test" {
  name           = "%s"
  local_area_id  = "Chinese-Mainland"
  remote_area_id = "Chinese-Mainland"
  charge_mode    = "bandwidth"
  billing_mode   = 3
  bandwidth      = 5
  description    = "This is an accaptance test"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func testDataSourceCcBandwidthPackageTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cc_bandwidth_package_tags" "test" {
  depends_on = [huaweicloud_cc_bandwidth_package.test]
}
`, testDataSourceCcBandwidthPackageTags_base(name))
}
