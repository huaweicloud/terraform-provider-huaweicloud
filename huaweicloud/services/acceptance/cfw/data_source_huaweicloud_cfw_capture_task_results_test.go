package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwCaptureTaskResults_basic(t *testing.T) {
	var obj interface{}
	dName := "data.huaweicloud_cfw_capture_task_results.test"
	dc := acceptance.InitDataSourceCheck(dName)
	rName := "huaweicloud_cfw_capture_task.test"
	name := acceptance.RandomAccResourceName()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCaptureTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwCaptureTaskResults_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "duration", "5"),
					resource.TestCheckResourceAttr(rName, "max_packets", "100000"),
				),
			},
			{
				Config: testDataSourceCfwCaptureTaskResults_update(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dName, "captcha"),
					resource.TestCheckResourceAttrSet(dName, "expires"),
					resource.TestCheckResourceAttrSet(dName, "url"),
				),
			},
		},
	})
}

func testDataSourceCfwCaptureTaskResults_basic(name string) string {
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
}`, acceptance.HW_CFW_INSTANCE_ID, name)
}

func testDataSourceCfwCaptureTaskResults_update(name string) string {
	return fmt.Sprintf(`
resource huaweicloud_cfw_capture_task "test" {
  fw_instance_id = "%[1]s"
  name           = "%[2]s"
  duration       = 5
  max_packets    = 100000
  stop_capture   = true
			  
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

data "huaweicloud_cfw_capture_task_results" "test" {
  fw_instance_id = "%[1]s"
  task_id        = huaweicloud_cfw_capture_task.test.task_id
}
`, acceptance.HW_CFW_INSTANCE_ID, name)
}
