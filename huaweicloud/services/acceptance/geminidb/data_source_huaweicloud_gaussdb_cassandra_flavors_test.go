package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCassandraFlavorsDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCassandraFlavorsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCassandraFlavorsDataSourceID("data.huaweicloud_gaussdb_cassandra_flavors.test"),
				),
			},
		},
	})
}

func testAccCheckCassandraFlavorsDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB cassandra flavors data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB cassandra flavors data source ID not set ")
		}

		return nil
	}
}

const testAccCassandraFlavorsDataSource_basic = `
data "huaweicloud_gaussdb_cassandra_flavors" "test" {
}
`
