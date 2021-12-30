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
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by acc test script."),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccEvsVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by acc test script."),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
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

func testAccEvsVolume_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  description       = "Created by acc test script."
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 100
  image_id          = data.huaweicloud_images_image.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, rName)
}

func testAccEvsVolume_update(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
  
resource "huaweicloud_evs_volume" "test" {
  name              = "%s_update"
  description       = "Updated by acc test script."
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 200
  image_id          = data.huaweicloud_images_image.test.id

  tags = {
    foo1 = "bar"
    key  = "value1"
  }
}
`, rName)
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
