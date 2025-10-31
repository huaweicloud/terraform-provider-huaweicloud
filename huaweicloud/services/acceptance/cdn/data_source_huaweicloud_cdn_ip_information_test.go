package cdn

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataIpInformation_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_ip_information.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDNIPAddresses(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataIpInformation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "information.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(rName, "information.0.ip"),
					resource.TestCheckResourceAttrSet(rName, "information.0.belongs"),
				),
			},
		},
	})
}

func testAccDataIpInformation_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_ip_information" "test" {
  ips = "%[1]s"
}
`, acceptance.HW_CDN_IP_ADDRESSES)
}
