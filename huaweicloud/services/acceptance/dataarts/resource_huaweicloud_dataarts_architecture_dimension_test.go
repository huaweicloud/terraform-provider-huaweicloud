package dataarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dataarts"
)

func getArchitectureDimensionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		workspaceId = state.Primary.Attributes["workspace_id"]
		id          = state.Primary.ID
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating DataArts Studio client: %s", err)
	}

	return dataarts.GetArchitectureDimensionById(client, workspaceId, id)
}

func TestAccResourceArchitectureDimension_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_dataarts_architecture_dimension.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getArchitectureDimensionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchitectureDimension_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name_ch", name),
					resource.TestCheckResourceAttr(rName, "name_en", name+"_en"),
					resource.TestCheckResourceAttr(rName, "dimension_type", "DIMENSION"),
					resource.TestCheckResourceAttr(rName, "owner", "test_owner"),
					resource.TestCheckResourceAttr(rName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_ch", "attr1"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_en", "attr1_en"),
					resource.TestCheckResourceAttr(rName, "attributes.0.data_type", "STRING"),
					resource.TestCheckResourceAttr(rName, "attributes.0.is_primary_key", "true"),
					resource.TestCheckResourceAttr(rName, "attributes.0.ordinal", "1"),
					resource.TestCheckResourceAttr(rName, "datasource.0.dw_type", "DWS"),
					resource.TestCheckResourceAttrSet(rName, "datasource.0.dw_id"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "env_type"),
				),
			},
			{
				Config: testAccArchitectureDimension_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", "updated by terraform"),
					resource.TestCheckResourceAttr(rName, "alias", "test_alias"),
					resource.TestCheckResourceAttr(rName, "attributes.#", "2"),
					resource.TestCheckResourceAttr(rName, "attributes.0.name_ch", "attr1_update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testDataArtsStudioImportState(rName),
			},
		},
	})
}

func testAccArchitectureDimension_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_dataarts_architecture_model" "test" {
  workspace_id = "%[1]s"
  name         = "%[2]s"
  type         = "THIRD_NF"
  physical     = true
  dw_type      = "DWS"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, name)
}

func testAccArchitectureDimension_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = "%[2]s"
  name_ch        = "%[3]s"
  name_en        = "%[3]s_en"
  dimension_type = "DIMENSION"
  l3_id          = huaweicloud_dataarts_architecture_model.test.id
  owner          = "test_owner"
  description    = "created by terraform"

  datasource {
    dw_id   = "%[4]s"
    dw_type = huaweicloud_dataarts_architecture_model.test.dw_type
  }

  attributes {
    name_ch        = "attr1"
    name_en        = "attr1_en"
    data_type      = "STRING"
    is_primary_key = true
    ordinal        = "1"
  }

  attributes {
    name_ch        = "attr2"
    name_en        = "attr2_en"
    data_type      = "INTEGER"
    is_primary_key = false
    ordinal        = "2"
  }
}
`,
		testAccArchitectureDimension_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_CONNECTION_ID,
	)
}

func testAccArchitectureDimension_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = "%[2]s"
  name_ch        = "%[3]s"
  name_en        = "%[3]s_en"
  dimension_type = "DIMENSION"
  l3_id          = huaweicloud_dataarts_architecture_model.test.id
  owner          = "test_owner"
  description    = "updated by terraform"
  alias          = "test_alias"

  datasource {
    dw_id   = "%[4]s"
    dw_type = huaweicloud_dataarts_architecture_model.test.dw_type
  }

  attributes {
    name_ch        = "attr1_update"
    name_en        = "attr1_en"
    data_type      = "STRING"
    is_primary_key = true
    ordinal        = "1"
  }

  attributes {
    name_ch        = "attr2"
    name_en        = "attr2_en"
    data_type      = "INTEGER"
    is_primary_key = false
    ordinal        = "2"
  }
}
`,
		testAccArchitectureDimension_base(name),
		acceptance.HW_DATAARTS_WORKSPACE_ID,
		name,
		acceptance.HW_DATAARTS_CONNECTION_ID,
	)
}
