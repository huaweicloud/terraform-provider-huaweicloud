package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGroupAlarmTemplateAsyncAssociate_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceGroupAlarmTemplateAsyncAssociate_basic(rName),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testResourceGroupAlarmTemplateAsyncAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

%[3]s

resource "huaweicloud_ces_resource_group_alarm_template_async_associate" "test" {
  group_id             = huaweicloud_ces_resource_group.test.id
  template_ids         = [huaweicloud_ces_alarm_template.test.id]
  notification_enabled = true
  alarm_notifications {
    type              = "contact"
    notification_list = [huaweicloud_smn_topic.test.id]
  }
  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
  effective_timezone      = "GMT+08:00"
  enterprise_project_id   = "0"
}
`, testResourceGroup_basic(name), testCesAlarmTemplate_basic(name), testCESAlarmRule_topicBase(name))
}
