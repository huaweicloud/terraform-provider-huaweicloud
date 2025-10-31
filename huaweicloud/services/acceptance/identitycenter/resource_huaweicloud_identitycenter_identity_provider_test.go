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

func getIdentityCenterIdentityProviderResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		listHttpUrl = "v1/identity-stores/{identity_store_id}/external-idp"
		listProduct = "identitystore"
	)
	client, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Identity Center Client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{identity_store_id}", state.Primary.Attributes["identity_store_id"])

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	idpId := state.Primary.ID

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Identity Center identity provider: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	association := utils.PathSearch(fmt.Sprintf("associations[?idp_id=='%s']|[0]", idpId), listRespBody, nil)
	if association == nil {
		return nil, fmt.Errorf("error get Identity Center identity provider")
	}
	return association, nil
}

func TestAccIdentityProvider_basic(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_identitycenter_identity_provider.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getIdentityCenterIdentityProviderResourceFunc,
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
				Config: testIdentityCenterIdentityProvider_create_without_metadata_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "entity_id", "https://create.entity.com"),
					resource.TestCheckResourceAttr(rName, "login_url", "https://create.login.com"),
					resource.TestCheckResourceAttrSet(rName, "is_enabled"),
					resource.TestCheckResourceAttrSet(rName, "want_request_signed"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate_ids.0.certificate_id"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate_ids.0.status"),
				),
			},
			{
				Config: testIdentityCenterIdentityProvider_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "identity_store_id",
						"data.huaweicloud_identitycenter_instance.test", "identity_store_id"),
					resource.TestCheckResourceAttr(rName, "entity_id", "https://update.entity.com"),
					resource.TestCheckResourceAttr(rName, "login_url", "https://update.login.com"),
					resource.TestCheckResourceAttrSet(rName, "is_enabled"),
					resource.TestCheckResourceAttrSet(rName, "want_request_signed"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate_ids.0.certificate_id"),
					resource.TestCheckResourceAttrSet(rName, "idp_certificate_ids.0.status"),
				),
			},
			{
				ResourceName:            rName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       testIdentityCenterIdentityProviderImportState(rName),
				ImportStateVerifyIgnore: []string{"idp_certificate"},
			},
		},
	})
}

func testIdentityCenterIdentityProvider_create_without_metadata_basic() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_identity_provider" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  idp_certificate   = "-----BEGIN CERTIFICATE-----MIIDUTCCAjmgAwIBAgIQAP9wc90YPLxcirh7/qTyBTANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAliYWlkdS5jb20wHhcNMjUwOTAzMDI0NzAzWhcNMzUwOTAxMDI0NzAzWjAUMRIwEAYDVQQDDAliYWlkdS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDZGZlyBZTE6I3RCsCdOTNvdT0rn0dzhO3squ0owQjpO5APLEVtu3rOeImkOpGIKrsx1s0Y17F+/DP2eB+YTkIOtkn01+XK9SUerL2kcAJQFMbNBtyihOxGPvq5S55NneZGY7UaPab+C9ugurmqi9xzH1NYdiDgb5Aa58t7XL2G0mnd4OkRYd43fY4lpo9jhmGhOUKZBFequ1TgV6EQg3FPTNpDzUN2skucmoU0Rz0Btz1bKdKtQ7tCHffFIPEI63dR6Rmi7CQT5WEm8f9j5UPhxFWC7f65PuaLbcTomccBbUXVm0JxaDI3L06yZyQ4U9ZKnLOSGSaZGs5hgtmeZlCNAgMBAAGjgZ4wgZswHQYDVR0OBBYEFIpds0vTmKvtGLhem/BxatDeyC+7MA4GA1UdDwEB/wQEAwIEsDAMBgNVHRMBAf8EAjAAMDsGA1UdJQQ0MDIGCCsGAQUFBwMCBggrBgEFBQcDAQYIKwYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcDCDAfBgNVHSMEGDAWgBSKXbNL05ir7Ri4XpvwcWrQ3sgvuzANBgkqhkiG9w0BAQsFAAOCAQEARbgGdia5b0QNDxQyitntXq+Gn5JAK5Lx/5JYL0/RsgJ7uT6kVveMN6ySyDclZW2/Pvf3weCJdiz0h8NIAY0TIKd9NNJ53YCAyGTzi/BNPXAwAFJztyEwGtTIBCSIHwNmHifYzfrEFBdy33LY6xBO+W98d9NyyOppstFbRtgz4WCEdGJxRDNQ2h4oZJcIloDj54WXFyEulibbieC4oIyVP58j2MXUZwXYrhfnlir/qtaQTudjcA43+YorkTP2CBDCONm9vjINy7mDF7dTdFDjuUMyWPqokuvqLVB7zHZpKu/QhfsOBNMKgxTiNfHgqQe+EFxwvhXnxXZnkmd7F1pXqw==-----END CERTIFICATE-----"
  entity_id         = "https://create.entity.com"
  login_url         = "https://create.login.com"

}
`, testAccDatasourceIdentityCenter_basic())
}

func testIdentityCenterIdentityProvider_basic_update() string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_identitycenter_identity_provider" "test" {
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
  idp_certificate   = "-----BEGIN CERTIFICATE-----MIIDUTCCAjmgAwIBAgIQAP9wc90YPLxcirh7/qTyBTANBgkqhkiG9w0BAQsFADAUMRIwEAYDVQQDDAliYWlkdS5jb20wHhcNMjUwOTAzMDI0NzAzWhcNMzUwOTAxMDI0NzAzWjAUMRIwEAYDVQQDDAliYWlkdS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDZGZlyBZTE6I3RCsCdOTNvdT0rn0dzhO3squ0owQjpO5APLEVtu3rOeImkOpGIKrsx1s0Y17F+/DP2eB+YTkIOtkn01+XK9SUerL2kcAJQFMbNBtyihOxGPvq5S55NneZGY7UaPab+C9ugurmqi9xzH1NYdiDgb5Aa58t7XL2G0mnd4OkRYd43fY4lpo9jhmGhOUKZBFequ1TgV6EQg3FPTNpDzUN2skucmoU0Rz0Btz1bKdKtQ7tCHffFIPEI63dR6Rmi7CQT5WEm8f9j5UPhxFWC7f65PuaLbcTomccBbUXVm0JxaDI3L06yZyQ4U9ZKnLOSGSaZGs5hgtmeZlCNAgMBAAGjgZ4wgZswHQYDVR0OBBYEFIpds0vTmKvtGLhem/BxatDeyC+7MA4GA1UdDwEB/wQEAwIEsDAMBgNVHRMBAf8EAjAAMDsGA1UdJQQ0MDIGCCsGAQUFBwMCBggrBgEFBQcDAQYIKwYBBQUHAwMGCCsGAQUFBwMEBggrBgEFBQcDCDAfBgNVHSMEGDAWgBSKXbNL05ir7Ri4XpvwcWrQ3sgvuzANBgkqhkiG9w0BAQsFAAOCAQEARbgGdia5b0QNDxQyitntXq+Gn5JAK5Lx/5JYL0/RsgJ7uT6kVveMN6ySyDclZW2/Pvf3weCJdiz0h8NIAY0TIKd9NNJ53YCAyGTzi/BNPXAwAFJztyEwGtTIBCSIHwNmHifYzfrEFBdy33LY6xBO+W98d9NyyOppstFbRtgz4WCEdGJxRDNQ2h4oZJcIloDj54WXFyEulibbieC4oIyVP58j2MXUZwXYrhfnlir/qtaQTudjcA43+YorkTP2CBDCONm9vjINy7mDF7dTdFDjuUMyWPqokuvqLVB7zHZpKu/QhfsOBNMKgxTiNfHgqQe+EFxwvhXnxXZnkmd7F1pXqw==-----END CERTIFICATE-----"
  entity_id         = "https://update.entity.com"
  login_url         = "https://update.login.com"
}
`, testAccDatasourceIdentityCenter_basic())
}

func testIdentityCenterIdentityProviderImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		identityStoreId := rs.Primary.Attributes["identity_store_id"]
		if identityStoreId == "" {
			return "", fmt.Errorf("attribute (identity_store_id) of Resource (%s) not found: %s", name, rs)
		}
		return identityStoreId + "/" + rs.Primary.ID, nil
	}
}
