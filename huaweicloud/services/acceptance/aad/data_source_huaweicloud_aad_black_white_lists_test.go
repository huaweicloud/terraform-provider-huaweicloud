package antiddos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBlackWhiteLists_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aad_black_white_lists.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a valid AAD instance ID and config it to the environment variable.
			acceptance.TestAccPreCheckAadInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAadBlackWhiteLists_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ips.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ips.0.ip"),
				),
			},
		},
	})
}

func testDataSourceAadBlackWhiteLists_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_aad_black_white_list" "test" {
  instance_id = "%[1]s"
  type        = "black"
  ips         = ["12.1.2.117"]
}

data "huaweicloud_aad_black_white_lists" "test" {
  depends_on = [huaweicloud_aad_black_white_list.test]

  instance_id = "%[1]s"
  type        = "black"
}
`, acceptance.HW_AAD_INSTANCE_ID)
}
