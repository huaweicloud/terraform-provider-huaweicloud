package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKmsV1Key_importBasic(t *testing.T) {
	resourceName := "huaweicloud_kms_key_v1.key_2"
	var keyAlias = fmt.Sprintf("kms_%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKmsV1KeyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccKmsV1Key_basic(keyAlias),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"key_usage",
					"pending_days",
				},
			},
		},
	})
}
