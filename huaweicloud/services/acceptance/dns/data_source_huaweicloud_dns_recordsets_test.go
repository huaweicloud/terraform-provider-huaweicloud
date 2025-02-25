package dns

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRecordsets_basic(t *testing.T) {
	var (
		name  = fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))
		rName = "huaweicloud_dns_recordset.test.0"

		dataSource = "data.huaweicloud_dns_recordsets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byRecordsetId           = "data.huaweicloud_dns_recordsets.filter_by_recordset_id"
		dcByRecordsetId         = acceptance.InitDataSourceCheck(byRecordsetId)
		byNotFoundRecordsetId   = "data.huaweicloud_dns_recordsets.filter_by_not_found_recordset_id"
		dcByNotFoundRecordsetId = acceptance.InitDataSourceCheck(byNotFoundRecordsetId)

		byNameFuzzy      = "data.huaweicloud_dns_recordsets.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNameExact      = "data.huaweicloud_dns_recordsets.filter_by_name_exact"
		dcByNameExact    = acceptance.InitDataSourceCheck(byNameExact)
		byNotFoundName   = "data.huaweicloud_dns_recordsets.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byLineId           = "data.huaweicloud_dns_recordsets.filter_by_line_id"
		dcByLineId         = acceptance.InitDataSourceCheck(byLineId)
		byNotFoundLineId   = "data.huaweicloud_dns_recordsets.filter_by_not_found_line_id"
		dcByNotFoundLineId = acceptance.InitDataSourceCheck(byNotFoundLineId)

		byStatus           = "data.huaweicloud_dns_recordsets.filter_by_status"
		dcByStatus         = acceptance.InitDataSourceCheck(byStatus)
		byNotFoundStatus   = "data.huaweicloud_dns_recordsets.filter_by_not_found_status"
		dcByNotFoundStatus = acceptance.InitDataSourceCheck(byNotFoundStatus)

		byType           = "data.huaweicloud_dns_recordsets.filter_by_type"
		dcByType         = acceptance.InitDataSourceCheck(byType)
		byNotFoundType   = "data.huaweicloud_dns_recordsets.filter_by_not_found_type"
		dcByNotFoundType = acceptance.InitDataSourceCheck(byNotFoundType)

		byTags           = "data.huaweicloud_dns_recordsets.filter_by_tags"
		dcByTags         = acceptance.InitDataSourceCheck(byTags)
		byNotFoundTags   = "data.huaweicloud_dns_recordsets.filter_by_not_found_tags"
		dcByNotFoundTags = acceptance.InitDataSourceCheck(byNotFoundTags)

		bySortAsc    = "data.huaweicloud_dns_recordsets.filter_by_sort_asc"
		dcBySortAsc  = acceptance.InitDataSourceCheck(bySortAsc)
		bySortDesc   = "data.huaweicloud_dns_recordsets.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataRecordsets_notFound(),
				ExpectError: regexp.MustCompile(`This zone does not exist`),
			},
			{
				Config: testAccDataRecordsets_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by recordset ID.
					dcByRecordsetId.CheckResourceExists(),
					resource.TestCheckOutput("is_recordset_id_filter_useful", "true"),
					dcByNotFoundRecordsetId.CheckResourceExists(),
					resource.TestCheckOutput("recordset_id_not_found_validation_pass", "true"),
					// Fuzzy search by recordset name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					// Exactly search by recordset name.
					dcByNameExact.CheckResourceExists(),
					resource.TestCheckOutput("is_name_exact_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by line ID.
					dcByLineId.CheckResourceExists(),
					resource.TestCheckOutput("is_line_id_filter_useful", "true"),
					dcByNotFoundLineId.CheckResourceExists(),
					resource.TestCheckOutput("line_id_not_found_validation_pass", "true"),
					// Filter by recordset status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByNotFoundStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_not_found_validation_pass", "true"),
					// Filter by recordset type.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByNotFoundType.CheckResourceExists(),
					resource.TestCheckOutput("type_not_found_validation_pass", "true"),
					// Filter by recordset tags.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNotFoundTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_not_found_validation_pass", "true"),
					// Check the sort results.
					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					// Check attributes.
					// The ID of the corresponding resource consists of zone ID and recordset ID.
					resource.TestCheckResourceAttrSet(byRecordsetId, "recordsets.0.id"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.zone_id", rName, "zone_id"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.zone_name", rName, "zone_name"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.type", rName, "type"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.ttl", rName, "ttl"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.records", rName, "records"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.weight", rName, "weight"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.description", rName, "description"),
					resource.TestMatchResourceAttr(byRecordsetId, "recordsets.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byRecordsetId, "recordsets.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// This value is false when created by the Terraform script.
					resource.TestCheckResourceAttr(byRecordsetId, "recordsets.0.default", "false"),
				),
			},
		},
	})
}

func testAccDataRecordsets_notFound() string {
	randomId, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_dns_recordsets" "test" {
  zone_id = "%[1]s"
}
`, randomId)
}

func testAccDataRecordsets_base() string {
	return fmt.Sprintf(`
data "huaweicloud_dns_recordsets" "test" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
}

# Filter by recordset ID.
locals {
  recordset_id = try(split("/", huaweicloud_dns_recordset.test.0.id)[1], "")
}

data "huaweicloud_dns_recordsets" "filter_by_recordset_id" {
  zone_id      = huaweicloud_dns_recordset.test.0.zone_id
  recordset_id = local.recordset_id
}

locals {
  recordset_id_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_recordset_id.recordsets[*].id :
  v == local.recordset_id]
}

output "is_recordset_id_filter_useful" {
  value = length(local.recordset_id_filter_result) > 0 && alltrue(local.recordset_id_filter_result)
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_recordset_id" {
  zone_id      = huaweicloud_dns_recordset.test.0.zone_id
  recordset_id = "recordset_id_not_found"
}

output "recordset_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_recordset_id.recordsets) == 0
}

# Filter by recordset name (fuzzy search).
locals {
  recordset_name = huaweicloud_dns_recordset.test.0.name
  name_suffix    = huaweicloud_dns_zone.test.name
}

data "huaweicloud_dns_recordsets" "filter_by_name_fuzzy" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  name    = local.name_suffix
}

locals {
  name_fuzzy_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_name_fuzzy.recordsets[*].name :
  strcontains(v, local.name_suffix)]
}

output "is_name_fuzzy_filter_useful" {
  value = length(local.name_fuzzy_filter_result) >= 2 && alltrue(local.name_fuzzy_filter_result)
}

# Filter by recordset name (exact search).
data "huaweicloud_dns_recordsets" "filter_by_name_exact" {
  zone_id     = huaweicloud_dns_recordset.test.0.zone_id
  name        = local.recordset_name
  search_mode = "equal"
}

locals {
  name_exact_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_name_exact.recordsets[*].name :
  v == local.recordset_name]
}

output "is_name_exact_filter_useful" {
  value = length(local.name_exact_filter_result) > 0 && alltrue(local.name_exact_filter_result)
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_name" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  name    = "recordset_name_not_found"
}

output "name_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_name.recordsets) == 0
}

# By filter recordset status.
data "huaweicloud_dns_recordsets" "filter_by_status" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  # In the corresponding resource, the value of status is ENABLE.
  status = "ACTIVE"
}

locals {
  status_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_status.recordsets[*].status : v == "ACTIVE"]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_status" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  status  = "status_not_found"
}

output "status_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_status.recordsets) == 0
}

# By filter recordset type.
locals {
  recordset_type = huaweicloud_dns_recordset.test.0.type
}

data "huaweicloud_dns_recordsets" "filter_by_type" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  type    = local.recordset_type
}

locals {
  type_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_type.recordsets[*].type : v == local.recordset_type]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_type" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  type    = "type_not_found"
}

output "type_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_type.recordsets) == 0
}

# Filter by recordset tags.
data "huaweicloud_dns_recordsets" "filter_by_tags" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  tags    = join("|", [for key, value in huaweicloud_dns_recordset.test.0.tags : format("%%v,%%v", key, value)])
}

output "is_tags_filter_useful" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_tags.recordsets) >= 2
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_tags" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  tags    = "not_found_tags,not_found"
}

output "tags_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_tags.recordsets) == 0
}

# Sort ascending order using name as the sort field.
data "huaweicloud_dns_recordsets" "filter_by_sort_asc" {
  zone_id  = huaweicloud_dns_recordset.test.0.zone_id
  sort_key = "name"
  sort_dir = "asc"
}

# Sort descending order using name as the sort field.
data "huaweicloud_dns_recordsets" "filter_by_sort_desc" {
  zone_id  = huaweicloud_dns_recordset.test.0.zone_id
  sort_key = "name"
  sort_dir = "desc"
}

locals {
  sort_desc_filter_result = data.huaweicloud_dns_recordsets.filter_by_sort_desc.recordsets
  sort_asc_first_name     = try(data.huaweicloud_dns_recordsets.filter_by_sort_asc.recordsets[0].name, "")
  sort_desc_last_name     = try(data.huaweicloud_dns_recordsets.filter_by_sort_desc.recordsets[length(local.sort_desc_filter_result) - 1].name, "")
}

output "sort_filter_is_useful" {
  value = length(local.sort_desc_filter_result) > 0 && local.sort_asc_first_name == local.sort_desc_last_name
}
`)
}

func testAccDataRecordsets_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dns_zone" "test" {
  name      = "%[1]s"
  zone_type = "public"
}

resource "huaweicloud_dns_recordset" "test" {
  count = 2

  zone_id     = huaweicloud_dns_zone.test.id
  name        = "${count.index}test.%[1]s"
  type        = "A"
  description = "Created recordset"
  ttl         = 300
  records     = ["10.1.0.0"]
  line_id     = "Dianxin_Shanxi"
  weight      = 3

  tags = {
    key1 = "value1"
    key2 = "value2"
  }
}

%[2]s

# By filter line ID.
locals {
  line_id = huaweicloud_dns_recordset.test.0.line_id
}

data "huaweicloud_dns_recordsets" "filter_by_line_id" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  line_id = local.line_id
}

locals {
  line_id_filter_result = [for v in data.huaweicloud_dns_recordsets.filter_by_line_id.recordsets[*].line_id :
  v == local.line_id]
}

output "is_line_id_filter_useful" {
  value = length(local.line_id_filter_result) > 0 && alltrue(local.line_id_filter_result)
}

data "huaweicloud_dns_recordsets" "filter_by_not_found_line_id" {
  zone_id = huaweicloud_dns_recordset.test.0.zone_id
  line_id = "line_id_not_found"
}

output "line_id_not_found_validation_pass" {
  value = length(data.huaweicloud_dns_recordsets.filter_by_not_found_line_id.recordsets) == 0
}
`, name, testAccDataRecordsets_base())
}

func TestAccDataRecordsets_private(t *testing.T) {
	var (
		name  = fmt.Sprintf("acpttest-recordset-%s.com.", acctest.RandString(5))
		rName = "huaweicloud_dns_recordset.test.0"

		dataSource = "data.huaweicloud_dns_recordsets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byRecordsetId           = "data.huaweicloud_dns_recordsets.filter_by_recordset_id"
		dcByRecordsetId         = acceptance.InitDataSourceCheck(byRecordsetId)
		byNotFoundRecordsetId   = "data.huaweicloud_dns_recordsets.filter_by_not_found_recordset_id"
		dcByNotFoundRecordsetId = acceptance.InitDataSourceCheck(byNotFoundRecordsetId)

		byNameFuzzy      = "data.huaweicloud_dns_recordsets.filter_by_name_fuzzy"
		dcByNameFuzzy    = acceptance.InitDataSourceCheck(byNameFuzzy)
		byNameExact      = "data.huaweicloud_dns_recordsets.filter_by_name_exact"
		dcByNameExact    = acceptance.InitDataSourceCheck(byNameExact)
		byNotFoundName   = "data.huaweicloud_dns_recordsets.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byStatus           = "data.huaweicloud_dns_recordsets.filter_by_status"
		dcByStatus         = acceptance.InitDataSourceCheck(byStatus)
		byNotFoundStatus   = "data.huaweicloud_dns_recordsets.filter_by_not_found_status"
		dcByNotFoundStatus = acceptance.InitDataSourceCheck(byNotFoundStatus)

		byType           = "data.huaweicloud_dns_recordsets.filter_by_type"
		dcByType         = acceptance.InitDataSourceCheck(byType)
		byNotFoundType   = "data.huaweicloud_dns_recordsets.filter_by_not_found_type"
		dcByNotFoundType = acceptance.InitDataSourceCheck(byNotFoundType)

		byTags           = "data.huaweicloud_dns_recordsets.filter_by_tags"
		dcByTags         = acceptance.InitDataSourceCheck(byTags)
		byNotFoundTags   = "data.huaweicloud_dns_recordsets.filter_by_not_found_tags"
		dcByNotFoundTags = acceptance.InitDataSourceCheck(byNotFoundTags)

		bySortAsc    = "data.huaweicloud_dns_recordsets.filter_by_sort_asc"
		dcBySortAsc  = acceptance.InitDataSourceCheck(bySortAsc)
		bySortDesc   = "data.huaweicloud_dns_recordsets.filter_by_sort_desc"
		dcBySortDesc = acceptance.InitDataSourceCheck(bySortDesc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataRecordsets_notFound(),
				ExpectError: regexp.MustCompile(`This zone does not exist`),
			},
			{
				Config: testAccDataRecordsets_private(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "recordsets.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by recordset ID.
					dcByRecordsetId.CheckResourceExists(),
					resource.TestCheckOutput("is_recordset_id_filter_useful", "true"),
					dcByNotFoundRecordsetId.CheckResourceExists(),
					resource.TestCheckOutput("recordset_id_not_found_validation_pass", "true"),
					// Fuzzy search by recordset name.
					dcByNameFuzzy.CheckResourceExists(),
					resource.TestCheckOutput("is_name_fuzzy_filter_useful", "true"),
					// Exactly search by recordset name.
					dcByNameExact.CheckResourceExists(),
					resource.TestCheckOutput("is_name_exact_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("name_not_found_validation_pass", "true"),
					// Filter by recordset status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					dcByNotFoundStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_not_found_validation_pass", "true"),
					// Filter by recordset type.
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByNotFoundType.CheckResourceExists(),
					resource.TestCheckOutput("type_not_found_validation_pass", "true"),
					// Filter by recordset tags.
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByNotFoundTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_not_found_validation_pass", "true"),
					// Check the sort results.
					dcBySortAsc.CheckResourceExists(),
					dcBySortDesc.CheckResourceExists(),
					resource.TestCheckOutput("sort_filter_is_useful", "true"),
					// Check attributes.
					// The ID of the corresponding resource consists of zone ID and recordset ID.
					resource.TestCheckResourceAttrSet(byRecordsetId, "recordsets.0.id"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.name", rName, "name"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.zone_id", rName, "zone_id"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.zone_name", rName, "zone_name"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.type", rName, "type"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.ttl", rName, "ttl"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.records", rName, "records"),
					resource.TestCheckResourceAttrPair(byRecordsetId, "recordsets.0.description", rName, "description"),
					resource.TestMatchResourceAttr(byRecordsetId, "recordsets.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byRecordsetId, "recordsets.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// This value is false when created by the Terraform script.
					resource.TestCheckResourceAttr(byRecordsetId, "recordsets.0.default", "false"),
				),
			},
		},
	})
}

func testAccDataRecordsets_private(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_dns_zone" "test" {
  name      = "%[2]s"
  zone_type = "private"

  router {
    router_id = huaweicloud_vpc.test.id
  }
}

resource "huaweicloud_dns_recordset" "test" {
  count =2

  zone_id     = huaweicloud_dns_zone.test.id
  name        = "${count.index}test.%[2]s"
  type        = "A"
  description = "Created a recordset by script"
  ttl         = 600
  records     = ["10.1.0.3"]

  tags = {
    foo = "bar_private"
  }
}

%[3]s
`, acceptance.RandomAccResourceName(), name, testAccDataRecordsets_base())
}
