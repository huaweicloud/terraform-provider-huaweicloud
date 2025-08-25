package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventBatchAction_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "huaweicloud_eg_event_batch_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testEventBatchAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "id"),
				),
			},
		},
	})
}

func testEventBatchAction_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_channel" "test" {
  name = "%[1]s"
}

resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = huaweicloud_eg_custom_event_channel.test.id
  type        = "APPLICATION"
  name        = "%[1]s"
  description = "Created by script"
}
`, name)
}

func testEventBatchAction_basic(name string) string {
	return fmt.Sprintf(`
%[2]s

resource "huaweicloud_eg_event_batch_action" "test" {
  channel_id = huaweicloud_eg_custom_event_channel.test.id

  events {
    id                = "%[1]s"
    source            = huaweicloud_eg_custom_event_source.test.id
    spec_version      = "1.0"
    type              = "com.example.object.created.v1"
    data_content_type = "application/json"
	time              = "2023-01-01T12:00:00Z"
    subject           = "test-object"

    data = jsonencode({
      object_id = "obj-123"
      timestamp = "2023-01-01T12:00:00Z"
    })
  }
}`, name, testEventBatchAction_base(name))
}
