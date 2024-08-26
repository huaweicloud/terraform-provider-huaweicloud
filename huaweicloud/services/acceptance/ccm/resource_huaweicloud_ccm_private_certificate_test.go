package ccm

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

func getCertificateResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("ccm", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}

	getCertificateHttpUrl := "v1/private-certificates/{certificate_id}"
	getCertificatePath := client.Endpoint + getCertificateHttpUrl
	getCertificatePath = strings.ReplaceAll(getCertificatePath, "{certificate_id}", state.Primary.ID)

	getCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getCertificateResp, err := client.Request("GET", getCertificatePath, &getCertificateOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCM private certificate: %s", err)
	}

	getCertificateRespBody, err := utils.FlattenResponse(getCertificateResp)
	if err != nil {
		return nil, fmt.Errorf("error prase CCM private certificate: %s", err)
	}

	return getCertificateRespBody, nil
}

func TestAccPrivateCertificate_basic(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_ccm_private_certificate.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getCertificateResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testPrivateCertificate_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.common_name", rName),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organization", "huawei"),

					resource.TestCheckResourceAttr(resourceName, "tags.fooo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.keye", "value"),

					resource.TestCheckResourceAttrSet(resourceName, "issuer_name"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "start_at"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testPrivateCertificate_tagsUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"validity",
					"key_usage",
					"server_auth",
					"client_auth",
					"code_signing",
					"email_protection",
					"time_stamping",
					"object_identifier",
					"object_identifier_value",
					"subject_alternative_names",
				},
			},
		},
	})
}

// lintignore:AT004
func testPrivateCertificate_basic(commonName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ccm_private_certificate" "test" {
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA256"
  key_usage           = ["nonRepudiation", "keyEncipherment"]

  distinguished_name {
    common_name         = "%s"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = "1"
  }

  subject_alternative_names {
    type  = "IP"
    value = "156.127.116.38"
  }

  tags = {
    fooo = "bar"
    keye = "value"
  }
}`, tesPrivateCA_postpaid_root(commonName), commonName)
}

// lintignore:AT004
func testPrivateCertificate_tagsUpdate(commonName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ccm_private_certificate" "test" {
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA256"
  key_usage           = ["nonRepudiation", "keyEncipherment"]

  distinguished_name {
    common_name         = "%s"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = "1"
  }

  subject_alternative_names {
    type  = "IP"
    value = "156.127.116.38"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, tesPrivateCA_postpaid_root(commonName), commonName)
}
