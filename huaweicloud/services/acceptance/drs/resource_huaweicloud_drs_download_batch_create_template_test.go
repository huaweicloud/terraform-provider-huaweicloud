package drs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDownloadBatchCreateTemplate_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDownloadBatchCreateTemplate_basic(),
			},
		},
	})
}

func testDownloadBatchCreateTemplate_basic() string {
	return `
resource "huaweicloud_drs_download_batch_create_template" "test" {
  engine_type        = "postgresql"
  template_file_name = "template_test"
}
`
}
