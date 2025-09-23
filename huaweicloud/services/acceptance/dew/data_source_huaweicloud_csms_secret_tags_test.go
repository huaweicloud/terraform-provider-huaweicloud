package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCsmsSecretTags_basic(t *testing.T) {
	var (
		dataSourcename = "data.huaweicloud_csms_secret_tags.test"
		dc             = acceptance.InitDataSourceCheck(dataSourcename)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare an import CSMS secret with tags and config it to the environment variable.
			acceptance.TestAccPreCheckCsmsSecretID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCsmsSecretTags_basic,
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

const testAccDataSourceCsmsSecretTags_basic = `data "huaweicloud_csms_secret_tags" "test" {}`
