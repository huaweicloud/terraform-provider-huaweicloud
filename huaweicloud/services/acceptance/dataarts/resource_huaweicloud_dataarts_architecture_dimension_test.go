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

func getArchitectureDimensionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("dataarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts client: %s", err)
	}

	return dataarts.GetArchitectureDimensionById(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccArchitectureDimension_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_dataarts_architecture_dimension.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getArchitectureDimensionFunc)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
			acceptance.TestAccPreCheckDataArtsArchitectureReviewer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.14.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureDimension_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name_ch", name),
					resource.TestCheckResourceAttr(rName, "name_en", fmt.Sprintf("dim_%s", name)),
					resource.TestCheckResourceAttr(rName, "dimension_type", "COMMON"),
					resource.TestCheckResourceAttr(rName, "owner", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrPair(rName, "l3_id", "huaweicloud_dataarts_architecture_subject.level3", "id"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_ch", "attr1_ch"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_en", "attr1_en"),
					resource.TestCheckResourceAttr(rName, "attributes.0.data_type", "STRING"),
					resource.TestCheckResourceAttr(rName, "attributes.0.is_primary_key", "true"),
					resource.TestCheckResourceAttr(rName, "attributes.0.ordinal", "1"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.id"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.domain_type"),
					resource.TestCheckResourceAttr(rName, "attributes.1.name_ch", "attr2_ch"),
					resource.TestCheckResourceAttr(rName, "attributes.1.name_en", "attr2_en"),
					resource.TestCheckResourceAttr(rName, "attributes.1.data_type", "BIGINT"),
					resource.TestCheckResourceAttr(rName, "attributes.1.is_primary_key", "false"),
					resource.TestCheckResourceAttr(rName, "attributes.1.ordinal", "2"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.create_by"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.create_time"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.update_time"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.id"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.domain_type"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.create_by"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.create_time"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.update_time"),
					resource.TestMatchResourceAttr(rName, "datasource.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "datasource.0.dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "datasource.0.dw_type", "DLI"),
					resource.TestCheckResourceAttr(rName, "datasource.0.db_name", "default"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.dw_name"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.biz_type"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.biz_id"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "l1_id"),
					resource.TestCheckResourceAttrSet(rName, "l2_id"),
					resource.TestCheckResourceAttrSet(rName, "l1_name"),
					resource.TestCheckResourceAttrSet(rName, "l2_name"),
					resource.TestCheckResourceAttrSet(rName, "l3_name"),
					resource.TestCheckResourceAttrSet(rName, "table_type"),
					resource.TestCheckResourceAttrSet(rName, "model_id"),
					resource.TestMatchResourceAttr(rName, "model.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttr(rName, "model.0.is_physical", "true"),
					resource.TestCheckResourceAttrSet(rName, "model.0.id"),
					resource.TestCheckResourceAttrSet(rName, "model.0.name"),
					resource.TestCheckResourceAttrSet(rName, "model.0.create_by"),
					resource.TestCheckResourceAttrSet(rName, "model.0.create_time"),
					resource.TestCheckResourceAttrSet(rName, "model.0.update_by"),
					resource.TestCheckResourceAttrSet(rName, "model.0.update_time"),
					resource.TestCheckResourceAttrSet(rName, "model.0.type"),
					resource.TestCheckResourceAttrSet(rName, "model.0.table_model_prefix"),
				),
			},
			{
				Config: testAccArchitectureDimension_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "name_ch", updateName),
					resource.TestCheckResourceAttr(rName, "name_en", fmt.Sprintf("dim_%s", updateName)),
					resource.TestCheckResourceAttr(rName, "dimension_type", "COMMON"),
					resource.TestCheckResourceAttr(rName, "owner", acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttrPair(rName, "l3_id", "huaweicloud_dataarts_architecture_subject.level3", "id"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_ch", "attr1_ch_updated"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_en", "attr1_en_updated"),
					resource.TestCheckResourceAttr(rName, "attributes.0.data_type", "STRING"),
					resource.TestCheckResourceAttr(rName, "attributes.0.is_primary_key", "false"),
					resource.TestCheckResourceAttr(rName, "attributes.0.ordinal", "1"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.id"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.domain_type"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.create_by"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.create_time"),
					resource.TestCheckResourceAttrSet(rName, "attributes.0.update_time"),
					resource.TestCheckResourceAttr(rName, "attributes.1.name_ch", "attr2_ch_updated"),
					resource.TestCheckResourceAttr(rName, "attributes.1.name_en", "attr2_en_updated"),
					resource.TestCheckResourceAttr(rName, "attributes.1.data_type", "BIGINT"),
					resource.TestCheckResourceAttr(rName, "attributes.1.is_primary_key", "true"),
					resource.TestCheckResourceAttr(rName, "attributes.1.ordinal", "2"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.id"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.domain_type"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.create_by"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.create_time"),
					resource.TestCheckResourceAttrSet(rName, "attributes.1.update_time"),
					resource.TestCheckResourceAttr(rName, "datasource.#", "1"),
					resource.TestCheckResourceAttr(rName, "datasource.0.dw_id", acceptance.HW_DATAARTS_CONNECTION_ID),
					resource.TestCheckResourceAttr(rName, "datasource.0.dw_type", "DLI"),
					resource.TestCheckResourceAttr(rName, "datasource.0.db_name", "default"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.dw_name"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.biz_type"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.biz_id"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "created_by"),
					resource.TestMatchResourceAttr(rName, "create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(rName, "update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(rName, "l1_id"),
					resource.TestCheckResourceAttrSet(rName, "l2_id"),
					resource.TestCheckResourceAttrSet(rName, "l1_name"),
					resource.TestCheckResourceAttrSet(rName, "l2_name"),
					resource.TestCheckResourceAttrSet(rName, "l3_name"),
					resource.TestCheckResourceAttrSet(rName, "table_type"),
					resource.TestCheckResourceAttrSet(rName, "model_id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccArchitectureDimensionImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{
					"is_delete_physical_table",
				},
			},
		},
	})
}

func testAccArchitectureDimensionImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found in state", rName)
		}

		workspaceId := rs.Primary.Attributes["workspace_id"]
		dimensionId := rs.Primary.ID
		if workspaceId == "" || dimensionId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, dimensionId)
		}
		return fmt.Sprintf("%s/%s", workspaceId, dimensionId), nil
	}
}

func testAccArchitectureDimension_base(name string) string {
	return fmt.Sprintf(`
# Need to create the parent subject first, and then create the child subject
resource "huaweicloud_dataarts_architecture_subject" "level1" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  code         = "%[1]s"
  owner        = "%[3]s"
  level        = 1
  description  = "level 1 created by terraform acc test"
}

resource "huaweicloud_dataarts_architecture_subject" "level2" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  code         = "%[1]s"
  owner        = "%[3]s"
  level        = 2
  parent_id    = huaweicloud_dataarts_architecture_subject.level1.id
  description  = "level 2 created by terraform acc test"
}

resource "huaweicloud_dataarts_architecture_subject" "level3" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
  code         = "%[1]s"
  owner        = "%[3]s"
  level        = 3
  department   = "%[1]s"
  parent_id    = huaweicloud_dataarts_architecture_subject.level2.id
  description  = "level 3 created by terraform acc test"
}

# The sub-subject can only be published after the parent subject is published
resource "huaweicloud_dataarts_architecture_batch_publishment" "publish1" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[4]s"
  approver_user_name = "%[3]s"
  fast_approval      = true

  biz_infos {
    biz_id   = huaweicloud_dataarts_architecture_subject.level1.id
    biz_type = "SUBJECT"
  }
}

resource "huaweicloud_dataarts_architecture_batch_publishment" "publish2" {
  workspace_id       = "%[2]s"
  approver_user_id   = "%[4]s"
  approver_user_name = "%[3]s"
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
  workspace_id       = "%[2]s"
  approver_user_id   = "%[4]s"
  approver_user_name = "%[3]s"
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
`, name, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME, acceptance.HW_DATAARTS_ARCHITECTURE_USER_ID)
}

func testAccArchitectureDimension_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = "%[2]s"
  name_ch        = "%[3]s"
  name_en        = "dim_%[3]s"
  l3_id          = huaweicloud_dataarts_architecture_subject.level3.id
  dimension_type = "COMMON"
  owner          = "%[4]s"
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
`, testAccArchitectureDimension_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
	)
}

func testAccArchitectureDimension_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = "%[2]s"
  name_ch        = "%[3]s"
  name_en        = "dim_%[3]s"
  l3_id          = huaweicloud_dataarts_architecture_subject.level3.id
  dimension_type = "COMMON"
  owner          = "%[4]s"
  description    = "Updated by terraform script"

  # delete physical table when deleting the dimension
  is_delete_physical_table = true

  attributes {
    data_type      = "STRING"
    is_primary_key = false
    ordinal        = 1
    name_en        = "attr1_en_updated"
    name_ch        = "attr1_ch_updated"
  }

  attributes {
    name_en        = "attr2_en_updated"
    name_ch        = "attr2_ch_updated"
    data_type      = "BIGINT"
    is_primary_key = true
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
`, testAccArchitectureDimension_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		updateName,
		acceptance.HW_DATAARTS_ARCHITECTURE_USER_NAME,
		acceptance.HW_DATAARTS_CONNECTION_ID,
	)
}
