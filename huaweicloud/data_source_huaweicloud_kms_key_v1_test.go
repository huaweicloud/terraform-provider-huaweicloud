package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var keyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))
var keyAlias_epsId = fmt.Sprintf("key_alias_%s", acctest.RandString(5))

func TestAccKmsKeyDataSource_Basic(t *testing.T) {
	var datasourceName = "data.huaweicloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckKms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_Basic(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "rotation_enabled", "false"),
					resource.TestCheckResourceAttr(datasourceName, "region", HW_REGION_NAME),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithTags(t *testing.T) {
	var datasourceName = "data.huaweicloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckKms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_WithTags(keyAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias),
					resource.TestCheckResourceAttr(datasourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(datasourceName, "tags.key", "value"),
				),
			},
		},
	})
}

func TestAccKmsKeyDataSource_WithEpsId(t *testing.T) {
	var datasourceName = "data.huaweicloud_kms_key.key_1"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckKms(t); testAccPreCheckEpsID(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsKeyDataSource_epsId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKmsKeyDataSourceID(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "key_alias", keyAlias_epsId),
					resource.TestCheckResourceAttr(datasourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckKmsKeyDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find Kms key data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Kms key data source ID not set")
		}

		return nil
	}
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

var testAccKmsKeyDataSource_epsId = fmt.Sprintf(`
resource "huaweicloud_kms_key_v1" "key_1" {
  key_alias       = "%s"
  key_description = "test description"
  pending_days    = "7"
  is_enabled      = true
  enterprise_project_id = "%s"
}

data "huaweicloud_kms_key_v1" "key_1" {
  key_alias       = huaweicloud_kms_key_v1.key_1.key_alias
  key_id          = huaweicloud_kms_key_v1.key_1.id
  key_description = "test description"
  key_state       = "2"
  enterprise_project_id = huaweicloud_kms_key_v1.key_1.enterprise_project_id
}
`, keyAlias_epsId, HW_ENTERPRISE_PROJECT_ID_TEST)
