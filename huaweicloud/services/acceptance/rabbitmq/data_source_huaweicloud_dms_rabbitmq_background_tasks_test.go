package rabbitmq

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDmsRabbitMQBackgroundTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dms_rabbitmq_background_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDmsRabbitMQBackgroundTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.params"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.user_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceDmsRabbitMQBackgroundTasks_basic(name string) string {
	beginTime := time.Now().UTC().Format("20060102150405")
	endTime := time.Now().Add(time.Hour).UTC().Format("20060102150405")

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dms_rabbitmq_background_tasks" "test" {
  instance_id = huaweicloud_dms_rabbitmq_instance.test.id
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, testAccDmsRabbitmqInstance_newFormat_single(name), beginTime, endTime)
}
