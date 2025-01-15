package cae

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cae"
)

func getCertificateFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("cae", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	environmentId := state.Primary.Attributes["environment_id"]
	return cae.GetCertificateById(client, environmentId, state.Primary.ID)
}

func TestAccCertificate_basic(t *testing.T) {
	var (
		obj interface{}

		name = acceptance.RandomAccResourceNameWithDash()

		rName = "huaweicloud_cae_certificate.test"
		rc    = acceptance.InitResourceCheck(
			rName,
			&obj,
			getCertificateFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironment(t)
			acceptance.TestAccPreCheckCertificateWithoutRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate_basic(name, acceptance.HW_CERTIFICATE_CONTENT, acceptance.HW_CERTIFICATE_PRIVATE_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "environment_id", acceptance.HW_CAE_ENVIRONMENT_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "crt", acceptance.HW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(rName, "key", acceptance.HW_CERTIFICATE_PRIVATE_KEY),
				),
			},
			{
				Config: testAccCertificate_basic(name, acceptance.HW_NEW_CERTIFICATE_CONTENT, acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "crt", acceptance.HW_NEW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(rName, "key", acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCertificateImportStateFunc(rName),
			},
		},
	})
}

func testAccCertificate_basic(name, content, privateKey string) string {
	return fmt.Sprintf(`

resource "huaweicloud_cae_certificate" "test" {
  environment_id = "%[1]s"
  name           = "%[2]s"

  # Base64 format corresponding to PEM encoding.
  crt = "%[3]s"
  key = "%[4]s"
}
 `, acceptance.HW_CAE_ENVIRONMENT_ID, name, content, privateKey)
}

func testAccCertificateImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		var (
			environmentId   = rs.Primary.Attributes["environment_id"]
			certificateName = rs.Primary.Attributes["name"]
		)
		if environmentId == "" || certificateName == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>', but got '%s/%s'",
				environmentId, certificateName)
		}

		return fmt.Sprintf("%s/%s", environmentId, certificateName), nil
	}
}
