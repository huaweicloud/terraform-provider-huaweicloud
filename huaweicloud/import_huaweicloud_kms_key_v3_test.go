package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
)

func TestAccKmsV3Key_importBasic(t *testing.T) {
	resourceName := "huaweicloud_kms_key_v3.key_1"
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV3KeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsV3Key_basic(keyAlias),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
