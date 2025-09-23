package accessanalyzer

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

func getArchiveRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("accessanalyzer", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Access Analyzer client: %s", err)
	}
	getArchiveRuleHttpUrl := "v5/analyzers/{analyzer_id}/archive-rules/{archive_rule_id}"
	getArchiveRulePath := client.Endpoint + getArchiveRuleHttpUrl
	getArchiveRulePath = strings.ReplaceAll(getArchiveRulePath, "{analyzer_id}", state.Primary.Attributes["analyzer_id"])
	getArchiveRulePath = strings.ReplaceAll(getArchiveRulePath, "{archive_rule_id}", state.Primary.ID)
	getArchiveRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getArchiveRuleResp, err := client.Request("GET", getArchiveRulePath, &getArchiveRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving archive rule: %s", err)
	}
	return utils.FlattenResponse(getArchiveRuleResp)
}

func TestAccArchiveRule_basic(t *testing.T) {
	var object interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_access_analyzer_archive_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&object,
		getArchiveRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccArchiveRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "analyzer_id",
						"huaweicloud_access_analyzer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "filters.0.key", "resource_type"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.criterion.0.eq.0", "iam:agency"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.criterion.0.eq.1", "obs:bucket"),
				),
			},
			{
				Config: testAccArchiveRule_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "analyzer_id",
						"huaweicloud_access_analyzer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "filters.0.key", "resource_type"),
					resource.TestCheckResourceAttr(resourceName, "filters.0.criterion.0.eq.0", "iam:agency"),
					resource.TestCheckResourceAttr(resourceName, "filters.1.key", "condition.g:SourceVpc"),
					resource.TestCheckResourceAttr(resourceName, "filters.1.criterion.0.exists", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccArchiveRuleImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccArchiveRuleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var analyzerID, ruleID string
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of archive rule is not found in the tfstate", resourceName)
		}
		analyzerID = rs.Primary.Attributes["analyzer_id"]
		ruleID = rs.Primary.ID
		if analyzerID == "" || ruleID == "" {
			return "", fmt.Errorf("the archive rule ID is not exist or analyzer ID is missing")
		}
		return fmt.Sprintf("%s/%s", analyzerID, ruleID), nil
	}
}

func testAccArchiveRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_access_analyzer_archive_rule" "test" {
  analyzer_id = huaweicloud_access_analyzer.test.id
  name        = "%s"

  filters {
    key = "resource_type"

    criterion {
        eq = ["iam:agency", "obs:bucket"]
    }
  }
}
`, testAccAnalyzer_basic(rName), rName)
}

func testAccArchiveRule_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_access_analyzer_archive_rule" "test" {
  analyzer_id = huaweicloud_access_analyzer.test.id
  name        = "%s"

  filters {
    key = "resource_type"

    criterion {
        eq = ["iam:agency"]
    }
  }

  filters {
    key = "condition.g:SourceVpc"

    criterion {
        exists = true
    }
  }
}
`, testAccAnalyzer_basic(rName), rName)
}
