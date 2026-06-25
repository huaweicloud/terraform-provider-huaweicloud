package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBInstanceAssociatedEip_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource := "data.huaweicloud_taurusdb_instance_associated_eip.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBInstanceAssociatedEip_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "can_enable_public_access"),
					resource.TestCheckResourceAttrSet(dataSource, "eip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "type"),
					resource.TestCheckResourceAttrSet(dataSource, "port_id"),
					resource.TestCheckResourceAttrSet(dataSource, "public_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "private_ip_address"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_name"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_size"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_share_type"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBInstanceAssociatedEip_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_taurusdb_instance_associated_eip" "test" {
  instance_id = huaweicloud_taurusdb_instance.test.id

  depends_on = [huaweicloud_taurusdb_eip_associate.test]
}
`, testTaurusDBEipAssociate_basic(rName))
}
