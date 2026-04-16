package rabbitmq

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRecycleInstances_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dms_rabbitmq_recycle_instances.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDMSRocketMQRecycleBinInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRecycleInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "retention_days"),
					resource.TestCheckResourceAttr(all, "default_use_recycle", "true"),
					resource.TestMatchResourceAttr(all, "instances.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(all, "instances.0.name"),
					resource.TestCheckResourceAttrSet(all, "instances.0.status"),
					resource.TestCheckResourceAttrSet(all, "instances.0.engine"),
					resource.TestCheckResourceAttrSet(all, "instances.0.in_recycle_time"),
					resource.TestCheckResourceAttrSet(all, "instances.0.save_time"),
					resource.TestCheckResourceAttrSet(all, "instances.0.auto_delete_time"),
					resource.TestCheckResourceAttrSet(all, "instances.0.product_id"),
					resource.TestCheckOutput("instance_id_validation", "true"),
				),
			},
		},
	})
}

func testDataSourceRecycleInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dms_rabbitmq_recycle_instances" "test" {}

output "instance_id_validation" {
  value = contains(data.huaweicloud_dms_rabbitmq_recycle_instances.test.instances[*].instance_id, "%[1]s")
}
`, acceptance.HW_DMS_RABBITMQ_RECYCLE_BIN_INSTANCE_ID)
}
