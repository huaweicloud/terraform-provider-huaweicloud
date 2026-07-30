package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceMetadataTags_basic(t *testing.T) {
	resourceName := "huaweicloud_dsc_metadata_tags.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscMetadataTagName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMetadataTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "names.0", acceptance.HW_DSC_METADATA_TAG_NAME),
				),
			},
			{
				Config: testAccResourceMetadataTags_update(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "names.0", acceptance.HW_DSC_METADATA_TAG_NAME_UPDATE),
				),
			},
		},
	})
}

func testAccResourceMetadataTags_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_metadata_tags" "test" {
  names = ["%s"]
}
`, acceptance.HW_DSC_METADATA_TAG_NAME)
}

func testAccResourceMetadataTags_update() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_metadata_tags" "test" {
  names = ["%s"]
}
`, acceptance.HW_DSC_METADATA_TAG_NAME_UPDATE)
}
