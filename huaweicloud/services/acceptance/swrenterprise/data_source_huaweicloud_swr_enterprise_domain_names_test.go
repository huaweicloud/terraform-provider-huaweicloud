package swrenterprise

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSwrEnterpriseDomainNames_basic(t *testing.T) {
	dataSource := "data.huaweicloud_swr_enterprise_domain_names.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSwrEnterpriseDomainNames_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.0.uid"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "domain_name_infos.0.updated_at"),

					resource.TestCheckOutput("uid_filter_is_useful", "true"),
					resource.TestCheckOutput("domain_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSwrEnterpriseDomainNames_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_swr_enterprise_domain_names" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
}

data "huaweicloud_swr_enterprise_domain_names" "filter_by_uid" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  uid         = data.huaweicloud_swr_enterprise_domain_names.test.domain_name_infos[0].uid
}

output "uid_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_domain_names.filter_by_uid.domain_name_infos) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_domain_names.filter_by_uid.domain_name_infos[*].uid : 
	  v == data.huaweicloud_swr_enterprise_domain_names.test.domain_name_infos[0].uid]
  )
}

data "huaweicloud_swr_enterprise_domain_names" "filter_by_domain_name" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  domain_name = data.huaweicloud_swr_enterprise_domain_names.test.domain_name_infos[0].domain_name
}

output "domain_name_filter_is_useful" {
  value = length(data.huaweicloud_swr_enterprise_domain_names.filter_by_domain_name.domain_name_infos) > 0 && alltrue(
	[for v in data.huaweicloud_swr_enterprise_domain_names.filter_by_domain_name.domain_name_infos[*].domain_name : 
	  v == data.huaweicloud_swr_enterprise_domain_names.test.domain_name_infos[0].domain_name]
  )
}
`, testAccSwrEnterpriseImageSignaturePolicy_basic(name))
}
