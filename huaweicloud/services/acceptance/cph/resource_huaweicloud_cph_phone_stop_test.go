package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCphPhoneStop_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCphPhoneStop_basic(name),
			},
		},
	})
}

func testCphPhoneStop_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cph_phones" "test" {
  server_id = huaweicloud_cph_server.test.id
}

resource "huaweicloud_cph_phone_stop" "test" {
  phone_id = data.huaweicloud_cph_phones.test.phones[0].phone_id
}
`, testCphServer_basic(name))
}
