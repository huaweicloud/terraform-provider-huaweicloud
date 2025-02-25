package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKafkaInstancesDataSource_basic(t *testing.T) {
	dc0 := acceptance.InitDataSourceCheck("data.huaweicloud_dms_kafka_instances.query_0")
	dc1 := acceptance.InitDataSourceCheck("data.huaweicloud_dms_kafka_instances.query_1")
	dc2 := acceptance.InitDataSourceCheck("data.huaweicloud_dms_kafka_instances.query_2")
	dc3 := acceptance.InitDataSourceCheck("data.huaweicloud_dms_kafka_instances.query_3")
	dc4 := acceptance.InitDataSourceCheck("data.huaweicloud_dms_kafka_instances.query_4")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKafkaInstancesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc0.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_dms_kafka_instances.query_0",
						"instances.#", regexp.MustCompile(`[1-9]\d*`)),
					dc1.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_dms_kafka_instances.query_1",
						"instances.#", regexp.MustCompile(`[1-9]\d*`)),
					dc2.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_dms_kafka_instances.query_2",
						"instances.#", regexp.MustCompile(`[1-9]\d*`)),
					dc3.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_dms_kafka_instances.query_3",
						"instances.#", regexp.MustCompile(`[1-9]\d*`)),
					dc4.CheckResourceExists(),
					resource.TestMatchResourceAttr("data.huaweicloud_dms_kafka_instances.query_4",
						"instances.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

func testAccKafkaInstancesDataSource_basic() string {
	rName := acceptance.RandomAccResourceNameWithDash()
	fuzzyName := rName[0 : len(rName)-1]

	return fmt.Sprintf(`
%s

data "huaweicloud_dms_kafka_instances" "query_0" {
  depends_on = [
    huaweicloud_dms_kafka_instance.test,
  ]

  name        = "%s"
  fuzzy_match = true
}

data "huaweicloud_dms_kafka_instances" "query_1" {
  depends_on = [
    huaweicloud_dms_kafka_instance.test,
  ]
  
  name = huaweicloud_dms_kafka_instance.test.name
}

data "huaweicloud_dms_kafka_instances" "query_2" {
  instance_id = huaweicloud_dms_kafka_instance.test.id
}

data "huaweicloud_dms_kafka_instances" "query_3" {
  depends_on = [
    huaweicloud_dms_kafka_instance.test,
  ]
  
  enterprise_project_id = huaweicloud_dms_kafka_instance.test.enterprise_project_id
}

data "huaweicloud_dms_kafka_instances" "query_4" {
  depends_on = [
    huaweicloud_dms_kafka_instance.test,
  ]

  status = "RUNNING"
}
`, testAccKafkaInstance_newFormat(rName), fuzzyName)
}
