package servicestage

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccV3ComponentRefresh_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Make sure at least one of node exist.
			acceptance.TestAccPreCheckCceClusterId(t)
			// At least one of JAR package must be provided.
			acceptance.TestAccPreCheckServiceStageJarPkgStorageURLs(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccV3ComponentRefresh_basic_step1(),
			},
		},
	})
}

func testAccV3ComponentRefresh_basic_step1() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_servicestagev3_component_refresh" "test" {
  application_id = huaweicloud_servicestagev3_application.test.id
  component_id   = huaweicloud_servicestagev3_component.test.id
}
`, testAccV3ComponentAction_basic_base())
}
