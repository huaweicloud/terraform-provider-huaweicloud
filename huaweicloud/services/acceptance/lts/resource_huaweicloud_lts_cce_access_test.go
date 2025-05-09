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

func getCceAccessConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}

	listCceAccessConfigHttpUrl := "v3/{project_id}/lts/access-config-list"
	listCceAccessConfigPath := ltsClient.Endpoint + listCceAccessConfigHttpUrl
	listCceAccessConfigPath = strings.ReplaceAll(listCceAccessConfigPath, "{project_id}", ltsClient.ProjectID)

	listCceAccessConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := state.Primary.Attributes["name"]
	listCceAccessConfigOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}

	listCceAccessConfigResp, err := ltsClient.Request("POST", listCceAccessConfigPath, &listCceAccessConfigOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCE access config: %s", err)
	}

	listCceAccessConfigRespBody, err := utils.FlattenResponse(listCceAccessConfigResp)
	if err != nil {
		return nil, fmt.Errorf("error flatten CCE access config response: %s", err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listCceAccessConfigRespBody = utils.PathSearch(jsonPath, listCceAccessConfigRespBody, nil)
	if listCceAccessConfigRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return listCceAccessConfigRespBody, nil
}

func TestAccCceAccessConfig_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		accessConfig        interface{}
		withContainerFile   = "huaweicloud_lts_cce_access.with_container_file"
		rcWithContainerFile = acceptance.InitResourceCheck(withContainerFile, &accessConfig, getCceAccessConfigResourceFunc)

		withContainerStdout   = "huaweicloud_lts_cce_access.with_container_stdout"
		rcWithContainerStdout = acceptance.InitResourceCheck(withContainerStdout, &accessConfig, getCceAccessConfigResourceFunc)

		withHostFile   = "huaweicloud_lts_cce_access.with_host_file"
		rcWithHostFile = acceptance.InitResourceCheck(withHostFile, &accessConfig, getCceAccessConfigResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSCCEAccess(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithContainerFile.CheckResourceDestroy(),
			rcWithContainerStdout.CheckResourceDestroy(),
			rcWithHostFile.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCceAccessConfig_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithContainerFile.CheckResourceExists(),
					// Check Required parameters.
					resource.TestCheckResourceAttr(withContainerFile, "name", name),
					resource.TestCheckResourceAttrPair(withContainerFile, "log_group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttrPair(withContainerFile, "log_stream_id", "huaweicloud_lts_stream.test", "id"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.path_type", "container_file"),
					// Check optional parameters.
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.paths.0", "/var"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.black_paths.0", "/var/a.log"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.name_space_regex", "namespace"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.pod_name_regex", "podname"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.container_name_regex", "containername"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_labels.loglabelkey1", "bar1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_labels.loglabelkey2", "bar"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels.includeKey1", "incval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels.includeKey2", "incval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels.excludeKey1", "excval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels.excludeKey2", "excval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_envs.envKey1", "envval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_envs.envKey2", "envval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs.inEnvKey1", "incval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs.inEnvKey2", "incval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs.exEnvKey1", "excval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs.exEnvKey2", "excval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_k8s.k8sKey1", "k8sval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_k8s.k8sKey2", "k8sval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels.ink8sKey1", "ink8sval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels.exk8sKey1", "exk8sval1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels_logical", "or"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.system_fields.#", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.system_fields.0", "pathFile"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.repeat_collect", "true"),
					resource.TestCheckResourceAttrPair(withContainerFile, "host_group_ids.0", "huaweicloud_lts_host_group.test", "id"),
					resource.TestCheckResourceAttr(withContainerFile, "tags.key", "value"),
					resource.TestCheckResourceAttr(withContainerFile, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(withContainerFile, "log_split", "true"),
					resource.TestCheckResourceAttrSet(withContainerFile, "demo_log"),
					resource.TestCheckResourceAttr(withContainerFile, "demo_fields.#", "2"),
					resource.TestCheckResourceAttrSet(withContainerFile, "demo_fields.0.field_name"),
					resource.TestCheckResourceAttrSet(withContainerFile, "demo_fields.0.field_value"),
					resource.TestCheckResourceAttr(withContainerFile, "processor_type", "SPLIT"),
					resource.TestCheckResourceAttr(withContainerFile, "processors.#", "2"),
					resource.TestCheckResourceAttr(withContainerFile, "encoding_format", "UTF-8"),
					resource.TestCheckResourceAttr(withContainerFile, "incremental_collect", "true"),
					// Check attributes.
					resource.TestCheckResourceAttr(withContainerFile, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttrSet(withContainerFile, "log_group_name"),
					resource.TestCheckResourceAttrSet(withContainerFile, "log_stream_name"),
					resource.TestCheckResourceAttrSet(withContainerFile, "created_at"),
					rcWithContainerStdout.CheckResourceExists(),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.path_type", "container_stdout"),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.stdout", "true"),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.multi_log_format.0.mode", "time"),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.multi_log_format.0.value", "YYYY-MM-DD hh:mm:ss"),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.system_fields.#", "2"),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.repeat_collect", "false"),
					resource.TestCheckResourceAttr(withContainerStdout, "encoding_format", "GBK"),
					resource.TestCheckResourceAttr(withContainerStdout, "incremental_collect", "false"),
					rcWithHostFile.CheckResourceExists(),
					resource.TestCheckResourceAttr(withHostFile, "access_config.0.paths.#", "2"),
					resource.TestCheckResourceAttr(withHostFile, "access_config.0.black_paths.#", "2"),
					resource.TestCheckResourceAttr(withHostFile, "access_config.0.path_type", "host_file"),
				),
			},
			{
				Config: testAccCceAccessConfig_basic_step2(name, updateName),
				Check: resource.ComposeTestCheckFunc(
					rcWithContainerFile.CheckResourceExists(),
					resource.TestCheckResourceAttr(withContainerFile, "name", updateName),
					resource.TestCheckResourceAttr(withContainerFile, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.name_space_regex", "namespace_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.pod_name_regex", "podname_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.container_name_regex", "containername_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_labels.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_labels.loglabelkey2", "bar_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_labels.includeKey2", "incval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_labels.excludeKey2", "excval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_envs.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_envs.envKey2", "envval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_envs.inEnvKey2", "incval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_envs.exEnvKey2", "excval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_k8s.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.log_k8s.k8sKey2", "k8sval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels_logical", "and"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels.%", "1"),
					resource.TestCheckResourceAttr(withContainerFile, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval_update"),
					resource.TestCheckResourceAttr(withContainerFile, "log_split", "false"),
					resource.TestCheckResourceAttr(withContainerFile, "demo_log", ""),
					resource.TestCheckResourceAttr(withContainerFile, "demo_fields.#", "0"),
					resource.TestCheckResourceAttr(withContainerFile, "processor_type", "SPLIT"),
					rcWithContainerStdout.CheckResourceExists(),
					resource.TestCheckResourceAttr(withContainerStdout, "access_config.0.repeat_collect", "true"),
					resource.TestCheckResourceAttr(withContainerStdout, "encoding_format", "UTF-8"),
					resource.TestCheckResourceAttr(withContainerStdout, "incremental_collect", "true"),
					rcWithHostFile.CheckResourceExists(),
					resource.TestCheckResourceAttr(withHostFile, "access_config.0.black_paths.#", "0"),
				),
			},
			{
				ResourceName:            withContainerFile,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testAccCceAccessConfigImportStateFunc(withContainerFile),
				ImportStateVerifyIgnore: []string{"processors"},
			},
			{
				ResourceName:      withContainerStdout,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCceAccessConfigImportStateFunc(withContainerStdout),
			},
			{
				ResourceName:      withHostFile,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCceAccessConfigImportStateFunc(withHostFile),
			},
		},
	})
}

func testAccCceAccessConfigImportStateFunc(rName string) resource.ImportStateIdFunc {
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

func testAccCceAccessConfig_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[1]s"
  ttl_in_days = 60
}

resource "huaweicloud_lts_host_group" "test" {
  name = "%[1]s"
  type = "linux"
}
`, name)
}

func testAccCceAccessConfig_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "with_container_file" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"
  log_split      = true

  access_config {
    path_type            = "container_file"
    paths                = ["/var"]
    black_paths          = ["/var/a.log"]
    name_space_regex     = "namespace"
    pod_name_regex       = "podname"
    container_name_regex = "containername"

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }

    log_labels = {
      loglabelkey1 = "bar1"
      loglabelkey2 = "bar"
    }

    include_labels = {
      includeKey1 = "incval1"
      includeKey2 = "incval"
    }

    exclude_labels = {
      excludeKey1 = "excval1"
      excludeKey2 = "excval"
    }

    log_envs = {
      envKey1 = "envval1"
      envKey2 = "envval"
    }

    include_envs = {
      inEnvKey1 = "incval1"
      inEnvKey2 = "incval"
    }

    exclude_envs = {
      exEnvKey1 = "excval1"
      exEnvKey2 = "excval"
    }

    log_k8s = {
      k8sKey1 = "k8sval1"
      k8sKey2 = "k8sval"
    }

    include_k8s_labels = {
      ink8sKey1 = "ink8sval1"
      ink8sKey2 = "ink8sval"
    }

    exclude_k8s_labels = {
      exk8sKey1 = "exk8sval1"
      exk8sKey2 = "exk8sval"
    }

    custom_key_value = {
      customkey = "custom_val"
    }
  }

  demo_log       = "2025-04-28 10:59:07.000 a.log:1 level:warn|error"
  processor_type = "SPLIT"

  demo_fields {
    field_name  = "field2"
    field_value = "error"
  }
  demo_fields {
    field_name  = "field1"
    field_value = "2025-04-28 10:59:07.000 a.log:1 level:warn"
  }

  processors {
    type = "processor_filter_regex"

    detail = jsonencode({
      "include" : {
        "name1" : "^terraform"
      },
      "exclude" : {
        "black" : "test"
      }
    })
  }
  processors {
    type = "processor_split_string"

    detail = jsonencode({
      "split_sep" : "|",
      "keys" : ["field1", "field2"],
      "keep_source" : false,
      "keep_source_if_parse_error" : false
    })
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}

resource "huaweicloud_lts_cce_access" "with_container_stdout" {
  name                = "%[2]s_container_file"
  log_group_id        = huaweicloud_lts_group.test.id
  log_stream_id       = huaweicloud_lts_stream.test.id
  cluster_id          = "%[3]s"
  encoding_format     = "GBK"
  incremental_collect = false

  access_config {
    path_type      = "container_stdout"
    stdout         = true
    repeat_collect = false
    system_fields  = ["pathFile", "hostName"]

    multi_log_format  {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }
  }
}

resource "huaweicloud_lts_cce_access" "with_host_file" {
  name          = "%[2]s_host_file"
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  cluster_id    = "%[3]s"

  access_config {
    path_type   = "host_file"
    paths       = ["/var", "/temp"]
    black_paths = ["/var/temp.log", "/var/a.log"]

    single_log_format {
      mode = "system"
    }
  }
}
`, testAccCceAccessConfig_base(name), name, acceptance.HW_LTS_CCE_CLUSTER_ID)
}

func testAccCceAccessConfig_basic_step2(name, updateName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "with_container_file" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"
  log_split      = false

  access_config {
    path_type                  = "container_file"
    paths                      = ["/var"]
    black_paths                = ["/var/a.log"]
    name_space_regex           = "namespace_update"
    pod_name_regex             = "podname_update"
    container_name_regex       = "containername_update"
    exclude_envs_logical       = "and"
    include_envs_logical       = "and"
    exclude_labels_logical     = "and"
    include_labels_logical     = "and"
    exclude_k8s_labels_logical = "and"
    include_k8s_labels_logical = "and"

    windows_log_info {
      categorys        = ["System", "Application"]
      event_level      = ["warning", "error"]
      time_offset_unit = "day"
      time_offset      = 7
    }

    single_log_format {
      mode = "system"
    }

    log_labels = {
      loglabelkey2 = "bar_update"
    }

    include_labels = {
      includeKey2 = "incval_update"
    }

    exclude_labels = {
      excludeKey2 = "excval_update"
    }

    log_envs = {
      envKey2 = "envval_update"
    }

    include_envs = {
      inEnvKey2 = "incval_update"
    }

    exclude_envs = {
      exEnvKey2 = "excval_update"
    }

    log_k8s = {
      k8sKey2 = "k8sval_update"
    }

    include_k8s_labels = {
      ink8sKey2 = "ink8sval_update"
    }

    exclude_k8s_labels = {
      exk8sKey2 = "exk8sval_update"
    }

    custom_key_value = {
      customkey = "custom_val"
    }
  }

  processor_type = "SPLIT"
  processors {}

  tags = {
    foo = "bar_update"
  }
}

resource "huaweicloud_lts_cce_access" "with_container_stdout" {
  name                = "%[2]s_container_file"
  log_group_id        = huaweicloud_lts_group.test.id
  log_stream_id       = huaweicloud_lts_stream.test.id
  cluster_id          = "%[3]s"
  encoding_format     = "UTF-8"
  incremental_collect = true

  access_config {
    path_type      = "container_stdout"
    stdout         = true
    repeat_collect = true

    multi_log_format  {
      mode  = "time"
      value = "YYYY-MM-DD hh:mm:ss"
    }
  }
}

resource "huaweicloud_lts_cce_access" "with_host_file" {
  name          = "%[2]s_host_file"
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
  cluster_id    = "%[3]s"

  access_config {
    path_type = "host_file"
    paths     = ["/var", "/temp"]

    single_log_format {
      mode = "system"
    }
  }
}
`, testAccCceAccessConfig_base(name), updateName, acceptance.HW_LTS_CCE_CLUSTER_ID)
}
