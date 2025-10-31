package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityAccessKeyDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_identity_access_key.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: AccIdentityAccessKeyDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credentials.0.status", "active"),
				),
			},
			{
				Config: AccIdentityAccessKeyDataSourceByAccessKey_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credentials.0.status", "active"),
				),
			},
			{
				Config: testAccIdentityAccessKeyDataSourceWithoutParams_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "credentials.#", "1"),
					resource.TestCheckResourceAttr(dataSource, "credentials.0.status", "active"),
				),
			},
		},
	})
}

const AccIdentityAccessKeyDataSource_basic = `
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  depends_on = [huaweicloud_identity_access_key.key_1]

  user_id = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  user_id = huaweicloud_identity_access_key.key_1.user_id
}
`

const AccIdentityAccessKeyDataSourceByAccessKey_basic = `
resource "huaweicloud_identity_user" "test_user" {
  name        = "my_test_huaweicloud"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}
resource "huaweicloud_identity_access_key" "key_1" {
  user_id = huaweicloud_identity_user.test_user.id
}
data "huaweicloud_identity_access_key" "test" {
  access_key = huaweicloud_identity_access_key.key_1.id
}
`

const testAccIdentityAccessKeyDataSourceWithoutParams_basic = `
data "huaweicloud_identity_access_key" "test" {}
`
