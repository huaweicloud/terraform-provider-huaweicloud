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

func getResourceSwrEnterpriseRetentionPolicy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<policy_id>', but got '%s'", state.Primary.ID)
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", namespaceName)
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

func TestAccSwrEnterpriseRetentionPolicy_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_retention_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseRetentionPolicy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseRetentionPolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "scheduled"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
				),
			},
			{
				Config: testAccSwrEnterpriseRetentionPolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "manual"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
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

func testAccSwrEnterpriseRetentionPolicy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_retention_policy" "test"{
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  name           = "%[2]s"
  algorithm      = "or"
  enabled        = true
  
  rules {
    priority        = 0
    action          = "retain"
    repo_scope_mode = "regular"
    disabled        = false
    template        = "latestPushedK"

    params = {
      latestPushedK = jsonencode(1)
    }

    scope_selectors {
      key = "repository"

      value {
        kind       = "doublestar"
        decoration = "repoMatches"
        pattern    = "**"
      }
    }
    
    tag_selectors {
      kind       = "doublestar"
      decoration = "matches"
      pattern    = "**"
    }
  }

  trigger {
    type = "scheduled"

    trigger_settings {
      cron = "0 0 0 1 * ?"
    }
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}

func testAccSwrEnterpriseRetentionPolicy_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_retention_policy" "test"{
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"
  name           = "%[2]s-update"
  algorithm      = "or"
  enabled        = true
  
  rules {
    priority        = 0
    action          = "retain"
    repo_scope_mode = "regular"
    disabled        = false
    template        = "latestPushedK"

    params = {
      latestPushedK = jsonencode(2)
    }

    scope_selectors {
      key = "repository"

      value {
        kind       = "doublestar"
        decoration = "repoMatches"
        pattern    = "**"
      }
    }
    
    tag_selectors {
      kind       = "doublestar"
      decoration = "matches"
      pattern    = "**"
    }
  }

  trigger {
    type = "manual"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}
