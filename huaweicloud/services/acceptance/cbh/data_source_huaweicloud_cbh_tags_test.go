package cbh

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_cbh_instance_tags.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please make sure that the cbh instance is created with tags(foo = "xxx") before running this test.
			acceptance.TestAccPreCheckCbhInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataTags_base(),
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

// testAccDataTags_base returns the base configuration for the data source test.
func testAccDataTags_base() string {
	return `
data "huaweicloud_cbh_instance_tags" "test" {}

output "tags_validation" {
  value = length([for t in data.huaweicloud_cbh_instance_tags.test.tags: t.key == "foo" &&
    alltrue([for k, v in data.huaweicloud_cbh_instance_tags.test.tags: contains(t.values, v) if k == "foo"])]) > 0
}
`
}
