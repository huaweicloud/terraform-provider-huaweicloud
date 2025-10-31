package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceIdentityKeystoneMetadataFile_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_keystone_metadata_file.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestDataSourceIdentityKeystoneMetadataFile,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata_file"),
				),
			},
			{
				Config: testTestDataSourceIdentityKeystoneMetadataFileWithUnsigned,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata_file"),
				),
			},
		},
	})
}

const testTestDataSourceIdentityKeystoneMetadataFile = `
data "huaweicloud_identity_keystone_metadata_file" "test" {}
`
const testTestDataSourceIdentityKeystoneMetadataFileWithUnsigned = `
data "huaweicloud_identity_keystone_metadata_file" "test" {
  unsigned = true
}
`
