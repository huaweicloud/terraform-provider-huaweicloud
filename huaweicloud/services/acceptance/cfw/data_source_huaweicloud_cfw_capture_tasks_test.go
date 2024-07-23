package cfw

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwCaptureTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_capture_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwCaptureTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.source_address"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dest_address"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.duration"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.remaining_days"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.max_packets"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.capture_size"),
					resource.TestMatchResourceAttr(dataSource, "records.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "records.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testDataSourceCfwCaptureTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cfw_capture_tasks" "test" {
  fw_instance_id = "%[2]s" 

  depends_on = [
    huaweicloud_cfw_capture_task.test
  ]
}
`, testDataSourceCfwCaptureTasks_base(name), acceptance.HW_CFW_INSTANCE_ID)
}

func testDataSourceCfwCaptureTasks_base(name string) string {
	return fmt.Sprintf(`
resource huaweicloud_cfw_capture_task "test" {
  fw_instance_id = "%[1]s"
  name           = "%[2]s"
  duration       = 5
  max_packets    = 100000
		
  destination {
    address      = "1.1.1.1"
    address_type = 0
  }
		
  source {
    address      = "2.2.2.2"
    address_type = 0
  }
		
  service {
    dest_port   = "80"
    protocol    =  6
    source_port = "80"
  }
}
`, acceptance.HW_CFW_INSTANCE_ID, name)
}
