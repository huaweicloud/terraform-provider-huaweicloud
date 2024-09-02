package ccm

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/go-uuid"
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

func TestAccCCMPrivateCA_postpaid_root(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_ccm_private_ca.test_root"
	)

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
				Config: tesPrivateCA_postpaid_root(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "pending_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "type", "ROOT"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "crl_configuration.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.common_name", fmt.Sprintf("%s-root", rName)),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.state", "GD"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.locality", "SZ"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organization", "huawei"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organizational_unit", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.type", "DAY"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.value", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVED"),

					// attributes
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "path_length"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_name"),
				),
			},
			{
				Config: tesPrivateCA_postpaid_rootUpdate1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.obs_bucket_name", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "0"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_dis_point", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_name", ""),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
				),
			},
			{
				Config: tesPrivateCA_postpaid_rootUpdate2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "crl_configuration.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "10"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days", "action", "auto_renew",
				},
			},
		},
	})
}

func tesPrivateCA_postpaid_root(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 5
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
    enabled         = true
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func tesPrivateCA_postpaid_rootUpdate1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  action              = "disable"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 5
  }

  crl_configuration {
    enabled = false
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, name)
}

func tesPrivateCA_postpaid_rootUpdate2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  action              = "disable"
  charging_mode       = "postPaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 5
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "10"
    enabled         = true
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, name)
}

func TestAccCCMPrivateCA_postpaid_subordinate(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_ccm_private_ca.test_subordinate"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesPrivateCA_postpaid_subordinate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "postPaid"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "pending_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "type", "SUBORDINATE"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttrPair(resourceName, "issuer_id", "huaweicloud_ccm_private_ca.test_root", "id"),
					resource.TestCheckResourceAttr(resourceName, "key_usages.0", "cRLSign"),
					resource.TestCheckResourceAttr(resourceName, "path_length", "5"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "crl_configuration.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.common_name", fmt.Sprintf("%s-subordinate", rName)),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.state", "GD"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.locality", "SZ"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organization", "huawei"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organizational_unit", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.type", "DAY"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVED"),

					// attributes
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "issuer_name"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_name"),
				),
			},
			{
				Config: tesPrivateCA_postpaid_subordinateUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.obs_bucket_name", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "0"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_dis_point", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_name", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days", "auto_renew",
				},
			},
		},
	})
}

func tesPrivateCA_postpaid_subordinate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  type                  = "SUBORDINATE"
  issuer_id             = huaweicloud_ccm_private_ca.test_root.id
  key_algorithm         = "RSA2048"
  signature_algorithm   = "SHA512"
  pending_days          = "7"
  charging_mode         = "postPaid"
  auto_renew            = false
  path_length           = 5
  enterprise_project_id = "%[2]s"
  key_usages            = ["cRLSign"]

  distinguished_name {
    common_name         = "%[3]s-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 1
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
    enabled         = true
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, tesPrivateCA_postpaid_root(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func tesPrivateCA_postpaid_subordinateUpdate(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ccm_private_ca" "test_subordinate" {
  type                  = "SUBORDINATE"
  issuer_id             = huaweicloud_ccm_private_ca.test_root.id
  key_algorithm         = "RSA2048"
  signature_algorithm   = "SHA512"
  pending_days          = "7"
  charging_mode         = "postPaid"
  auto_renew            = false
  path_length           = 5
  enterprise_project_id = "%[2]s"
  key_usages            = ["cRLSign"]

  distinguished_name {
    common_name         = "%[3]s-subordinate"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 1
  }

  crl_configuration {
    enabled = false
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, tesPrivateCA_postpaid_root(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func TestAccCCMPrivateCA_prepaid_root(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_ccm_private_ca.test_root"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: tesPrivateCA_prepaid_root(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "auto_renew", "false"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
					resource.TestCheckResourceAttr(resourceName, "key_algorithm", "RSA2048"),
					resource.TestCheckResourceAttr(resourceName, "pending_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "signature_algorithm", "SHA512"),
					resource.TestCheckResourceAttr(resourceName, "type", "ROOT"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "crl_configuration.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.common_name", fmt.Sprintf("%s-root", rName)),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.country", "CN"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.state", "GD"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.locality", "SZ"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organization", "huawei"),
					resource.TestCheckResourceAttr(resourceName, "distinguished_name.0.organizational_unit", "cloud"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.type", "MONTH"),
					resource.TestCheckResourceAttr(resourceName, "validity.0.value", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVED"),

					// attributes
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "expired_at"),
					resource.TestCheckResourceAttrSet(resourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(resourceName, "gen_mode"),
					resource.TestCheckResourceAttrSet(resourceName, "path_length"),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "enterprise_project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_name"),
				),
			},
			{
				Config: tesPrivateCA_prepaid_rootUpdate1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.obs_bucket_name", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "0"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_dis_point", ""),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.crl_name", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo_update", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
				),
			},
			{
				Config: tesPrivateCA_prepaid_rootUpdate2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "action", "enable"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVED"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.enabled", "true"),
					resource.TestCheckResourceAttrPair(resourceName, "crl_configuration.0.obs_bucket_name",
						"huaweicloud_obs_bucket.test", "bucket"),
					resource.TestCheckResourceAttr(resourceName, "crl_configuration.0.valid_days", "30"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_dis_point"),
					resource.TestCheckResourceAttrSet(resourceName, "crl_configuration.0.crl_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days", "action", "auto_renew",
				},
			},
		},
	})
}

func tesPrivateCA_prepaid_root(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                  = "ROOT"
  key_algorithm         = "RSA2048"
  signature_algorithm   = "SHA512"
  pending_days          = "7"
  charging_mode         = "prePaid"
  auto_renew            = false
  enterprise_project_id = "%[2]s"

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "MONTH"
    value = 1
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
    enabled         = true
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func tesPrivateCA_prepaid_rootUpdate1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                  = "ROOT"
  key_algorithm         = "RSA2048"
  signature_algorithm   = "SHA512"
  pending_days          = "7"
  action                = "disable"
  charging_mode         = "prePaid"
  auto_renew            = false
  enterprise_project_id = "%[2]s"

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "MONTH"
    value = 1
  }

  crl_configuration {
    enabled = false
  }

  tags = {
    foo_update = "bar_update"
    key_update = "value_update"
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func tesPrivateCA_prepaid_rootUpdate2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                  = "ROOT"
  key_algorithm         = "RSA2048"
  signature_algorithm   = "SHA512"
  pending_days          = "7"
  action                = "enable"
  charging_mode         = "prePaid"
  auto_renew            = false
  enterprise_project_id = "%[2]s"

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "MONTH"
    value = 1
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "30"
    enabled         = true
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

// TestAccCCMPrivateCA_other using to test some special code
func TestAccCCMPrivateCA_other(t *testing.T) {
	var (
		obj                 interface{}
		rName               = acceptance.RandomAccResourceNameWithDash()
		prepaidResourceName = "huaweicloud_ccm_private_ca.test_root"
	)

	rc := acceptance.InitResourceCheck(
		prepaidResourceName,
		&obj,
		getPrivateCAResourceFunc,
	)

	// test prepaid special code
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      tesPrivateCA_prepaid_other1(rName),
				ExpectError: regexp.MustCompile("only `YEAR` or `MONTH` is supported when creating a prepaid private CA"),
			},
			{
				Config: tesPrivateCA_prepaid_other2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
				),
			},
			{
				ResourceName:      prepaidResourceName,
				ExpectError:       regexp.MustCompile("Cannot import non-existent remote object"),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testOtherCaseImportIdFunc(),
				ImportStateVerifyIgnore: []string{
					"validity", "key_usages", "pending_days", "action", "auto_renew",
				},
			},
		},
	})
}

// testOtherCaseImportIdFunc using to import a non-exist resource ID
func testOtherCaseImportIdFunc() resource.ImportStateIdFunc {
	return func(_ *terraform.State) (string, error) {
		randUUID, err := uuid.GenerateUUID()
		if err != nil {
			return "", fmt.Errorf("error generating uuid: %s", err)
		}
		return randUUID, nil
	}
}

func tesPrivateCA_prepaid_other1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ccm_private_ca" "prepaid_root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "prePaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "DAY"
    value = 2
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}

func tesPrivateCA_prepaid_other2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%[1]s-crl-bucket"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ccm_private_ca" "test_root" {
  type                = "ROOT"
  key_algorithm       = "RSA2048"
  signature_algorithm = "SHA512"
  pending_days        = "7"
  charging_mode       = "prePaid"
  auto_renew          = false

  distinguished_name {
    common_name         = "%[1]s-root"
    country             = "CN"
    state               = "GD"
    locality            = "SZ"
    organization        = "huawei"
    organizational_unit = "cloud"
  }

  validity {
    type  = "MONTH"
    value = 1
  }

  crl_configuration {
    obs_bucket_name = huaweicloud_obs_bucket.test.bucket
    valid_days      = "7"
    enabled         = true
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, name)
}
