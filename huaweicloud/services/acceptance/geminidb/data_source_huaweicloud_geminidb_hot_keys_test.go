package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHotKeys_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_hot_keys.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Because there is no available data for testing, the test case is only
			// used to verify that the API can be invoked.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHotKeys_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keys.#"),
				),
			},
		},
	})
}

func testAccDataSourceHotKeys_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_hot_keys" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
