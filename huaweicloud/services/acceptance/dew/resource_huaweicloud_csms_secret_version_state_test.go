package dew

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getSecretVersionStateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region                 = acceptance.HW_REGION_NAME
		getVersionStatehttpUrl = "v1/{project_id}/secrets/{secret_name}/stages/{stage_name}"
		product                = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS client: %s", err)
	}

	getVersionStatePath := client.Endpoint + getVersionStatehttpUrl
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{project_id}", client.ProjectID)
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{secret_name}", state.Primary.Attributes["secret_name"])
	getVersionStatePath = strings.ReplaceAll(getVersionStatePath, "{stage_name}", state.Primary.ID)
	getVersionStateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getVersionStatePath, &getVersionStateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving secret version state: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccSecretVersionState_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_csms_secret_version_state.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getSecretVersionStateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSecretVersionState_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "secret_name", "huaweicloud_csms_secret.test", "name"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttrPair(rName, "version_id", "huaweicloud_csms_secret.test", "latest_version"),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testAccSecretVersionState_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "version_id", "v2"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestMatchResourceAttr(rName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSecretVersionStateImportStateFunc(rName),
			},
		},
	})
}

func testAccSecretVersionStateImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var secretName, stateId string
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("the resource (%s) not found", name)
		}

		secretName = rs.Primary.Attributes["secret_name"]
		stateId = rs.Primary.ID
		if secretName == "" || stateId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<secret_name>/<id>', but got '%s/%s'",
				secretName, stateId)
		}
		return fmt.Sprintf("%s/%s", secretName, stateId), nil
	}
}

func testAccSecretVersionState_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name        = "%[1]s"
  secret_text = "secret version"
  description = "acc test"
}

resource "huaweicloud_csms_secret_version_state" "test" {
  secret_name = huaweicloud_csms_secret.test.name
  name        = "%[1]s"
  version_id  = huaweicloud_csms_secret.test.latest_version
}
`, name)
}

func testAccSecretVersionState_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_csms_secret" "test" {
  name        = "%[1]s"
  secret_text = "version state"
  description = "acc test"
}

resource "huaweicloud_csms_secret_version_state" "test" {
  secret_name = huaweicloud_csms_secret.test.name
  name        = "%[1]s"
  version_id  = "v2"
}
`, name)
}
