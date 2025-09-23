package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAlertRuleSimulation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterInstanceID(t)
			acceptance.TestAccPreCheckSecMasterPipelineID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAlertRuleSimulation_basic(),
			},
		},
	})
}

func testAlertRuleSimulation_basic() string {
	// This query rule only applies to this pipeline (sec-cts-audit)
	queryRule := "trace_name = 'addNic' and service_type = 'ECS' and trace_rating = 'normal' | select src_ip__,resource_name,__time," +
		"resource_id,request,ops.rgn, user.name, trim('cts') as data_source_product_feature"
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_alert_rule_simulation" "test" {
  workspace_id   = "%[1]s"
  pipeline_id    = "%[2]s"
  query_rule     = "%[3]s"
  query_type     = "SQL"
  from_time      = "2025-01-16 19:04:05"
  to_time        = "2025-01-16 19:06:05"
  event_grouping = false

  triggers {
    mode              = "COUNT"
    operator          = "GT"
    expression        = "0"
    severity          = "TIPS"
    accumulated_times = 1
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_PIPELINE_ID, queryRule)
}
