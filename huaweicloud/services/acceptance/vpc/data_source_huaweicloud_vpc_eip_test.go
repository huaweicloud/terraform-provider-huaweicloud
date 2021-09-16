package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVpcEipDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_vpc_eip.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcEipConfig_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "status", "DOWN"),
					resource.TestCheckResourceAttr(dataSourceName, "type", "5_bgp"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_size", "8"),
					resource.TestCheckResourceAttr(dataSourceName, "bandwidth_share_type", "PER"),
				),
			},
		},
	})
}

func testAccDataSourceVpcEipConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    name        = "%s"
    size        = 8
    share_type  = "PER"
    charge_mode = "traffic"
  }
}

data "huaweicloud_vpc_eip" "test" {
  public_ip = huaweicloud_vpc_eip.test.address
}
`, rName)
}
