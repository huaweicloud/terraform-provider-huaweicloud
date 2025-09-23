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

func getKmsGrantResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getGrant: Query the KMS manual Grant
	var (
		getGrantHttpUrl = "v1.0/{project_id}/kms/list-grants"
		getGrantProduct = "kms"
	)
	getGrantClient, err := cfg.NewServiceClient(getGrantProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS Client: %s", err)
	}

	getGrantPath := getGrantClient.Endpoint + getGrantHttpUrl
	getGrantPath = strings.ReplaceAll(getGrantPath, "{project_id}", getGrantClient.ProjectID)

	getGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		JSONBody: map[string]interface{}{
			"key_id": state.Primary.Attributes["key_id"],
			"limit":  100,
		},
	}
	getGrantResp, err := getGrantClient.Request("POST", getGrantPath, &getGrantOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	grantResponseBody, err := utils.FlattenResponse(getGrantResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	searchPath := fmt.Sprintf("grants[?grant_id=='%s']|[0]", state.Primary.ID)
	r := utils.PathSearch(searchPath, grantResponseBody, nil)
	if r == nil {
		return nil, fmt.Errorf("error retrieving KMS grant: %s", err)
	}

	return r, nil
}

func TestAccKmsGrant_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_kms_grant.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getKmsGrantResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKmsGrant_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "operations.#", "8"),
					resource.TestCheckResourceAttrPair(rName, "key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "grantee_principal", "huaweicloud_identity_user.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "creator"),
					resource.TestCheckResourceAttrPair(rName, "retiring_principal", "huaweicloud_identity_user.test", "id"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testKmsGrantImportState(rName),
			},
		},
	})
}

func testKmsGrant_basic(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_identity_user" "test" {
  name        = "%s"
  password    = "password123@!"
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_kms_grant" "test" {
  key_id             = huaweicloud_kms_key.test.id
  grantee_principal  = huaweicloud_identity_user.test.id
  retiring_principal = huaweicloud_identity_user.test.id

  operations = [
    "create-datakey",
    "create-datakey-without-plaintext",
    "describe-key",
    "encrypt-data",
    "decrypt-data",
    "decrypt-datakey",
    "retire-grant",
    "encrypt-datakey"
  ]
}
`, name, name)
}

func testKmsGrantImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.Attributes["key_id"] == "" {
			return "", fmt.Errorf("attribute (key_id) of Resource (%s) not found: %s", name, rs)
		}
		if rs.Primary.ID == "" {
			return "", fmt.Errorf("attribute (ID) of Resource (%s) not found: %s", name, rs)
		}

		return rs.Primary.Attributes["key_id"] + "/" +
			rs.Primary.ID, nil
	}
}
