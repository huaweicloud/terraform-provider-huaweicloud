package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCphPhoneProperty_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCphPhoneProperty_basic(name),
			},
		},
	})
}

func testCphPhoneProperty_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cph_phones" "test" {
  server_id = huaweicloud_cph_server.test.id
}

resource "huaweicloud_cph_phone_property" "test" {
  phones {
    phone_id = data.huaweicloud_cph_phones.test.phones[0].phone_id
    property = jsonencode({
      "com.cph.mainkeys":0,
      "disable.status.bar":0,
      "ro.permission.changed":0,
      "ro.horizontal.screen":0,
      "ro.install.auto":0,
      "ro.com.cph.sfs_enable":0,
      "ro.product.manufacturer":"Huawei",
      "ro.product.name":"monbox",
      "ro.com.cph.notification_disable":0
    })
  }
}
`, testCphServer_build(name))
}

func testCphServer_build(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cph_server" "test" {
  name          = "%s"
  server_flavor = "physical.rx1.xlarge"
  phone_flavor  = "rx1.cp.c15.d46.e1v1"
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
