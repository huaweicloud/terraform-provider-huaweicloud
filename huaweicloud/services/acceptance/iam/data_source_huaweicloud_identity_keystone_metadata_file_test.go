package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataKeystoneMetadataFile_basic(t *testing.T) {
	var (
		check = "data.huaweicloud_identity_keystone_metadata_file.test"
		dc    = acceptance.InitDataSourceCheck(check)

		byUnsigned   = "data.huaweicloud_identity_keystone_metadata_file.test_with_unsigned"
		dcByUnsigned = acceptance.InitDataSourceCheck(byUnsigned)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataKeystoneMetadataFile_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(check, "metadata_file"),
					dcByUnsigned.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byUnsigned, "metadata_file"),
				),
			},
		},
	})
}

const testAccDataKeystoneMetadataFile_basic = `
data "huaweicloud_identity_keystone_metadata_file" "test" {}

data "huaweicloud_identity_keystone_metadata_file" "test_with_unsigned" {
  unsigned = true
}
`
