package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/as"
)

func getASBandWidthPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("autoscaling", region)
	if err != nil {
		return nil, fmt.Errorf("error creating AS bandWidth policy client: %s", err)
	}

	return as.GetBandwidthPolicy(client, state.Primary.ID)
}

func TestAccASBandWidthPolicy_basic(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_bandwidth_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASBandWidthPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASBandWidthPolicy_scheduled(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(resourceName, "scaling_resource_type", "BANDWIDTH"),
					resource.TestCheckResourceAttr(resourceName, "action", "pause"),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSED"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "300"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.size", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
				),
			},
			{
				Config: testASBandWidthPolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", rName+"-updated"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "SCHEDULED"),
					resource.TestCheckResourceAttr(resourceName, "scaling_resource_type", "BANDWIDTH"),
					resource.TestCheckResourceAttr(resourceName, "action", "resume"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "900"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.limits", "300"),
				),
			},
			{
				Config: testASBandWidthPolicy_recurrence(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "RECURRENCE"),
					resource.TestCheckResourceAttr(resourceName, "action", "pause"),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSED"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.launch_time", "07:00"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.recurrence_type", "Weekly"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"action"},
			},
		},
	})
}

func TestAccASBandWidthPolicy_alarm(t *testing.T) {
	var obj interface{}

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_bandwidth_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASBandWidthPolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testASBandWidthPolicy_alarm(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "ALARM"),
					resource.TestCheckResourceAttr(resourceName, "scaling_resource_type", "BANDWIDTH"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttrPair(resourceName, "bandwidth_id", "huaweicloud_vpc_bandwidth.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "alarm_id", "huaweicloud_ces_alarmrule.alarmrule_1", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testASBandWidthPolicy_scheduled(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "%[1]s"
  scaling_policy_type = "SCHEDULED"
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  action              = "pause"

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }

  scheduled_policy {
    launch_time = "2088-09-30T12:00Z"
  }
}
`, name)
}

func testASBandWidthPolicy_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "%[1]s-updated"
  scaling_policy_type = "SCHEDULED"
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  cool_down_time      = 900
  action              = "resume"

  scaling_policy_action {
    operation = "ADD"
    size      = 2
    limits    = 300
  }

  scheduled_policy {
    launch_time = "2099-09-30T12:00Z"
  }
}
`, name)
}

func testASBandWidthPolicy_recurrence(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "%[1]s"
  scaling_policy_type = "RECURRENCE"
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  cool_down_time      = 600
  action              = "pause"

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }

  scheduled_policy {
    launch_time      = "07:00"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    start_time       = "2022-09-30T12:00Z"
    end_time         = "2122-12-30T12:00Z"
  }
}
`, name)
}

func testASBandWidthPolicy_alarm(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[1]s"
  size = 5
}

resource "huaweicloud_ces_alarmrule" "alarmrule_1" {
  alarm_name           = "rule-%[1]s"
  alarm_description    = "autoScaling"
  alarm_action_enabled = true
  alarm_enabled        = true

  metric {
    namespace   = "SYS.VPC"
    metric_name = "downstream_bandwidth"

    dimensions {
      name  = "bandwidth_id"
      value = huaweicloud_vpc_bandwidth.test.id
    }
  }

  condition  {
    period              = 300
    filter              = "max"
    comparison_operator = ">"
    value               = 3600
    unit                = "bit/s"
    count               = 2
    suppress_duration   = 300
  }

  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = "%[1]s"
  scaling_policy_type = "ALARM"
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  alarm_id            = huaweicloud_ces_alarmrule.alarmrule_1.id

  scaling_policy_action {
    operation = "ADD"
    size      = 2
    limits    = 300
  }
}
`, name)
}
