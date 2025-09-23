package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iotda"
)

func getBatchTaskResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region    = acceptance.HW_REGION_NAME
		isDerived = iotda.WithDerivedAuth(cfg, region)
		product   = "iotda"
		taskId    = state.Primary.ID
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA client: %s", err)
	}

	// There is no need to handle pagination related parameters here, as it only applies to the structure of subtasks.
	return iotda.GetBatchTaskById(client, taskId)
}

func TestAccBatchTask_basic(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_batchtask.test_freeze"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBatchTaskResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBatchTask_basic(name, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "freezeDevices"),
					resource.TestCheckResourceAttrPair(resourceName, "space_id", "data.huaweicloud_iotda_spaces.test", "spaces.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "task_progress.0.total", "2"),
					resource.TestCheckResourceAttr(resourceName, "task_details.#", "2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id", "targets", "targets_filter", "targets_file"},
			},
		},
	})
}

func TestAccBatchTask_withTargetsFilterField(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_batchtask.test_unfreeze"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBatchTaskResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBatchTask_withTargetsFilterField(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "unfreezeDevices"),
					resource.TestCheckResourceAttr(resourceName, "status", "Fail"),

					resource.TestCheckResourceAttrSet(resourceName, "status_desc"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttr(resourceName, "task_details.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id", "targets", "targets_filter", "targets_file"},
			},
		},
	})
}

func testBatchTask_base(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_iotda_spaces" "test" {
  is_default = true
}

resource "huaweicloud_iotda_product" "test" {
  name        = "%[2]s"
  device_type = "test"
  protocol    = "MQTT"
  space_id    = data.huaweicloud_iotda_spaces.test.spaces[0].id
  data_type   = "json"

  services {
    id   = "service_1"
    type = "serv_type"
  }
}

resource "huaweicloud_iotda_device" "test1" {
  name       = "%[2]s"
  node_id    = "%[3]s"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
}

resource "huaweicloud_iotda_device" "test2" {
  name       = "%[2]s_2"
  node_id    = "%[3]s_2"
  space_id   = data.huaweicloud_iotda_spaces.test.spaces[0].id
  product_id = huaweicloud_iotda_product.test.id
  secret     = "1234567890"
}
`, buildIoTDAEndpoint(), name, nodeId)
}

func testBatchTask_basic(name, nodeId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_batchtask" "test_freeze" {
  name     = "%[2]s"
  type     = "freezeDevices"
  space_id = data.huaweicloud_iotda_spaces.test.spaces[0].id
  targets  = [huaweicloud_iotda_device.test1.id, huaweicloud_iotda_device.test2.id]
}
`, testBatchTask_base(name, nodeId), name, nodeId)
}

func testBatchTask_withTargetsFilterField(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_batchtask" "test_unfreeze" {
  name = "%[2]s"
  type = "unfreezeDevices"

  # The status of batch task created with non-existent group ID must be failed.
  targets_filter {
    group_ids = ["123456789", "987654321"]
  }
}
`, buildIoTDAEndpoint(), name)
}

func TestAccBatchTask_withTargetsFileField(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_batchtask.test_create"
		name         = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getBatchTaskResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIOTDABatchTaskFilePath(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testBatchTask_withTargetsFileField(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "createDevices"),

					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "task_details.#"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"space_id", "targets", "targets_filter", "targets_file"},
			},
		},
	})
}

func testBatchTask_withTargetsFileField(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_batchtask" "test_create" {
  name         = "%[2]s"
  type         = "createDevices"
  targets_file = "%[3]s"
}
`, buildIoTDAEndpoint(), name, acceptance.HW_IOTDA_BATCHTASK_FILE_PATH)
}
