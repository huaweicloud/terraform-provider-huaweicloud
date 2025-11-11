package swrenterprise

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

func getResourceSwrEnterpriseReplicationPolicy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, want '<instance_id>/<policy_id>', but got '%s'", state.Primary.ID)
	}
	instanceId := parts[0]
	id := parts[1]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func TestAccSwrEnterpriseReplicationPolicy_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_replication_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseReplicationPolicy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSwrSignatureWithImageEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseReplicationPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "description", "demo"),
					resource.TestCheckResourceAttr(resourceName, "src_registry.0.id", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "dest_registry.0.id",
						"huaweicloud_swr_enterprise_instance_registry.test", "registry_id"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "scheduled"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				Config: testAccSwrEnterpriseReplicationPolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrPair(resourceName, "src_registry.0.id",
						"huaweicloud_swr_enterprise_instance_registry.test", "registry_id"),
					resource.TestCheckResourceAttr(resourceName, "dest_registry.0.id", "0"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSwrEnterpriseReplicationPolicy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_replication_policy" "test"{
  instance_id     = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name            = "%[2]s"
  enabled         = true
  repo_scope_mode = "regular"
  description     = "demo"

  dest_registry {
    id = huaweicloud_swr_enterprise_instance_registry.test.registry_id
  }

  filters {
    type  = "name"
    value = "**/**"
  }
    
  filters {
    type  = "tag"
    value = "**"
  }

  trigger {
    trigger_settings {
      cron = "0 0 0 1 * ?"
    }
    type = "scheduled"
  }
}
`, testAccSwrEnterpriseInstanceRegistry_enterprise(rName), rName)
}

func testAccSwrEnterpriseReplicationPolicy_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_replication_policy" "test"{
  instance_id    = data.huaweicloud_swr_enterprise_instances.test.instances[0].id
  name            = "%[2]s-update"
  enabled         = false
  repo_scope_mode = "regular"
  description     = ""

  src_registry {
    id = huaweicloud_swr_enterprise_instance_registry.test.registry_id
  }

  filters {
    type  = "name"
    value = "**/**"
  }
    
  filters {
    type  = "tag"
    value = "**"
  }

  trigger {
    type = "manual"
  }
}
`, testAccSwrEnterpriseInstanceRegistry_enterprise(rName), rName)
}
