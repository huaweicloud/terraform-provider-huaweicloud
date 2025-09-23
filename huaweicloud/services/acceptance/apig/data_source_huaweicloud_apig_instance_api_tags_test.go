package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceApiTags_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_apig_instance_api_tags.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceApiTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "tags.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_tags_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceApiTags_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

resource "huaweicloud_apig_api" "test" {
  count            = 2
  instance_id      = "%[1]s"
  group_id         = huaweicloud_apig_group.test.id
  name             = "%[2]s${count.index}"
  type             = "Private"
  request_protocol = "HTTP"
  request_method   = "GET"
  request_path     = "/mock/test${count.index}"

  mock {
    status_code = 200
  }

  tags = count.index == 0 ? ["foo", "bar"] : ["TF"]
}

data "huaweicloud_apig_instance_api_tags" "test" {
  instance_id = "%[1]s"

  depends_on = [huaweicloud_apig_api.test]
}

locals {
  filter_result = [for v in flatten(huaweicloud_apig_api.test[*].tags) : contains(data.huaweicloud_apig_instance_api_tags.test.tags, v)]
}

output "is_tags_set_and_valid" {
  value = length(local.filter_result) > 0 && alltrue(local.filter_result)
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
