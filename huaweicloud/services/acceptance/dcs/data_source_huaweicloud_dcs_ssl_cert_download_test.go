package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsSslCertDownload_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_dcs_ssl_cert_download.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	name := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDcsSslCertDownload_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "file_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "link"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bucket_name"),
				),
			},
		},
	})
}

func testAccDataSourceDcsSslCertDownload_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_dcs_flavors" "test" {
  engine         = "Redis"
  engine_version = "6.0"
  capacity       = 1
  name           = "redis.ha.au1.large.r4.1"
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_dcs_instance" "test" {
  count = 2

  name               = "%[1]s_${count.index}"
  engine_version     = "6.0"
  password           = "Huawei_test"
  engine             = "Redis"
  capacity           = 1
  vpc_id             = data.huaweicloud_vpc.test.id
  subnet_id          = data.huaweicloud_vpc_subnet.test.id
  availability_zones = [data.huaweicloud_availability_zones.test.names[0]]
  flavor             = data.huaweicloud_dcs_flavors.test.flavors[0].name
}
`, name)
}

func testAccDataSourceDcsSslCertDownload_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_ssl_cert_download" "test" {
  instance_id = huaweicloud_dcs_instance.test[0].id
}
`, testAccDataSourceDcsSslCertDownload_base(name))
}
