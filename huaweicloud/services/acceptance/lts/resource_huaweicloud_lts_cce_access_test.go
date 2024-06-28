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

func TestAccCceAccessConfig_containerFile(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_cce_access.container_file"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCceAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSCCEAccess(t)
			acceptance.TestAccPreCheckLTSHostGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigContainerFile(name, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "container_file"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey1", "bar1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey1", "envval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey1", "k8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey1", "ink8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey1", "exk8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval"),
				),
			},
			{
				Config: testCceAccessConfigContainerFile(name, "-update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCceAccessConfigImportStateFunc(rName),
			},
		},
	})
}

func TestAccCceAccessConfig_containerStdout(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_cce_access.container_stdout"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCceAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSCCEAccess(t)
			acceptance.TestAccPreCheckLTSHostGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigContainerStdout(name, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "container_stdout"),
					resource.TestCheckResourceAttr(rName, "access_config.0.stdout", "true"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey1", "bar1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey1", "envval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey1", "incval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey1", "excval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval"),

					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey1", "k8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey1", "ink8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey1", "exk8sval1"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval"),
				),
			},
			{
				Config: testCceAccessConfigContainerStdout(name, "-update"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.name_space_regex", "namespace-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.pod_name_regex", "podname-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.container_name_regex", "containername-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_labels.loglabelkey2", "bar-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_labels.includeKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_labels.excludeKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_envs.envKey2", "envval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_envs.inEnvKey2", "incval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_envs.exEnvKey2", "excval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.log_k8s.k8sKey2", "k8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.include_k8s_labels.ink8sKey2", "ink8sval-update"),
					resource.TestCheckResourceAttr(rName, "access_config.0.exclude_k8s_labels.exk8sKey2", "exk8sval-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCceAccessConfigImportStateFunc(rName),
			},
		},
	})
}

func TestAccCceAccessConfig_hostFile(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_lts_cce_access.host_file"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getCceAccessConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLTSCCEAccess(t)
			acceptance.TestAccPreCheckLTSHostGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCceAccessConfigHostFile(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrSet(rName, "log_group_name"),
					resource.TestCheckResourceAttrSet(rName, "log_stream_name"),
					resource.TestCheckResourceAttr(rName, "tags.key", "value"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/a.log"),
					resource.TestCheckResourceAttr(rName, "access_config.0.path_type", "host_file"),
				),
			},
			{
				Config: testCceAccessConfigHostFileUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "tags.key", "value-updated"),
					resource.TestCheckResourceAttr(rName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(rName, "access_type", "K8S_CCE"),
					resource.TestCheckResourceAttr(rName, "access_config.0.paths.0", "/var/logs"),
					resource.TestCheckResourceAttr(rName, "access_config.0.black_paths.0", "/var/logs/a.log"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCceAccessConfigImportStateFunc(rName),
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

func testCceAccessConfig_base(name string) string {
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
  name     = "%[1]s"
  type     = "linux"
  host_ids = split(",", "%[2]s")
}
`, name, acceptance.HW_LTS_HOST_IDS)
}

func testCceAccessConfigContainerFile(name, suffix string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "container_file" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"

  access_config {
    path_type            = "container_file"
    paths                = ["/var"]
    black_paths          = ["/var/a.log"]
    name_space_regex     = "namespace%[4]s"
    pod_name_regex       = "podname%[4]s"
    container_name_regex = "containername%[4]s"

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
      loglabelkey2 = "bar%[4]s"
    }

    include_labels = {
      includeKey1 = "incval1"
      includeKey2 = "incval%[4]s"
    }

    exclude_labels = {
      excludeKey1 = "excval1"
      excludeKey2 = "excval%[4]s"
    }

    log_envs = {
      envKey1 = "envval1"
      envKey2 = "envval%[4]s"
    }

    include_envs = {
      inEnvKey1 = "incval1"
      inEnvKey2 = "incval%[4]s"
    }

    exclude_envs = {
      exEnvKey1 = "excval1"
      exEnvKey2 = "excval%[4]s"
    }

    log_k8s = {
      k8sKey1 = "k8sval1"
      k8sKey2 = "k8sval%[4]s"
    }

    include_k8s_labels = {
      ink8sKey1 = "ink8sval1"
      ink8sKey2 = "ink8sval%[4]s"
    }

    exclude_k8s_labels = {
      exk8sKey1 = "exk8sval1"
      exk8sKey2 = "exk8sval%[4]s"
    }
  }

  tags = {
    key = "value"
    foo = "bar%[4]s"
  }
}
`, testCceAccessConfig_base(name), name, acceptance.HW_LTS_CCE_CLUSTER_ID, suffix)
}

func testCceAccessConfigContainerStdout(name, suffix string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "container_stdout" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"

  access_config {
    path_type            = "container_stdout"
    stdout               = true
    name_space_regex     = "namespace%[4]s"
    pod_name_regex       = "podname%[4]s"
    container_name_regex = "containername%[4]s"

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
      loglabelkey2 = "bar%[4]s"
    }

    include_labels = {
      includeKey1 = "incval1"
      includeKey2 = "incval%[4]s"
    }

    exclude_labels = {
      excludeKey1 = "excval1"
      excludeKey2 = "excval%[4]s"
    }

    log_envs = {
      envKey1 = "envval1"
      envKey2 = "envval%[4]s"
    }

    include_envs = {
      inEnvKey1 = "incval1"
      inEnvKey2 = "incval%[4]s"
    }

    exclude_envs = {
      exEnvKey1 = "excval1"
      exEnvKey2 = "excval%[4]s"
    }

    log_k8s = {
      k8sKey1 = "k8sval1"
      k8sKey2 = "k8sval%[4]s"
    }

    include_k8s_labels = {
      ink8sKey1 = "ink8sval1"
      ink8sKey2 = "ink8sval%[4]s"
    }

    exclude_k8s_labels = {
      exk8sKey1 = "exk8sval1"
      exk8sKey2 = "exk8sval%[4]s"
    }
  }

  tags = {
    key = "value"
    foo = "bar%[4]s"
  }
}
`, testCceAccessConfig_base(name), name, acceptance.HW_LTS_CCE_CLUSTER_ID, suffix)
}

func testCceAccessConfigHostFile(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "host_file" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"

  access_config {
    path_type   = "host_file"
    paths       = ["/var"]
    black_paths = ["/var/a.log"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key = "value"
    foo = "bar"
  }
}
`, testCceAccessConfig_base(name), name, acceptance.HW_LTS_CCE_CLUSTER_ID)
}

func testCceAccessConfigHostFileUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_cce_access" "host_file" {
  name           = "%[2]s"
  log_group_id   = huaweicloud_lts_group.test.id
  log_stream_id  = huaweicloud_lts_stream.test.id
  host_group_ids = [huaweicloud_lts_host_group.test.id]
  cluster_id     = "%[3]s"

  access_config {
    path_type   = "host_file"
    paths       = ["/var/logs"]
    black_paths = ["/var/logs/a.log"]

    single_log_format {
      mode = "system"
    }
  }

  tags = {
    key   = "value-updated"
    owner = "terraform"
  }
}
`, testCceAccessConfig_base(name), name, acceptance.HW_LTS_CCE_CLUSTER_ID)
}
