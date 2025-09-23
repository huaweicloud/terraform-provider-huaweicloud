package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEventReport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCesTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testEventReport_basic(),
			},
		},
	})
}

func testEventReport_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_event_report" "test" {
  name   = "test"
  source = "test.System"
  time   = "%[1]s"

  detail {
    state         = "normal"
    level         = "Major"
    content       = "test content"
    resource_id   = "ecs001"
    resource_name = "test resource name"
    user          = "test user"

    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }

    dimensions {
      name  = "instance_name"
      value = "test_instance_name"
    }
  }
}
`, acceptance.HW_CES_START_TIME)
}
