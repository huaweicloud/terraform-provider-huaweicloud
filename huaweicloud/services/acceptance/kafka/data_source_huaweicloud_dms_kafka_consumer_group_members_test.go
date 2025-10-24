package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, make sure that the consumer group status is `STABLE`.
func TestAccDataConsumerGroupMembers_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_consumer_group_members.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byHost   = "data.huaweicloud_dms_kafka_consumer_group_members.filter_by_host"
		dcByHost = acceptance.InitDataSourceCheck(byHost)

		byMemberId   = "data.huaweicloud_dms_kafka_consumer_group_members.filter_by_member_id"
		dcByMemberId = acceptance.InitDataSourceCheck(byMemberId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataConsumerGroupMembers_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataConsumerGroupMembers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "members.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "members.0.id"),
					resource.TestCheckResourceAttrSet(all, "members.0.host"),
					resource.TestCheckResourceAttrSet(all, "members.0.client_id"),
					dcByHost.CheckResourceExists(),
					resource.TestCheckOutput("is_host_filter_useful", "true"),
					dcByMemberId.CheckResourceExists(),
					resource.TestCheckOutput("is_member_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataConsumerGroupMembers_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_consumer_group_members" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
}
`, randomId, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}

func testAccDataConsumerGroupMembers_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_consumer_group_members" "test" {
  instance_id = "%[1]s"
  group       = "%[2]s"
}

locals {
  host      = try(data.huaweicloud_dms_kafka_consumer_group_members.test.members[0].host, null)
  member_id = try(data.huaweicloud_dms_kafka_consumer_group_members.test.members[0].id, null)
}

data "huaweicloud_dms_kafka_consumer_group_members" "filter_by_host" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  host        = local.host
}

locals {
  filter_by_host_result = [for v in data.huaweicloud_dms_kafka_consumer_group_members.filter_by_host.members : strcontains(v.host, local.host)]
}

output "is_host_filter_useful" {
  value = length(local.filter_by_host_result) >= 0 && alltrue(local.filter_by_host_result)
}

data "huaweicloud_dms_kafka_consumer_group_members" "filter_by_member_id" {
  instance_id = "%[1]s"
  group       = "%[2]s"
  member_id   = local.member_id
}

locals {
  filter_by_member_id_result = [for v in data.huaweicloud_dms_kafka_consumer_group_members.filter_by_member_id.members :
  strcontains(v.id, local.member_id)]
}

output "is_member_id_filter_useful" {
  value = length(local.filter_by_member_id_result) >= 0 && alltrue(local.filter_by_member_id_result)
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME)
}
