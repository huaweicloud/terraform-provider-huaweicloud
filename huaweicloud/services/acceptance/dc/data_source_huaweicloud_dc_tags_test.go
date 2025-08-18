package dc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceDcTags_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_dc_tags.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDcTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceDcTags_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dc_virtual_gateway" "test" {
  vpc_id      = huaweicloud_vpc.test.id
  name        = "%[2]s"
  description = "Created by acc test"

  local_ep_group = [
    huaweicloud_vpc.test.cidr,
  ]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestVpc(name), name)
}

func testDataSourceDcTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dc_tags" "test" {
  depends_on = [huaweicloud_dc_virtual_gateway.test]

  resource_type = "dc-vgw"
}
`, testDataSourceDcTags_base(name))
}
