package huaweicloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAvailabilityZonesV2_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDeprecated(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.huaweicloud_compute_availability_zones_v2.zones", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccAvailabilityZonesConfig = `
data "huaweicloud_compute_availability_zones_v2" "zones" {}
`
