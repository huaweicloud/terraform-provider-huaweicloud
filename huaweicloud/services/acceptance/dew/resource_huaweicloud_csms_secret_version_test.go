package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDewCsmsSecretVersion_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	resourceName := "data.huaweicloud_csms_secret_version.version_1"

	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDewCsmsSecretVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "secret_name", name),
					resource.TestCheckResourceAttrSet(resourceName, "version"),
					resource.TestCheckResourceAttr(resourceName, "secret_text", "this is a password"),
				),
			},
			{
				Config: testAccDewCsmsSecretVersion_version(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "secret_name", name),
					resource.TestCheckResourceAttr(resourceName, "version", "v1"),
					resource.TestCheckResourceAttr(resourceName, "secret_text", "this is a password"),
				),
			},
		},
	})
}

func testAccDewCsmsSecretVersion_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "secret_1" {
  name        = "%s"
  secret_text = "this is a password"
}

data "huaweicloud_csms_secret_version" "version_1" {
  secret_name = "%s"

  depends_on = [huaweicloud_csms_secret.secret_1]
}
`, name, name)
}

func testAccDewCsmsSecretVersion_version(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "secret_1" {
  name        = "%s"
  secret_text = "this is a new password"
}

data "huaweicloud_csms_secret_version" "version_1" {
  secret_name = "%s"
  version     = "v1"

  depends_on = [huaweicloud_csms_secret.secret_1]
}
`, name, name)
}
