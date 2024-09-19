package ims

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccImageExport_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource will export an image file to the OBS bucket,
			// so here we need to set a URL in the format **OBS bucket name:image file name**.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccImageExport_basic(),
			},
		},
	})
}

func TestAccImageExport_with_isQuickExport(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy
	// method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This resource will export an image file to the OBS bucket,
			// so here we need to set a URL in the format **OBS bucket name:image file name**.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// One-time action resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccImageExport_with_isQuickExport(),
			},
		},
	})
}

func testAccImageExport_basic() string {
	rName := acceptance.RandomAccResourceName()
	rNameSuffix := acctest.RandString(4)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_image_export" "test" {
  image_id    = huaweicloud_ims_ecs_system_image.test.id
  bucket_url  = "%[2]s_%[3]s"
  file_format = "qcow2"
}
`, testAccEcsSystemImage_basic(rName), acceptance.HW_IMS_IMAGE_URL, rNameSuffix)
}

func testAccImageExport_with_isQuickExport() string {
	rName := acceptance.RandomAccResourceName()
	rNameSuffix := acctest.RandString(4)

	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_image_export" "test" {
  image_id        = huaweicloud_ims_ecs_system_image.test.id
  bucket_url      = "%[2]s_%[3]s"
  is_quick_export = true
}
`, testAccEcsSystemImage_basic(rName), acceptance.HW_IMS_IMAGE_URL, rNameSuffix)
}
