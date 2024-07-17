package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEnvironments_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_apig_environments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironments_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "environments.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccDataSourceEnvironments_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_environment" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  description = "Created by script"
}

data "huaweicloud_apig_environments" "test" {
  depends_on = [huaweicloud_apig_environment.test]

  instance_id = "%[1]s"
  name        = huaweicloud_apig_environment.test.name
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}
