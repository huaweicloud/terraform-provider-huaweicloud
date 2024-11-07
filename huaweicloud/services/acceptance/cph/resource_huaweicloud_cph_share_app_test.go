package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccShareApp_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckCphObsBucketName(t)
			acceptance.TestAccPrecheckCphAdbObjectPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testShareApp_basic(name),
			},
		},
	})
}

func testShareApp_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cph_share_app" "test" {
  server_id       = huaweicloud_cph_server.test.id
  package_name    = "com.cph.config"
  bucket_name     = %[2]s
  object_path     = %[3]s
  pre_install_app = 0
}
`, testCphServer_supportShareApp(name), acceptance.HW_CPH_OBS_BUCKET_NAME, acceptance.HW_CPH_OBS_OBJECT_PATH)
}

func testCphServer_supportShareApp(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cph_server" "test" {
  name          = "%s"
  server_flavor = "physical.kg1.4xlarge.cp"
  phone_flavor  = "rs2.max"
  image_id      = data.huaweicloud_cph_phone_images.test.images[0].id
  keypair_name  = huaweicloud_kps_keypair.test.name

  vpc_id    = huaweicloud_vpc.test.id
  subnet_id = huaweicloud_vpc_subnet.test.id
  eip_type  = "5_bgp"

  bandwidth {
    share_type  = "0"
    charge_mode = "1"
    size        = 300
  }

  period_unit = "month"
  period      = 1
  auto_renew  = "true"

  lifecycle {
    ignore_changes = [
      image_id, auto_renew, period, period_unit,
    ]
  }
}
`, testCphServerBase(name), name)
}
