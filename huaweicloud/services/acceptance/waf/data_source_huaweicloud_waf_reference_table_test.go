package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDataSourceReferenceTablesV1_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_reference_tables.ref_table"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccReferenceTablesV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReferenceTablesId(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.conditions.0"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tables.0.creation_time"),
				),
			},
		},
	})
}

func testAccCheckReferenceTablesId(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find WAF reference tables data source: %s.", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The WAF reference tables data source ID not set.")
		}
		return nil
	}
}

func testAccReferenceTablesV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_reference_tables" "ref_table" {
  depends_on = [huaweicloud_waf_reference_table.ref_table]
}
`, testAccWafReferenceTableV1_conf(name))
}
