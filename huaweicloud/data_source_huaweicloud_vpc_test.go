package huaweicloud

import (
	"strconv"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcV1DataSource_basic(t *testing.T) {
	rName := fmtp.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	tmp := strconv.Itoa(acctest.RandIntRange(1, 254))
	cidr := fmtp.Sprintf("172.16.%s.0/24", tmp)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcV1Config(rName, cidr),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcV1Check("data.huaweicloud_vpc.by_id", rName),
					testAccDataSourceVpcV1Check("data.huaweicloud_vpc.by_cidr", rName),
					testAccDataSourceVpcV1Check("data.huaweicloud_vpc.by_name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc.by_id", "status", "OK"),
				),
			},
		},
	})
}

func testAccDataSourceVpcV1Check(n, rName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("root module has no resource called %s", n)
		}

		vpcRs, ok := s.RootModule().Resources["huaweicloud_vpc.test"]
		if !ok {
			return fmtp.Errorf("can't find huaweicloud_vpc.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != vpcRs.Primary.Attributes["id"] {
			return fmtp.Errorf(
				"id is %s; want %s",
				attr["id"],
				vpcRs.Primary.Attributes["id"],
			)
		}

		if attr["name"] != rName {
			return fmtp.Errorf("bad vpc name %s", attr["name"])
		}

		return nil
	}
}

func testAccDataSourceVpcV1Config(rName, cidr string) string {
	return fmtp.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

data "huaweicloud_vpc" "by_id" {
  id = huaweicloud_vpc.test.id
}

data "huaweicloud_vpc" "by_cidr" {
  cidr = huaweicloud_vpc.test.cidr
}

data "huaweicloud_vpc" "by_name" {
  name = huaweicloud_vpc.test.name
}
`, rName, cidr)
}
