package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityDomainQuotaDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_domain_quota.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityDomainQuotaDataSource,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "11"),
				),
			},
			{
				Config: testAccIdentityDomainQuotaDataSourceWithType,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "resources.0.type", "user"),
				),
			},
		},
	})
}

const testAccIdentityDomainQuotaDataSource = `
data "huaweicloud_identity_domain_quota" "test" {}
`

const testAccIdentityDomainQuotaDataSourceWithType = `
data "huaweicloud_identity_domain_quota" "test" {
  type = "user"
}
`
