package dataarts

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureAggregationLogicTables_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameCh   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name_ch"
		dcByNameCh = acceptance.InitDataSourceCheck(byNameCh)

		byNameEn   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name_en"
		dcByNameEn = acceptance.InitDataSourceCheck(byNameEn)

		byCreateBy   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_create_by"
		dcByCreateBy = acceptance.InitDataSourceCheck(byCreateBy)

		byApprover   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_approver"
		dcByApprover = acceptance.InitDataSourceCheck(byApprover)

		byOwner   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_owner"
		dcByOwner = acceptance.InitDataSourceCheck(byOwner)

		byStatus   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		bySyncStatus   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_sync_status"
		dcBySyncStatus = acceptance.InitDataSourceCheck(bySyncStatus)

		byBeginTime   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_begin_time"
		dcByBeginTime = acceptance.InitDataSourceCheck(byBeginTime)

		byBizCatalogId   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_biz_catalog_id"
		dcByBizCatalogId = acceptance.InitDataSourceCheck(byBizCatalogId)

		byAutoGenerate   = "data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_auto_generate"
		dcByAutoGenerate = acceptance.InitDataSourceCheck(byAutoGenerate)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
			acceptance.TestAccPreCheckDataArtsArchitectureSecrecyLevelIds(t, 1)
		},
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.9.1",
			},
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataArchitectureAggregationLogicTables_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tables.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'name' parameter fuzzy search.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by 'name_ch' parameter exact match.
					dcByNameCh.CheckResourceExists(),
					resource.TestCheckOutput("is_name_ch_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.id",
						"huaweicloud_dataarts_architecture_aggregation_logic_table.test", "id"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.tb_name",
						"huaweicloud_dataarts_architecture_aggregation_logic_table.test", "tb_name"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.tb_logic_name",
						"huaweicloud_dataarts_architecture_aggregation_logic_table.test", "tb_logic_name"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.dw_type", "DLI"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.owner", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.description", "Created by terraform script"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.alias",
						"huaweicloud_dataarts_architecture_aggregation_logic_table.test", "alias"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.queue_name", acceptance.HW_DATAARTS_DLI_QUEUE_NAME),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.table_attributes.#", "1"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.id"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.name_ch"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.name_en"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.data_type"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.attribute_type"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.is_primary_key"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.is_partition_key"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.not_null"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.description"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.domain_type"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.table_attributes.0.ref_id",
						"huaweicloud_dataarts_architecture_business_metric.test", "id"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.table_attributes.0.stand_row_id",
						"huaweicloud_dataarts_architecture_data_standard.test", "id"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.alias"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.table_attributes.0.secrecy_levels.#", "1"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.secrecy_levels.0.id"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.secrecy_levels.0.uuid"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.secrecy_levels.0.name"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.secrecy_levels.0.slevel"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.table_attributes.0.ordinal"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.table_type", "EXTERNAL"),
					resource.TestCheckResourceAttrPair(byNameCh, "tables.0.obs_location",
						"huaweicloud_dli_table.test", "bucket_location"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.configs"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.dimension_group"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.sql"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.partition_conf"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.dirty_out_switch", "true"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.dirty_out_database"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.dirty_out_prefix"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.dirty_out_suffix"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.self_defined_fields.#", "1"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.self_defined_fields.0.fd_name_ch"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.self_defined_fields.0.fd_name_en"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.self_defined_fields.0.not_null"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.self_defined_fields.0.fd_value"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.tb_logic_guid"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.status", "PUBLISHED"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.model_id"),
					resource.TestCheckResourceAttr(byNameCh, "tables.0.create_by", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestMatchResourceAttr(byNameCh, "tables.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byNameCh, "tables.0.update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.env_type"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.physical_table"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.technical_asset"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.meta_data_link"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.data_quality"),
					resource.TestCheckResourceAttrSet(byNameCh, "tables.0.summary_status"),
					// Filter by 'name_en' parameter exact match.
					dcByNameEn.CheckResourceExists(),
					resource.TestCheckOutput("is_name_en_filter_useful", "true"),
					// Filter by 'create_by' parameter exact match.
					dcByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),
					// Filter by 'approver' parameter.
					dcByApprover.CheckResourceExists(),
					resource.TestCheckOutput("is_approver_filter_useful", "true"),
					// Filter by 'owner' parameter.
					dcByOwner.CheckResourceExists(),
					resource.TestCheckOutput("is_owner_filter_useful", "true"),
					// Filter by 'status' parameter.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					// Filter by 'sync_status' and 'sync_key' parameter.
					dcBySyncStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_sync_status_filter_useful", "true"),
					// Filter by 'begin_time' and 'end_time' parameter.
					dcByBeginTime.CheckResourceExists(),
					resource.TestCheckOutput("is_begin_time_filter_useful", "true"),
					// Filter by 'biz_catalog_id' parameter.
					dcByBizCatalogId.CheckResourceExists(),
					resource.TestCheckOutput("is_biz_catalog_id_filter_useful", "true"),
					// Filter by 'auto_generate' parameter.
					dcByAutoGenerate.CheckResourceExists(),
					resource.TestCheckOutput("is_auto_generate_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureAggregationLogicTables_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  name = "%[1]s"
}

resource "huaweicloud_obs_bucket" "test" {
  bucket = replace("%[1]s", "_", "-")
  acl    = "private"
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "user/data/user.csv"
  content      = "Jason,Tokyo"
  content_type = "text/plain"
}

resource "huaweicloud_dli_table" "test" {
  database_name   = huaweicloud_dli_database.test.name
  name            = "%[1]s"
  data_location   = "OBS"
  data_format     = "csv"
  bucket_location = "obs://${huaweicloud_obs_bucket_object.test.bucket}/user/data"

  columns {
    name = "name"
    type = "string"
  }
  columns {
    name         = "address"
    type         = "string"
    is_partition = true
  }
}

resource "huaweicloud_dataarts_architecture_data_standard_template" "test" {
  workspace_id = "%[2]s"

  custom_fields {
    fd_name    = "%[1]s"
    required   = false
    searchable = false
  }
}

resource "huaweicloud_dataarts_architecture_directory" "test" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  type         = "STANDARD_ELEMENT"

  depends_on = [huaweicloud_dataarts_architecture_data_standard_template.test]
}

resource "huaweicloud_dataarts_architecture_data_standard" "test" {
  workspace_id = "%[2]s"
  directory_id = huaweicloud_dataarts_architecture_directory.test.id

  values {
    fd_name  = "nameCh"
    fd_value = "%[1]s"
  }
  values {
    fd_name  = "nameEn"
    fd_value = "%[1]s"
  }
}

resource "huaweicloud_dataarts_architecture_subject" "test" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  code         = "%[1]s"
  owner        = "%[3]s"
  level        = 1
}

resource "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  owner        = "%[3]s"
}

resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  name             = "%[1]s"
  workspace_id     = "%[2]s"
  biz_catalog_id   = huaweicloud_dataarts_architecture_process.test.id
  owner            = "%[3]s"
  owner_department = "test-department"
  time_filters     = "双周"
  interval_type    = "HOUR"
  destination      = "test destination"
  definition       = "test definition"
  expression       = "a+b+c"
}

locals {
  biz_infos = [
    {
      biz_id   = huaweicloud_dataarts_architecture_data_standard.test.id
      biz_type = "STANDARD_ELEMENT"
    },
    {
      biz_id   = huaweicloud_dataarts_architecture_business_metric.test.id
      biz_type = "BIZ_METRIC"
    },
    {
      biz_id   = huaweicloud_dataarts_architecture_subject.test.id
      biz_type = "SUBJECT"
    },
  ]
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "test" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[4]s"
  approver_user_name = "%[3]s"
  fast_approval      = true

  dynamic "biz_infos" {
    for_each = local.biz_infos

    content {
      biz_id   = biz_infos.value.biz_id
      biz_type = biz_infos.value.biz_type
    }
  }
}

resource "huaweicloud_dataarts_architecture_aggregation_logic_table" "test" {
  workspace_id  = "%[2]s"
  tb_name       = "dws_%[1]s"
  tb_logic_name = "%[1]s"
  l3_id         = huaweicloud_dataarts_architecture_subject.test.id
  dw_id         = "%[5]s"
  dw_type       = "DLI"
  db_name       = huaweicloud_dli_database.test.name
  owner         = "%[3]s"
  alias         = "%[1]s_alias"
  queue_name    = "%[6]s"
  table_type    = "EXTERNAL"
  obs_location  = huaweicloud_dli_table.test.bucket_location

  table_attributes {
    name_ch        = "key3_ch"
    name_en        = "key3_en"
    alias          = "key_alias"
    data_type      = "STRING"
    attribute_type = "BIZ_METRIC"
    ref_id         = huaweicloud_dataarts_architecture_business_metric.test.id
    is_primary_key = true
    not_null       = true
    description    = "Created attribute by terraform script"
    stand_row_id   = huaweicloud_dataarts_architecture_data_standard.test.id

    secrecy_levels {
      id = try(split(",", "%[7]s")[0], null)
    }
  }

  configs = jsonencode({
    custom = "value"
  })

  description        = "Created by terraform script"
  sql                = "select * from ${huaweicloud_dli_table.test.name}"
  partition_conf     = "name is not null"
  dirty_out_switch   = true
  dirty_out_database = huaweicloud_dli_database.test.name
  dirty_out_prefix   = "tf"
  dirty_out_suffix   = "test"
  secret_type        = "PUBLIC"

  self_defined_fields {
    fd_name_ch = "key_ch1"
    fd_name_en = "key_en1"
    not_null   = true
    fd_value   = "value"
  }

  depends_on = [huaweicloud_dataarts_architecture_batch_publishment.test]
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "aggregation_logic_table" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[4]s"
  approver_user_name = "%[3]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_aggregation_logic_table.test.id
    biz_type = "AGGREGATION_LOGIC_TABLE"
  }
}

resource "time_sleep" "wait_after_aggregation_logic_table_publishment" {
  depends_on      = [huaweicloud_dataarts_architecture_batch_publishment.aggregation_logic_table]
  create_duration = "60s"
}
`,
		name,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME,
		acceptance.HW_DATAARTS_ARCHITECTURE_SECRECY_LEVEL_IDS,
	)
}

func testAccDataArchitectureAggregationLogicTables_basic(name string) string {
	currentTime := time.Now().UTC()
	beginTime := currentTime.Add(-2 * time.Minute).Format("2006-01-02T15:04:05.000Z")
	endTime := currentTime.Add(3 * time.Minute).Format("2006-01-02T15:04:05.000Z")
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_architecture_aggregation_logic_table.test]
}

# Filter by 'name' parameter fuzzy search.
locals {
  name_ch = huaweicloud_dataarts_architecture_aggregation_logic_table.test.tb_logic_name
  name    = substr(local.name_ch, 1, 4)
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.name

  depends_on = [huaweicloud_dataarts_architecture_batch_publishment.aggregation_logic_table]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name.tables[*].tb_logic_name : strcontains(v, local.name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by 'name_ch' parameter exact match.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_name_ch" {
  workspace_id = "%[2]s"
  name_ch      = local.name_ch

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

locals {
  name_ch_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name_ch.tables[*].tb_logic_name : v == local.name_ch
  ]
}

output "is_name_ch_filter_useful" {
  value = length(local.name_ch_filter_result) > 0 && alltrue(local.name_ch_filter_result)
}

# Filter by 'name_en' parameter exact match.
locals {
  name_en = huaweicloud_dataarts_architecture_aggregation_logic_table.test.tb_name
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_name_en" {
  workspace_id = "%[2]s"
  name_en      = local.name_en

  depends_on = [huaweicloud_dataarts_architecture_aggregation_logic_table.test]
}

locals {
  name_en_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_name_en.tables[*].tb_name : v == local.name_en
  ]
}

output "is_name_en_filter_useful" {
  value = length(local.name_en_filter_result) > 0 && alltrue(local.name_en_filter_result)
}

# Filter by 'create_by' parameter exact match.
locals {
  create_by = huaweicloud_dataarts_architecture_aggregation_logic_table.test.created_by
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_create_by" {
  workspace_id = "%[2]s"
  create_by    = local.create_by

  depends_on = [huaweicloud_dataarts_architecture_aggregation_logic_table.test]
}

locals {
  create_by_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_create_by.tables[*].create_by : v == local.create_by
  ]
}

output "is_create_by_filter_useful" {
  value = length(local.create_by_filter_result) > 0 && alltrue(local.create_by_filter_result)
}

# Filter by 'approver' parameter.
locals {
  approver = "%[3]s"
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_approver" {
  workspace_id = "%[2]s"
  approver     = local.approver

  depends_on = [huaweicloud_dataarts_architecture_batch_publishment.aggregation_logic_table]
}

locals {
  approver_filter_result = [
    for v in try(flatten(data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_approver.tables[*].approval_info[*])[*], []) :
    v.approver == local.approver
  ]
}

output "is_approver_filter_useful" {
  value = length(local.approver_filter_result) > 0 && alltrue(local.approver_filter_result)
}

# Filter by 'owner' parameter.
locals {
  owner = huaweicloud_dataarts_architecture_aggregation_logic_table.test.owner
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_owner" {
  workspace_id = "%[2]s"
  owner        = local.owner

  depends_on = [huaweicloud_dataarts_architecture_aggregation_logic_table.test]
}

locals {
  owner_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_owner.tables[*].owner : v == local.owner
  ]
}

output "is_owner_filter_useful" {
  value = length(local.owner_filter_result) > 0 && alltrue(local.owner_filter_result)
}

# Filter by 'status' parameter.
locals {
  status = "PUBLISHED"
}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_status" {
  workspace_id = "%[2]s"
  status       = local.status

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_status.tables[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by 'sync_status' and 'sync_key' parameter.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_sync_status" {
  workspace_id = "%[2]s"
  sync_status  = "SUMMARY_SUCCESS"
  sync_key     = ["BUSINESS_ASSET"]

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

locals {
  sync_status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_sync_status.tables[*].business_asset :
    v == "CREATE_SUCCESS"
  ]
}

output "is_sync_status_filter_useful" {
  value = length(local.sync_status_filter_result) > 0 && alltrue(local.sync_status_filter_result)
}

# Filter by 'begin_time' and 'end_time' parameter.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_begin_time" {
  workspace_id = "%[2]s"
  begin_time   = "%[4]s"
  end_time     = "%[5]s"

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

locals {
  begin_time_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_begin_time.tables[*].update_time :
    timecmp("%[4]s", formatdate("YYYY-MM-DD'T'HH:mm:ss'Z'", timeadd(v, "-8h"))) <= 0 &&
    timecmp(formatdate("YYYY-MM-DD'T'HH:mm:ss'Z'", timeadd(v, "-8h")), "%[5]s") <= 0
  ]
}

output "is_begin_time_filter_useful" {
  value = length(local.begin_time_filter_result) > 0 && alltrue(local.begin_time_filter_result)
}

# Filter by 'biz_catalog_id' parameter.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_biz_catalog_id" {
  workspace_id   = "%[2]s"
  biz_catalog_id = huaweicloud_dataarts_architecture_aggregation_logic_table.test.l3_id

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

output "is_biz_catalog_id_filter_useful" {
  value = length(data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_biz_catalog_id.tables) > 0
}

# Filter by 'auto_generate' parameter.
data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "filter_by_auto_generate" {
  workspace_id  = "%[2]s"
  auto_generate = false

  depends_on = [time_sleep.wait_after_aggregation_logic_table_publishment]
}

output "is_auto_generate_filter_useful" {
  value = length(data.huaweicloud_dataarts_architecture_aggregation_logic_tables.filter_by_auto_generate.tables) > 0
}
`, testAccDataArchitectureAggregationLogicTables_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		beginTime,
		endTime,
	)
}
