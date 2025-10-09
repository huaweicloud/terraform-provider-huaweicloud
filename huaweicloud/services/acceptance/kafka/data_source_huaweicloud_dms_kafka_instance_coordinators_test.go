package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceCoordinators_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_kafka_instance_coordinators.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataInstanceCoordinators_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataInstanceCoordinators_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "coordinators.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "coordinators.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "coordinators.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "coordinators.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "coordinators.0.group_id"),
				),
			},
		},
	})
}

func testAccDataInstanceCoordinators_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instance_coordinators" "test" {
  instance_id = "%[1]s"
}
`, randomId)
}

func testAccDataInstanceCoordinators_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
}

data "huaweicloud_dms_kafka_instance_coordinators" "test" {
  instance_id = "%[1]s"

  depends_on = [huaweicloud_dms_kafka_consumer_group.test]
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID, acceptance.RandomAccResourceName())
}
