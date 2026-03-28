package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIntelligentSessionKill_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccIntelligentSessionKill_basic(rName),
			},
		},
	})
}

func testAccIntelligentSessionKill_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_rds_intelligent_session_kill" "test" {
  instance_id             = huaweicloud_rds_instance.test.id
  auto_add_sql_limit_rule = "true"
}
`, testAccRdsInstance_mysql(rName))
}
