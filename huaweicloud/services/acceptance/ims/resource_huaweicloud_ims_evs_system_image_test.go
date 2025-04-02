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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getEvsSystemImageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "ims"
		httpUrl = "v2/cloudimages"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += fmt.Sprintf("?id=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS EVS system image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, then return `404` error code.
	if image == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return image, nil
}

func TestAccEvsSystemImage_basic(t *testing.T) {
	var (
		image        interface{}
		rName        = acceptance.RandomAccResourceName()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_ims_evs_system_image.test"
		defaultEpsId = "0"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getEvsSystemImageResourceFunc,
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

func testAccEvsSystemImage_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  volume_type       = "GPSSD"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  server_id         = huaweicloud_compute_instance.test.id
  size              = 100
  charging_mode     = "postPaid"
}
`, common.TestBaseNetwork(rName), rName)
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
`, testAccEvsSystemImage_base(rName), rName)
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
`, testAccEvsSystemImage_base(rName), rNameUpdate, migrateEpsId, maxRAM, minRAM)
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
`, testAccEvsSystemImage_base(rName), rNameUpdate, defaultEpsId, maxRAM, minRAM)
}
