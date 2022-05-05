package eip

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBandWidthDataSource_basic(t *testing.T) {
	randName := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_vpc_bandwidth.test"
	eipResourceName := "huaweicloud_vpc_eip.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBandWidthDataSource_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", randName),
					resource.TestCheckResourceAttr(dataSourceName, "size", "10"),
					resource.TestCheckResourceAttr(dataSourceName, "publicips.#", "1"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.id",
						eipResourceName, "id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "publicips.0.ip_address",
						eipResourceName, "address"),
				),
			},
		},
	})
}

func TestAccBandWidthDataSource_basic2(t *testing.T) {
	err := b(0)
	if _, ok := err.(golangsdk.ErrDefault404); ok {
		logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCI network")

	} else {
		logp.Printf("[DEBUG] Successfully deleted HuaweiCloud CCI network2222")
	}
}
func b(i int) error {
	if i == 1 {
		return golangsdk.ErrDefault404{}
	}
	return nil
}

func testAccBandWidthDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%s"
  size = 10
}

resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }
  bandwidth {
    share_type = "WHOLE"
    id         = huaweicloud_vpc_bandwidth.test.id
  }
}

data "huaweicloud_vpc_bandwidth" "test" {
  depends_on = [huaweicloud_vpc_eip.test]

  name = huaweicloud_vpc_bandwidth.test.name
}
`, rName)
}
