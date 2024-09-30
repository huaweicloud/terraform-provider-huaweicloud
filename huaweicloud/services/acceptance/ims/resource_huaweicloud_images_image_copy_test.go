package ims

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ims"
)

func getImageCopyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	targetRegion := state.Primary.Attributes["target_region"]
	if targetRegion != "" {
		region = targetRegion
	}

	imsClient, err := cfg.ImageV2Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS v2 client: %s", err)
	}

	imageList, err := ims.GetImageList(imsClient, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS images: %s", err)
	}

	if len(imageList) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return imageList[0], nil
}

func TestAccImageCopy_basic(t *testing.T) {
	var (
		obj             interface{}
		sourceImageName = acceptance.RandomAccResourceName()
		name            = sourceImageName + "-copy"
		updateName      = name + "-update"
		defaultEpsId    = "0"
		migrateEpsId    = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
		rName           = "huaweicloud_images_image_copy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a non default enterprise project ID.
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImageCopy_within_region_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(rName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttrPair(rName, "source_image_id", "huaweicloud_ims_ecs_system_image.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "kms_key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "os_version"),
					resource.TestCheckResourceAttrSet(rName, "visibility"),
					resource.TestCheckResourceAttrSet(rName, "min_disk"),
					resource.TestCheckResourceAttrSet(rName, "data_origin"),
					resource.TestCheckResourceAttrSet(rName, "disk_format"),
					resource.TestCheckResourceAttrSet(rName, "image_size"),
					resource.TestMatchResourceAttr(rName, "active_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testImageCopy_within_region_update(sourceImageName, updateName, migrateEpsId, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", migrateEpsId),
				),
			},
			{
				Config: testImageCopy_within_region_update(sourceImageName, updateName, defaultEpsId, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "min_ram", "0"),
					resource.TestCheckResourceAttr(rName, "max_ram", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", defaultEpsId),
				),
			},
		},
	})
}

func TestAccImageCopy_cross_region_basic(t *testing.T) {
	var (
		obj             interface{}
		sourceImageName = acceptance.RandomAccResourceName()
		name            = sourceImageName + "-copy"
		updateName      = name + "-update"
		rName           = "huaweicloud_images_image_copy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a region different from the region where the source image is located.
			acceptance.TestAccPreCheckDestRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImageCopy_cross_region_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(rName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrPair(rName, "source_image_id", "huaweicloud_ims_ecs_system_image.test", "id"),
				),
			},
			{
				Config: testImageCopy_cross_region_update(sourceImageName, updateName, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
			{
				Config: testImageCopy_cross_region_update(sourceImageName, updateName, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "min_ram", "0"),
					resource.TestCheckResourceAttr(rName, "max_ram", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func TestAccImageCopy_cross_region_withVaultId_basic(t *testing.T) {
	var (
		obj             interface{}
		sourceImageName = acceptance.RandomAccResourceName()
		name            = sourceImageName + "-copy"
		updateName      = name + "-update"
		rName           = "huaweicloud_images_image_copy.test"
	)

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getImageCopyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a region different from the region where the source image is located.
			acceptance.TestAccPreCheckDestRegion(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testImageCopy_cross_region_withVaultId_basic(sourceImageName, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", "description test"),
					resource.TestCheckResourceAttr(rName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(rName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "tags.key2", "value2"),
					resource.TestCheckResourceAttrPair(rName, "source_image_id", "huaweicloud_ims_ecs_whole_image.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "vault_id", "huaweicloud_cbr_vault.test_replication", "id"),
				),
			},
			{
				Config: testImageCopy_cross_region_withVaultId_update(sourceImageName, updateName, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(rName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
			{
				Config: testImageCopy_cross_region_withVaultId_update(sourceImageName, updateName, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "status", "active"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "min_ram", "0"),
					resource.TestCheckResourceAttr(rName, "max_ram", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key1", "value1_update"),
					resource.TestCheckResourceAttr(rName, "tags.key3", "value3"),
					resource.TestCheckResourceAttr(rName, "tags.key4", "value4"),
				),
			},
		},
	})
}

func testImageCopy_within_region_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[2]s"
  pending_days = "7"
}
`, testAccEcsSystemImage_basic(name), name)
}

func testImageCopy_within_region_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = huaweicloud_ims_ecs_system_image.test.id
  name            = "%[2]s"
  kms_key_id      = huaweicloud_kms_key.test.id
  min_ram         = 1024
  max_ram         = 4096
  description     = "description test"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}
`, testImageCopy_within_region_base(baseImageName), copyImageName)
}

func testImageCopy_within_region_update(baseImageName, copyImageName, migrateEpsId string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id       = huaweicloud_ims_ecs_system_image.test.id
  name                  = "%[2]s"
  kms_key_id            = huaweicloud_kms_key.test.id
  enterprise_project_id = "%[3]s"
  min_ram               = %[4]d
  max_ram               = %[5]d

  tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
  }
}
`, testImageCopy_within_region_base(baseImageName), copyImageName, migrateEpsId, minRAM, maxRAM)
}

func testImageCopy_cross_region_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = huaweicloud_ims_ecs_system_image.test.id
  name            = "%[2]s"
  target_region   = "%[3]s"
  agency_name     = "ims_admin_agency"
  min_ram         = 1024
  max_ram         = 4096
  description     = "description test"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}`, testAccEcsSystemImage_basic(baseImageName), copyImageName, acceptance.HW_DEST_REGION)
}

func testImageCopy_cross_region_update(baseImageName, copyImageName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = huaweicloud_ims_ecs_system_image.test.id
  name            = "%[2]s"
  target_region   = "%[3]s"
  agency_name     = "ims_admin_agency"
  min_ram         = %[4]d
  max_ram         = %[5]d

  tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
  }
}
`, testAccEcsSystemImage_basic(baseImageName), copyImageName, acceptance.HW_DEST_REGION, minRAM, maxRAM)
}

func testImageCopy_cross_region_withVaultId_base(baseImageName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cbr_vault" "test_replication" {
  region           = "%[3]s"
  name             = "%[2]s_replication"
  type             = "server"
  consistent_level = "crash_consistent"
  protection_type  = "replication"
  size             = 200
}

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "huaweicloud_ims_ecs_whole_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id
  vault_id    = huaweicloud_cbr_vault.test.id
}
`, testAccEcsSystemImage_base(baseImageName), baseImageName, acceptance.HW_DEST_REGION)
}

func testImageCopy_cross_region_withVaultId_basic(baseImageName, copyImageName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = huaweicloud_ims_ecs_whole_image.test.id
  name            = "%[2]s"
  target_region   = "%[3]s"
  agency_name     = "ims_admin_agency"
  vault_id        = huaweicloud_cbr_vault.test_replication.id
  min_ram         = 1024
  max_ram         = 4096
  description     = "description test"

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}`, testImageCopy_cross_region_withVaultId_base(baseImageName), copyImageName, acceptance.HW_DEST_REGION)
}

func testImageCopy_cross_region_withVaultId_update(baseImageName, copyImageName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = huaweicloud_ims_ecs_whole_image.test.id
  name            = "%[2]s"
  target_region   = "%[3]s"
  agency_name     = "ims_admin_agency"
  vault_id        = huaweicloud_cbr_vault.test_replication.id
  min_ram         = %[4]d
  max_ram         = %[5]d

  tags = {
    key1 = "value1_update"
    key3 = "value3"
    key4 = "value4"
  }
}
`, testImageCopy_cross_region_withVaultId_base(baseImageName), copyImageName, acceptance.HW_DEST_REGION, minRAM, maxRAM)
}
