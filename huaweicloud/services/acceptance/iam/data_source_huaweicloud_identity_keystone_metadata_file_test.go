package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityKeystoneMetadataFile_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_identity_keystone_metadata_file.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testIdentityKeystoneMetadataFile,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "metadata_file"),
				),
			},
			{
				Config: testIdentityKeystoneMetadataFileWithUnsigned,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "metadata_file"),
				),
			},
		},
	})
}

const testIdentityKeystoneMetadataFile = `
data "huaweicloud_identity_keystone_metadata_file" "test" {}
`
const testIdentityKeystoneMetadataFileWithUnsigned = `
data "huaweicloud_identity_keystone_metadata_file" "test" {
  unsigned = true
}
`
