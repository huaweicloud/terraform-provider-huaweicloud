package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestResourceSwrV2NamespacesCreate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSwrDestory,
		Steps: []resource.TestStep{
			{
				Config: testSwrConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("huaweicloud_swr_v2_namespace.test", "namespace", "taobao"),
				),
			},
		},
	})
}

var testSwrConfig = `
resource "huaweicloud_swr_v2_namespace" "test" {
  	namespace	  = "taobao"
}
`

func testAccCheckSwrDestory(s *terraform.State) error {
	return nil
}
