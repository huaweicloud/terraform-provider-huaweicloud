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

func getResourceSwrEnterpriseImageSignaturePolicy(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
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

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}"
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

func TestAccSwrEnterpriseImageSignaturePolicy_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_image_signature_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseImageSignaturePolicy,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseImageSignaturePolicy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "desc"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "scope_rules.0.repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "event_based"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
				),
			},
			{
				Config: testAccSwrEnterpriseImageSignaturePolicy_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttr(resourceName, "scope_rules.0.repo_scope_mode", "regular"),
					resource.TestCheckResourceAttr(resourceName, "trigger.0.type", "event_based"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
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

func testAccKmsKey_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias             = "%s"
  key_algorithm         = "SM2"
  key_usage             = "SIGN_VERIFY"
  origin                = "kms"
  enterprise_project_id = "0"
  pending_days          = "7"
}
`, name)
}

func testAccSwrEnterpriseImageSignaturePolicy_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_swr_enterprise_image_signature_policy" "test" {
  instance_id         = huaweicloud_swr_enterprise_instance.test.id
  namespace_name      = "library"
  name                = "%[3]s"
  description         = "desc"
  enabled             = true
  signature_method    = "KMS"
  signature_algorithm = "SM2DSA_SM3"
  signature_key       = huaweicloud_kms_key.test.id
  
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

  trigger {
    type = "event_based"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), testAccKmsKey_basic(rName), rName)
}

func testAccSwrEnterpriseImageSignaturePolicy_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

%[2]s

resource "huaweicloud_swr_enterprise_image_signature_policy" "test" {
  instance_id         = huaweicloud_swr_enterprise_instance.test.id
  namespace_name      = "library"
  name                = "%[3]s-update"
  enabled             = false
  signature_method    = "KMS"
  signature_algorithm = "SM2DSA_SM3"
  signature_key       = huaweicloud_kms_key.test.id
  
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
      extras     = jsonencode({
        "untagged": true
      })
    }
  }

  trigger {
    type = "event_based"
  }
}
`, testAccSwrEnterpriseInstance_update(rName), testAccKmsKey_basic(rName), rName)
}
