package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceRdsSslCertDownloadLinks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_ssl_cert_download_links.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsSslCertDownloadLinks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "cert_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cert_info_list.0.download_link"),
					resource.TestCheckResourceAttrSet(dataSource, "cert_info_list.0.category"),
				),
			},
		},
	})
}

func testDataSourceRdsSslCertDownloadLinks_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_flavors" "test" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.test.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id
  availability_zone = slice(sort(data.huaweicloud_rds_flavors.test.flavors[0].availability_zones), 0, 1)

  db {
    password = "Huangwei!120521"
    type     = "MySQL"
    version  = "8.0"
    port     = 3306
  }
    
  volume {
    type = "CLOUDSSD"
    size = 40
  }
}
`, common.TestBaseNetwork(name), name)
}

func testDataSourceRdsSslCertDownloadLinks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_rds_ssl_cert_download_links" "test" {
  instance_id = huaweicloud_rds_instance.test.id
}
`, testDataSourceRdsSslCertDownloadLinks_base(name))
}
