package cfw

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDownloadImportedAclRule_basic(t *testing.T) {
	exportFileName := fmt.Sprintf("./%s.xlsx", acceptance.RandomAccResourceName())
	defer func() {
		if err := os.Remove(exportFileName); err != nil {
			log.Printf("error deleting testing file %s: %s", exportFileName, err)
		}
	}()

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDownloadImportedAclRule_basic(exportFileName),
			},
		},
	})
}

func testDownloadImportedAclRule_basic(exportFileName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_download_imported_acl_rule" "test" {
  object_id        = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  export_file_name = "%[2]s"
}
`, testAccDatasourceFirewalls_basic(), exportFileName)
}
