package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccInstanceNodeDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceNodeDelete_basic(),
			},
		},
	})
}

func TestAccInstanceNodeDelete_prepaid(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceNodeDelete_prpaid(),
			},
		},
	})
}

func testAccInstanceNodeDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_instance_node_delete" "test" {
  instance_id = "%s"
  num         = "2"
}
`, acceptance.HW_DDS_INSTANCE_ID)
}

func testAccInstanceNodeDelete_prpaid() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_instance_node_delete" "test" {
  instance_id = "%s"
  num         = "2"
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
