package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceDNSRecordsets_basic(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))

		rName = "data.huaweicloud_dns_recordsets.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_dns_recordsets.recordset_id_filter"
		dcById = acceptance.InitDataSourceCheck(byId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSRecordsets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckOutput("line_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("recordset_id_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.id"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.name"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.zone_id"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.zone_name"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.type"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.ttl"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.records.#"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.status"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.default"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.line_id"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.weight"),
					resource.TestMatchResourceAttr(byId, "recordsets.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "recordsets.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr("data.huaweicloud_dns_recordsets.tags_filter", "recordsets.#", "1"),
				),
			},
		},
	})
}

func TestAccDatasourceDNSRecordsets_private(t *testing.T) {
	var (
		name = fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))

		rName = "data.huaweicloud_dns_recordsets.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byId   = "data.huaweicloud_dns_recordsets.recordset_id_filter"
		dcById = acceptance.InitDataSourceCheck(byId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceDNSRecordsets_private(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(rName, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("recordset_id_filter_is_useful", "true"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.id"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.name"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.zone_id"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.zone_name"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.type"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.ttl"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.records.#"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.status"),
					resource.TestCheckResourceAttrSet(byId, "recordsets.0.default"),
					resource.TestMatchResourceAttr(byId, "recordsets.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "recordsets.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttr("data.huaweicloud_dns_recordsets.tags_filter", "recordsets.#", "1"),
				),
			},
		},
	})
}

func testAccDatasourceDNSRecordsets_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dns_recordsets" "test" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
}

data "huaweicloud_dns_recordsets" "line_id_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  line_id = huaweicloud_dns_recordset.test.line_id
}
data "huaweicloud_dns_recordsets" "status_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  status  = "ACTIVE"
}
data "huaweicloud_dns_recordsets" "type_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  type    = huaweicloud_dns_recordset.test.type
}
data "huaweicloud_dns_recordsets" "name_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  name    = huaweicloud_dns_recordset.test.name
}
data "huaweicloud_dns_recordsets" "recordset_id_filter" {
  zone_id      = huaweicloud_dns_recordset.test.zone_id
  recordset_id = split("/", huaweicloud_dns_recordset.test.id).1
}

locals {
  line_id_filter_result = [for v in data.huaweicloud_dns_recordsets.line_id_filter.recordsets[*].line_id :
v == huaweicloud_dns_recordset.test.line_id]
  status_filter_result = [for v in data.huaweicloud_dns_recordsets.status_filter.recordsets[*].status : v == "ACTIVE"]
  type_filter_result = [for v in data.huaweicloud_dns_recordsets.type_filter.recordsets[*].type : v == huaweicloud_dns_recordset.test.type]
  name_filter_result = [for v in data.huaweicloud_dns_recordsets.name_filter.recordsets[*].name : v == huaweicloud_dns_recordset.test.name]
  recordset_id_filter_result = [for v in data.huaweicloud_dns_recordsets.recordset_id_filter.recordsets[*].id :
v == split("/", huaweicloud_dns_recordset.test.id).1]
}

output "line_id_filter_is_useful" {
  value = alltrue(local.line_id_filter_result) && length(local.line_id_filter_result) > 0
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "recordset_id_filter_is_useful" {
  value = alltrue(local.recordset_id_filter_result) && length(local.recordset_id_filter_result) > 0
}

data "huaweicloud_dns_recordsets" "tags_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  tags    = "key1,value1"
}
`, testDNSRecordset_basic(name))
}

func testAccDatasourceDNSRecordsets_private(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dns_zone" "zone_1" {
  name      = "%[2]s"
  zone_type = "private"

  router {
    router_id = huaweicloud_vpc.test.id
  }
}

resource "huaweicloud_dns_recordset" "test" {
  zone_id     = huaweicloud_dns_zone.zone_1.id
  name        = "%[2]s"
  type        = "A"
  description = "Created a record set by script"
  ttl         = 600
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar_private"
  }
}

data "huaweicloud_dns_recordsets" "test" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
}

data "huaweicloud_dns_recordsets" "status_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  status  = "ACTIVE"
}
data "huaweicloud_dns_recordsets" "type_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  type    = huaweicloud_dns_recordset.test.type
}
data "huaweicloud_dns_recordsets" "name_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  name    = huaweicloud_dns_recordset.test.name
}
data "huaweicloud_dns_recordsets" "recordset_id_filter" {
  zone_id      = huaweicloud_dns_recordset.test.zone_id
  recordset_id = split("/", huaweicloud_dns_recordset.test.id).1
}

locals {
  status_filter_result = [for v in data.huaweicloud_dns_recordsets.status_filter.recordsets[*].status : v == "ACTIVE"]
  type_filter_result = [for v in data.huaweicloud_dns_recordsets.type_filter.recordsets[*].type : v == huaweicloud_dns_recordset.test.type]
  name_filter_result = [for v in data.huaweicloud_dns_recordsets.name_filter.recordsets[*].name : v == huaweicloud_dns_recordset.test.name]
  recordset_id_filter_result = [for v in data.huaweicloud_dns_recordsets.recordset_id_filter.recordsets[*].id :
v == split("/", huaweicloud_dns_recordset.test.id).1]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "recordset_id_filter_is_useful" {
  value = alltrue(local.recordset_id_filter_result) && length(local.recordset_id_filter_result) > 0
}

data "huaweicloud_dns_recordsets" "tags_filter" {
  zone_id = huaweicloud_dns_recordset.test.zone_id
  tags    = "foo,bar_private"
}
`, acceptance.RandomAccResourceName(), name)
}
