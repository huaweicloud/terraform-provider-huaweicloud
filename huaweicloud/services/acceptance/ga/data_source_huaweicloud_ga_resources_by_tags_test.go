package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaResourcesByTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ga_resources_by_tags.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaResourceByTags_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),

					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaResourceByTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ga_resources_by_tags" "test" {
  depends_on = [huaweicloud_ga_accelerator.test]

  resource_type = "ga-accelerators"
}

data "huaweicloud_ga_resources_by_tags" "filter_by_tags" {
  depends_on = [huaweicloud_ga_accelerator.test]

  resource_type = "ga-accelerators"

  tags {
    key    = "foo"
    values = ["bar1"]
  }
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_ga_resources_by_tags.filter_by_tags.resources) > 0
}
`, testAccGaResourceByTags_base(name))
}

func testAccGaResourceByTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ga_accelerator" "test" {
  name        = "%[1]s"
  description = "terraform test"

  ip_sets {
    area = "CM"
  }

  tags = {
    foo = "bar1"
    key = "value1"
  }
}
`, name)
}
