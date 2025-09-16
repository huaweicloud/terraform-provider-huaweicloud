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

func getResourceSwrEnterpriseLongTermCredential(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("swr", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SWR client: %s", err)
	}

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/long-term-credentials"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", state.Primary.Attributes["instance_id"])

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

	searchPath := fmt.Sprintf("auth_tokens[?token_id=='%s']|[0]", state.Primary.ID)
	token := utils.PathSearch(searchPath, getRespBody, nil)
	if token == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return token, nil
}

func TestAccSwrEnterpriseLongTermCredential_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_swr_enterprise_long_term_credential.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getResourceSwrEnterpriseLongTermCredential,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSwrEnterpriseLongTermCredential_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable", "false"),
				),
			},
			{
				Config: testAccSwrEnterpriseLongTermCredential_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enable", "true"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testSwrEnterpriseLongTermCredentialImportState(resourceName),
				ImportStateVerifyIgnore: []string{"auth_token"},
			},
		},
	})
}

func testAccSwrEnterpriseLongTermCredential_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_long_term_credential" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = "%s"
  enable      = false
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}

func testAccSwrEnterpriseLongTermCredential_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_swr_enterprise_long_term_credential" "test" {
  instance_id = huaweicloud_swr_enterprise_instance.test.id
  name        = "%s"
  enable      = true
}
`, testAccSwrEnterpriseInstance_update(rName), rName)
}

func testSwrEnterpriseLongTermCredentialImportState(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}
		if rs.Primary.ID == "" || rs.Primary.Attributes["instance_id"] == "" {
			return "", fmt.Errorf("resource (%s) instance ID not found: %s", rName, rs)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["instance_id"], rs.Primary.ID), nil
	}
}
