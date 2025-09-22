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

func getResourceSwrEnterpriseTrigger(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<id>', but got '%s'", state.Primary.ID)
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}"
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

func TestAccSwrEnterpriseTrigger_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_trigger.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseTrigger,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseTrigger_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "scope_rules.0.repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.address", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.address_type", "internal"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.auth_header", "Test:Header"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.skip_cert_verify", "false"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "http"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "trigger_id"),
				),
			},
			{
				Config: testAccSwrEnterpriseTrigger_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "scope_rules.0.repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.address", "https://test.com"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.address_type", "internal"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.auth_header", "TestUpdate:Header"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.skip_cert_verify", "true"),
					resource.TestCheckResourceAttr(resourceName, "targets.0.type", "http"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "trigger_id"),
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

func testAccSwrEnterpriseTrigger_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_trigger" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  description    = "desc"
  enabled        = true
  event_types    = ["PUSH_ARTIFACT"]
  name           = "%[2]s"
  namespace_name = "library"
  
  scope_rules {
    repo_scope_mode = "regular"

    scope_selectors {
      key = "repository"

      value {
        decoration = "repoMatches"
        kind       = "doublestar"
        pattern    = "**"
	  }
    }

    tag_selectors {
      decoration = "matches"
      kind       = "doublestar"
      pattern    = "**"
    }
  }

  targets {
    address          = "https://test.com"
    address_type     = "internal"
    auth_header      = "Test:Header"
    skip_cert_verify = false
    type             = "http"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}

func testAccSwrEnterpriseTrigger_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_trigger" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  description    = ""
  enabled        = false
  event_types    = ["PUSH_ARTIFACT"]
  name           = "%[2]s-update"
  namespace_name = "library"
  
  scope_rules {
    repo_scope_mode = "regular"

    scope_selectors {
      key = "repository"

      value {
        decoration = "repoMatches"
        kind       = "doublestar"
        pattern    = "nginx-*"
      }
	}

    tag_selectors {
      decoration = "matches"
      kind       = "doublestar"
      pattern    = "**"
	}
  }

  targets {
    address          = "https://test.com"
    address_type     = "internal"
    auth_header      = "TestUpdate:Header"
    skip_cert_verify = true
    type             = "http"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}
