package ims

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEvsSystemImage_basic(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = acceptance.RandomAccResourceName()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_ims_evs_system_image.test"
		defaultEpsId = "0"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImsImageResourceFunc,
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
				Config: testAccEvsSystemImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "os_version", "Windows Server 2019 Standard 64bit"),
					resource.TestCheckResourceAttr(resourceName, "type", "ECS"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform description test"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "2048"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "visibility"),
					resource.TestCheckResourceAttrSet(resourceName, "image_size"),
					resource.TestCheckResourceAttrSet(resourceName, "os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "min_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(resourceName, "data_origin"),
					resource.TestMatchResourceAttr(resourceName, "active_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "huaweicloud_evs_volume.test", "id"),
				),
			},
			{
				Config: testAccEvsSystemImage_update1(rName, rNameUpdate, migrateEpsId, 2048, 1024),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "2048"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", migrateEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				Config: testAccEvsSystemImage_update2(rName, rNameUpdate, defaultEpsId, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"type",
				},
			},
		},
	})
}

func testAccEvsSystemImage_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_evs_system_image" "test" {
  name        = "%[2]s"
  volume_id   = huaweicloud_evs_volume.test.id
  os_version  = "Windows Server 2019 Standard 64bit"
  type        = "ECS"
  description = "terraform description test"
  max_ram     = 4096
  min_ram     = 2048

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsDataImage_base(rName), rName)
}

func testAccEvsSystemImage_update1(rName, rNameUpdate, migrateEpsId string, maxRAM, minRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_evs_system_image" "test" {
  name                  = "%[2]s"
  volume_id             = huaweicloud_evs_volume.test.id
  os_version            = "Windows Server 2019 Standard 64bit"
  type                  = "ECS"
  max_ram               = %[4]d
  min_ram               = %[5]d
  enterprise_project_id = "%[3]s"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccEvsDataImage_base(rName), rNameUpdate, migrateEpsId, maxRAM, minRAM)
}

func testAccEvsSystemImage_update2(rName, rNameUpdate, defaultEpsId string, maxRAM, minRAM int) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_evs_system_image" "test" {
  name                  = "%[2]s"
  volume_id             = huaweicloud_evs_volume.test.id
  os_version            = "Windows Server 2019 Standard 64bit"
  type                  = "ECS"
  max_ram               = %[4]d
  min_ram               = %[5]d
  enterprise_project_id = "%[3]s"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccEvsDataImage_base(rName), rNameUpdate, defaultEpsId, maxRAM, minRAM)
}
