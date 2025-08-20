package dns

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDNSV21PtrRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dnsv21_ptrrecords.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := fmt.Sprintf("acpttest-ptr-%s.com", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDNSV21PtrRecords_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.names.#"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.publicip_id"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.ttl"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.address"),
					resource.TestCheckResourceAttrSet(dataSource, "ptrrecords.0.tags.%"),

					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDNSV21PtrRecords_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dnsv21_ptrrecords" "test" {
  depends_on = [huaweicloud_dnsv21_ptrrecord.test]
}

// filter by eps_id
data "huaweicloud_dnsv21_ptrrecords" "filter_by_eps_id" {
  enterprise_project_id = huaweicloud_dnsv21_ptrrecord.test.enterprise_project_id
}

locals {
  filter_result_by_eps_id = [for v in data.huaweicloud_dnsv21_ptrrecords.filter_by_eps_id.ptrrecords[*].enterprise_project_id :
    v == huaweicloud_dnsv21_ptrrecord.test.enterprise_project_id]
}

output "is_eps_id_filter_useful" {
  value = length(local.filter_result_by_eps_id) > 0 && alltrue(local.filter_result_by_eps_id)
}

// filter by status
data "huaweicloud_dnsv21_ptrrecords" "filter_by_status" {
  status = huaweicloud_dnsv21_ptrrecord.test.status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_dnsv21_ptrrecords.filter_by_status.ptrrecords[*].status :
    v == huaweicloud_dnsv21_ptrrecord.test.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status)
}

// filter by tags
data "huaweicloud_dnsv21_ptrrecords" "filter_by_tags" {
  tags = huaweicloud_dnsv21_ptrrecord.test.tags
}

locals {
  filter_result_by_tags = [for v in data.huaweicloud_dnsv21_ptrrecords.filter_by_tags.ptrrecords[*].tags :
    v == huaweicloud_dnsv21_ptrrecord.test.tags]
}

output "is_tags_filter_useful" {
  value = length(local.filter_result_by_tags) > 0 && alltrue(local.filter_result_by_tags)
}
`, testAccDNSV21PtrRecord_basic(name))
}
