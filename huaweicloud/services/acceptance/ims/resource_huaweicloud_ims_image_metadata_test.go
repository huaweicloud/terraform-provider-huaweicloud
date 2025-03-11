package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImageMetadata_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccImageMetadata_basic(),
			},
		},
	})
}

func testAccImageMetadata_basic() string {
	rName := acceptance.RandomAccResourceName()

	return fmt.Sprintf(`
resource "huaweicloud_ims_image_metadata" "test" {
  __os_version     = "Ubuntu 20.04 server 64bit"
  visibility       = "private"
  name1            = "%s"
  protected        = false
  container_format = "bare"
  disk_format      = "vhd"

  tags = [
    "test=testvalue",
    "image=imagevalue"
  ]

  min_ram  = 1024
  min_disk = 80
}
`, rName)
}
