package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_ga_tags.test"
		name           = acceptance.RandomAccResourceNameWithDash()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTags_basic(name),
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

func testDataSourceTags_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%s"
  description = "terraform test"

  ip_sets {
    ip_type = "IPV4"
    area    = "CM"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}

data "huaweicloud_ga_tags" "test" {
  depends_on = [huaweicloud_ga_accelerator.test]

  resource_type = "ga-accelerators"
}
`, name)
}
