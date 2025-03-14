package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceNatPrivateTransitIpTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_nat_private_transit_ip_tags.test"
		name           = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePrivateTransitIpTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourcePrivateTransitIpTags_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_nat_private_transit_ip" "test" {
  subnet_id             = huaweicloud_vpc_subnet.test.id
  ip_address            = "192.168.0.68"
  enterprise_project_id = "0"

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_nat_private_transit_ip_tags" "test" {
  depends_on = [huaweicloud_nat_private_transit_ip.test]
}
`, common.TestBaseNetwork(name))
}
