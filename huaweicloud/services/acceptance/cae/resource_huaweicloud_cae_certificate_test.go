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
	return cae.GetCertificateById(
		client,
		state.Primary.Attributes["environment_id"],
		state.Primary.ID,
		state.Primary.Attributes["enterprise_project_id"],
	)
}

func TestAccCertificate_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		certificate interface{}
		rName       = "huaweicloud_cae_certificate.test.0"
		rc          = acceptance.InitResourceCheck(rName, &certificate, getCertificateFunc)

		withNotDefaultEpsId   = "huaweicloud_cae_certificate.test.1"
		rcWithNotDefaultEpsId = acceptance.InitResourceCheck(withNotDefaultEpsId, &certificate, getCertificateFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCaeEnvironmentIds(t, 2)
			acceptance.TestAccPreCheckCertificateWithoutRootCA(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCertificate_basic(name, acceptance.HW_CERTIFICATE_CONTENT, acceptance.HW_CERTIFICATE_PRIVATE_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "environment_id"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s0", name)),
					resource.TestCheckResourceAttr(rName, "crt", acceptance.HW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(rName, "key", acceptance.HW_CERTIFICATE_PRIVATE_KEY),
					resource.TestCheckNoResourceAttr(rName, "enterprise_project_id"),
					rcWithNotDefaultEpsId.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(withNotDefaultEpsId, "environment_id"),
					resource.TestCheckResourceAttr(withNotDefaultEpsId, "name", fmt.Sprintf("%s1", name)),
					resource.TestCheckResourceAttr(withNotDefaultEpsId, "crt", acceptance.HW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(withNotDefaultEpsId, "key", acceptance.HW_CERTIFICATE_PRIVATE_KEY),
					resource.TestCheckResourceAttrPair(withNotDefaultEpsId, "enterprise_project_id",
						"data.huaweicloud_cae_environments.test", "environments.0.annotations.enterprise_project_id"),
				),
			},
			{
				Config: testAccCertificate_basic(name, acceptance.HW_NEW_CERTIFICATE_CONTENT, acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "crt", acceptance.HW_NEW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(rName, "key", acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY),
					rcWithNotDefaultEpsId.CheckResourceExists(),
					resource.TestCheckResourceAttr(withNotDefaultEpsId, "crt", acceptance.HW_NEW_CERTIFICATE_CONTENT),
					resource.TestCheckResourceAttr(withNotDefaultEpsId, "key", acceptance.HW_NEW_CERTIFICATE_PRIVATE_KEY),
				),
			},
			{
				ResourceName:      "huaweicloud_cae_certificate.test[0]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCertificateImportStateFunc(rName, true),
			},
			{
				ResourceName:      "huaweicloud_cae_certificate.test[1]",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccCertificateImportStateFunc(withNotDefaultEpsId, false),
			},
		},
	})
}

func testAccCertificate_basic(name, content, privateKey string) string {
	return fmt.Sprintf(`
locals {
  env_ids = split(",", "%[1]s")
}

data "huaweicloud_cae_environments" "test" {
  environment_id = local.env_ids[1]
}

resource "huaweicloud_cae_certificate" "test" {
  count = 2

  environment_id        = local.env_ids[count.index]
  name                  = "%[2]s${count.index}"
  enterprise_project_id = count.index == 1 ? try(data.huaweicloud_cae_environments.test.environments[0].annotations.enterprise_project_id,
  null) : null

  # Base64 format corresponding to PEM encoding.
  crt = "%[3]s"
  key = "%[4]s"
}
 `, acceptance.HW_CAE_ENVIRONMENT_IDS, name, content, privateKey)
}

func testAccCertificateImportStateFunc(name string, isDefaultEpsId bool) resource.ImportStateIdFunc {
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
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<name>' or "+
				"'<environment_id>/<name>/<enterprise_project_id>', but got '%s/%s'", environmentId, certificateName)
		}

		if isDefaultEpsId {
			return fmt.Sprintf("%s/%s", environmentId, certificateName), nil
		}

		epsId := rs.Primary.Attributes["enterprise_project_id"]
		if epsId == "" {
			return "", fmt.Errorf("some import IDs are missing, want '<environment_id>/<id>/<enterprise_project_id>', but got '%s/%s/%s'",
				environmentId, certificateName, epsId)
		}
		return fmt.Sprintf("%s/%s/%s", environmentId, certificateName, epsId), nil
	}
}
