package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/kafka"
)

func getTopicQuotaResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}
	return kafka.GetTopicQuotaByTopicName(client, state.Primary.Attributes["instance_id"], state.Primary.Attributes["topic"])
}

// Before running this test, ensure the instance is not a single-node instance.
func TestAccTopicQuota_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dms_kafka_topic_quota.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getTopicQuotaResourceFunc)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccTopicQuota_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "topic", name),
					resource.TestCheckResourceAttr(rName, "producer_byte_rate", "2097152"),
					resource.TestCheckResourceAttr(rName, "consumer_byte_rate", "3145728"),
				),
			},
			{
				Config: testAccTopicQuota_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "producer_byte_rate", "3145728"),
					resource.TestCheckResourceAttr(rName, "consumer_byte_rate", "0"),
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

func testAccTopicQuota_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  partitions  = 20
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccTopicQuota_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic_quota" "test" {
  instance_id = "%[2]s"
  topic       = huaweicloud_dms_kafka_topic.test.name

  # 2097152 B/s corresponds to 2 MB/s.
  producer_byte_rate = 2097152
  consumer_byte_rate = 3145728
}
`, testAccTopicQuota_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}

func testAccTopicQuota_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dms_kafka_topic_quota" "test" {
  instance_id        = "%[2]s"
  topic              = huaweicloud_dms_kafka_topic.test.name
  producer_byte_rate = 3145728
}
`, testAccTopicQuota_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
