package bms

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getVolumeAttach(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/os-volume_attachments"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating BMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", state.Primary.Attributes["server_id"])

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	searchPath := fmt.Sprintf("volumeAttachments[?volumeId=='%s']|[0]", state.Primary.Attributes["volume_id"])
	volumeAttach := utils.PathSearch(searchPath, getRespBody, nil)
	if volumeAttach == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func TestAccBmsVolumeAttach_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_bms_volume_attach.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getVolumeAttach,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBmsVolumeAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_bms_instance.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id",
						"huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "device", "/dev/sdb"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccBmsVolumeAttach_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_bms_flavors" "test" {
  cpu_arch          = "x86_64"
  memory            = "192"
  vcpus             = "56"
  availability_zone = try(element(data.huaweicloud_availability_zones.test.names, 0), "")
}

data "huaweicloud_images_images" "test" {
  name_regex = "x86"
  os         = "CentOS"
  image_type = "Ironic"
}

locals {
  x86_images = [for v in data.huaweicloud_images_images.test.images: v.id if v.container_format == "bare"]
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_bms_instance" "test" {
  security_groups   = [huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  flavor_id         = data.huaweicloud_bms_flavors.test.flavors[0].id
  key_pair          = huaweicloud_kps_keypair.test.name
  image_id          = try(local.x86_images[0], "")
  name              = "%[2]s"
  user_id           = "%[3]s"
  system_disk_type  = "GPSSD"
  system_disk_size  = 150

  nics {
    subnet_id = huaweicloud_vpc_subnet.test.id
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = "1"
  auto_renew    = "false"
}

resource "huaweicloud_evs_volume" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 100
  volume_type       = "GPSSD2"
  iops              = 3000
  throughput        = 125
  device_type       = "SCSI"
  multiattach       = false
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_USER_ID)
}

func testAccBmsVolumeAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_bms_volume_attach" "test" {
  server_id = huaweicloud_bms_instance.test.id
  volume_id = huaweicloud_evs_volume.test.id
  device    = "/dev/sdb"
}
`, testAccBmsVolumeAttach_base(rName))
}
