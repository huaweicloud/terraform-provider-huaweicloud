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

func getV3VolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "evs"
	)

	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating EVS client: %s", err)
	}

	return evs.GetV3VolumeDetail(client, state.Primary.ID)
}

func TestAccV3Volume_basic(t *testing.T) {
	var (
		volume       interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_evsv3_volume.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getV3VolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Volume_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrPair(resourceName, "image_id",
						"data.huaweicloud_images_image.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "links.#"),
					resource.TestCheckResourceAttrSet(resourceName, "attachments.#"),
					resource.TestCheckResourceAttrSet(resourceName, "bootable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_image_metadata.%"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccV3Volume_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				Config: testAccV3Volume_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"volume_type",
					"disaster_recovery_azs",
					"dedicated_storage_id",
					"cascade",
				},
			},
		},
	})
}

func testAccV3Volume_base() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
`
}

func testAccV3Volume_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = "GPSSD"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  description       = "test description"
  image_id          = data.huaweicloud_images_image.test.id

  metadata = {
    volume_owner = "tf test"
  }

  multiattach = false
  name        = "%[2]s"
  size        = 100
  cascade     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccV3Volume_base(), name)
}

func testAccV3Volume_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = "GPSSD"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  description       = "test description update"
  image_id          = data.huaweicloud_images_image.test.id

  metadata = {
    volume_owner = "tf test update"
    data_test    = "value_test"
  }

  multiattach = false
  name        = "%[2]s_update"
  size        = 200
  cascade     = true

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, testAccV3Volume_base(), name)
}

func testAccV3Volume_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = "GPSSD"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  description       = ""
  image_id          = data.huaweicloud_images_image.test.id

  metadata = {
    volume_owner = "tf test update"
    data_test    = "value_test"
  }

  multiattach = false
  name        = "%[2]s_update"
  size        = 200
  cascade     = true

  tags = {}
}
`, testAccV3Volume_base(), name)
}

func TestAccV3Volume_volumeType(t *testing.T) {
	var (
		volume       interface{}
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_evsv3_volume.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&volume,
		getV3VolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV3Volume_volumeType_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "GPSSD2"),
					resource.TestCheckResourceAttrPair(resourceName, "availability_zone",
						"data.huaweicloud_availability_zones.test", "names.0"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttrPair(resourceName, "image_id",
						"data.huaweicloud_images_image.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "iops", "3000"),
					resource.TestCheckResourceAttr(resourceName, "throughput", "125"),
					resource.TestCheckResourceAttr(resourceName, "cascade", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "links.#"),
					resource.TestCheckResourceAttrSet(resourceName, "attachments.#"),
					resource.TestCheckResourceAttrSet(resourceName, "bootable"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_image_metadata.%"),
					resource.TestCheckResourceAttrSet(resourceName, "iops_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "throughput_attribute.#"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccV3Volume_volumeType_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"volume_type",
					"disaster_recovery_azs",
					"dedicated_storage_id",
					"cascade",
				},
			},
		},
	})
}

func testAccV3Volume_volumeType_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = "GPSSD2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  description       = "test description"
  image_id          = data.huaweicloud_images_image.test.id

  metadata = {
    volume_owner = "tf test"
  }

  multiattach = false
  name        = "%[2]s"
  size        = 100
  iops        = 3000
  throughput  = 125
  cascade     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccV3Volume_base(), name)
}

func testAccV3Volume_volumeType_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_evsv3_volume" "test" {
  volume_type       = "GPSSD2"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  description       = "test description update"
  image_id          = data.huaweicloud_images_image.test.id

  metadata = {
    volume_owner = "tf test update"
    data_test    = "value_test"
  }

  multiattach = false
  name        = "%[2]s_update"
  size        = 200
  iops        = 3000
  throughput  = 125
  cascade     = true

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, testAccV3Volume_base(), name)
}
