package cph

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCphPhones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cph_phones.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCphPhones_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "phones.0.phone_name"),
					resource.TestCheckResourceAttrSet(dataSource, "phones.0.server_id"),
					resource.TestCheckResourceAttrSet(dataSource, "phones.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "phones.0.status"),

					resource.TestCheckOutput("is_phone_name_filter_useful", "true"),
					resource.TestCheckOutput("is_server_id_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
			{
				Config: testCphServerBase(rName),
				Check: resource.ComposeTestCheckFunc(
					waitForDeletionCooldownComplete(),
				),
			},
		},
	})
}

func testDataSourceCphPhones_basic(name string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_cph_phones" "test" {
  depends_on = [huaweicloud_cph_server.test]
}

locals {
  phone_name = data.huaweicloud_cph_phones.test.phones[0].phone_name
  server_id  = data.huaweicloud_cph_phones.test.phones[0].server_id
  status     = data.huaweicloud_cph_phones.test.phones[0].status
  type       = tostring(data.huaweicloud_cph_phones.test.phones[0].type)
}

data "huaweicloud_cph_phones" "filter_by_phone_name" {
  phone_name = local.phone_name
}

data "huaweicloud_cph_phones" "filter_by_server_id" {
  server_id = local.server_id
}

data "huaweicloud_cph_phones" "filter_by_type" {
  type = local.type
}

data "huaweicloud_cph_phones" "filter_by_status" {
  status = local.status
}

locals {
  list_by_phone_name = data.huaweicloud_cph_phones.filter_by_phone_name.phones
  list_by_server_id  = data.huaweicloud_cph_phones.filter_by_server_id.phones
  list_by_type       = data.huaweicloud_cph_phones.filter_by_type.phones
  list_by_status     = data.huaweicloud_cph_phones.filter_by_status.phones
}

output "is_phone_name_filter_useful" {
  value = length(local.list_by_phone_name) > 0 && alltrue(
    [for v in local.list_by_phone_name[*].phone_name : v == local.phone_name]
  )
}

output "is_server_id_filter_useful" {
  value = length(local.list_by_server_id) > 0 && alltrue(
    [for v in local.list_by_server_id[*].server_id : v == local.server_id]
  )
}

output "is_type_filter_useful" {
  value = length(local.list_by_type) > 0 && alltrue(
    [for v in local.list_by_type[*].type : tostring(v) == local.type]
  )
}

output "is_status_filter_useful" {
  value = length(local.list_by_status) > 0 && alltrue(
    [for v in local.list_by_status[*].status : v == local.status]
  )
}
`, testCphServer_basic(name))
}
