package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsNodeStatusChange_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsNodeStatusChange_basic(name),
			},
		},
	})
}

func testAccDcsNodeStatusChange_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dcs_backup" "test" {
  instance_id   = huaweicloud_dcs_instance.instance_1.id
  description   = "test DCS backup remark"
  backup_format = "rdb"
}

resource "huaweicloud_dcs_node_status_change" "test" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  action      = "start"
	
  depends_on = [huaweicloud_dcs_backup.test]
}
`, testAccDcsInstance_base(name))
}
