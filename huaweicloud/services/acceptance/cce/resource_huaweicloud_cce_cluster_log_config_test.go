package cce

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

func getClusterLogConfigResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getClusterLogConfigHttpUrl = "api/v3/projects/{project_id}/cluster/{cluster_id}/log-configs"
		getClusterLogConfigProduct = "cce"
	)

	clusterID := state.Primary.ID
	client, err := cfg.NewServiceClient(getClusterLogConfigProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE Client: %s", err)
	}

	getClusterLogConfigPath := client.Endpoint + getClusterLogConfigHttpUrl
	getClusterLogConfigPath = strings.ReplaceAll(getClusterLogConfigPath, "{project_id}", client.ProjectID)
	getClusterLogConfigPath = strings.ReplaceAll(getClusterLogConfigPath, "{cluster_id}", clusterID)

	getClusterLogConfigOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getClusterLogConfigResp, err := client.Request("GET", getClusterLogConfigPath, &getClusterLogConfigOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving log config: %s", err)
	}

	return utils.FlattenResponse(getClusterLogConfigResp)
}

func TestAccClusterLogConfig_basic(t *testing.T) {
	var logConfig interface{}
	resourceName := "huaweicloud_cce_cluster_log_config.test"
	randName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		resourceName,
		&logConfig,
		getClusterLogConfigResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterLogConfig_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "7"),
					resource.TestCheckOutput("log_configs", "true"),
				),
			},
			{
				Config: testAccClusterLogConfig_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "id",
						"huaweicloud_cce_cluster.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "3"),
					resource.TestCheckOutput("log_configs", "true"),
				),
			},
			{
				Config:            testAccClusterLogConfig_config(randName),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClusterLogConfig_basic(randName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_log_config" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id

  log_configs {
    name   = "kube-apiserver"
    enable = true
  }

  log_configs {
    name   = "kube-controller-manager"
    enable = true
  }

  log_configs {
    name   = "kube-scheduler"
    enable = true
  }
}

variable "configs" {
  type = map

  default = {
    "kube-apiserver"          = true,
    "kube-controller-manager" = true,
    "kube-scheduler"          = true
  }
}

output "log_configs" {
  value = alltrue([for s in huaweicloud_cce_cluster_log_config.test.log_configs : var.configs[s.name] == s.enable])
}
`, testAccCluster_basic(randName))
}

func testAccClusterLogConfig_update(randName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_log_config" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  ttl_in_days = 3

  log_configs {
    name   = "kube-apiserver"
    enable = true
  }

  log_configs {
    name   = "kube-controller-manager"
    enable = false
  }

  log_configs {
    name   = "kube-scheduler"
    enable = false
  }

  log_configs {
    name   = "audit"
    enable = true
  }
}

variable "configs" {
  type = map

  default = {
    "kube-apiserver"          = true,
    "kube-controller-manager" = false,
    "kube-scheduler"          = false,
    "audit"                   = true
  }
}

output "log_configs" {
  value = alltrue([for s in huaweicloud_cce_cluster_log_config.test.log_configs : var.configs[s.name] == s.enable])
}
`, testAccCluster_basic(randName))
}

func testAccClusterLogConfig_config(randName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_cluster_log_config" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  ttl_in_days = 3

  log_configs {
    name   = "kube-apiserver"
    enable = true
  }

  log_configs {
    name   = "kube-controller-manager"
    enable = false
  }

  log_configs {
    name   = "kube-scheduler"
    enable = false
  }

  log_configs {
    name   = "audit"
    enable = true
  }
}
`, testAccCluster_basic(randName))
}
