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

func getPersistentVolumeClaimResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CciV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCI client: %s", err)
	}
	return cci.GetV2PersistentVolumeClaimDetail(c, state.Primary.Attributes["namespace"], state.Primary.Attributes["name"])
}

func TestAccV2PersistentVolumeClaim_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_cciv2_persistent_volume_claim.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPersistentVolumeClaimResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV2PersistentVolumeClaim_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "namespace", "huaweicloud_cci_namespace.test", "name"),
					resource.TestCheckResourceAttrPair(resourceName, "storage_class_name",
						"huaweicloud_cciv2_persistent_volume.test", "name"),
					resource.TestCheckResourceAttr(resourceName, "access_modes.0", "ReadWriteMany"),
					resource.TestCheckResourceAttr(resourceName, "volume_mode", "Filesystem"),
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

func testAccV2PersistentVolumeClaim_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_cciv2_persistent_volume_claim" test {
  name               = "%[3]s"
  namespace          = huaweicloud_cciv2_namespace.test.name
  access_modes       = ["ReadWriteMany"]
  storage_class_name = huaweicloud_cciv2_persistent_volume.test.name
  volume_mode        = "Filesystem"

  annotations = {
    "everest.io/obs-volume-type"       = "STANDARD"
    "csi.storage.k8s.io/fstype"        = "s3fs"
    "everest.io/enterprise-project-id" = "0"
  }

  resources {
    requests = {
      storage = "1Gi"
    }
  }
}
`, testAccV2Namespace_basic(rName), testAccV2PersistentVolume_basic(rName), rName)
}
