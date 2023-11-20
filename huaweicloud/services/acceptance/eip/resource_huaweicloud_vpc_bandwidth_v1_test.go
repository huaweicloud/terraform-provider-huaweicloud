package eip

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcBandWidthV1_update(t *testing.T) {
	var bandwidth bandwidths.BandWidth

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_bandwidth.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&bandwidth,
		getBandwidthResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandWidthV1_update(randName, 5),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "size", "5"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "publicips.#", "0"),
				),
			},
			{
				Config: testAccVpcBandWidthV1_update(randName+"_update", 6),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "size", "6"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "WHOLE"),
					resource.TestCheckResourceAttr(resourceName, "status", "NORMAL"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccVpcBandWidthV1_update(rName string, size int) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%s"
  size = "%d"
}
`, rName, size)
}
