package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCloudVendorUserResourcesSync_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCloudVendorUserResourcesSync_basic(),
			},
		},
	})
}

func testCloudVendorUserResourcesSync_basic() string {
	return `
resource "huaweicloud_coc_cloud_vendor_user_resources_sync" "test" {
  vendor = "HCS"
}
`
}
