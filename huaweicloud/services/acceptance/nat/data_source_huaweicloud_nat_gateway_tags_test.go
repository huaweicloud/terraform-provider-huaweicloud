package nat

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceGatewayTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_nat_gateway_tags.test"
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
				Config: testDataSourceGatewayTags_basic(name),
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

func testDataSourceGatewayTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_nat_gateway" "test" {
  name      = "%[2]s"
  spec      = "1"
  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_nat_gateway_tags" "test" {
  depends_on = [huaweicloud_nat_gateway.test]
}
`, common.TestBaseNetwork(name), name)
}
