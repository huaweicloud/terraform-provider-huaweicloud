package hss

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Because there is a lack of scenarios for testing the API, the test case only tests one failure error.
func TestAccEventDeleteIsolatedFile_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testEventDeleteIsolatedFile_basic,
				ExpectError: regexp.MustCompile(`error deleting HSS isolated file`),
			},
		},
	})
}

const testEventDeleteIsolatedFile_basic string = `
resource "huaweicloud_hss_event_delete_isolated_file" "test" {
  data_list {
    host_id   = "not-exist-host_id"
    file_hash = "not-exist-file_hash"
    file_path = "not-exist-file_path"
    file_attr = "not-exist-file_attr"
  }

  enterprise_project_id = "0"
}
`
