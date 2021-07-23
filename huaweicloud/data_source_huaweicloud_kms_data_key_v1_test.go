package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var datakeyAlias = fmt.Sprintf("key_alias_%s", acctest.RandString(5))

func TestAccKmsDataKeyV1DataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckKms(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccKmsDataKeyV1DataSource_key,
			},
			{
				Config: testAccKmsDataKeyV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_kms_data_key_v1.kms_datakey1", "plain_text"),
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_kms_data_key_v1.kms_datakey1", "cipher_text"),
				),
			},
		},
	})
}

var testAccKmsDataKeyV1DataSource_key = fmt.Sprintf(`
resource "huaweicloud_kms_key_v1" "key1" {
  key_alias    = "%s"
  pending_days = "7"
}`, datakeyAlias)

var testAccKmsDataKeyV1DataSource_basic = fmt.Sprintf(`
%s
data "huaweicloud_kms_data_key_v1" "kms_datakey1" {
  key_id           =   "${huaweicloud_kms_key_v1.key1.id}"
  datakey_length   =   "512"
}
`, testAccKmsDataKeyV1DataSource_key)
