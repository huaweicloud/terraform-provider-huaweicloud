package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccOfficialEventBatchAction_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_eg_official_event_batch_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgOfficialChannelName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testOfficialEventBatchAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "id"),
				),
			},
		},
	})
}

func testOfficialEventBatchAction_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_official_event_batch_action" "test" {
  source_name = "HC.CES"

  events {
    id                = "%[1]s"
    source            = "HC.CES"
    spec_version      = "1.0"
    type              = "CES:Event:SYS.RMS-configurationNoncomplianceNotification"
    data_content_type = "application/json"
	subject           = "CES:cn-north-4:0970dd7a1300f5672ff2c003c60ae115:Event:ev1756103124686vZvwRJE1Y",
	time              = "2023-01-01T12:00:00Z"

    data = jsonencode({
      object_id = "obj-123"
      timestamp = "2023-01-01T12:00:00Z"

	  "eventId": "ev1756103124686vZvwRJE1Y",
	  "eventInfo": {
	    "event_name": "configurationNoncomplianceNotification",
	    "event_source": "SYS.RMS",
	    "detail": {
          "sub_event_type": "SUB_EVENT.OPS",
          "event_state": "warning",
          "event_type": "EVENT.SYS",
          "event_level": "Major",
          "resource_id": "cdf1169a-5a7a-4c23-a253-6ab0a461d887",
          "resource_name": "pms-cluster-02",
          "content": ""
		},
		"time": 1756103124298
		},
    })
  }
}`, name)
}
