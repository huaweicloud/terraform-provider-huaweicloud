package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsKeyDataSource_Basic(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	datasourceName := "data.huaweicloud_kms_key.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "rotation_enabled", "false"),
					resource.TestCheckResourceAttr(datasourceName, "region", acceptance.HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithTags(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	var datasourceName = "data.huaweicloud_kms_key.key_1"
	dc := acceptance.InitDataSourceCheck(datasourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(datasourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithEpsId(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	var datasourceName = "data.huaweicloud_kms_key.key_1"
	dc := acceptance.InitDataSourceCheck(datasourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t); acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_epsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccKmsKeyDataSource_Basic(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_key" "key_1" {
  key_alias = huaweicloud_kms_key.key_1.key_alias
  key_id    = huaweicloud_kms_key.key_1.id
  key_state = "2"
}
`, testAccKmsKey_Basic(keyAlias))
}

func testAccKmsKeyDataSource_WithTags(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_key" "key_1" {
  key_alias = huaweicloud_kms_key.key_1.key_alias
  key_id    = huaweicloud_kms_key.key_1.id
  key_state = "2"
}
`, testAccKmsKey_WithTags(keyAlias))
}

func testAccKmsKeyDataSource_epsId(keyAlias string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias             = "%s"
  key_description       = "test description"
  pending_days          = "7"
  is_enabled            = true
  enterprise_project_id = "%s"
}

data "huaweicloud_kms_key" "key_1" {
  key_alias             = huaweicloud_kms_key.key_1.key_alias
  key_id                = huaweicloud_kms_key.key_1.id
  key_description       = "test description"
  key_state             = "2"
  enterprise_project_id = huaweicloud_kms_key.key_1.enterprise_project_id
}
`, keyAlias, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
