package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPostPaidOrder_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterPostPaidOrder(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPostPaidOrder_basic(),
			},
		},
	})
}

func testPostPaidOrder_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_post_paid_order" "name" {
  product_list {
    id                 = "%[1]s"
    product_id         = "%[2]s"
    cloud_service_type = "hws.service.type.sa"
    resource_type      = "hws.resource.type.secmaster.typical"
    resource_spec_code = "secmaster.professional"
    usage_measure_id   = 4
    usage_value        = 1
    resource_size      = 4
    usage_factor       = "duration"    
  }
}
`, acceptance.HW_SECMASTER_ORDER_ID, acceptance.HW_SECMASTER_PRODUCT_ID)
}
