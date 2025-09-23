package lts

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getHostAccessConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	listHostAccessConfigHttpUrl := "v3/{project_id}/lts/access-config-list"
	listHostAccessConfigPath := ltsClient.Endpoint + listHostAccessConfigHttpUrl
	listHostAccessConfigPath = strings.ReplaceAll(listHostAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listHostAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := state.Primary.Attributes["name"]
	listHostAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}

	listHostAccessConfigResp, err := ltsClient.Request("POST", listHostAccessConfigPath, &listHostAccessConfigOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving host access config: %s", err)
	}

	listHostAccessConfigRespBody, err := utils.FlattenResponse(listHostAccessConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten host access config response: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listHostAccessConfigRespBody = utils.PathSearch(jsonPath, listHostAccessConfigRespBody, nil)
	if listHostAccessConfigRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return listHostAccessConfigRespBody, nil
}

func TestAccHostAccessConfig_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		hostAccess interface{}
		rName      = "huaweicloud_lts_host_access.test"
		rc         = acceptance.InitResourceCheck(
			rName,
			&hostAccess,
			getHostAccessConfigResourceFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// Check required Parameter.
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "2"),
					// Check optional Parameter.
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.#", "2"),
					resource.TestCheckResourceAttr(rName, "access_config.0.repeat_collect", "false"),
					resource.TestCheckResourceAttr(rName, "access_config.0.custom_key_value.%", "1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.custom_key_value.flag", "terraform"),
					resource.TestCheckResourceAttr(rName, "access_config.0.system_fields.0", "pathFile"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "0"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "2"),
					resource.TestCheckResourceAttr(rName, "processor_type", "SPLIT"),
					resource.TestCheckResourceAttr(rName, "binary_collect", "true"),
					resource.TestCheckResourceAttr(rName, "encoding_format", "GBK"),
					resource.TestCheckResourceAttr(rName, "incremental_collect", "false"),
					resource.TestCheckResourceAttr(rName, "log_split", "true"),
					// Check attributes.
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestMatchResourceAttr(rName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testHostAccessConfig_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					// Check required Parameter.
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.#", "1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/log/*/*.log"),
					// Check optional Parameter.
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/log/*/a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.repeat_collect", "true"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value-updated"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttrPair(rName, "host_group_ids.0", "huaweicloud_lts_host_group.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "demo_log"),
					resource.TestCheckResourceAttr(rName, "demo_fields.#", "1"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.name", "field1"),
					resource.TestCheckResourceAttr(rName, "demo_fields.0.value", "level:warn1"),
					resource.TestCheckResourceAttr(rName, "processor_type", "SPLIT"),
					resource.TestCheckResourceAttr(rName, "encoding_format", "UTF-8"),
					resource.TestCheckResourceAttr(rName, "incremental_collect", "true"),
					resource.TestCheckResourceAttr(rName, "log_split", "false"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccHostAccessConfigImportStateFunc(rName),
				ImportStateVerifyIgnore: []string{"processors"},
			},
		},
	})
}

func TestAccHostAccessConfig_windows(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_host_access.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getHostAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testHostAccessConfig_windows_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					// Check required parameters.
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "D:\\data\\log\\*"),
					// Check optional parameters.
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "D:\\data\\log\\a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset", "7"),
					resource.TestCheckResourceAttr(rName, "access_config.0.windows_log_info.0.time_offset_unit", "day"),
					resource.TestCheckResourceAttr(rName, "host_group_ids.#", "1"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "binary_collect", "false"),
					resource.TestCheckResourceAttr(rName, "encoding_format", "UTF-8"),
					resource.TestCheckResourceAttr(rName, "incremental_collect", "true"),
					resource.TestCheckResourceAttr(rName, "log_split", "false"),
					// Check attributes.
					resource.TestCheckResourceAttr(rName, "access_type", "AGENT"),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccHostAccessConfigImportStateFunc(rName),
			},
		},
	})
}

func testAccHostAccessConfigImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		name := rs.Primary.Attributes["name"]
		if name == "" {
			return "", fmt.Errorf("can not find the 'name' parameter in %s", rName)
		}
		return name, nil
	}
}

func testHostAccessConfig_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
}
`, name)
}

func testHostAccessConfig_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "huaweicloud_lts_host_group" "test" {
  name = "%[2]s"
  type = "linux"
}

resource "huaweicloud_lts_host_access" "test" {
  name          = "%[2]s"
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id

  access_config {
    paths          = ["/var/temp", "/var/log/*"]
    black_paths    = ["/var/temp", "/var/log/*/a.log"]
    repeat_collect = false

    single_log_format {
      mode = "system"
    }

    custom_key_value = {
      "flag": "terraform"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }

  demo_log = "2024-10-11 10:59:07.000 a.log:1 level:warn|error"

  demo_fields {
   name  = "field2"
   value = "error"
  } 
  demo_fields {
   name  = "field1"
   value = "2024-10-11-10:59:07.000 a.log:1 level:warn"
  }

  processor_type = "SPLIT"

  processors {
    type   = "processor_filter_regex"
    detail = jsonencode({
      "include": {
        "name1": "^terraform"
      },
      "exclude": {
        "black": "test"
      }
    })
  }
  processors {
    type   = "processor_split_string"
    detail = jsonencode({
      "split_sep": "|",
      "keys": ["field1", "field2"],
      "keep_source": false,
      "keep_source_if_parse_error": false
    })
  }

  binary_collect      = true
  encoding_format     = "GBK"
  incremental_collect = false
  log_split           = true
}
`, testHostAccessConfig_base(name), name)
}

func testHostAccessConfig_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "huaweicloud_lts_host_group" "test" {
  name = "%[2]s"
  type = "linux"
}

resource "huaweicloud_lts_host_access" "test" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]

  access_config {
    paths       = ["/var/log/*/*.log"]
    black_paths = ["/var/log/*/a.log"]

    multi_log_format {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }

    custom_key_value = {
      "flag": "terraform"
    }
  }

  tags = {
    key   = "value-updated"
    owner = "terraform"
  }

  demo_log = "2024-10-11-10:59:07.000 level:warn1"

  demo_fields {
   name  = "field1"
   value = "level:warn1"
  }

  processor_type = "SPLIT"

  processors {
    type = "processor_split_string"
    detail = jsonencode({
      "split_sep":" ",
      "keys": ["field1"],
      "keep_source": true,
      "keep_source_if_parse_error": true
    })
  }

  encoding_format  = "UTF-8"
  binary_collect   = true
}
`, testHostAccessConfig_base(name), updateName)
}

func testHostAccessConfig_windows_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# A host group without hosts
resource "huaweicloud_lts_host_group" "test" {
  name = "%[2]s"
  type = "windows"
}

resource "huaweicloud_lts_host_access" "test" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]

  access_config {
    paths       = ["D:\\data\\log\\*"]
    black_paths = ["D:\\data\\log\\a.log"]

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testHostAccessConfig_base(name), name)
}
