package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dns_tags.test"
	name := fmt.Sprintf("acpttest-zone-%s.com", acctest.RandString(5))
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDNSTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDNSTags_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dns_tags" "test" {
  depends_on = [huaweicloud_dns_zone.test]
  
  resource_type = "DNS-public_zone"
}
`, testAccDNSZone_basic(name))
}
