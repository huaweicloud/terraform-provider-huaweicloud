package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV2Namespaces_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cciv2_namespaces.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccV2Namespaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "namespaces.#", "1"),
				),
			},
		},
	})
}

func testAccV2Namespaces_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cci_namespaces" "test" {
  depends_on = [huaweicloud_cci_network.test]

  name = "%s"
}
`, testAccDataCciNamespaces_base(rName), rName)
}
