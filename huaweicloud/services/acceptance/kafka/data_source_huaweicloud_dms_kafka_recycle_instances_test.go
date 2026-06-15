package kafka

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRecycleInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_kafka_recycle_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSKafkaRecycleBinInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRecycleInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "retention_days"),
					resource.TestCheckResourceAttr(all, "default_use_recycle", "true"),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.status"),
					resource.TestCheckResourceAttrSet(all, "instances.0.engine"),
					resource.TestMatchResourceAttr(all, "instances.0.in_recycle_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.save_time"),
					resource.TestMatchResourceAttr(all, "instances.0.auto_delete_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.product_id"),
					resource.TestCheckOutput("instance_id_validation", "true"),
				),
			},
		},
	})
}

func testAccDataRecycleInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_kafka_recycle_instances" "test" {}

output "instance_id_validation" {
  value = contains(data.huaweicloud_dms_kafka_recycle_instances.test.instances[*].id, "%[1]s")
}
`, acceptance.HW_DMS_KAFKA_RECYCLE_BIN_INSTANCE_ID)
}
