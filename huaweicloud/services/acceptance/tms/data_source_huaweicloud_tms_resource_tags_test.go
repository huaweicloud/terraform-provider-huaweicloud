package tms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTmsResourceTags_basic(t *testing.T) {
	dataSource := "data.huaweicloud_tms_resource_tags.test"
	name := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTmsResourceTags_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tags.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "tags.0.values.#"),
					resource.TestCheckOutput("huaweicloud_tms_resource_tags_length", "true"),
				),
			},
		},
	})
}

func testDataSourceTmsResourceTags_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_identity_projects" "test" {
  name = "ap-southeast-1"
}

resource "huaweicloud_vpc" "vpc_1" {
  name = "%[1]s_1"
  cidr = "192.168.0.0/16"
  tags = {
    k1 = "v1"
    k2 = "v2"
  }
}

resource "huaweicloud_vpc" "vpc_2" {
  name = "%[1]s_2"
  cidr = "192.168.0.0/16"
  tags = {
    k1 = "v1"
    k3 = "v3"
  }
}

`, name)
}

func testDataSourceTmsResourceTags_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_tms_resource_tags" "test" {
  resource_types = "vpc"
  project_id     = data.huaweicloud_identity_projects.test.projects[0].id
}
output "huaweicloud_tms_resource_tags_length" {
  value = length(data.huaweicloud_tms_resource_tags.test.tags) >= 3
}
`, testDataSourceTmsResourceTags_base(name))
}
