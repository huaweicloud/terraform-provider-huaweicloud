package geminidb

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDatabaseVersions_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_geminidb_database_versions.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDatabaseVersions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "versions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

const testAccDataSourceDatabaseVersions_basic = `
data "huaweicloud_geminidb_database_versions" "test" {
  datastore_name = "redis"
}
`
