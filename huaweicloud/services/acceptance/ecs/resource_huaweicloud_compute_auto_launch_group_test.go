package ecs

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getAutoLaunchGroupResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CMS client: %s", err)
	}

	getAutoLaunchGroupHttpUrl := "v2/{domain_id}/auto-launch-groups/{auto_launch_group_id}"
	getAutoLaunchGroupPath := client.Endpoint + getAutoLaunchGroupHttpUrl
	getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{domain_id}", cfg.DomainID)
	getAutoLaunchGroupPath = strings.ReplaceAll(getAutoLaunchGroupPath, "{auto_launch_group_id}", state.Primary.ID)
	getAutoLaunchGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getAutoLaunchGroupResp, err := client.Request("GET", getAutoLaunchGroupPath, &getAutoLaunchGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving auto launch groups: %s", err)
	}

	return utils.FlattenResponse(getAutoLaunchGroupResp)
}

func TestAccResourceAutoLaunchGroup_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_compute_auto_launch_group.test"
	rName := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getAutoLaunchGroupResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSLaunchTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAutoLaunchGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "target_capacity", "1"),
					resource.TestCheckResourceAttr(resourceName, "launch_template_id", acceptance.HW_ECS_LAUNCH_TEMPLATE_ID),
					resource.TestCheckResourceAttr(resourceName, "launch_template_version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "stable_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "excess_fulfilled_capacity_behavior"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "instances_behavior_with_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_since"),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "spot_price"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "current_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "current_stable_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "task_state"),
				),
			},
			{
				Config: testAccAutoLaunchGroup_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "target_capacity", "2"),
					resource.TestCheckResourceAttr(resourceName, "launch_template_id", acceptance.HW_ECS_LAUNCH_TEMPLATE_ID),
					resource.TestCheckResourceAttr(resourceName, "launch_template_version", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "stable_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "excess_fulfilled_capacity_behavior"),
					resource.TestCheckResourceAttrSet(resourceName, "type"),
					resource.TestCheckResourceAttrSet(resourceName, "instances_behavior_with_expiration"),
					resource.TestCheckResourceAttrSet(resourceName, "valid_since"),
					resource.TestCheckResourceAttrSet(resourceName, "allocation_strategy"),
					resource.TestCheckResourceAttrSet(resourceName, "spot_price"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "current_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "current_stable_capacity"),
					resource.TestCheckResourceAttrSet(resourceName, "task_state"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"task_state",
				},
			},
		},
	})
}

const testAccAutoLaunchGroupConfigbasic = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}
`

func testAccAutoLaunchGroup_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_auto_launch_group" "test" {
  name                    = "%s"
  target_capacity         = 1
  launch_template_id      = "%s"
  launch_template_version = "1"

  overrides {
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
    flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  }
}
`, testAccAutoLaunchGroupConfigbasic, name, acceptance.HW_ECS_LAUNCH_TEMPLATE_ID)
}

func testAccAutoLaunchGroup_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_auto_launch_group" "test" {
  name                    = "%s-update"
  target_capacity         = 2
  launch_template_id      = "%s"
  launch_template_version = "1"

  overrides {
    availability_zone = data.huaweicloud_availability_zones.test.names[0]
    flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  }
  overrides {
    availability_zone = data.huaweicloud_availability_zones.test.names[1]
    flavor_id         = data.huaweicloud_compute_flavors.test.ids[1]
  }
}
`, testAccAutoLaunchGroupConfigbasic, name, acceptance.HW_ECS_LAUNCH_TEMPLATE_ID)
}
