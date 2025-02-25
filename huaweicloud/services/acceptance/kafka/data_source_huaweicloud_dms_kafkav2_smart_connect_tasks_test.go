package kafka

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsKafkav2SmartConnectTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafkav2_smart_connect_tasks.test"
	rName := acceptance.RandomAccResourceNameWithDash()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkav2SmartConnectTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkav2SmartConnectTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafkav2_smart_connect_tasks" "test" {
  depends_on = [huaweicloud_dms_kafkav2_smart_connect_task.test]

  instance_id = huaweicloud_dms_kafka_instance.test.id
}
`, testDmsKafkav2SmartConnectTask_basic(name))
}

func TestAccDataSourceDmsKafkav2SmartConnectTasks_kafkaToKafka(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafkav2_smart_connect_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDmsKafkav2SmartConnectTasks_kafkaToKafka(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDmsKafkav2SmartConnectTasks_kafkaToKafka(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafkav2_smart_connect_tasks" "test" {
  depends_on = [huaweicloud_dms_kafkav2_smart_connect_task.test]

  instance_id = huaweicloud_dms_kafka_instance.test1.id
}
`, testDmsKafkav2SmartConnectTask_kafkaToKafka(name))
}
