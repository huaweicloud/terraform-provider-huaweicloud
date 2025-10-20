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

func getResourceSwrEnterpriseImmutableTagRule(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(state.Primary.ID, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<immutable_rule_id>', but got '%s'",
			state.Primary.ID)
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/immutabletagrules?limit=100"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	var rule interface{}
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%v", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return nil, fmt.Errorf("error querying SWR instance immutable tag rule: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, fmt.Errorf("error flattening SWR instance immutable tag rule response: %s", err)
		}

		rules := utils.PathSearch("immutable_rules", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}

		searchPath := fmt.Sprintf("immutable_rules[?namespace_name=='%s'&&id==`%s`]|[0]", namespaceName, id)
		rule = utils.PathSearch(searchPath, getRespBody, nil)
		if rule != nil {
			break
		}

		// offset must be the multiple of limit
		offset += 100
	}

	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func TestAccSwrEnterpriseImmutableTagRule_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_immutable_tag_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseImmutableTagRule,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseImmutableTagRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "immutable_rule_id"),
				),
			},
			{
				Config: testAccSwrEnterpriseImmutableTagRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id", "huaweicloud_swr_enterprise_instance.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "namespace_name", "library"),
					resource.TestCheckResourceAttrSet(resourceName, "namespace_id"),
					resource.TestCheckResourceAttrSet(resourceName, "immutable_rule_id"),
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

func testAccSwrEnterpriseImmutableTagRule_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_immutable_tag_rule" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"

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

  priority = 0
  action   = "immutable"
  disabled = false
  template = "immutable_template"
}
`, testAccSwrEnterpriseInstance_update(rName))
}

func testAccSwrEnterpriseImmutableTagRule_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_swr_enterprise_immutable_tag_rule" "test" {
  instance_id    = huaweicloud_swr_enterprise_instance.test.id
  namespace_name = "library"

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

  priority = 0
  action   = "immutable"
  disabled = true
  template = "immutable_template"
}
`, testAccSwrEnterpriseInstance_update(rName))
}
