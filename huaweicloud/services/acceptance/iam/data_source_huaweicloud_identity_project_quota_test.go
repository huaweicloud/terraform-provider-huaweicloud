package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataProjectQuota_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_project_quota.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataProjectQuota_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "resources.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "resources.0.type", "project"),
				),
			},
		},
	})
}

func testAccDataProjectQuota_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_project" "test" {
  name        = "%[1]s_%[2]s"
}

data "huaweicloud_identity_project_quota" "test" {
  project_id = huaweicloud_identity_project.test.id
  
  # Waiting for the project to be created.
  depends_on = [huaweicloud_identity_project.test]
}
`, acceptance.HW_REGION_NAME, name)
}
