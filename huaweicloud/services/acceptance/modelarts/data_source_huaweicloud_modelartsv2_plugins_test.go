package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV2Plugins_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_modelartsv2_plugins.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV2Plugins_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "plugins.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.kind"),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.api_version"),

					resource.TestMatchResourceAttr(dcName, "plugins.0.metadata.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.metadata.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.metadata.0.created_at"),

					resource.TestMatchResourceAttr(dcName, "plugins.0.spec.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestMatchResourceAttr(dcName, "plugins.0.spec.0.template.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.spec.0.template.0.name"),

					resource.TestMatchResourceAttr(dcName, "plugins.0.status.#",
						regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.status.0.phase"),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.status.0.version"),
					resource.TestCheckResourceAttrSet(dcName, "plugins.0.status.0.reason"),
				),
			},
		},
	})
}

func testAccDataV2Plugins_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_plugins" "test" {
  pool_id = "%[1]s"
}
`, acceptance.HW_MODELARTS_RESOURCE_POOL_ID)
}
