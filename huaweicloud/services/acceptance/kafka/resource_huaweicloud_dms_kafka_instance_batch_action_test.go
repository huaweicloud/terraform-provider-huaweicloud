package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccInstanceBatchAction_basic(t *testing.T) {
	name := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testInstanceBatchAction_basic_step1(name),
			},
			{
				Config: testInstanceBatchAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("is_delete_success", "true"),
				),
				// After deleting an instance using huaweicloud_dms_kafka_instance_batch_action, terraform will detect
				// the change in the next plan because the instance status is modified. Set ExpectNonEmptyPlan to true to expect
				// this non-empty plan as the normal behavior.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testInstanceBatchAction_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_dms_kafka_flavors" "test" {
  type = "single"
}

locals {
  flavor = data.huaweicloud_dms_kafka_flavors.test.flavors[0]
}

resource "huaweicloud_dms_kafka_instance" "test" {
  count = 2

  name               = "%[2]s${count.index}"
  vpc_id             = huaweicloud_vpc.test.id
  network_id         = huaweicloud_vpc_subnet.test.id
  security_group_id  = huaweicloud_networking_secgroup.test.id
  flavor_id          = local.flavor.id
  storage_spec_code  = "dms.physical.storage.extreme"
  engine_version     = "3.x"
  broker_num         = 1
  storage_space      = try(local.flavor.properties[0].min_broker * local.flavor.properties[0].min_storage_per_node, null)
  availability_zones = slice(data.huaweicloud_availability_zones.test.names, 0, 1)
}
`, common.TestBaseNetwork(name), name)
}

func testInstanceBatchAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_batch_action" "test" {
  instances = huaweicloud_dms_kafka_instance.test[*].id
  action    = "restart"
}
`, testInstanceBatchAction_base(name))
}

func testInstanceBatchAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_batch_action" "delete" {
  instances    = huaweicloud_dms_kafka_instance.test[*].id
  action       = "delete"
  force_delete = true

  lifecycle {
    ignore_changes = [instances]
  }
}

data "huaweicloud_dms_kafka_instances" "test" {
  count = 2

  instance_id = huaweicloud_dms_kafka_instance.test[count.index].id

  depends_on = [huaweicloud_dms_kafka_instance_batch_action.delete]
}

output "is_delete_success" {
  value = length(flatten(data.huaweicloud_dms_kafka_instances.test[*].instances)) == 0
}
`, testInstanceBatchAction_base(name))
}
