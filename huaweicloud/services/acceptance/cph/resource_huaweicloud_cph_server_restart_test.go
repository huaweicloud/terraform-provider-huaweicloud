package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCphServerRestart_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCphServerRestart_basic(name),
			},
		},
	})
}

func testCphServerRestart_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cph_server_restart" "test" {
  server_id = huaweicloud_cph_server.test.id
}
`, testCphServer_basic(name))
}
