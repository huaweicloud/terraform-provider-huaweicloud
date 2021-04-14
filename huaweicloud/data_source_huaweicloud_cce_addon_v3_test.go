package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCCEAddonV3DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "data.huaweicloud_cce_addon.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCCEAddonV3DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEAddonV3DataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "template_name", "metrics-server"),
				),
			},
		},
	})
}

func testAccCheckCCEAddonV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find Addons data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Addon data source ID not set ")
		}

		return nil
	}
}

func testAccCCEAddonV3DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_addon" "test" {
  cluster_id    = huaweicloud_cce_cluster.test.id
  template_name = huaweicloud_cce_addon.test.template_name
}
`, testAccCCEAddonV3_basic(rName))
}
