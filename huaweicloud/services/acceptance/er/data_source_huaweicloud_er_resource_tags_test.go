package er

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataResourceTags_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_er_resource_tags.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataResourceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "tags.%", "2"),
					resource.TestCheckResourceAttr(all, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(all, "tags.owner", "terraform"),
					resource.TestCheckOutput("tags_validation", "true"),
				),
			},
		},
	})
}

func testAccDataResourceTags_base() string {
	var (
		name     = acceptance.RandomAccResourceName()
		bgpAsNum = acctest.RandIntRange(64512, 65534)
	)

	return fmt.Sprintf(`
data "huaweicloud_er_availability_zones" "test" {}

resource "huaweicloud_er_instance" "test" {
  availability_zones = try(slice(data.huaweicloud_er_availability_zones.test.names, 0, 1), [])
  name               = "%[1]s"
  asn                = %[2]d

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  lifecycle {
    ignore_changes = [
      availability_zones, # Available instances in an availability zone may be sold out, so its value may be empty.
    ]
  }
}
`, name, bgpAsNum)
}

func testAccDataResourceTags_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_er_resource_tags" "test" {
  resource_type = "instance"
  resource_id   = huaweicloud_er_instance.test.id
}

output "tags_validation" {
  value = length(keys(data.huaweicloud_er_resource_tags.test.tags)) == length(keys(huaweicloud_er_instance.test.tags)) && alltrue(
    [for k, v in data.huaweicloud_er_resource_tags.test.tags: lookup(huaweicloud_er_instance.test.tags, k, "NONE") == v])
}
`, testAccDataResourceTags_base())
}
