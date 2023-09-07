package eg

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/eg/v1/source/custom"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getCustomEventSourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.EgV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating EG v1 client: %s", err)
	}

	return custom.Get(client, state.Primary.ID)
}

func TestAccCustomEventSource_basic(t *testing.T) {
	var (
		obj custom.Source

		rName = "huaweicloud_eg_custom_event_source.test"
		name  = acceptance.RandomAccResourceName()
		rc    = acceptance.InitResourceCheck(rName, &obj, getCustomEventSourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEgChannelId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCustomEventSource_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "channel_id", acceptance.HW_EG_CHANNEL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(rName, "description", "Created by acceptance test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testAccCustomEventSource_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "channel_id", acceptance.HW_EG_CHANNEL_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "type", "APPLICATION"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCustomEventSource_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id  = "%[1]s"
  name        = "%[2]s"
  description = "Created by acceptance test"
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}

func testAccCustomEventSource_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_eg_custom_event_source" "test" {
  channel_id = "%[1]s"
  name       = "%[2]s"
}
`, acceptance.HW_EG_CHANNEL_ID, name)
}
