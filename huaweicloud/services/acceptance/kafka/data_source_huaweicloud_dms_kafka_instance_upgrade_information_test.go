package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataInstanceUpgradeInformation_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dms_kafka_instance_upgrade_information.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataInstanceUpgradeInformation_instanceNotFound(),
				ExpectError: regexp.MustCompile(`This DMS instance does not exist`),
			},
			{
				Config: testAccDataInstanceUpgradeInformation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "current_version"),
					resource.TestCheckResourceAttrSet(dataSource, "latest_version"),
				),
			},
		},
	})
}

func testAccDataInstanceUpgradeInformation_instanceNotFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instance_upgrade_information" "test" {
  instance_id = "%[1]s"
}
`, randomId)
}

func testAccDataInstanceUpgradeInformation_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_instance_upgrade_information" "test" {
  instance_id = "%[1]s"
}
`, acceptance.HW_DMS_KAFKA_INSTANCE_ID)
}
