package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/coc"
)

func getDiagnosisTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("coc", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating COC client: %s", err)
	}

	return coc.GetDiagnosisTask(client, state.Primary.ID, state.Primary.Attributes["resource_id"])
}

func TestAccResourceCocDiagnosisTask_basic(t *testing.T) {
	var obj interface{}
	resourceName := "huaweicloud_coc_diagnosis_task.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getDiagnosisTaskResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocDiagnosisTask_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "user_name"),
					resource.TestCheckResourceAttrSet(resourceName, "work_order_id"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "start_time"),
					resource.TestCheckResourceAttrSet(resourceName, "os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "region"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.#"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.0.code"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.0.diagnosis_task_id"),
					resource.TestCheckResourceAttrSet(resourceName, "node_list.0.status"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"node_list.0.status", "progress", "status"},
				ImportStateIdFunc:       testDiagnosisTaskImportState(resourceName),
			},
		},
	})
}

func testCocDiagnosisTask_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_diagnosis_task" "test" {
  resource_id = "%s"
  type        = "ECS"
}
`, acceptance.HW_COC_INSTANCE_ID)
}

func testDiagnosisTaskImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		resourceID := rs.Primary.Attributes["resource_id"]
		if resourceID == "" {
			return "", fmt.Errorf("attribute (resource_id) of resource (%s) not found", name)
		}

		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of resource (%s) not found: %s", name, rs)
		}

		return fmt.Sprintf("%s/%s", resourceID, rs.Primary.ID), nil
	}
}
