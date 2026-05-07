package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataWorkspaceQuotas_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_modelarts_workspace_quotas.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		rName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataWorkspaceQuotas_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "quotas.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.resource"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.quota"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.min_quota"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.max_quota"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.used_quota"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.name_cn"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.name_en"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.unit_cn"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.unit_en"),
					resource.TestCheckResourceAttrSet(dcName, "quotas.0.updated_at"),
				),
			},
		},
	})
}

func testAccDataWorkspaceQuotas_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_modelarts_workspace" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"
}

data "huaweicloud_modelarts_workspace_quotas" "test" {
  workspace_id = huaweicloud_modelarts_workspace.test.id
}
`, name)
}
