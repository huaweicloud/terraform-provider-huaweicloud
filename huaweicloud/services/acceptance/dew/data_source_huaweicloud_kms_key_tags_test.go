package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKmsKeyTags_basic(t *testing.T) {
	var (
		dataSourcename = "data.huaweicloud_kms_key_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourcename)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare a import KMS key ID with tag and config it to the environment variable.
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKmsKeyTags_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSourcename, "tags.0.values.#"),
				),
			},
		},
	})
}

const testAccDataSourceKmsKeyTags_basic = `data "huaweicloud_kms_key_tags" "test" {}`
