package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccDataSourceRdsTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_tags.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsTags_basic(rName),
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

func testDataSourceRdsTags_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_rds_flavors" "instance" {
  db_type       = "MySQL"
  db_version    = "8.0"
  instance_mode = "single"
  group_type    = "dedicated"
  vcpus         = 4
}

resource "huaweicloud_rds_instance" "test" {
  name              = "%[2]s"
  flavor            = data.huaweicloud_rds_flavors.instance.flavors[0].name
  security_group_id = huaweicloud_networking_secgroup.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  vpc_id            = huaweicloud_vpc.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  db {
    password = "test_1234"
    type     = "MySQL"
    version  = "8.0"
  }

  volume {
    type = "CLOUDSSD"
    size = 40
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testDataSourceRdsTags_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rds_tags" "test" {
  depends_on = [huaweicloud_rds_instance.test]
}
`, testDataSourceRdsTags_base(rName))
}
