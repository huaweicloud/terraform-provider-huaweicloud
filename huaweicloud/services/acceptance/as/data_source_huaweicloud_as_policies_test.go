package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccPoliciesDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		baseConfig     = testACCPolicy_base()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSource_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.operation"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.action.0.instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.cool_down_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.created_at"),
					resource.TestCheckOutput("is_scaling_policy_id_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("is_scaling_policy_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testACCASGroup_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_as_configuration" "test" {
  scaling_configuration_name = "%[2]s"

  instance_config {
    image    = data.huaweicloud_images_image.test.id
    flavor   = data.huaweicloud_compute_flavors.test.ids[0]
    key_name = huaweicloud_kps_keypair.test.id

    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "huaweicloud_as_group" "acc_as_group" {
  scaling_group_name       = "%[2]s"
  scaling_configuration_id = huaweicloud_as_configuration.test.id
  vpc_id                   = huaweicloud_vpc.test.id

  networks {
    id = huaweicloud_vpc_subnet.test.id
  }

  security_groups {
    id = huaweicloud_networking_secgroup.test.id
  }
}
`, common.TestBaseComputeResources(rName), rName)
}

func testACCPolicy_base() string {
	name := acceptance.RandomAccResourceNameWithDash()

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ces_alarmrule" "alarm_rule" {
  alarm_name = "%[2]s"

  metric {
    namespace   = "SYS.AS"
    metric_name = "cpu_util"

    dimensions {
      name  = "AutoScalingGroup"
      value = huaweicloud_as_group.acc_as_group.id
    }
  }

  condition {
    period              = 300
    filter              = "average"
    comparison_operator = ">="
    value               = 60
    unit                = "%%"
    count               = 1
    suppress_duration   = 300
  }

  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "huaweicloud_as_policy" "policy_alarm" {
  scaling_policy_name = "%[2]s"
  scaling_policy_type = "ALARM"
  scaling_group_id    = huaweicloud_as_group.acc_as_group.id
  alarm_id            = huaweicloud_ces_alarmrule.alarm_rule.id
  cool_down_time      = 600

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
}
`, testACCASGroup_base(name), name)
}

func testAccPoliciesDataSource_basic(baseConfig string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_as_policies" "test" {
  depends_on = [huaweicloud_as_policy.policy_alarm]

  scaling_group_id = huaweicloud_as_group.acc_as_group.id
}

// Filter using scaling policy ID.
locals {
  scaling_policy_id = data.huaweicloud_as_policies.test.policies[0].id
}

data "huaweicloud_as_policies" "scaling_policy_id_filter" {
  scaling_group_id  = huaweicloud_as_group.acc_as_group.id
  scaling_policy_id = local.scaling_policy_id
}

output "is_scaling_policy_id_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_id_filter.policies[*].id : v == local.scaling_policy_id]
  )
}

// Filter using scaling policy name.
locals {
  scaling_policy_name = data.huaweicloud_as_policies.test.policies[0].name
}

data "huaweicloud_as_policies" "scaling_policy_name_filter" {
  scaling_group_id    = huaweicloud_as_group.acc_as_group.id
  scaling_policy_name = local.scaling_policy_name
}

output "is_scaling_policy_name_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_name_filter.policies[*].name : v == local.scaling_policy_name]
  )
}

// Filter using scaling policy type.
locals {
  scaling_policy_type = data.huaweicloud_as_policies.test.policies[0].type
}

data "huaweicloud_as_policies" "scaling_policy_type_filter" {
  scaling_group_id    = huaweicloud_as_group.acc_as_group.id
  scaling_policy_type = local.scaling_policy_type
}

output "is_scaling_policy_type_filter_useful" {
  value = length(data.huaweicloud_as_policies.scaling_policy_type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_as_policies.scaling_policy_type_filter.policies[*].type : v == local.scaling_policy_type]
  )
}
`, baseConfig)
}
