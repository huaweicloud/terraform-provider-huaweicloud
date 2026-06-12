package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBigKeys_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_big_keys.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBigKeys_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "keys.#"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.db_id"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_type"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_name"),
					resource.TestCheckResourceAttrSet(dataSource, "keys.0.key_size"),

					resource.TestCheckOutput("key_types_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceBigKeys_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_big_keys" "test" {
  instance_id = "%[1]s"
  key_types   = ["string","hash"]
}

locals {
  key_type = lower(data.huaweicloud_geminidb_big_keys.test.keys[0].key_type)
}

data "huaweicloud_geminidb_big_keys" "key_types_filter" {
  instance_id = "%[1]s"
  key_types   = [local.key_type]
}

output "key_types_filter_useful" {
  value = length(data.huaweicloud_geminidb_big_keys.key_types_filter.keys) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID)
}
