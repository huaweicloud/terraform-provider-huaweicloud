package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsServerTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_server_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsServerTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.volume_type"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.vpc.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.vpc.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.vpc.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.nics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.nics.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.nics.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.security_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.security_groups.0.name"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
					resource.TestCheckOutput("region_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsServerTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_sms_server_templates" "test" {
  depends_on = [huaweicloud_sms_server_template.test]
}

data "huaweicloud_sms_server_templates" "name_filter" {
  name = huaweicloud_sms_server_template.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_sms_server_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_sms_server_templates.name_filter.templates[*].name :
      v == huaweicloud_sms_server_template.test.name]
  )
}

data "huaweicloud_sms_server_templates" "availability_zone_filter" {
  availability_zone = huaweicloud_sms_server_template.test.availability_zone
}

output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_sms_server_templates.availability_zone_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_sms_server_templates.availability_zone_filter.templates[*].availability_zone :
      v == huaweicloud_sms_server_template.test.availability_zone]
  )
}

data "huaweicloud_sms_server_templates" "region_filter" {
  region = huaweicloud_sms_server_template.test.region
}

output "region_filter_is_useful" {
  value = length(data.huaweicloud_sms_server_templates.region_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_sms_server_templates.region_filter.templates[*].region :
      v == huaweicloud_sms_server_template.test.region]
  )
}
`, testAccServerTemplate_basic(name))
}
