package cbr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cbr_tags.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "tags.0.key"),
					resource.TestMatchResourceAttr(all, "tags.0.values.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("tags_validation", "true"),
				),
			},
		},
	})
}

func testAccDataTags_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cbr_vault" "test" {
  count = 2

  name                  = format("%[1]s_%%d", count.index)
  type                  = "server"
  consistent_level      = "crash_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "0"

  tags = {
    foo = format("bar%%d", count.index)
  }
}
`, name)
}

func testAccDataTags_basic() string {
	name := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_tags" "test" {
  depends_on = [huaweicloud_cbr_vault.test]

  resource_type = "vault"
}

output "tags_validation" {
  value = length([for t in data.huaweicloud_cbr_tags.test.tags: t.key == "foo" &&
    alltrue([for k, v in huaweicloud_cbr_vault.test[*].tags: contains(t.values, v) if k == "foo"])]) > 0
}
`, testAccDataTags_base(name))
}
