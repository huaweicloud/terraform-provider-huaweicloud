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

func getKafkaMessageDiagnosisTaskReportResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dmsv2", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DMS client: %s", err)
	}

	return kafka.GetKafkaMessageDiagnosisTaskReport(client, state.Primary.Attributes["instance_id"], state.Primary.ID)
}

func TestAccKafkaMessageDiagnosisTaskReport_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_dms_kafka_message_diagnosis_task.test"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getKafkaMessageDiagnosisTaskReportResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
			acceptance.TestAccPreCheckDMSKafkaTopicName(t)
			acceptance.TestAccPreCheckDMSKafkaConsumerGroupName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKafkaMessageDiagnosisTaskReport_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "instance_id", acceptance.HW_DMS_KAFKA_INSTANCE_ID),
					resource.TestCheckResourceAttr(resourceName, "group_name", acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME),
					resource.TestCheckResourceAttr(resourceName, "topic_name", acceptance.HW_DMS_KAFKA_TOPIC_NAME),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "begin_time"),
					resource.TestCheckResourceAttrSet(resourceName, "end_time"),
					resource.TestCheckResourceAttrSet(resourceName, "accumulated_partitions"),
					resource.TestCheckResourceAttrSet(resourceName, "diagnosis_dimension_list.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccResourceMessageDiagnosisTaskImportStateIDFunc(resourceName),
			},
		},
	})
}

func testKafkaMessageDiagnosisTaskReport_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_message_diagnosis_task" "test" {
  instance_id = "%[1]s"
  group_name  = "%[2]s"
  topic_name  = "%[3]s"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.HW_DMS_KAFKA_CONSUMER_GROUP_NAME, acceptance.HW_DMS_KAFKA_TOPIC_NAME)
}

func testAccResourceMessageDiagnosisTaskImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", resourceName, rs)
		}

		instanceID := rs.Primary.Attributes["instance_id"]
		reportID := rs.Primary.ID
		if instanceID == "" || reportID == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<report_id>', but got '%s/%s'",
				instanceID, reportID)
		}
		return fmt.Sprintf("%s/%s", instanceID, reportID), nil
	}
}
