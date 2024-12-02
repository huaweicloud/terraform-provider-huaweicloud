package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBDehResourceDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBDehResourceDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGeminiDBDehResourceDataSourceID("data.huaweicloud_gaussdb_cassandra_dedicated_resource.test"),
				),
			},
		},
	})
}

func testAccCheckGeminiDBDehResourceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB cassandra dedicated resource data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB cassandra dedicated resource data source ID not set ")
		}

		return nil
	}
}

const testAccGeminiDBDehResourceDataSource_basic = `
data "huaweicloud_gaussdb_cassandra_dedicated_resource" "test" {
}
`
