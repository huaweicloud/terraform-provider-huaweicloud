package evs

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/evs/v2/cloudvolumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getVolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.BlockStorageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud block storage v2 client: %s", err)
	}
	return cloudvolumes.Get(c, state.Primary.ID).Extract()
}

func TestAccEvsVolume_basic(t *testing.T) {
	var volume cloudvolumes.Volume
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_evs_volume.test"

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
					rc.CheckMultiResourcesExists(4),
					// Common configuration
					resource.TestCheckResourceAttrPair("huaweicloud_evs_volume.test.0", "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "description", "Created by acc test script."),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "volume_type", "SSD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "size", "100"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "tags.foo", "bar"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "tags.key", "value"),
					// Personalized configuration
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "name", rName+"_vbd_normal_volume"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "device_type", "VBD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "multiattach", "false"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "name", rName+"_vbd_share_volume"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "device_type", "VBD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "multiattach", "true"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "name", rName+"_scsi_normal_volume"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "device_type", "SCSI"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "multiattach", "false"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "name", rName+"_scsi_share_volume"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "device_type", "SCSI"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "multiattach", "true"),
				),
			},
			{
				Config: testAccEvsVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(4),
					// Common configuration
					resource.TestCheckResourceAttrPair("huaweicloud_evs_volume.test.0", "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "description", "Updated by acc test script."),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "volume_type", "SSD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "size", "200"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "tags.foo1", "bar"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "tags.key", "value1"),
					// Personalized configuration
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "name", rName+"_vbd_normal_volume_update"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "device_type", "VBD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.0", "multiattach", "false"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "name", rName+"_vbd_share_volume_update"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "device_type", "VBD"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.1", "multiattach", "true"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "name", rName+"_scsi_normal_volume_update"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "device_type", "SCSI"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.2", "multiattach", "false"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "name", rName+"_scsi_share_volume_update"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "device_type", "SCSI"),
					resource.TestCheckResourceAttr("huaweicloud_evs_volume.test.3", "multiattach", "true"),
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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
				),
			},
			{
				Config: testAccEvsVolume_prePaid(rName, true),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "true"),
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
				},
			},
		},
	})
}

func testAccEvsVolume_base() string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    suffix      = string
    device_type = string
    multiattach = bool
  }))
  default = [
    {suffix = "vbd_normal_volume",  device_type = "VBD",  multiattach = false},
    {suffix = "vbd_share_volume",   device_type = "VBD",  multiattach = true},
    {suffix = "scsi_normal_volume", device_type = "SCSI", multiattach = false},
    {suffix = "scsi_share_volume",  device_type = "SCSI", multiattach = true},
  ]
}

data "huaweicloud_availability_zones" "test" {}
`)
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
  volume_type       = "SSD"
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach

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
  volume_type       = "SSD"
  device_type       = var.volume_configuration[count.index].device_type
  multiattach       = var.volume_configuration[count.index].multiattach

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

func testAccEvsVolume_prePaid(rName string, isAutoRenew bool) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_evs_volume" "test" {
  name              = "%[1]s"
  description       = "test volume for charging mode"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SSD"
  size              = 100

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = "%v"
}
`, rName, isAutoRenew)
}
