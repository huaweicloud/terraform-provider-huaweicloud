package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPublicAccessSwitch_basic(t *testing.T) {
	var (
		rName = "huaweicloud_dms_kafka_instance_public_access_switch.test"
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccPublicAccessSwitch_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "instance_id"),
					resource.TestCheckResourceAttrSet(rName, "publicip_id"),
				),
			},
		},
	})
}

func testAccPublicAccessSwitch_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_eip" "test" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    name        = "%[1]s"
    size        = 5
    share_type  = "PER"
    charge_mode = "traffic"
  }
}`, name)
}

func testAccPublicAccessSwitch_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_public_access_switch" "test" {
  instance_id = "%[2]s"

  public_boundwidth = 0
  publicip_id       = huaweicloud_vpc_eip.test.id
}`, testAccPublicAccessSwitch_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
