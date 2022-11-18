package as

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getASBandWidthPolicyResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getBandwidthPolicy: Query the AS bandwidth scaling policy
	var (
		getBandwidthPolicyHttpUrl = "autoscaling-api/v2/{project_id}/scaling_policy/{id}"
		getBandwidthPolicyProduct = "autoscaling"
	)
	getBandwidthPolicyClient, err := config.NewServiceClient(getBandwidthPolicyProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ASBandWidthPolicy Client: %s", err)
	}

	getBandwidthPolicyPath := getBandwidthPolicyClient.Endpoint + getBandwidthPolicyHttpUrl
	getBandwidthPolicyPath = strings.ReplaceAll(getBandwidthPolicyPath, "{project_id}", getBandwidthPolicyClient.ProjectID)
	getBandwidthPolicyPath = strings.ReplaceAll(getBandwidthPolicyPath, "{id}", state.Primary.ID)

	getBandwidthPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBandwidthPolicyResp, err := getBandwidthPolicyClient.Request("GET", getBandwidthPolicyPath, &getBandwidthPolicyOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ASBandWidthPolicy: %s", err)
	}
	return utils.FlattenResponse(getBandwidthPolicyResp)
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
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
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
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "900"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.size", "2"),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_action.0.limits", "300"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testASBandWidthPolicy_recurrence(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_name", rName),
					resource.TestCheckResourceAttr(resourceName, "scaling_policy_type", "RECURRENCE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.launch_time", "07:00"),
					resource.TestCheckResourceAttr(resourceName, "scheduled_policy.0.recurrence_type", "Weekly"),
				),
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

  scaling_policy_action {
    operation = "ADD"
    size      = 1
  }
  scheduled_policy {
    launch_time      = "07:00"
    recurrence_type  = "Weekly"
    recurrence_value = "1,3,5"
    start_time       = "2022-09-30T12:00Z"
    end_time         = "2022-12-30T12:00Z"
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
