package dew

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

	basePath := getAliasClient.Endpoint + getAliasHttpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", getAliasClient.ProjectID)

	getAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	basePath += fmt.Sprintf("?key_id=%s&limit=50", state.Primary.Attributes["key_id"])
	for {
		getAliasPath := basePath
		if marker != "" {
			getAliasPath = basePath + fmt.Sprintf("&marker=%s", marker)
		}

		getAliasResp, err := getAliasClient.Request("GET", getAliasPath, &getAliasOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving KMS alias: %s", err)
		}

		getAliasRespBody, err := utils.FlattenResponse(getAliasResp)
		if err != nil {
			return nil, err
		}
		aliases := utils.PathSearch("aliases", getAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) == 0 {
			break
		}

		allAliases = append(allAliases, aliases...)

		marker = utils.PathSearch("page_info.next_marker", getAliasRespBody, "").(string)
		if marker == "" {
			break
		}
	}
	searchPath := fmt.Sprintf("[?alias=='%s']|[0]", state.Primary.Attributes["alias"])
	aliasDetail := utils.PathSearch(searchPath, allAliases, nil)
	if aliasDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return aliasDetail, nil
}

func TestAccKmsAlias_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	aliasName := fmt.Sprintf("alias/%s", name)
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
				Config: testKmsAlias_basic(aliasName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr(rName, "alias", aliasName),
					resource.TestCheckResourceAttrSet(rName, "alias_urn"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
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
  key_id = "%s"
  alias  = "%s"
}
`, acceptance.HW_KMS_KEY_ID, name)
}

func testKmsAliasImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		if rs.Primary.Attributes["key_id"] == "" {
			return "", fmt.Errorf("attribute (key_id) of Resource (%s) not found", name)
		}

		if rs.Primary.Attributes["alias"] == "" {
			return "", fmt.Errorf("attribute (alias) of Resource (%s) not found", name)
		}
		return rs.Primary.Attributes["key_id"] + "?" + rs.Primary.Attributes["alias"], nil
	}
}
