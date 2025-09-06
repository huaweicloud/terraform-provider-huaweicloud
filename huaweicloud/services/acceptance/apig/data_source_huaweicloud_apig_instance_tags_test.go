package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, please make sure at least two tags that the instance have.
func TestAccDataSourceInstanceTags_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_instance_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceTags_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "tags.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.0.value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.1.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.1.value"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceTags_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_apig_instance_tags" "test" {
  instance_id = "%s"
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
