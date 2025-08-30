package iam

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"testing"
)

func TestAccIdentityAccessKeyDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_key.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentityAccessKeyDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "credentials.#"),
				),
			},
		},
	})
}

const AccIdentityAccessKeyDataSource_basic = `
data "huaweicloud_identity_key" "test" {}`
