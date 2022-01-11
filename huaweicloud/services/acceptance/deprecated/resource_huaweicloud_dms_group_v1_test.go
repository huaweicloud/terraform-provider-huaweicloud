package deprecated

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/dms/v1/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDmsGroupFunc(c *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := c.DmsV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud DMS client(V1): %s", err)
	}
	queueID := state.Primary.Attributes["queue_id"]
	page, err := groups.List(client, queueID, false).AllPages()
	if err == nil {
		groupsList, err := groups.ExtractGroups(page)
		if err != nil {
			return nil, fmtp.Errorf("Error getting groups in queue %s: %s", queueID, err)
		}
		if len(groupsList) > 0 {
			for _, group := range groupsList {
				if group.ID == state.Primary.ID {
					return group, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("not found dms group: %s", err)
}

func TestAccDmsGroupsV1_basic(t *testing.T) {
	var group groups.Group
	var groupName = acceptance.RandomAccResourceName()
	var queueName = acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_dms_group_v1.group_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getDmsGroupFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckDeprecated(t)
			acceptance.TestAccPreCheckDms(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDmsV1Group_basic(groupName, queueName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", groupName),
				),
			},
		},
	})
}

func testAccDmsV1Group_basic(groupName string, queueName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_queue_v1" "queue_1" {
  name = "%s"
}
resource "huaweicloud_dms_group_v1" "group_1" {
  name     = "%s"
  queue_id = "${huaweicloud_dms_queue_v1.queue_1.id}"
}`, queueName, groupName)
}
