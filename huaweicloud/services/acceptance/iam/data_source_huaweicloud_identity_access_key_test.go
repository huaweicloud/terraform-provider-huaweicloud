package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccessKey_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identity_access_key.test"

	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAccessKey_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.status", "active"),
				),
			},
			{
				Config: testAccDataSourceAccessKey_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.status", "active"),
				),
			},
			{
				Config: testAccDataSourceAccessKey_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "credentials.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credentials.0.status", "active"),
				),
			},
		},
	})
}

func testAccDataSourceAccessKey_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "password@123!"
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_identity_access_key" "test" {
  user_id = huaweicloud_identity_user.test.id
}
`, name)
}

func testAccDataSourceAccessKey_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_access_key" "test" {
  user_id = huaweicloud_identity_access_key.test.user_id
}
`, testAccDataSourceAccessKey_base(name))
}

func testAccDataSourceAccessKey_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_access_key" "test" {
  access_key = huaweicloud_identity_access_key.test.id
}
`, testAccDataSourceAccessKey_base(name))
}

func testAccDataSourceAccessKey_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_access_key" "test" {}
`, testAccDataSourceAccessKey_base(name))
}
