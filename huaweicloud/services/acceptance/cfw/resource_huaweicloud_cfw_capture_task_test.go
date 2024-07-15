package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cfw"
)

func getCaptureTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	product := "cfw"

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CFW client: %s", err)
	}

	return cfw.GetCaptureTask(client, state.Primary.ID, state.Primary.Attributes["fw_instance_id"])
}

func TestAccCaptureTask_basic(t *testing.T) {
	var obj interface{}

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
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCaptureTask_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "duration", "5"),
					resource.TestCheckResourceAttr(rName, "max_packets", "100000"),
					resource.TestCheckResourceAttr(rName, "destination.0.address", "1.1.1.1"),
					resource.TestCheckResourceAttr(rName, "destination.0.address_type", "0"),
					resource.TestCheckResourceAttr(rName, "source.0.address", "2.2.2.2"),
					resource.TestCheckResourceAttr(rName, "source.0.address_type", "0"),
					resource.TestCheckResourceAttr(rName, "service.0.dest_port", "80"),
					resource.TestCheckResourceAttr(rName, "service.0.protocol", "6"),
					resource.TestCheckResourceAttr(rName, "service.0.source_port", "80"),
					resource.TestCheckResourceAttrSet(rName, "task_id"),
					resource.TestCheckResourceAttrSet(rName, "created_at"),
					resource.TestCheckResourceAttrSet(rName, "updated_at"),
				),
			},
			{
				Config: testCaptureTask_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "duration", "5"),
					resource.TestCheckResourceAttr(rName, "max_packets", "100000"),
					resource.TestCheckResourceAttr(rName, "stop_capture", "true"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testCaptureTaskImportState(rName),
				ImportStateVerifyIgnore: []string{
					"stop_capture",
				},
			},
		},
	})
}

func testCaptureTask_basic(name string) string {
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

func testCaptureTask_update(name string) string {
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
`, acceptance.HW_CFW_INSTANCE_ID, name)
}

func testCaptureTaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["fw_instance_id"] == "" {
			return "", fmt.Errorf("attribute (fw_instance_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found: %s", name, rs)
		}
		return rs.Primary.Attributes["fw_instance_id"] + "/" + rs.Primary.ID, nil
	}
}
