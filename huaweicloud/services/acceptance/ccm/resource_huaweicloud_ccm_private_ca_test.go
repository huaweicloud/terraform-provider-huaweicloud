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

func getPrivateCAResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("ccm", state.Primary.Attributes["region"])
	if err != nil {
		return nil, fmt.Errorf("error creating CCM client: %s", err)
	}

	getPrivateCAHttpUrl := "v1/private-certificate-authorities/{id}"
	getPrivateCAPath := client.Endpoint + getPrivateCAHttpUrl
	getPrivateCAPath = strings.ReplaceAll(getPrivateCAPath, "{id}", state.Primary.ID)
	getPrivateCAOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPrivateCAResp, err := client.Request("GET", getPrivateCAPath, &getPrivateCAOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CCM private CA: %s", err)
	}
	getPrivateCARespBody, err := utils.FlattenResponse(getPrivateCAResp)
	if err != nil {
		return nil, fmt.Errorf("error prase CCM private CA: %s", err)
	}

	status := utils.PathSearch("status", getPrivateCARespBody, nil).(string)
	if status == "DELETED" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getPrivateCARespBody, nil
}

func TestAccCCMPrivateCA_basic(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ccm_private_ca.test_subordinate"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesPrivateCA_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "SUBORDINATE"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "charging_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
				),
			},
			{
				Config: testPrivateCA_updateTags(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days",
				},
			},
		},
	})
}

// lintignore:AT004
func tesPrivateCA_base(commonName string) string {
	return fmt.Sprintf(`
provider "huaweicloud" {
  endpoints = {
    ccm = "https://ccm.cn-north-4.myhuaweicloud.com/"
  }
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type = "ROOT"
  distinguished_name {
    common_name         = "%s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  validity {
    type  = "DAY"
    value = 5
  }
}`, commonName)
}

// lintignore:AT004
func tesPrivateCA_basic(commonName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  type = "SUBORDINATE"
  distinguished_name {
    common_name         = "%[2]s-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  signature_algorithm = "SHA512"
  pending_days        = "7"
  validity {
    type  = "DAY"
    value = 1
  }
  tags = {
    foo = "bar"
    key = "value"
  }
  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
  }
}`, tesPrivateCA_base(commonName), commonName)
}

// lintignore:AT004
func testPrivateCA_updateTags(commonName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[2]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}
  
resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  type = "SUBORDINATE"
  distinguished_name {
    common_name         = "%[2]s-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  signature_algorithm = "SHA512"
  pending_days        = "7"
  validity {
    type  = "DAY"
    value = 1
  }
  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
  }
}`, tesPrivateCA_base(commonName), commonName)
}

func TestAccCCMPrivateCA_prePaid(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_ccm_private_ca.test_subordinate"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesPrivateCA_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "SUBORDINATE"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer_name"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days",
				},
			},
		},
	})
}

// lintignore:AT004
func tesPrivateCA_prePaid(commonName string) string {
	return fmt.Sprintf(`
provider "huaweicloud" {
  endpoints = {
    ccm = "https://ccm.cn-north-4.myhuaweicloud.com/"
  }
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type = "ROOT"
  distinguished_name {
    common_name         = "%s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  validity {
    type  = "MONTH"
    value = 2
  }
  charging_mode = "prePaid"
}

resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  type = "SUBORDINATE"
  distinguished_name {
    common_name         = "%s-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  issuer_id           = huaweicloud_ccm_private_ca.test_root.id
  signature_algorithm = "SHA512"
  pending_days        = "7"
  validity {
    type  = "MONTH"
    value = 1
  }
  charging_mode = "prePaid"
}`, commonName, commonName)
}

func TestAccCCMPrivateCA_withoutParentCA(t *testing.T) {
	var obj interface{}
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_ccm_private_ca.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesPrivateCA_withoutParentCA(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "ROOT"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "charging_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
				),
			},
			{
				Config: tesPrivateCA_withoutParentCA_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "tags.foo1", "bar1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "action", "enable"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVED"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days", "action",
				},
			},
		},
	})
}

// lintignore:AT004
func tesPrivateCA_withoutParentCA(commonName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "test" {
  type = "ROOT"
  distinguished_name {
    common_name         = "%s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  action              = "disable"
  validity {
    type  = "DAY"
    value = 5
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}`, commonName)
}

// lintignore:AT004
func tesPrivateCA_withoutParentCA_update(commonName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "test" {
  type = "ROOT"
  distinguished_name {
    common_name         = "%s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  action              = "enable"
  validity {
    type  = "DAY"
    value = 5
  }

  tags = {
    foo1 = "bar1"
    key1 = "value1"
  }
}`, commonName)
}
