package huaweicloud

import (
	"fmt"
	"math/rand"
	"testing"

	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcIdsV1DataSource_basic(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	rInt := rand.Intn(50)
	cidr := fmt.Sprintf("172.16.%d.0/24", rInt)
	name := fmt.Sprintf("terraform-testacc-vpc-data-source-%d", rInt)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcIdsV1Config(name, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccVpcIdsV2DataSourceID("data.huaweicloud_vpc_ids_v1.vpc_ids"),
				),
			},
		},
	})
}

func testAccVpcIdsV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceVpcIdsV1Config(name, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "vpc_1" {
	name = "%s"
	cidr = "%s"
}

data "huaweicloud_vpc_ids_v1" "vpc_ids" {
}
`, name, cidr)
}
