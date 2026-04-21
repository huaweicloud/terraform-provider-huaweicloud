package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsTemplateVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rfs_template_versions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		name       = acceptance.RandomAccResourceName()
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRfsTemplateVersions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.version_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.template_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.create_time"),
				),
			},
		},
	})
}

func testAccDataSourceRfsTemplateVersions_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_rfs_template" "test" {
  template_name        = "%s"
  template_description = "template for versions datasource test"
  version_description  = "v1"
  template_body        = <<-EOF
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}
EOF
}

resource "huaweicloud_rfs_template_version" "test" {
  template_name       = huaweicloud_rfs_template.test.template_name
  template_id         = huaweicloud_rfs_template.test.template_id
  version_description = "v2 - add subnet"
  template_body       = <<-EOF
variable "vpc_name" {
  type    = string
  default = "my-vpc"
}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name       = "my-subnet"
  vpc_id     = huaweicloud_vpc.vpc.id
  cidr       = "172.16.1.0/24"
  gateway_ip = "172.16.1.1"
}
EOF
}
`, name)
}

func testAccDataSourceRfsTemplateVersions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_rfs_template_versions" "test" {
  template_name = huaweicloud_rfs_template.test.template_name

  depends_on = [
    huaweicloud_rfs_template_version.test
  ]
}
`, testAccDataSourceRfsTemplateVersions_base(name))
}
