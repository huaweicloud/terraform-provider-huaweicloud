package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataTopicQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dms_kafka_topic_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byKeyword   = "data.huaweicloud_dms_kafka_topic_quotas.filter_by_keyword"
		dcByKeyword = acceptance.InitDataSourceCheck(byKeyword)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataTopicQuotas_instanceNotFound(),
				ExpectError: regexp.MustCompile("This DMS instance does not exist"),
			},
			{
				Config: testAccDataTopicQuotas_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "quotas.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByKeyword.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "quotas.0.topic"),
					resource.TestMatchResourceAttr(dataSource, "quotas.0.consumer_byte_rate", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dataSource, "quotas.0.producer_byte_rate", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckOutput("is_keyword_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataTopicQuotas_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_topic_quotas" "test" {
  instance_id = "%[1]s"
}
`, randomId)
}

func testAccDataTopicQuotas_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_topic" "test" {
  count = 2

  instance_id = "%[1]s"
  name        = "%[2]s_${count.index}"
  partitions  = 10
}

resource "huaweicloud_dms_kafka_topic_quota" "test" {
  count = 2

  instance_id        = "%[1]s"
  topic              = huaweicloud_dms_kafka_topic.test[count.index].name
  consumer_byte_rate = 2097152
  producer_byte_rate = 2097152
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}

func testAccDataTopicQuotas_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_topic_quotas" "test" {
  instance_id = "%[2]s"

  depends_on = [huaweicloud_dms_kafka_topic_quota.test]
}

data "huaweicloud_dms_kafka_topic_quotas" "filter_by_keyword" {
  instance_id = "%[2]s"
  keyword     = "%[3]s"

  depends_on = [huaweicloud_dms_kafka_topic_quota.test]
}

locals {
  keyword_filter_result = [for v in data.huaweicloud_dms_kafka_topic_quotas.filter_by_keyword.quotas :
  strcontains(v.topic, "%[3]s")]
}

output "is_keyword_filter_useful" {
  value = length(local.keyword_filter_result) >= 2 && alltrue(local.keyword_filter_result)
}
`, testAccDataTopicQuotas_base(name), acceptance.HW_DMS_KAFKA_INSTANCE_ID, name)
}
