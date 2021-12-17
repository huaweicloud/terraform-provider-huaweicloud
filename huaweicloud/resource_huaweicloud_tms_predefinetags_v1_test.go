package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestTmsTagsV1(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTMSTagsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testTmsTagsV1Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_tms_tags.test", "tag.#", "1"),
				),
			},
			{
				Config: testTmsTagsV1ConfigUpdates,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_tms_tags.test", "tag.#", "1"),
				),
			},
		},
	})
}

var testTmsTagsV1Config = `
resource "huaweicloud_tms_tags" "test" {
	tag {
		key = "xxn"
		value = "11111"
	}
}
`

var testTmsTagsV1ConfigUpdates = `
resource "huaweicloud_tms_tags" "test" {
	tag {
		key = "nxx"
		value = "11111"
	}
}
`

func testAccCheckTMSTagsDestroy(s *terraform.State) error {
	return nil
}
