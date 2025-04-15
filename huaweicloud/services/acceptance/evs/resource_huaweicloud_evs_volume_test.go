package evs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/evs"
)

func getVolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetVolumeDetail(client, state.Primary.ID)
}

// This test case is used to test the `GPSSD2` type cloud hard disk, which can be used in cn-north-4.
func TestAccEvsVolume_postPaidWithoutServer(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
				Config: testAccEvsVolume_postPaidWithoutServer(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "all_metadata.%"),
					resource.TestCheckResourceAttrSet(resourceName, "bootable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.0.frozened"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.0.total_val"),
					resource.TestCheckResourceAttrSet(resourceName, "links.#"),
					resource.TestCheckResourceAttrSet(resourceName, "links.0.href"),
					resource.TestCheckResourceAttrSet(resourceName, "links.0.rel"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "service_type"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.0.frozened"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.0.total_val"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidWithoutServer_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "150"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidWithoutServer_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_type", "SCSI"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "150"),

					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
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

func testAccEvsVolume_postPaidWithoutServer(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s"
  size                  = 100
  description           = "test description"
  volume_type           = "GPSSD2"
  iops                  = 3000
  throughput            = 125
  device_type           = "SCSI"
  multiattach           = false
  enterprise_project_id = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidWithoutServer_update1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD2"
  iops                  = 4000
  throughput            = 150
  device_type           = "SCSI"
  multiattach           = false
  enterprise_project_id = "%s"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidWithoutServer_update2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  volume_type           = "GPSSD2"
  iops                  = 4000
  throughput            = 150
  device_type           = "SCSI"
  multiattach           = false
  enterprise_project_id = "%s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// This test case is used to test the `GPSSD` type cloud hard disk, which can be used in cn-north-4.
func TestAccEvsVolume_postPaidWithServer(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
				Config: testAccEvsVolume_postPaidWithServer(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.#"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.attached_at"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.device"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.id"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.instance_id"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.attached_volume_id"),
					resource.TestCheckResourceAttrSet(resourceName, "attachment.0.volume_id"),
					resource.TestCheckResourceAttrSet(resourceName, "bootable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "links.#"),
					resource.TestCheckResourceAttrSet(resourceName, "links.0.href"),
					resource.TestCheckResourceAttrSet(resourceName, "links.0.rel"),
					resource.TestCheckResourceAttrSet(resourceName, "all_metadata.%"),
					resource.TestCheckResourceAttrSet(resourceName, "service_type"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidWithServer_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidWithServer_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
					"charging_mode",
					"server_id",
				},
			},
		},
	})
}

func testAccEvsVolume_postPaidWithServer_base(name string) string {
	return fmt.Sprintf(`
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

resource "huaweicloud_compute_instance" "test" {
  name                = "%s"
  description         = "terraform test"
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
}
`, name)
}

func testAccEvsVolume_postPaidWithServer(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s"
  size                  = 100
  description           = "test description"
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "postPaid"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsVolume_postPaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidWithServer_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "postPaid"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, testAccEvsVolume_postPaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidWithServer_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s_update"
  size                  = 200
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "postPaid"
}
`, testAccEvsVolume_postPaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// This test case is used to test the `GPSSD2` type cloud hard disk, which can be used in cn-north-4.
func TestAccEvsVolume_prePaidWithoutServer(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaidWithoutServer(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidWithoutServer_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "150"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidWithoutServer_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "4000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "150"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
					"period_unit",
					"period",
					"server_id",
					"auto_renew",
					"charging_mode",
				},
			},
		},
	})
}

func testAccEvsVolume_prePaidWithoutServer(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s"
  size                  = 100
  description           = "test description"
  volume_type           = "GPSSD2"
  iops                  = 3000
  throughput            = 125
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  auto_renew            = "true"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidWithoutServer_update1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD2"
  iops                  = 4000
  throughput            = 150
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  auto_renew            = "false"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidWithoutServer_update2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  volume_type           = "GPSSD2"
  iops                  = 4000
  throughput            = 150
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
  auto_renew            = "true"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// This test case is used to test the `GPSSD` type cloud hard disk, which can be used in cn-north-4.
func TestAccEvsVolume_prePaidWithServer(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaidWithServer(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidWithServer_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidWithServer_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttrPair(resourceName, "server_id",
						"huaweicloud_compute_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "period_unit", "month"),
					resource.TestCheckResourceAttr(resourceName, "period", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
					"charging_mode",
					"server_id",
					"period_unit",
					"period",
				},
			},
		},
	})
}

func testAccEvsVolume_prePaidWithServer_base(name string) string {
	return fmt.Sprintf(`
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

resource "huaweicloud_compute_instance" "test" {
  name                = "%s"
  description         = "terraform test"
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

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
}
`, name)
}

func testAccEvsVolume_prePaidWithServer(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s"
  size                  = 100
  description           = "test description"
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccEvsVolume_prePaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidWithServer_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, testAccEvsVolume_prePaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidWithServer_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%[2]s_update"
  size                  = 200
  volume_type           = "GPSSD"
  multiattach           = false
  enterprise_project_id = "%[3]s"
  server_id             = huaweicloud_compute_instance.test.id
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
}
`, testAccEvsVolume_prePaidWithServer_base(name), name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Changing the disk type takes a long time. This test case may take several hours, so a separate test case is provided.
// Before executing this test case, please submit a work order to apply for public beta qualification to change the disk type.
func TestAccEvsVolume_postPaidEditDiskType(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
			acceptance.TestAccPreCheckEVSFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_postPaidEditDiskType(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidEditDiskType_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccEvsVolume_postPaidEditDiskType_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
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

func testAccEvsVolume_postPaidEditDiskType(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s"
  size                  = 100
  description           = "test description"
  volume_type           = "SAS"
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidEditDiskType_update1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD2"
  iops                  = 3000
  throughput            = 125
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_postPaidEditDiskType_update2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  volume_type           = "GPSSD"
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// Changing the disk type takes a long time. This test case may take several hours, so a separate test case is provided.
// Before executing this test case, please submit a work order to apply for public beta qualification to change the disk type.
func TestAccEvsVolume_prePaidEditDiskType(t *testing.T) {
	var volume interface{}
	name := acceptance.RandomAccResourceName()
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
			acceptance.TestAccPreCheckEVSFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccEvsVolume_prePaidEditDiskType(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidEditDiskType_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccEvsVolume_prePaidEditDiskType_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "device_type", "VBD"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SSD"),
					resource.TestCheckResourceAttr(resourceName, "iops", "0"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "0"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttrSet(resourceName, "wwn"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"cascade",
					"period_unit",
					"period",
					"auto_renew",
					"charging_mode",
				},
			},
		},
	})
}

func testAccEvsVolume_prePaidEditDiskType(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s"
  size                  = 100
  description           = "test description"
  volume_type           = "GPSSD"
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidEditDiskType_update1(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  description           = "test description update"
  volume_type           = "GPSSD2"
  iops                  = 3000
  throughput            = 125
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccEvsVolume_prePaidEditDiskType_update2(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  name                  = "%s_update"
  size                  = 200
  volume_type           = "SSD"
  device_type           = "VBD"
  multiattach           = false
  enterprise_project_id = "%s"
  charging_mode         = "prePaid"
  period_unit           = "month"
  period                = 1
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
