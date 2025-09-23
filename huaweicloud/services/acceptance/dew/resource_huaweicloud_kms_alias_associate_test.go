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

func getKmsAliasAssociateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                   = acceptance.HW_REGION_NAME
		getAssociateAliasHttpUrl = "v1.0/{project_id}/kms/aliases"
		getAssociateAliasProduct = "kms"
	)
	getAssociateAliasClient, err := cfg.NewServiceClient(getAssociateAliasProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	basePath := getAssociateAliasClient.Endpoint + getAssociateAliasHttpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", getAssociateAliasClient.ProjectID)

	getAssociateAliasOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allAliases := make([]interface{}, 0)
	marker := ""
	basePath += fmt.Sprintf("?key_id=%s&limit=50", state.Primary.Attributes["target_key_id"])
	for {
		getAssociateAliasPath := basePath
		if marker != "" {
			getAssociateAliasPath = basePath + fmt.Sprintf("&marker=%s", marker)
		}

		getAssociateAliasResp, err := getAssociateAliasClient.Request("GET", getAssociateAliasPath, &getAssociateAliasOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving KMS alias: %s", err)
		}

		getAssociateAliasRespBody, err := utils.FlattenResponse(getAssociateAliasResp)
		if err != nil {
			return nil, err
		}
		aliases := utils.PathSearch("aliases", getAssociateAliasRespBody, make([]interface{}, 0)).([]interface{})
		if len(aliases) == 0 {
			break
		}

		allAliases = append(allAliases, aliases...)

		marker = utils.PathSearch("page_info.next_marker", getAssociateAliasRespBody, "").(string)
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

func TestAccKmsAliasAssociate_basic(t *testing.T) {
	var obj interface{}
	rName := "huaweicloud_kms_alias_associate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsAliasAssociateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckKmsKeyID(t)
			// an existing alias must be set for association
			acceptance.TestAccPreCheckKmsAlias(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKmsAliasAssociate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "target_key_id", acceptance.HW_KMS_KEY_ID),
					resource.TestCheckResourceAttr(rName, "alias", acceptance.HW_KMS_ALIAS),
					resource.TestCheckResourceAttrSet(rName, "alias_urn"),
					resource.TestCheckResourceAttrSet(rName, "update_time"),
					resource.TestCheckResourceAttrSet(rName, "create_time"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKmsAliasAssociateImportState(rName),
			},
		},
	})
}

func testKmsAliasAssociate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_alias_associate" "test" {
  target_key_id = "%s"
  alias         = "%s"
}
`, acceptance.HW_KMS_KEY_ID, acceptance.HW_KMS_ALIAS)
}

func testKmsAliasAssociateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found", name)
		}

		if rs.Primary.Attributes["target_key_id"] == "" {
			return "", fmt.Errorf("attribute (target_key_id) of Resource (%s) not found", name)
		}

		if rs.Primary.Attributes["alias"] == "" {
			return "", fmt.Errorf("attribute (alias) of Resource (%s) not found", name)
		}
		return rs.Primary.Attributes["target_key_id"] + "?" + rs.Primary.Attributes["alias"], nil
	}
}
