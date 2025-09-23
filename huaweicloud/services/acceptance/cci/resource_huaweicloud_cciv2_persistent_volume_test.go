package cci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cci"
)

func getPersistentVolumeResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetPersistentVolume(c, state.Primary.ID)
}

func TestAccV2PersistentVolume_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_persistent_volume.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPersistentVolumeResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2PersistentVolume_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "access_modes.0", "ReadWriteMany"),
					resource.TestCheckResourceAttr(resourceName, "csi.0.driver", "a-a"),
					resource.TestCheckResourceAttr(resourceName, "csi.0.volume_handle", "aaa"),
					resource.TestCheckResourceAttr(resourceName, "capacity.storage", "2Gi"),
				),
			},
			{
				Config: testAccV2PersistentVolume_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "capacity.storage", "4Gi"),
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

func testAccV2PersistentVolume_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cciv2_persistent_volume" test {
  name         = "%[1]s"
  access_modes = ["ReadWriteMany"]

  capacity = {
    storage = "2Gi"
  }

  csi {
    driver        = "a-a"
    volume_handle = "aaa"
  }
}
`, rName)
}

func testAccV2PersistentVolume_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_cciv2_persistent_volume" test {
  name         = "%[1]s"
  access_modes = ["ReadWriteMany"]

  capacity = {
    storage = "4Gi"
  }

  csi {
    driver        = "a-a"
    volume_handle = "aaa"
  }
}
`, rName)
}
