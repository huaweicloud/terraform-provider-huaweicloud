package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkaMessageDiagnosisTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_message_diagnosis_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkaMessageDiagnosisTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.report_id"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.group_name"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.topic_name"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.accumulated_partitions"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "report_list.0.end_time"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkaMessageDiagnosisTasks_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_kafka_message_diagnosis_tasks" "test" {
  depends_on = [huaweicloud_dms_kafka_message_diagnosis_task.test]

  instance_id = "%[2]s"
}
`, testKafkaMessageDiagnosisTaskReport_basic(), acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
