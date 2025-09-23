package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ims/v2/cloudimages"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/deprecated"
)

func getImageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	imsClient, err := cfg.ImageV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS v2 client: %s", err)
	}

	img, err := deprecated.GetCloudImage(imsClient, state.Primary.ID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving image: %s", err)
	}

	return &img, nil
}

func TestAccImsImage_basic(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_images_image.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 1024, 4096),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_update(rName, rNameUpdate, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "description",
						"created by Terraform AccTest for update"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
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

func TestAccImsImage_withEpsId(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		resourceName = "huaweicloud_images_image.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccImsImage_wholeImage_withServer(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_images_image.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_wholeImage_withServer(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 1024, 4096),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 4096, 8192),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "8192"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withServer_update(rName, rNameUpdate, 0, 0),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vault_id"},
			},
		},
	})
}

func TestAccImsImage_wholeImage_withBackup(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_images_image.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckImsBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_wholeImage_withBackup(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccImsImage_wholeImage_withBackup_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func TestAccImsImage_dataImage_withVolumeId(t *testing.T) {
	var (
		image        cloudimages.Image
		rName        = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_images_image.image_test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccImsImage_dataImage_withVolumeId_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "description test"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "huaweicloud_evs_volume.test", "id"),
				),
			},
			{
				Config: testAccImsImage_dataImage_withVolumeId_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", "description update"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "min_ram", "1024"),
					resource.TestCheckResourceAttr(resourceName, "max_ram", "4096"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
			},
		},
	})
}

func testAccImsImage_basic(rName string) string {
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

resource "huaweicloud_images_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccImsImage_update(rName, rNameUpdate string, minRAM, maxRAM int) string {
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

resource "huaweicloud_images_image" "test" {
  name        = "%[3]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest for update"
  min_ram     = %[4]d
  max_ram     = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, common.TestBaseNetwork(rName), rName, rNameUpdate, minRAM, maxRAM)
}

func testAccImsImage_withEpsId(rName string) string {
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

resource "huaweicloud_images_image" "test" {
  name                  = "%[2]s"
  instance_id           = huaweicloud_compute_instance.test.id
  description           = "created by Terraform AccTest"
  enterprise_project_id = "%[3]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccImsImage_wholeImage_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
  image_name         = "Ubuntu 18.04 server 64bit"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  data_disks {
    type = "SAS"
    size = "10"
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[1]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccImsImage_wholeImage_withServer(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_images_image" "test" {
  name        = "%[3]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"
  vault_id    = huaweicloud_cbr_vault.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), testAccImsImage_wholeImage_base(rName), rName)
}

func testAccImsImage_wholeImage_withServer_update(rName, updateName string, minRAM, maxRAM int) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_images_image" "test" {
  name        = "%[3]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "created by Terraform AccTest"
  vault_id    = huaweicloud_cbr_vault.test.id
  min_ram     = %[4]d
  max_ram     = %[5]d

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, common.TestBaseNetwork(rName), testAccImsImage_wholeImage_base(rName), updateName, minRAM, maxRAM)
}

func testAccImsImage_wholeImage_withBackup(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_images_image" "test" {
  name        = "%[1]s"
  backup_id   = "%[2]s"
  description = "created by Terraform AccTest"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName, acceptance.HW_IMS_BACKUP_ID)
}

func testAccImsImage_wholeImage_withBackup_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_images_image" "test" {
  name        = "%[1]s"
  backup_id   = "%[2]s"
  description = "created by Terraform AccTest"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, rName, acceptance.HW_IMS_BACKUP_ID)
}

func testAccImsImage_dataImage_withVolumeId_base(rName string) string {
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
  size              = 60
  charging_mode     = "postPaid"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccImsImage_dataImage_withVolumeId_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image" "image_test" {
  name        = "%[2]s"
  volume_id   = huaweicloud_evs_volume.test.id
  description = "description test"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccImsImage_dataImage_withVolumeId_base(rName), rName)
}

func testAccImsImage_dataImage_withVolumeId_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_images_image" "image_test" {
  name        = "%[2]s"
  volume_id   = huaweicloud_evs_volume.test.id
  description = "description update"
  min_ram     = 1024
  max_ram     = 4096

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccImsImage_dataImage_withVolumeId_base(rName), rName)
}
