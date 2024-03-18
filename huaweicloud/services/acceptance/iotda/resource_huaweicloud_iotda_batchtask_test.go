package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getBatchTaskResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcIoTdaV5Client(acceptance.HW_REGION_NAME, WithDerivedAuth())
	if err != nil {
		return nil, fmt.Errorf("error creating IoTDA v5 client: %s", err)
	}

	// There is no need to handle pagination related parameters here, as it only applies to the structure of subtasks.
	resp, err := client.ShowBatchTask(&model.ShowBatchTaskRequest{
		TaskId: state.Primary.ID,
	})

	if err != nil {
		return nil, fmt.Errorf("error querying IoTDA batch task: %s", err)
	}

	return resp, nil
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
					resource.TestCheckResourceAttrPair(resourceName, "space_id", "huaweicloud_iotda_space.test", "id"),

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

func TestAccBatchTask_derived(t *testing.T) {
	var (
		obj          interface{}
		resourceName = "huaweicloud_iotda_batchtask.test_derived"
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
				Config: testBatchTask_derived(name),
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
	productBasic := testProduct_basic(name)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_device" "test1" {
  name       = "%[2]s"
  node_id    = "%[3]s"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
}

resource "huaweicloud_iotda_device" "test2" {
  name       = "%[2]s_2"
  node_id    = "%[3]s_2"
  space_id   = huaweicloud_iotda_space.test.id
  product_id = huaweicloud_iotda_product.test.id
  secret     = "1234567890"
}

`, productBasic, name, nodeId)
}

func testBatchTask_basic(name, nodeId string) string {
	batchTaskBase := testBatchTask_base(name, nodeId)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_batchtask" "test_freeze" {
  name     = "%[2]s"
  type     = "freezeDevices"
  space_id = huaweicloud_iotda_space.test.id
  targets  = [huaweicloud_iotda_device.test1.id, huaweicloud_iotda_device.test2.id]
}
`, batchTaskBase, name, nodeId)
}

func testBatchTask_withTargetsFilterField(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_iotda_batchtask" "test_unfreeze" {
  name = "%s"
  type = "unfreezeDevices"

  # The status of batch task created with non-existent group ID must be failed.
  targets_filter {
    group_ids = ["123456789", "987654321"]
  }
}
`, name)
}

func testBatchTask_withTargetsFileField(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_iotda_batchtask" "test_create" {
  name         = "%[1]s"
  type         = "createDevices"
  targets_file = "%[2]s"
}
`, name, acceptance.HW_IOTDA_BATCHTASK_FILE_PATH)
}

func testBatchTask_derived(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_iotda_batchtask" "test_derived" {
  name = "%[2]s"
  type = "unfreezeDevices"

  # The status of batch task created with non-existent group ID must be failed.
  targets_filter {
    group_ids = ["123456789", "987654321"]
  }
}
`, buildIoTDAEndpoint(), name)
}
