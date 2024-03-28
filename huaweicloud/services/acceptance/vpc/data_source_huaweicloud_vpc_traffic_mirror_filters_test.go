package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcTrafficMirrorFilters_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpc_traffic_mirror_filters.test1"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcTrafficMirrorFilters_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "name", rName),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcTrafficMirrorFilters_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_traffic_mirror_filter" "test1" {
  name        = "%[1]s"
  description = "tf acc test filter"
}

data "huaweicloud_vpc_traffic_mirror_filters" "test1" {
  name = "%[1]s"
}
`, name)
}
