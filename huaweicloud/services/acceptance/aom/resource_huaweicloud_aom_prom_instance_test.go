package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/aom"
)

func getPromInstanceResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("aom", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating AOM client: %s", err)
	}

	return aom.GetPrometheusInstanceById(client, state.Primary.ID)
}

func TestAccPromInstance_basic(t *testing.T) {
	var (
		obj interface{}

		updateAfterCreate  = "huaweicloud_aom_prom_instance.update_limits_after_create"
		rcAfterCreate      = acceptance.InitResourceCheck(updateAfterCreate, &obj, getPromInstanceResourceFunc)
		updateDuringUpdate = "huaweicloud_aom_prom_instance.update_limits_during_update"
		rcDuringUpdate     = acceptance.InitResourceCheck(updateDuringUpdate, &obj, getPromInstanceResourceFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcAfterCreate.CheckResourceDestroy(),
			rcDuringUpdate.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAOMPromInstance_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcAfterCreate.CheckResourceExists(),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_name", fmt.Sprintf("%s_after_create", name)),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_type", "VPC"),
					resource.TestCheckResourceAttr(updateAfterCreate, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_version", "1.5"),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_limits.#", "1"),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_limits.0.compactor_blocks_retention_period", "360h"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "created_at"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "remote_write_url"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "remote_read_url"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "prom_http_api_endpoint"),
					rcDuringUpdate.CheckResourceExists(),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_name", fmt.Sprintf("%s_during_update", name)),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_type", "VPC"),
					resource.TestCheckResourceAttr(updateDuringUpdate, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_version", "1.5"),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_limits.#", "1"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "prom_limits.0.compactor_blocks_retention_period"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "created_at"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "remote_write_url"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "remote_read_url"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "prom_http_api_endpoint"),
				),
			},
			{
				Config: testAOMPromInstance_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rcAfterCreate.CheckResourceExists(),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_name", fmt.Sprintf("%s_after_create", updateName)),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_type", "VPC"),
					resource.TestCheckResourceAttr(updateAfterCreate, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_version", "1.5"),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_limits.#", "1"),
					resource.TestCheckResourceAttr(updateAfterCreate, "prom_limits.0.compactor_blocks_retention_period", "360h"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "created_at"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "remote_write_url"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "remote_read_url"),
					resource.TestCheckResourceAttrSet(updateAfterCreate, "prom_http_api_endpoint"),
					rcDuringUpdate.CheckResourceExists(),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_name", fmt.Sprintf("%s_during_update", updateName)),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_type", "VPC"),
					resource.TestCheckResourceAttr(updateDuringUpdate, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_version", "1.5"),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_limits.#", "1"),
					resource.TestCheckResourceAttr(updateDuringUpdate, "prom_limits.0.compactor_blocks_retention_period", "1440h"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "created_at"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "remote_write_url"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "remote_read_url"),
					resource.TestCheckResourceAttrSet(updateDuringUpdate, "prom_http_api_endpoint"),
				),
			},
			{
				ResourceName:      updateAfterCreate,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      updateDuringUpdate,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAOMPromInstance_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "update_limits_after_create" {
  prom_name             = "%[1]s_after_create"
  prom_type             = "VPC"
  enterprise_project_id = "%[2]s"
  prom_version          = "1.5"

  prom_limits {
    compactor_blocks_retention_period = "360h"
  }
}

resource "huaweicloud_aom_prom_instance" "update_limits_during_update" {
  prom_name             = "%[1]s_during_update"
  prom_type             = "VPC"
  enterprise_project_id = "%[2]s"
  prom_version          = "1.5"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAOMPromInstance_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_prom_instance" "update_limits_after_create" {
  prom_name             = "%[1]s_after_create"
  prom_type             = "VPC"
  enterprise_project_id = "%[2]s"
  prom_version          = "1.5"

  prom_limits {
    compactor_blocks_retention_period = "360h"
  }
}

resource "huaweicloud_aom_prom_instance" "update_limits_during_update" {
  prom_name             = "%[1]s_during_update"
  prom_type             = "VPC"
  enterprise_project_id = "%[2]s"
  prom_version          = "1.5"

  prom_limits {
    compactor_blocks_retention_period = "1440h"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
