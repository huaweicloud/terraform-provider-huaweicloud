package dew

import (
	"fmt"
	"strings"
	"testing"

	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getKmsAliasResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region          = acceptance.HW_REGION_NAME
		getAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		getAliasProduct = "kms"
	)
	getAliasClient, err := cfg.NewServiceClient(getAliasProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	getAliasPath := getAliasClient.Endpoint + getAliasHttpUrl
	getAliasPath = strings.ReplaceAll(getAliasPath, "{project_id}", getAliasClient.ProjectID)

	getAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	nextMarker := ""
	for {
		params := url.Values{}
		params.Add("key_id", state.Primary.Attributes["key_id"])
		params.Add("limit", "50")
		if marker != "" {
			params.Add("marker", marker)
		}
		getAliasPath := fmt.Sprintf("%s?%s", getAliasPath, params.Encode())

		getAliasResp, err := getAliasClient.Request("GET", getAliasPath, &getAliasOpt)

		getAliasRespBody, err := utils.FlattenResponse(getAliasResp)
		if err != nil {
			return nil, fmt.Errorf("error creating KMS client: %s", err)
		}
		aliases := utils.PathSearch("aliases", getAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) > 0 {
			allAliases = append(allAliases, aliases...)
		}

		nextMarker = utils.PathSearch("page_info.next_marker", getAliasRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		marker = nextMarker
	}
	searchPath := fmt.Sprintf("[?alias=='%s']|[0]", state.Primary.Attributes["alias"])
	aliasDetail := utils.PathSearch(searchPath, allAliases, nil)
	if aliasDetail == nil {
		return nil, fmt.Errorf("error retrieving KMS alias")
	}

	return aliasDetail, nil
}

func TestAccKmsAlias_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_kms_alias.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsAliasResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKmsAlias_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttrSet(rName, "alias_urn"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
					resource.TestCheckResourceAttrSet(rName, "alias"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKmsAliasImportState(rName),
			},
		},
	})
}

func testKmsAlias_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_alias" "test" {
  key_id             = "%s"
  alias              = "alias/%s"
}
`, acceptance.HW_KMS_KEY_ID, name)
}

func testKmsAliasImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		if rs.Primary.Attributes["key_id"] == "" {
			return "", fmt.Errorf("attribute (key_id) of Resource (%s) not found: %s", name, rs)
		}

		if rs.Primary.Attributes["alias"] == "" {
			return "", fmt.Errorf("attribute (alias) of Resource (%s) not found: %s", name, rs)
		}
		return rs.Primary.Attributes["key_id"] + "?" + rs.Primary.Attributes["alias"], nil
	}
}
