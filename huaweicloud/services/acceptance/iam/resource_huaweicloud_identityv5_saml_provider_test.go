package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/iam"
)

func getV5SamlProviderFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("iam", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	return iam.GetV5SamlProviderById(client, state.Primary.ID)
}

// Please ensure that the user executing the acceptance test has 'admin' permission,
// and the account has enabled the IAM v5 (new IAM) feature.
func TestAccV5SamlProvider_basic(t *testing.T) {
	var (
		obj    interface{}
		rName  = "huaweicloud_identityv5_saml_provider.test"
		dsName = "data.huaweicloud_identityv5_saml_providers.test"
		rc     = acceptance.InitResourceCheck(rName, &obj, getV5SamlProviderFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckIAMV5(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccV5SamlProvider_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", "acc-saml-provider"),
					resource.TestCheckResourceAttr(rName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(rName, "provider_id"),
				),
			},
			{
				Config: testAccV5SamlProvider_basic_step2() + testAccV5SamlProviders_dataSource(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(dsName, "providers.#", "1"),
					resource.TestCheckResourceAttr(dsName, "providers.0.name", "acc-saml-provider"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccV5SamlProvider_basic_step1() string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_saml_provider" "test" {
  name        = "acc-saml-provider"
  description = "Created by terraform script"

  saml_metadata_document = <<EOF
%[1]s
EOF
}
`, testAccV5SamlMetadataDocument())
}

func testAccV5SamlProvider_basic_step2() string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_saml_provider" "test" {
  name        = "acc-saml-provider"
  description = "Updated by terraform script"

  saml_metadata_document = <<EOF
%[1]s
EOF
}
`, testAccV5SamlMetadataDocumentUpdated())
}

func testAccV5SamlProviders_dataSource() string {
	return `
data "huaweicloud_identityv5_saml_providers" "test" {
  name       = "acc-saml-provider"
  depends_on = [huaweicloud_identityv5_saml_provider.test]
}
`
}

func testAccV5SamlMetadataDocument() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor entityID="https://terraform-acc-test.example.com/saml" xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata">
  <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:KeyDescriptor use="signing">
      <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <ds:X509Data>
          <ds:X509Certificate>%s</ds:X509Certificate>
        </ds:X509Data>
      </ds:KeyInfo>
    </md:KeyDescriptor>
    <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://terraform-acc-test.example.com/sso/saml"/>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://terraform-acc-test.example.com/sso/saml"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`, testAccV5SamlCertificate())
}

func testAccV5SamlMetadataDocumentUpdated() string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<md:EntityDescriptor entityID="https://terraform-acc-test2.example.com/saml" xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata">
  <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <md:KeyDescriptor use="signing">
      <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
        <ds:X509Data>
          <ds:X509Certificate>%s</ds:X509Certificate>
        </ds:X509Data>
      </ds:KeyInfo>
    </md:KeyDescriptor>
    <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://terraform-acc-test2.example.com/sso/saml"/>
    <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://terraform-acc-test2.example.com/sso/saml"/>
  </md:IDPSSODescriptor>
</md:EntityDescriptor>`, testAccV5SamlCertificate())
}

// The self-signed certificate used for testing.
func testAccV5SamlCertificate() string {
	return "MIIDJTCCAg2gAwIBAgIUFlMB5ZIh68pot5nLOxT8/lVS+V0wDQYJKoZIhvcNAQELBQAwIjEgMB4GA1UEAwwXdGVycmFmb3JtLWFjYy1zYW1sLXRlc3QwHhcNMjYxMDI2MDMzNzUwWhcNMzYxMDIzMDMzNzUwWjAiMSAwHgYDVQQDDBd0ZXJyYWZvcm0tYWNjLXNhbWwtdGVzdDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMbUZqE9LNsdKHA4IjWmXjK8NO8AwDuKOgCC135OJsF+jMsJ0aNkHRo7jnxESD+L4TIQn9trERSTJB9cPLK0d+p1uzCyqiM0+qobGka2sqFrGUU+TlXmCVktNRPI+TkNQzCJZiIR9hZlU2Xy+EMFdI9cZbCN08wstqaY/Jv4Eku7TMGTaAhcGhif+tqo3i9hpsV3XwC3ab2A50CtpVKFnnSFe1O0L/TfwSCqdT4KwChLnVTROZhy0U/1iBlzWShg24/8SEVuoPoN/Qi0ioTnhBnKhQA2Ca8joBqzdUb7aY3vv5AzALPUskWn2qCRBBlsaHzjmkFPjQULcu4wrzO0FBECAwEAAaNTMFEwHQYDVR0OBBYEFPUz59O/KXcGfakZU5Bq/OoJO2lwMB8GA1UdIwQYMBaAFPUz59O/KXcGfakZU5Bq/OoJO2lwMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQELBQADggEBAD7SLMMcg+Zl4HzlgxQcysFUAI+MXqmf/OiSfCGQExtXTBkAiNGOGdUEPVXENCbgfYmkO5anYMeR7b7WaX8qI7S+M1goyCZW4VXg4B8SVmywNoqF6LogkUnyJS9SFtbSq0TJotik8/OTF3qeJrjV4MBdgWunZU6++drzMU52253mE7E1CQA9oKndq4PXB5rtN66CApYh+rrTZtqa0c0bz5EUrLalcfA06M8iHTGqcu1O3LQZYFtmPcIVcLMOyMgiqwuJnA/G5YYKH0cvNluD0GAnwITG8wA0W1dF9IYOPB5ExPl/AevgoH2/rbuY5GorAIXfmrdU+XoWwqpn7QThQNQ="
}
