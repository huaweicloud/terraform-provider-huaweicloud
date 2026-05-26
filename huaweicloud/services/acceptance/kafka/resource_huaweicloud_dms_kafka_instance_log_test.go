package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

func getInstanceLogResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetInstanceLog(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["log_type"])
}

func TestAccInstanceLog_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dms_kafka_instance_log.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getInstanceLogResourceFunc)

		rNameBalance = "huaweicloud_dms_kafka_instance_log.balance"
		rcBalance    = acceptance.InitResourceCheck(rNameBalance, &obj, getInstanceLogResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rc.CheckResourceDestroy(),
			rcBalance.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceLog_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "instance_id", acceptance.HW_DMS_KAFKA_INSTANCE_ID),
					resource.TestCheckResourceAttr(rName, "log_type", "topic_log"),
					resource.TestCheckResourceAttrSet(rName, "log_file_name"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttrSet(rName, "log_group_id"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_id"),
					resource.TestCheckResourceAttr(rName, "status", "OPEN"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					rcBalance.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameBalance, "log_type", "REBALANCE"),
					resource.TestCheckResourceAttr(rNameBalance, "log_file_name", "coordinator.log"),
					resource.TestCheckResourceAttrPair(rNameBalance, "log_group_name",
						"huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(rNameBalance, "log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
					resource.TestCheckResourceAttrPair(rNameBalance, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rNameBalance, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
				),
			},
			{
				Config: testAccInstanceLog_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "log_type", "topic_log"),
					resource.TestCheckResourceAttr(rName, "log_file_name", "topic.log"),
					resource.TestCheckResourceAttrPair(rName, "log_group_name",
						"huaweicloud_lts_group.test", "group_name"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_name",
						"huaweicloud_lts_stream.test", "stream_name"),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),

					rcBalance.CheckResourceExists(),
					resource.TestCheckResourceAttr(rNameBalance, "log_type", "REBALANCE"),
					resource.TestCheckResourceAttrSet(rNameBalance, "log_file_name"),
					resource.TestCheckResourceAttrSet(rNameBalance, "log_group_name"),
					resource.TestCheckResourceAttrSet(rNameBalance, "log_stream_name"),
					resource.TestCheckResourceAttrSet(rNameBalance, "log_group_id"),
					resource.TestCheckResourceAttrSet(rNameBalance, "log_stream_id"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_force_new"},
			},
			{
				ResourceName:            rNameBalance,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"enable_force_new"},
			},
		},
	})
}

func testAccInstanceLog_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  type = string
  default = "%[1]s"
}

resource "huaweicloud_lts_group" "test" {
  group_name            = "%[2]s"
  ttl_in_days           = 7
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = "%[2]s"
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccInstanceLog_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_log" "test" {
  instance_id = "%[2]s"
  log_type    = "topic_log"

  depends_on = [
    huaweicloud_lts_group.test,
    huaweicloud_lts_stream.test,
  ]
}

resource "huaweicloud_dms_kafka_instance_log" "balance" {
  instance_id     = "%[2]s"
  log_type        = "REBALANCE"
  log_file_name   = "coordinator.log"
  log_group_name  = huaweicloud_lts_group.test.group_name
  log_stream_name = huaweicloud_lts_stream.test.stream_name

  depends_on = [
    huaweicloud_lts_group.test,
    huaweicloud_lts_stream.test,
  ]
}
`, testAccInstanceLog_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccInstanceLog_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_instance_log" "test" {
  instance_id      = "%[2]s"
  log_type         = "topic_log"
  log_file_name    = "topic.log"
  log_group_name   = huaweicloud_lts_group.test.group_name
  log_stream_name  = huaweicloud_lts_stream.test.stream_name
  enable_force_new = "true"
}

resource "huaweicloud_dms_kafka_instance_log" "balance" {
  instance_id      = "%[2]s"
  log_type         = "REBALANCE"
  enable_force_new = "true"
}
`, testAccInstanceLog_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
