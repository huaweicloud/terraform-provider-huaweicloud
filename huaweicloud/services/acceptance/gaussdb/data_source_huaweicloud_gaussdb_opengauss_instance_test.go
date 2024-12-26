package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccOpenGaussInstanceDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "sharding_num", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "coordinator_num", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstanceDataSource_haModeCentralized(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
			acceptance.TestAccPreCheckHighCostAllow(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstanceDataSource_haModeCentralized(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstanceDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "replica_num", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "volume.0.size", "40"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstanceDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB opengauss instance data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB opengauss data source ID not set ")
		}

		return nil
	}
}

func testAccOpenGaussInstanceDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_instance" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, testAccOpenGaussInstance_basic(rName, fmt.Sprintf("%s@123", acctest.RandString(5)), 2))
}

func testAccOpenGaussInstanceDataSource_haModeCentralized(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_instance" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, testAccOpenGaussInstance_haModeCentralized(rName, fmt.Sprintf("%s@123", acctest.RandString(5))))
}
