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

func TestAccEvsVolume_GPSSD2(t *testing.T) {
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
			// Type GPSSD2 is only supported in part availability_zones under the certain region.
			acceptance.TestAccPreCheckAvailabilityZoneGPSSD2(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_GPSSD2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", acceptance.HW_EVS_AVAILABILITY_ZONE_GPSSD2),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
				),
			},
			{
				Config: testAccEvsVolume_GPSSD2_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "150"),
				),
			},
		},
	})
}

func TestAccEvsVolume_ESSD2(t *testing.T) {
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
			// Type ESSD2 is only supported in part availability_zones under the certain region.
			acceptance.TestAccPreCheckAvailabilityZoneESSD2(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_ESSD2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "availability_zone", acceptance.HW_EVS_AVAILABILITY_ZONE_ESSD2),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "ESSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
				),
			},
			{
				Config: testAccEvsVolume_ESSD2_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "ESSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
				),
			},
		},
	})
}

func TestAccEvsVolume_prePaid_withoutServerId(t *testing.T) {
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
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaid_withoutServerId(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description",
						"test volume when charging_mode is prePaid and the volume is not attached"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
				),
			},
			{
				Config: testAccEvsVolume_prePaid_withoutServerId(rName, false),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
		},
	})
}

func testAccEvsVolume_prePaid_withoutServerId(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume when charging_mode is prePaid and the volume is not attached"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  size              = 100
  volume_type       = "SSD"

  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
  auto_renew        = "%[2]v"
}`, rName, isAutoRenew)
}

func TestAccEvsVolume_prePaid_withServerId(t *testing.T) {
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
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaid_withServerId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description",
						"test volume when charging_mode is prePaid and the volume is attached"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SSD"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
			{
				Config: testAccEvsVolume_prePaid_withServerId_updateSize(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
				),
			},
		},
	})
}

func testAccEvsVolume_prePaid_withServerId(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  description       = "test volume when charging_mode is prePaid and the volume is attached"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  size              = 100
  volume_type       = "SSD"
  server_id         = huaweicloud_compute_instance.test.id

  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}`, testAccComputeInstance_prePaid(rName), rName)
}

func testAccEvsVolume_prePaid_withServerId_updateSize(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  description       = "test volume when charging_mode is prePaid and the volume is attached"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  size              = 200
  volume_type       = "SSD"
  server_id         = huaweicloud_compute_instance.test.id

  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}`, testAccComputeInstance_prePaid(rName), rName)
}

func testAccComputeInstance_prePaid(rName string) string {
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

  charging_mode     = "prePaid"
  period_unit       = "month"
  period            = 1
}
`, testAccCompute_data, rName)
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

func testAccEvsVolume_GPSSD2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume for updating QoS when volume_type is GPSSD2"
  availability_zone = "%[2]s"
  size              = 100
  volume_type       = "GPSSD2"
  iops              = 3000
  throughput        = 125
}
`, rName, acceptance.HW_EVS_AVAILABILITY_ZONE_GPSSD2)
}

func testAccEvsVolume_GPSSD2_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume for updating QoS when volume_type is GPSSD2"
  availability_zone = "%[2]s"
  size              = 100
  volume_type       = "GPSSD2"
  iops              = 4000
  throughput        = 150
}
`, rName, acceptance.HW_EVS_AVAILABILITY_ZONE_GPSSD2)
}

func testAccEvsVolume_ESSD2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume for updating QoS when volume_type is ESSD2"
  availability_zone = "%[2]s"
  size              = 100
  volume_type       = "ESSD2"
  iops              = 3000
}
`, rName, acceptance.HW_EVS_AVAILABILITY_ZONE_ESSD2)
}

func testAccEvsVolume_ESSD2_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume for updating QoS when volume_type is ESSD2"
  availability_zone = "%[2]s"
  size              = 100
  volume_type       = "ESSD2"
  iops              = 4000
}
`, rName, acceptance.HW_EVS_AVAILABILITY_ZONE_ESSD2)
}
