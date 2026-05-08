package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getArchitectureAggregationLogicTableFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}

	return dataarts.GetArchitectureAggregationLogicTableById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccArchitectureAggregationLogicTable_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		obj   interface{}
		rName = "huaweicloud_dataarts_architecture_aggregation_logic_table.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getArchitectureAggregationLogicTableFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
			acceptance.TestAccPreCheckDataArtsRelatedDliQueueName(t)
			acceptance.TestAccPreCheckDataArtsArchitectureSecrecyLevelIds(t, 2)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureAggregationLogicTable_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "tb_name", fmt.Sprintf("dws_%s", name)),
					resource.TestCheckResourceAttr(rName, "tb_logic_name", name),
					resource.TestCheckResourceAttr(rName, "dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "dw_type", "DLI"),
					resource.TestCheckResourceAttrPair(rName, "db_name", "huaweicloud_dli_database.test.0", "name"),
					resource.TestCheckResourceAttr(rName, "owner", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestCheckResourceAttrPair(rName, "l3_id", "huaweicloud_dataarts_architecture_subject.test.0", "id"),
					resource.TestCheckResourceAttr(rName, "alias", fmt.Sprintf("%s_alias", name)),
					resource.TestCheckResourceAttr(rName, "queue_name", acceptance.HW_DATAARTS_DLI_QUEUE_NAME),
					resource.TestCheckResourceAttr(rName, "table_type", "MANAGED"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(rName, "sql", fmt.Sprintf("SELECT * FROM %s", name)),
					resource.TestCheckResourceAttr(rName, "partition_conf", "name is not null"),
					resource.TestCheckResourceAttr(rName, "dirty_out_switch", "true"),
					resource.TestCheckResourceAttrPair(rName, "dirty_out_database", "huaweicloud_dli_database.test.1", "name"),
					resource.TestCheckResourceAttr(rName, "dirty_out_prefix", name),
					resource.TestCheckResourceAttr(rName, "dirty_out_suffix", name),
					resource.TestCheckResourceAttr(rName, "secret_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet(rName, "configs"),
					resource.TestCheckResourceAttr(rName, "table_attributes.#", "3"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.name_ch", "key3_ch"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.name_en", "key3_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.data_type", "STRING"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.name_ch", "key1_ch"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.name_en", "key1_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.data_type", "DATE"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.attribute_type", "SUMMARY_TIME"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.is_primary_key", "true"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.is_partition_key", "true"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.not_null", "true"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.description", "Created attribute by terraform script"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.alias", "key1_alias"),
					resource.TestCheckResourceAttrPair(rName, "table_attributes.1.stand_row_id",
						"huaweicloud_dataarts_architecture_data_standard.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.id"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.secrecy_levels.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.secrecy_levels.0.uuid"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.secrecy_levels.0.name"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.domain_type"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.stand_row_name"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.1.ordinal"),
					resource.TestCheckResourceAttr(rName, "table_attributes.2.name_ch", "key2_ch"),
					resource.TestCheckResourceAttr(rName, "table_attributes.2.name_en", "key2_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.2.data_type", "BIGINT"),
					resource.TestCheckResourceAttr(rName, "table_attributes.2.attribute_type", "BIZ_METRIC"),
					resource.TestCheckResourceAttrPair(rName, "table_attributes.2.ref_id",
						"huaweicloud_dataarts_architecture_business_metric.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.2.ref_name_ch"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.2.ref_name_en"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.fd_name_ch", "key_ch2"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.fd_name_en", "key_en2"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.1.fd_name_ch", "key_ch1"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.1.fd_name_en", "key_en1"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.1.not_null", "true"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.1.fd_value", "value"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(rName, "dw_name"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccArchitectureAggregationLogicTable_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "tb_name", fmt.Sprintf("dws_%s", updateName)),
					resource.TestCheckResourceAttr(rName, "tb_logic_name", updateName),
					resource.TestCheckResourceAttr(rName, "dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "dw_type", "DLI"),
					resource.TestCheckResourceAttrPair(rName, "db_name", "huaweicloud_dli_database.test.1", "name"),
					resource.TestCheckResourceAttr(rName, "owner", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestCheckResourceAttrPair(rName, "l3_id", "huaweicloud_dataarts_architecture_subject.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "alias", fmt.Sprintf("%s_alias", updateName)),
					resource.TestCheckResourceAttr(rName, "queue_name", acceptance.HW_DATAARTS_DLI_QUEUE_NAME),
					resource.TestCheckResourceAttr(rName, "table_type", "EXTERNAL"),
					resource.TestCheckResourceAttrPair(rName, "obs_location", "huaweicloud_dli_table.test", "bucket_location"),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(rName, "sql", fmt.Sprintf("SELECT name FROM %s", name)),
					resource.TestCheckResourceAttr(rName, "partition_conf", "key1_en is not null"),
					resource.TestCheckResourceAttr(rName, "dirty_out_switch", "false"),
					resource.TestCheckResourceAttr(rName, "dirty_out_database", ""),
					resource.TestCheckResourceAttr(rName, "dirty_out_prefix", ""),
					resource.TestCheckResourceAttr(rName, "dirty_out_suffix", ""),
					resource.TestCheckResourceAttr(rName, "secret_type", "PUBLIC"),
					resource.TestCheckResourceAttrSet(rName, "configs"),
					resource.TestCheckResourceAttr(rName, "table_attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.name_ch", "key2_ch"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.name_en", "key2_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.data_type", "STRING"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.attribute_type", "BIZ_METRIC"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.is_primary_key", "true"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.is_partition_key", "true"),
					resource.TestCheckResourceAttr(rName, "table_attributes.0.not_null", "true"),
					resource.TestCheckResourceAttrPair(rName, "table_attributes.0.stand_row_id",
						"huaweicloud_dataarts_architecture_data_standard.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "table_attributes.0.ref_id",
						"huaweicloud_dataarts_architecture_business_metric.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.0.stand_row_name"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.0.ref_name_ch"),
					resource.TestCheckResourceAttrSet(rName, "table_attributes.0.ref_name_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.name_ch", "key1_ch"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.name_en", "key1_en"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.data_type", "DATE"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.attribute_type", "SUMMARY_TIME"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.alias", "key1_alias"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.is_primary_key", "false"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.is_partition_key", "false"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.not_null", "false"),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.description", ""),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.stand_row_id", ""),
					resource.TestCheckResourceAttr(rName, "table_attributes.1.secrecy_levels.#", "0"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.#", "1"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.fd_name_ch", "key_ch1"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.fd_name_en", "key_en1"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.not_null", "true"),
					resource.TestCheckResourceAttr(rName, "self_defined_fields.0.fd_value", "value"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccArchitectureAggregationLogicTableImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{"table_attributes_origin", "l3_id", "secret_type", "del_type"},
			},
		},
	})
}

func testAccArchitectureAggregationLogicTableImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found in state", rName)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		aggregationLogicTableId := rs.Primary.ID
		if workspaceId == "" || aggregationLogicTableId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, aggregationLogicTableId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, aggregationLogicTableId), nil
	}
}

func testAccArchitectureAggregationLogicTable_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dli_database" "test" {
  count = 2

  name = "%[1]s${count.index}"
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
  database_name   = try(huaweicloud_dli_database.test[0].name, "NOT_FOUND")
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
  count = 2

  workspace_id = "%[2]s"
  name         = "%[1]s${count.index}"
  code         = "%[1]s${count.index}"
  owner        = "%[4]s"
  level        = 1
}

resource "huaweicloud_dataarts_architecture_process" "test" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  owner        = "%[4]s"
}

resource "huaweicloud_dataarts_architecture_business_metric" "test" {
  name             = "%[1]s"
  workspace_id     = "%[2]s"
  biz_catalog_id   = huaweicloud_dataarts_architecture_process.test.id
  owner            = "%[4]s"
  owner_department = "test-department"
  time_filters     = "双周"
  interval_type    = "HOUR"
  destination      = "test destination"
  definition       = "test definition"
  expression       = "a+b+c"
}

locals {
  biz_infos = concat(
    [
      {
        biz_id   = huaweicloud_dataarts_architecture_data_standard.test.id
        biz_type = "STANDARD_ELEMENT"
      },
      {
        biz_id   = huaweicloud_dataarts_architecture_business_metric.test.id
        biz_type = "BIZ_METRIC"
      }
    ],
    [
      for id in huaweicloud_dataarts_architecture_subject.test[*].id : {
        biz_id   = id
        biz_type = "SUBJECT"
      }
    ]
  )
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "test" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[3]s"
  approver_user_name = "%[4]s"
  fast_approval      = true

  dynamic "biz_infos" {
    for_each = local.biz_infos
  
    content {
      biz_id   = biz_infos.value.biz_id
      biz_type = biz_infos.value.biz_type
    }
  }
}
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME)
}

func testAccArchitectureAggregationLogicTable_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_aggregation_logic_table" "test" {
  workspace_id  = "%[2]s"
  tb_name       = "dws_%[3]s"
  tb_logic_name = "%[3]s"
  l3_id         = try(huaweicloud_dataarts_architecture_subject.test[0].id, "NOT_FOUND")
  dw_id         = "%[4]s"
  dw_type       = "DLI"
  db_name       = try(huaweicloud_dli_database.test[0].name, "NOT_FOUND")
  owner         = "%[5]s"
  alias         = "%[3]s_alias"
  queue_name    = "%[6]s"
  table_type    = "MANAGED"

  table_attributes {
    name_ch   = "key3_ch"
    name_en   = "key3_en"
    data_type = "STRING"
  }
  table_attributes {
    name_ch          = "key1_ch"
    name_en          = "key1_en"
    alias            = "key1_alias"
    data_type        = "DATE"
    attribute_type   = "SUMMARY_TIME"
    is_primary_key   = true
    is_partition_key = true
    not_null         = true
    description      = "Created attribute by terraform script"
    stand_row_id     = huaweicloud_dataarts_architecture_data_standard.test.id

	dynamic "secrecy_levels" {
      for_each = reverse(sort(split(",", "%[7]s")))

      content {
        id = secrecy_levels.value
      }
    }
  }
  table_attributes {
    name_ch        = "key2_ch"
    name_en        = "key2_en"
    data_type      = "BIGINT"
    attribute_type = "BIZ_METRIC"
    ref_id         = huaweicloud_dataarts_architecture_business_metric.test.id
  }

  configs = jsonencode({
    custom = "value"
  })

  description        = "Created by terraform script"
  sql                = "SELECT * FROM ${huaweicloud_dli_table.test.name}"
  partition_conf     = "name is not null"
  dirty_out_switch   = true
  dirty_out_database = try(huaweicloud_dli_database.test[1].name, "NOT_FOUND")
  dirty_out_prefix   = "%[3]s"
  dirty_out_suffix   = "%[3]s"
  secret_type        = "PUBLIC"

  self_defined_fields {
    fd_name_ch = "key_ch2"
    fd_name_en = "key_en2"
  }
  self_defined_fields {
    fd_name_ch = "key_ch1"
    fd_name_en = "key_en1"
    not_null   = true
    fd_value   = "value"
  }

  depends_on = [
    huaweicloud_dataarts_architecture_batch_publishment.test,
  ]
}
`, testAccArchitectureAggregationLogicTable_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME,
		acceptance.HW_DATAARTS_ARCHITECTURE_SECRECY_LEVEL_IDS,
	)
}

func testAccArchitectureAggregationLogicTable_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_aggregation_logic_table" "test" {
  workspace_id  = "%[2]s"
  tb_name       = "dws_%[3]s"
  tb_logic_name = "%[3]s"
  l3_id         = try(huaweicloud_dataarts_architecture_subject.test[1].id, "NOT_FOUND")
  dw_id         = "%[4]s"
  dw_type       = "DLI"
  db_name       = try(huaweicloud_dli_database.test[1].name, "NOT_FOUND")
  owner         = "%[6]s"
  alias         = "%[3]s_alias"
  queue_name    = "%[5]s"
  table_type    = "EXTERNAL"
  obs_location  = huaweicloud_dli_table.test.bucket_location

  table_attributes {
    name_ch          = "key2_ch"
    name_en          = "key2_en"
    data_type        = "STRING"
    attribute_type   = "BIZ_METRIC"
    is_primary_key   = true
    is_partition_key = true
    not_null         = true
    stand_row_id     = huaweicloud_dataarts_architecture_data_standard.test.id
    ref_id           = huaweicloud_dataarts_architecture_business_metric.test.id
  }
  table_attributes {
    name_ch        = "key1_ch"
    name_en        = "key1_en"
    alias          = "key1_alias"
    data_type      = "DATE"
    attribute_type = "SUMMARY_TIME"
  }

  configs = jsonencode({
    stage  = 2
  })

  description        = "Updated by terraform script"
  sql                = "SELECT name FROM ${huaweicloud_dli_table.test.name}"
  partition_conf     = "key1_en is not null"
  secret_type        = "PUBLIC"

  self_defined_fields {
    fd_name_ch = "key_ch1"
    fd_name_en = "key_en1"
    not_null   = true
    fd_value   = "value"
  }

  del_type = "PHYSICAL_TABLE"

  depends_on = [
    huaweicloud_dataarts_architecture_batch_publishment.test,
  ]
}
`, testAccArchitectureAggregationLogicTable_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		updateName,
		acceptance.HW_DATAARTS_CONNECTION_ID,
		acceptance.HW_DATAARTS_DLI_QUEUE_NAME,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
	)
}
