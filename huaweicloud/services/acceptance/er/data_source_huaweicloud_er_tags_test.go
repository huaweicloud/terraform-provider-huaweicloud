package er

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_tags.test"
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

func testAccDataTags_base() string {
	var (
		name = acceptance.RandomAccResourceName()
		// The valid value is range from 64512 to 65534, maintain an upper limit as the boundary value of asn+count.
		bgpAsNum = acctest.RandIntRange(64512, 65533)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  count = 2

  availability_zones = try(slice(data.huaweicloud_er_availability_zones.test.names, 0, 1), [])
  name               = format("%[1]s_%%d", count.index)
  asn                = %[2]d+count.index

  tags = {
    foo = format("bar%%d", count.index)
  }

  lifecycle {
    ignore_changes = [
      availability_zones, # Available instances in an availability zone may be sold out, so its value may be empty.
    ]
  }
}
`, name, bgpAsNum)
}

func testAccDataTags_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_tags" "test" {
  depends_on = [huaweicloud_er_instance.test]

  resource_type = "instance"
}

output "tags_validation" {
  value = length([for t in data.huaweicloud_er_tags.test.tags: t.key == "foo" &&
    alltrue([for k, v in huaweicloud_er_instance.test[*].tags: contains(t.values, v) if k == "foo"])]) > 0
}
`, testAccDataTags_base())
}
