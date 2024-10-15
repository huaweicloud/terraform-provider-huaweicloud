package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCphAdbCommand_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCphAdbObjectPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCphAdbCommand_basic(name),
			},
		},
	})
}

func testCphAdbCommand_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cph_adb_command" "test" {
  command    = "push"
  content    = "%[2]s"
  server_ids = [huaweicloud_cph_server.test.id]
}
`, testCphServer_basic(name), acceptance.HW_CPH_OBS_OBJECT_PATH)
}
