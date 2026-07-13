package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataArchitectureDimensions_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_architecture_dimensions.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNameCh   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_ch"
		dcByNameCh = acceptance.InitDataSourceCheck(byNameCh)

		byNameEn   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_en"
		dcByNameEn = acceptance.InitDataSourceCheck(byNameEn)

		byCreateBy   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_create_by"
		dcByCreateBy = acceptance.InitDataSourceCheck(byCreateBy)

		byStatus   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byDimensionType   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_dimension_type"
		dcByDimensionType = acceptance.InitDataSourceCheck(byDimensionType)

		byBeginTime   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_begin_time"
		dcByBeginTime = acceptance.InitDataSourceCheck(byBeginTime)

		byEndTime   = "data.huaweicloud_dataarts_architecture_dimensions.filter_by_end_time"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.14.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDataArchitectureDimensions_nonExistentWorkspace(),
				ExpectError: regexp.MustCompile("error querying DataArts Architecture dimensions"),
			},
			{
				Config: testAccDataArchitectureDimensions_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "dimensions.#", regexp.MustCompile(`^[2-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.id"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.name_ch"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.name_en"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.dimension_type"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.status"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.created_by"),
					resource.TestMatchResourceAttr(all, "dimensions.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "dimensions.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.l1_id"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.l2_id"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.l1_name"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.l2_name"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.l3_name"),
					resource.TestCheckResourceAttrSet(all, "dimensions.0.model_id"),
					resource.TestMatchResourceAttr(all, "dimensions.0.datasource.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestMatchResourceAttr(all, "dimensions.0.attributes.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestMatchResourceAttr(all, "dimensions.0.model.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),

					// Filter by name (fuzzy).
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// Filter by name_ch (exact).
					dcByNameCh.CheckResourceExists(),
					resource.TestCheckOutput("is_name_ch_filter_useful", "true"),

					// Filter by name_en (exact).
					dcByNameEn.CheckResourceExists(),
					resource.TestCheckOutput("is_name_en_filter_useful", "true"),

					// Filter by create_by.
					dcByCreateBy.CheckResourceExists(),
					resource.TestCheckOutput("is_create_by_filter_useful", "true"),

					// Filter by status.
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					// Filter by dimension_type.
					dcByDimensionType.CheckResourceExists(),
					resource.TestCheckOutput("is_dimension_type_filter_useful", "true"),

					// Filter by begin_time.
					dcByBeginTime.CheckResourceExists(),
					resource.TestCheckOutput("is_begin_time_filter_useful", "true"),

					// Filter by end_time.
					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataArchitectureDimensions_nonExistentWorkspace() string {
	randUUID := uuid.New().String()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_architecture_dimensions" "test" {
  workspace_id = "%[1]s"
}
`, randUUID)
}

func testAccDataArchitectureDimension_base(name string) string {
	return fmt.Sprintf(`
# Need to create the parent subject first, and then create the child subject
resource "huaweicloud_dataarts_architecture_subject" "level1" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s"
  owner        = "%[3]s"
  level        = 1
  description  = "level 1 created by terraform acc test"
}

resource "huaweicloud_dataarts_architecture_subject" "level2" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s"
  owner        = "%[3]s"
  level        = 2
  parent_id    = huaweicloud_dataarts_architecture_subject.level1.id
  description  = "level 2 created by terraform acc test"
}

resource "huaweicloud_dataarts_architecture_subject" "level3" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  code         = "%[2]s"
  department   = "%[2]s"
  owner        = "%[3]s"
  level        = 3
  parent_id    = huaweicloud_dataarts_architecture_subject.level2.id
  description  = "level 3 created by terraform acc test"
}

# The sub-subject can only be published after the parent subject is published
resource "huaweicloud_dataarts_architecture_batch_publishment" "publish1" {
  workspace_id       = "%[1]s"
  approver_user_name = "%[3]s"
  approver_user_id   = "%[4]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.level1.id
    biz_type = "SUBJECT"
  }
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "publish2" {
  workspace_id       = "%[1]s"
  approver_user_name = "%[3]s"
  approver_user_id   = "%[4]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.level2.id
    biz_type = "SUBJECT"
  }

  depends_on = [
    huaweicloud_dataarts_architecture_batch_publishment.publish1,
  ]
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "publish3" {
  workspace_id       = "%[1]s"
  approver_user_name = "%[3]s"
  approver_user_id   = "%[4]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.level3.id
    biz_type = "SUBJECT"
  }

  depends_on = [
    huaweicloud_dataarts_architecture_batch_publishment.publish2,
  ]
}

# Need to wait for the dimension database in the cloud service backend to obtain the subject publishment information
resource "time_sleep" "wait_10_second" {
  create_duration = "10s"

  depends_on = [
    huaweicloud_dataarts_architecture_batch_publishment.publish3
  ]
}

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  count = 2

  workspace_id   = "%[1]s"
  name_ch        = "dim_%[2]s_${count.index + 1}"
  name_en        = "dim_%[2]s_${count.index + 1}"
  l3_id          = huaweicloud_dataarts_architecture_subject.level3.id
  dimension_type = "COMMON"
  owner          = "%[3]s"
  description    = "Created by terraform script"

  # delete physical table when deleting the dimension
  is_delete_physical_table = true

  attributes {
    name_en        = "attr1_en"
    name_ch        = "attr1_ch"
    data_type      = "STRING"
    is_primary_key = true
    ordinal        = 1
  }

  attributes {
    name_en        = "attr2_en"
    name_ch        = "attr2_ch"
    data_type      = "BIGINT"
    is_primary_key = false
    ordinal        = 2
  }

  datasource {
    dw_id   = "%[5]s"
    dw_type = "DLI"
    db_name = "default"
  }

  depends_on = [
    time_sleep.wait_10_second
  ]
}
`,
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID,
		acceptance.HW_DATAARTS_CONNECTION_ID)
}

func testAccDataArchitectureDimensions_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_architecture_dimensions" "all" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

# Filter by name (fuzzy).
locals {
  dimension_name = huaweicloud_dataarts_architecture_dimension.test[0].name_ch
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_name" {
  workspace_id = "%[2]s"
  name         = local.dimension_name

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

locals {
  name_filter_has_one   = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[0].id)
  name_filter_has_other = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[1].id)
  name_filter_is_useful = local.name_filter_has_one && !local.name_filter_has_other
}

output "is_name_filter_useful" {
  value = local.name_filter_is_useful
}

# Filter by name_ch (exact).
locals {
  dimension_name_ch = huaweicloud_dataarts_architecture_dimension.test[0].name_ch
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_name_ch" {
  workspace_id = "%[2]s"
  name_ch      = local.dimension_name_ch

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

locals {
  name_ch_filter_has_one   = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_ch.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[0].id)
  name_ch_filter_has_other = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_ch.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[1].id)
  is_name_ch_filter_useful = local.name_ch_filter_has_one && !local.name_ch_filter_has_other
}

output "is_name_ch_filter_useful" {
  value = local.is_name_ch_filter_useful
}

# Filter by name_en (exact).
locals {
  dimension_name_en = try(huaweicloud_dataarts_architecture_dimension.test[0].name_en, "NOT_FOUND")
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_name_en" {
  workspace_id = "%[2]s"
  name_en      = local.dimension_name_en

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

locals {
  name_en_filter_has_one   = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_en.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[0].id)
  name_en_filter_has_other = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_en.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[1].id)
  is_name_en_filter_useful = local.name_en_filter_has_one && !local.name_en_filter_has_other
}

output "is_name_en_filter_useful" {
  value = local.is_name_en_filter_useful
}

# Filter by create_by.
locals {
  dimension_creator = try(huaweicloud_dataarts_architecture_dimension.test[0].created_by, "NOT_FOUND")
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_create_by" {
  workspace_id = "%[2]s"
  create_by    = local.dimension_creator

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

output "is_create_by_filter_useful" {
  value = length(data.huaweicloud_dataarts_architecture_dimensions.filter_by_create_by.dimensions) >= 1
}

# Filter by status.
locals {
  dimension_status = try(huaweicloud_dataarts_architecture_dimension.test[0].status, "NOT_FOUND")
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_status" {
  workspace_id = "%[2]s"
  status       = local.dimension_status

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dataarts_architecture_dimensions.filter_by_status.dimensions :
      v.status == local.dimension_status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by dimension_type.
data "huaweicloud_dataarts_architecture_dimensions" "filter_by_dimension_type" {
  workspace_id    = "%[2]s"
  dimension_type  = "COMMON"

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

locals {
  type_filter_has_one   = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_en.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[0].id)
  type_filter_has_other = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_name_en.dimensions[*].id,
  huaweicloud_dataarts_architecture_dimension.test[1].id)
  type_filter_useful    = local.type_filter_has_one && !local.type_filter_has_other
}

output "is_dimension_type_filter_useful" {
  value = local.type_filter_useful
}

# Filter by begin_time.
locals {
  dimension_begin_time = timeadd(huaweicloud_dataarts_architecture_dimension.test[0].create_time, "-10m")
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_begin_time" {
  workspace_id = "%[2]s"
  begin_time   = local.dimension_begin_time

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

output "is_begin_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_begin_time.dimensions[*].id,
    huaweicloud_dataarts_architecture_dimension.test[0].id)
}

# Filter by end_time.
locals {
  dimension_end_time = timeadd(huaweicloud_dataarts_architecture_dimension.test[0].create_time, "10m")
}

data "huaweicloud_dataarts_architecture_dimensions" "filter_by_end_time" {
  workspace_id = "%[2]s"
  end_time     = local.dimension_end_time

  depends_on = [huaweicloud_dataarts_architecture_dimension.test]
}

output "is_end_time_filter_useful" {
  value = contains(data.huaweicloud_dataarts_architecture_dimensions.filter_by_end_time.dimensions[*].id,
    huaweicloud_dataarts_architecture_dimension.test[0].id)
}
`, testAccDataArchitectureDimension_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
