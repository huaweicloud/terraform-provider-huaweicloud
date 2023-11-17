package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/scheduledtasks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPlannedTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AS v1 client: %s", err)
	}

	listOpts := scheduledtasks.ListOpts{
		GroupID: state.Primary.Attributes["scaling_group_id"],
	}
	tasks, err := scheduledtasks.List(client, listOpts)
	if err == nil && len(tasks) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return tasks[0], err
}

func TestAccPlannedTask_basic(t *testing.T) {
	var (
		task scheduledtasks.ScheduledTask

		name           = acceptance.RandomAccResourceName()
		taskName       = acceptance.RandomAccResourceNameWithDash()
		taskUpdateName = acceptance.RandomAccResourceNameWithDash()
		rName          = "huaweicloud_as_planned_task.test"

		rc = acceptance.InitResourceCheck(
			rName,
			&task,
			getPlannedTaskResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlannedTask_basic(name, taskName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", taskName),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.launch_time", "10:00"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.recurrence_type", "WEEKLY"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.start_time", "2025-11-30T12:00Z"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.end_time", "2025-12-30T12:00Z"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.recurrence_value", "6"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.max", "3"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.min", "1"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.desire", "2"),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id",
						"huaweicloud_as_group.acc_as_group", "id"),
				),
			},
			{
				Config: testAccPlannedTask_update(name, taskUpdateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", taskUpdateName),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.launch_time", "12:00"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.recurrence_type", "MONTHLY"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.start_time", "2025-12-30T12:00Z"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.end_time", "2026-01-30T12:00Z"),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.recurrence_value", "25"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.max", "2"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.min", "0"),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id",
						"huaweicloud_as_group.acc_as_group", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testPlannedTaskImportState(rName),
			},
		},
	})
}

func TestAccPlannedTask_timedTask(t *testing.T) {
	var (
		task scheduledtasks.ScheduledTask

		name     = acceptance.RandomAccResourceName()
		taskName = acceptance.RandomAccResourceNameWithDash()
		rName    = "huaweicloud_as_planned_task.test"

		rc = acceptance.InitResourceCheck(
			rName,
			&task,
			getPlannedTaskResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccPlannedTask_timedTask(name, taskName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", taskName),
					resource.TestCheckResourceAttr(rName, "scheduled_policy.0.launch_time", "2025-11-30T12:00Z"),
					resource.TestCheckResourceAttr(rName, "instance_number.0.max", "1"),
					resource.TestCheckResourceAttrPair(rName, "scaling_group_id",
						"huaweicloud_as_group.acc_as_group", "id"),
				),
			},
		},
	})
}

func testAccPlannedTask_basic(name, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  name             = "%[2]s"

  scheduled_policy {
    launch_time      = "10:00"
    recurrence_type  = "WEEKLY"
    start_time       = "2025-11-30T12:00Z"
    end_time         = "2025-12-30T12:00Z"
    recurrence_value = "6"
  }

  instance_number {
    max    = "3"
    min    = "1"
    desire = "2"
  }
}
`, testASGroup_basic(name), rName)
}

func testAccPlannedTask_update(name, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  name             = "%[2]s"

  scheduled_policy {
    launch_time      = "12:00"
    recurrence_type  = "MONTHLY"
    start_time       = "2025-12-30T12:00Z"
    end_time         = "2026-01-30T12:00Z"
    recurrence_value = "25"
  }

  instance_number {
    max = "2"
    min = "0"
  }
}
`, testASGroup_basic(name), rName)
}

func testAccPlannedTask_timedTask(name, rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  name             = "%[2]s"

  scheduled_policy {
    launch_time = "2025-11-30T12:00Z"
  }

  instance_number {
    max = "1"
  }
}
`, testASGroup_basic(name), rName)
}

func testPlannedTaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		id := rs.Primary.ID
		groupID := rs.Primary.Attributes["scaling_group_id"]
		if groupID == "" || id == "" {
			return "", fmt.Errorf("the planned task is not exist or related scaling group ID is missing")
		}

		return fmt.Sprintf("%s/%s", groupID, id), nil
	}
}
