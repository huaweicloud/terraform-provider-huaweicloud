package lts

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getStructConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region      = acceptance.HW_REGION_NAME
		httpUrl     = "v2/{project_id}/lts/struct/template"
		product     = "lts"
		logGroupId  = state.Primary.Attributes["log_group_id"]
		logStreamId = state.Primary.Attributes["log_stream_id"]
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += fmt.Sprintf("?logGroupId=%s&logStreamId=%s", logGroupId, logStreamId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	rawString, isString := getRespBody.(string)
	if !isString {
		return nil, fmt.Errorf("the detail API response is not string")
	}

	if rawString == "" {
		// the structuring configuration is not exist
		return nil, golangsdk.ErrDefault404{}
	}
	return getRespBody, nil
}

func TestAccStructConfig_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_structing_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStructConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testStructConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "CTS"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_fields.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				Config: testStructConfig_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "GAUSSDB_MYSQL_SLOW"),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "true"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_fields.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStructConfigImportState(rName),
				ImportStateVerifyIgnore: []string{
					"template_type",
					"template_id",
					"demo_fields",
					"tag_fields",
					"quick_analysis",
				},
			},
		},
	})
}

func TestAccStructConfig_customTemplate(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_structing_template.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getStructConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsStructConfigCustom(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testStructConfig_custom(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "template_type", "custom"),
					resource.TestCheckResourceAttr(rName, "template_id", acceptance.HW_LTS_STRUCT_CONFIG_TEMPLATE_ID),
					resource.TestCheckResourceAttr(rName, "template_name", acceptance.HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "true"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				Config: testStructConfig_custom_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "log_group_id",
						"huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id",
						"huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "template_type", "built_in"),
					resource.TestCheckResourceAttr(rName, "template_name", "GAUSSDB_MYSQL_SLOW"),
					resource.TestCheckResourceAttr(rName, "quick_analysis", "false"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "tag_fields.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStructConfigImportState(rName),
				ImportStateVerifyIgnore: []string{
					"template_type",
					"template_id",
					"demo_fields",
					"tag_fields",
					"quick_analysis",
				},
			},
		},
	})
}

func testStructConfig_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_structing_template" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  template_type = "built_in"
  template_name = "CTS"

  demo_fields {
    field_name  = "event_type"
    is_analysis = true
  }

  demo_fields {
    field_name  = "resource_type"
    is_analysis = false
  }

  tag_fields {
    field_name  = "hostIP"
    is_analysis = true
  }
}
`, testAccLtsStream_basic(name), name)
}

func testStructConfig_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_structing_template" "test" {
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  template_type  = "built_in"
  template_name  = "GAUSSDB_MYSQL_SLOW"
  quick_analysis = true

  demo_fields {
    field_name  = "query_time"
    is_analysis = true
  }

  demo_fields {
    field_name  = "rows_examined"
    is_analysis = false
  }

  tag_fields {
    field_name  = "hostName"
    is_analysis = false
  }
}
`, testAccLtsStream_basic(name))
}

func testStructConfig_custom(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_structing_template" "test" {
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  template_type  = "custom"
  template_id    = "%[2]s"
  template_name  = "%[3]s"
  quick_analysis = true
}
`, testAccLtsStream_basic(name), acceptance.HW_LTS_STRUCT_CONFIG_TEMPLATE_ID, acceptance.HW_LTS_STRUCT_CONFIG_TEMPLATE_NAME)
}

func testStructConfig_custom_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_structing_template" "test" {
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  template_type  = "built_in"
  template_name  = "GAUSSDB_MYSQL_SLOW"
  quick_analysis = false

  demo_fields {
    field_name  = "query_time"
    is_analysis = true
  }

  demo_fields {
    field_name  = "rows_examined"
    is_analysis = true
  }

  tag_fields {
    field_name  = "hostName"
    is_analysis = true
  }
}
`, testAccLtsStream_basic(name), name)
}

func testLtsStructConfigImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		logGroupId := rs.Primary.Attributes["log_group_id"]
		logStreamId := rs.Primary.Attributes["log_stream_id"]
		if logGroupId == "" || logStreamId == "" {
			return "", fmt.Errorf("invalid format specified for import ID (LTS structuring configuration),"+
				" want '<log_group_id>/<log_stream_id>', but got '%s/%s'",
				logGroupId, logStreamId)
		}
		return fmt.Sprintf("%s/%s", logGroupId, logStreamId), nil
	}
}
