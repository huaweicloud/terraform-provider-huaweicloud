package identitycenter

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

func getIdentityCenterIdentityProviderCertificateResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp/{idp_id}/certificate"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}",
		state.Primary.Attributes["identity_store_id"])
	listPath = strings.ReplaceAll(listPath, "{idp_id}",
		state.Primary.Attributes["idp_id"])

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	certificateId := state.Primary.ID

	listResp, err := client.Request("GET",
		listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center identity provider certificate: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	certificate := utils.PathSearch(fmt.Sprintf("idp_certificates[?certificate_id =='%s']|[0]", certificateId), listRespBody, nil)
	if certificate == nil {
		return nil, fmt.Errorf("error get Identity Center identity provider certificate")
	}
	return certificate, nil
}

func TestAccIdentityProviderCertificate_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_identity_provider_certificate.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterIdentityProviderCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testIdentityCenterIdentityProviderCertificate_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttrSet(rName, "issuer_name"),
					resource.TestCheckResourceAttrSet(rName, "not_after"),
					resource.TestCheckResourceAttrSet(rName, "not_before"),
					resource.TestCheckResourceAttrSet(rName, "public_key"),
					resource.TestCheckResourceAttrSet(rName, "serial_number_string"),
					resource.TestCheckResourceAttrSet(rName, "subject_name"),
					resource.TestCheckResourceAttrSet(rName, "version"),
					resource.TestCheckResourceAttrSet(rName, "signature_algorithm_name"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterIdentityProviderCertificateImportState(rName),
				ImportStateVerifyIgnore: []string{"certificate_use"},
			},
		},
	})
}

func testIdentityCenterIdentityProviderCertificate_basic() string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_identitycenter_identity_provider_certificate" "test"{
  identity_store_id       = data.huaweicloud_identitycenter_instance.test.identity_store_id
  idp_id                  = huaweicloud_identitycenter_identity_provider.test.id
  x509_certificate_in_pem = "-----BEGIN CERTIFICATE-----MIIDUTCCAjmgAwIBAgIQAP9wc90YPLxcirh7/qTyBTANBgkqhkiG9w0BA`+
		`QsFADAUMRIwEAYDVQQDDAliYWlkdS5jb20wHhcNMjUwOTAzMDI0NzAzWhcNMzUwOTAxMDI0NzAzWjAUMRIwEAYDVQQDDAliYWlkdS5jb`+
		`20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDZGZlyBZTE6I3RCsCdOTNvdT0rn0dzhO3squ0owQjpO5APLEVtu3rOeImkOp`+
		`GIKrsx1s0Y17F+/DP2eB+YTkIOtkn01+XK9SUerL2kcAJQFMbNBtyihOxGPvq5S55NneZGY7UaPab+C9ugurmqi9xzH1NYdiDgb5Aa58t7`+
		`XL2G0mnd4OkRYd43fY4lpo9jhmGhOUKZBFequ1TgV6EQg3FPTNpDzUN2skucmoU0Rz0Btz1bKdKtQ7tCHffFIPEI63dR6Rmi7CQT5WEm8f`+
		`9j5UPhxFWC7f65PuaLbcTomccBbUXVm0JxaDI3L06yZyQ4U9ZKnLOSGSaZGs5hgtmeZlCNAgMBAAGjgZ4wgZswHQYDVR0OBBYEFIpds0v`+
		`TmKvtGLhem/BxatDeyC+7MA4GA1UdDwEB/wQEAwIEsDAMBgNVHRMBAf8EAjAAMDsGA1UdJQQ0MDIGCCsGAQUFBwMCBggrBgEFBQcDAQYIK`+
		`wYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcDCDAfBgNVHSMEGDAWgBSKXbNL05ir7Ri4XpvwcWrQ3sgvuzANBgkqhkiG9w0BAQsFAAOCAQEA`+
		`RbgGdia5b0QNDxQyitntXq+Gn5JAK5Lx/5JYL0/RsgJ7uT6kVveMN6ySyDclZW2/Pvf3weCJdiz0h8NIAY0TIKd9NNJ53YCAyGTzi/BNPX`+
		`AwAFJztyEwGtTIBCSIHwNmHifYzfrEFBdy33LY6xBO+W98d9NyyOppstFbRtgz4WCEdGJxRDNQ2h4oZJcIloDj54WXFyEulibbieC4oIyVP`+
		`58j2MXUZwXYrhfnlir/qtaQTudjcA43+YorkTP2CBDCONm9vjINy7mDF7dTdFDjuUMyWPqokuvqLVB7zHZpKu/QhfsOBNMKgxTiNfHgqQe`+
		`+EFxwvhXnxXZnkmd7F1pXqw==-----END CERTIFICATE-----"
  certificate_use         = "SIGNING"
}
`, testIdentityCenterIdentityProvider_create_without_metadata_basic())
}

func testIdentityCenterIdentityProviderCertificateImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		idpID := rs.Primary.Attributes["idp_id"]
		if idpID == "" {
			return "", fmt.Errorf("attribute (idp_id) of Resource (%s) not found: %s", name, rs)
		}
		certificateID := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", identityStoreId, idpID, certificateID), nil
	}
}
