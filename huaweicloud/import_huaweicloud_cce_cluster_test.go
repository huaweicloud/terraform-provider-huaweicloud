package huaweicloud

import (
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccCCEClusterV3_importBasic(t *testing.T) {
	resourceName := "huaweicloud_cce_cluster_v3.cluster_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCE(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV3Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCEClusterV3_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
