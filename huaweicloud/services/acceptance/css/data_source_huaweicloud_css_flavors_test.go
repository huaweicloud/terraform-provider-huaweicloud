package css

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssFlavorsDataSource_basic(t *testing.T) {
	var (
		typeFilter    = acceptance.InitDataSourceCheck("data.huaweicloud_css_flavors.type_filter")
		versionFilter = acceptance.InitDataSourceCheck("data.huaweicloud_css_flavors.version_filter")
		vcpusFilter   = acceptance.InitDataSourceCheck("data.huaweicloud_css_flavors.vcpus_filter")
		memoryFilter  = acceptance.InitDataSourceCheck("data.huaweicloud_css_flavors.memory_filter")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_basic,
				Check: resource.ComposeTestCheckFunc(
					typeFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					versionFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					vcpusFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_vcpus_filter_useful", "true"),
					memoryFilter.CheckResourceExists(),
					resource.TestCheckOutput("is_memory_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_basic = `

data "huaweicloud_css_flavors" "type_filter" {
  type = "ess"
}

output "is_type_filter_useful" {
  value = !contains([for v in data.huaweicloud_css_flavors.type_filter.flavors[*].type : v == "ess"], "false")
}

data "huaweicloud_css_flavors" "version_filter" {
  version = "7.9.3"
}

output "is_version_filter_useful" {
  value = !contains([for v in data.huaweicloud_css_flavors.version_filter.flavors[*].version : v == "7.9.3"], "false")
}

data "huaweicloud_css_flavors" "vcpus_filter" {
  vcpus = 32
}

output "is_vcpus_filter_useful" {
  value = !contains([for v in data.huaweicloud_css_flavors.vcpus_filter.flavors[*].vcpus : v == 32], "false")
}

data "huaweicloud_css_flavors" "memory_filter" {
  memory = 256
}

output "is_memory_filter_useful" {
  value = !contains([for v in data.huaweicloud_css_flavors.memory_filter.flavors[*].memory : v == 256], "false")
}
`

func TestAccCssFlavorsDataSource_all(t *testing.T) {
	dataSourceName := "data.huaweicloud_css_flavors.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCssFlavors_all,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.type", "ess"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.version", "7.9.3"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.id"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.region", "cn-north-4"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.name", "ess.spec-ds.8xlarge.8"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.memory", "256"),
					resource.TestCheckResourceAttr(dataSourceName, "flavors.0.vcpus", "32"),
					resource.TestCheckResourceAttrSet(dataSourceName, "flavors.0.disk_range"),
				),
			},
		},
	})
}

const testAccDataSourceCssFlavors_all = `
data "huaweicloud_css_flavors" "test" {
  type    = "ess"
  version = "7.9.3"
  vcpus   = 32
  memory  = 256
  region  = "cn-north-4"
  name    = "ess.spec-ds.8xlarge.8"
}
`
