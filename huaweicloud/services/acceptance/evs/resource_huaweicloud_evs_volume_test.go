package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.BlockStorageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating block storage v2 client: %s", err)
	}
	return cloudvolumes.Get(c, state.Primary.ID).Extract()
}

func TestAccEvsVolume_basic(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_evs_volume.test"
	resourceName1 := "huaweicloud_evs_volume.test.0"
	resourceName2 := "huaweicloud_evs_volume.test.1"
	resourceName3 := "huaweicloud_evs_volume.test.2"
	resourceName4 := "huaweicloud_evs_volume.test.3"
	resourceName5 := "huaweicloud_evs_volume.test.4"
	resourceName6 := "huaweicloud_evs_volume.test.5"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(6),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Created by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName1, "size", "100"),
					resource.TestCheckResourceAttr(resourceName1, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName1, "tags.key", "value"),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume"),
					resource.TestCheckResourceAttr(resourceName1, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName1, "multiattach", "false"),

					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume"),
					resource.TestCheckResourceAttr(resourceName2, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName2, "multiattach", "true"),

					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume"),
					resource.TestCheckResourceAttr(resourceName3, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName3, "multiattach", "false"),

					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume"),
					resource.TestCheckResourceAttr(resourceName4, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName4, "multiattach", "true"),

					resource.TestCheckResourceAttr(resourceName5, "name", rName+"_gpssd2_normal_volume"),
					resource.TestCheckResourceAttr(resourceName5, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName5, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName5, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName5, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName5, "throughput", "500"),

					resource.TestCheckResourceAttr(resourceName6, "name", rName+"_essd2_normal_volume"),
					resource.TestCheckResourceAttr(resourceName6, "volume_type", "ESSD2"),
					resource.TestCheckResourceAttr(resourceName6, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName6, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName6, "iops", "3000"),
				),
			},
			{
				Config: testAccEvsVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(6),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"Updated by acc test script."),
					resource.TestCheckResourceAttr(resourceName1, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName1, "size", "200"),
					resource.TestCheckResourceAttr(resourceName1, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName1, "tags.key", "value1"),
					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_vbd_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_vbd_share_volume_update"),
					resource.TestCheckResourceAttr(resourceName3, "name", rName+"_scsi_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName4, "name", rName+"_scsi_share_volume_update"),
					resource.TestCheckResourceAttr(resourceName5, "name", rName+"_gpssd2_normal_volume_update"),
					resource.TestCheckResourceAttr(resourceName6, "name", rName+"_essd2_normal_volume_update"),
				),
			},
		},
	})
}

func TestAccEvsVolume_withEpsId(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_evs_volume.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_epsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
				},
			},
		},
	})
}

func TestAccEvsVolume_prePaid(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_evs_volume.test"
	resourceName1 := "huaweicloud_evs_volume.test.0"
	resourceName2 := "huaweicloud_evs_volume.test.1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaid(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					// Common configuration
					resource.TestCheckResourceAttrPair(resourceName1, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName1, "description",
						"test volume for charging mode"),
					resource.TestCheckResourceAttr(resourceName1, "size", "100"),

					// Personalized configuration
					resource.TestCheckResourceAttr(resourceName1, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName1, "name", rName+"_ssd_volume"),
					resource.TestCheckResourceAttr(resourceName1, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName1, "auto_renew", "false"),

					resource.TestCheckResourceAttr(resourceName2, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName2, "name", rName+"_gpssd2_volume"),
					resource.TestCheckResourceAttr(resourceName2, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName2, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName2, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName2, "throughput", "500"),
				),
			},
			{
				Config: testAccEvsVolume_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttr(resourceName1, "auto_renew", "true"),
					resource.TestCheckResourceAttr(resourceName2, "auto_renew", "true"),
				),
			},
		},
	})
}

func TestAccEvsVolume_withServerId(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_evs_volume.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_serverId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "attachment.0.instance_id", "huaweicloud_compute_instance.test", "id")),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade", "server_id", "charging_mode", "period", "period_unit",
				},
			},
		},
	})
}

func testAccEvsVolume_base() string {
	return `
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    device_type = string
    volume_type = string
    multiattach = bool
    iops        = number
    throughput  = number
  }))
  default = [
    {
      suffix = "vbd_normal_volume",
      device_type = "VBD",
      volume_type = "SSD",
      multiattach = false,
      iops = 0,
      throughput = 0
    },
    {
      suffix = "vbd_share_volume",
      device_type = "VBD",
      volume_type = "SSD",
      multiattach = true,
      iops = 0,
      throughput = 0
    },
    {
      suffix = "scsi_normal_volume",
      device_type = "SCSI",
      volume_type = "SSD",
      multiattach = false,
      iops = 0,
      throughput = 0
    },
    {
      suffix = "scsi_share_volume",
      device_type = "SCSI",
      volume_type = "SSD",
      multiattach = true,
      iops = 0,
      throughput = 0
    },
    {
      suffix = "gpssd2_normal_volume",
      device_type = "SCSI",
      volume_type = "GPSSD2",
      multiattach = false,
      iops = 3000,
      throughput = 500
    },
    {
      suffix = "essd2_normal_volume",
      device_type = "SCSI",
      volume_type = "ESSD2",
      multiattach = false,
      iops = 3000,
      throughput = 0
    },
  ]
}

data "huaweicloud_availability_zones" "test" {}
`
}

func testAccEvsVolume_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}"
  size              = 100
  description       = "Created by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach
  iops              = var.volume_configuration[count.index].iops
  throughput        = var.volume_configuration[count.index].throughput

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsVolume_base(), rName)
}

func testAccEvsVolume_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%s_${var.volume_configuration[count.index].suffix}_update"
  size              = 200
  description       = "Updated by acc test script."
  volume_type       = var.volume_configuration[count.index].volume_type
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach
  iops              = var.volume_configuration[count.index].iops
  throughput        = var.volume_configuration[count.index].throughput

  tags = {
    foo1 = "bar"
    key  = "value1"
  }
}
`, testAccEvsVolume_base(), rName)
}

func testAccEvsVolume_epsId(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name                  = "%s"
  description           = "test volume for epsID"
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  volume_type           = "SSD"
  size                  = 100
  enterprise_project_id = "%s"
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prepaid_base() string {
	return `
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    volume_type = string
    iops        = number
    throughput  = number
  }))
  default = [
    {
      suffix = "ssd_volume",
      volume_type = "SSD",
      iops = 0,
      throughput = 0
    },
    {
      suffix = "gpssd2_volume",
      volume_type = "GPSSD2",
      iops = 3000,
      throughput = 500
    },
  ]
}

data "huaweicloud_availability_zones" "test" {}
`
}

func testAccEvsVolume_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  name              = "%s_${var.volume_configuration[count.index].suffix}"
  description       = "test volume for charging mode"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  size              = 100
  volume_type       = var.volume_configuration[count.index].volume_type
  iops              = var.volume_configuration[count.index].iops
  throughput        = var.volume_configuration[count.index].throughput

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%[3]v"
}
`, testAccEvsVolume_prepaid_base(), rName, isAutoRenew)
}

const testAccCompute_data = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name                = "%[2]s"
  description         = "terraform test"
  hostname            = "hostname-test"
  image_id            = data.huaweicloud_images_image.test.id
  flavor_id           = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids  = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  stop_before_destroy = true
  agency_name         = "test111"
  agent_list          = "hss"

  network {
    uuid              = data.huaweicloud_vpc_subnet.test.id
    source_dest_check = false
  }

  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  metadata = {
    foo = "bar"
    key = "value"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccEvsVolume_serverId(rName string) string {
	return fmt.Sprintf(`
%[1]s
	
resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  volume_type       = "GPSSD"
  description       = "test volume for charging mode"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  server_id         = huaweicloud_compute_instance.test.id
  size              = 100
  charging_mode     = "postPaid"
}
`, testAccComputeInstance_basic(rName), rName)
}
