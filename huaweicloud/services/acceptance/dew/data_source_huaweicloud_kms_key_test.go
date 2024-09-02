package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKmsKeyDataSource_basic(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	datasourceName := "data.huaweicloud_kms_key.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "rotation_enabled", "false"),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_withTags(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	var datasourceName = "data.huaweicloud_kms_key.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_withTags(keyAlias),
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

func TestAccKmsKeyDataSource_withEpsId(t *testing.T) {
	keyAlias := acceptance.RandomAccResourceName()
	var datasourceName = "data.huaweicloud_kms_key.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckKms(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_withEpsId(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func testAccKmsKeyDataSource_basic(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_key" "test" {
  key_alias = huaweicloud_kms_key.test.key_alias
  key_id    = huaweicloud_kms_key.test.id
  key_state = "2"
}
`, testAccKmsKey_basic(keyAlias))
}

func testAccKmsKeyDataSource_withTags(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_key" "test" {
  key_alias = huaweicloud_kms_key.test.key_alias
  key_id    = huaweicloud_kms_key.test.id
  key_state = "2"
}
`, testAccKmsKey_basic(keyAlias))
}

func testAccKmsKeyDataSource_withEpsId(keyAlias string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_kms_key" "test" {
  key_alias             = huaweicloud_kms_key.test.key_alias
  key_id                = huaweicloud_kms_key.test.id
  key_state             = "2"
  enterprise_project_id = huaweicloud_kms_key.test.enterprise_project_id
}
`, testAccKmsKey_basic(keyAlias))
}
