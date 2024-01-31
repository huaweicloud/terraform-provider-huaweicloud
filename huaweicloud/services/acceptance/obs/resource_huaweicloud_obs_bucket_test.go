package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccObsBucket_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucket_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "bucket_domain_name", testAccObsBucketDomainName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "multi_az", "false"),
					resource.TestCheckResourceAttr(resourceName, "parallel_fs", "false"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "region", acceptance.HW_REGION_NAME),
					resource.TestCheckResourceAttr(resourceName, "bucket_version", "3.0"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "storage_info.0.size", "0"),
					resource.TestCheckResourceAttr(resourceName, "storage_info.0.object_number", "0"),
				),
			},
			{
				Config: testAccObsBucket_basic_update(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "WARM"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"acl",
					"force_destroy",
				},
			},
		},
	})
}

func TestAccObsBucket_withEpsId(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
			acceptance.TestAccPreCheckMigrateEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucket_epsId(rInt, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccObsBucket_epsId(rInt, acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_MIGRATE_PROJECT_ID_TEST),
				),
			},
			{
				Config: testAccObsBucket_epsId(rInt, "0"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "enterprise_project_id", "0"),
				),
			},
		},
	})
}

func TestAccObsBucket_multiAZ(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigMultiAZ(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "parallel_fs", "false"),
					resource.TestCheckResourceAttr(resourceName, "tags.multi_az", "3az"),
				),
			},
		},
	})
}

func TestAccObsBucket_parallelFS(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigParallelFS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "parallel_fs", "true"),
					resource.TestCheckResourceAttr(resourceName, "multi_az", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.parallel_fs", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.multi_az", "3az"),
					resource.TestCheckNoResourceAttr(resourceName, "bucket_version"),
				),
			},
		},
	})
}

func TestAccObsBucket_encryption(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucket_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "bucket_domain_name", testAccObsBucketDomainName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "false"),
				),
			},
			{
				Config: testAccObsBucket_encryption(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "kms_key_id"),
				),
			},
			{
				Config: testAccObsBucket_encryptionByDefaultKMSMasterKey(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "kms_key_id", ""),
				),
			},
			{
				Config: testAccObsBucket_encryptionBySSEObsMode(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "sse_algorithm", "AES256"),
				),
			},
			{
				Config: testAccObsBucket_updateEncryptionToSSEKMSMode(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "encryption", "true"),
					resource.TestCheckResourceAttr(resourceName, "sse_algorithm", "kms"),
				),
			},
			{
				Config: testAccObsBucket_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "bucket_domain_name", testAccObsBucketDomainName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "encryption", "false"),
					resource.TestCheckResourceAttr(resourceName, "sse_algorithm", ""),
				),
			},
		},
	})
}

func TestAccObsBucket_versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "versioning", "true"),
				),
			},
			{
				Config: testAccObsBucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "versioning", "false"),
				),
			},
		},
	})
}

func TestAccObsBucket_logging(t *testing.T) {
	rInt := acctest.RandInt()
	targetBucket := fmt.Sprintf("tf-test-log-bucket-%d", rInt)
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithLogging(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					testAccCheckObsBucketLogging(resourceName, targetBucket, "log/", "live_to_obs"),
				),
			},
		},
	})
}

func TestAccObsBucket_quota(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithQuota(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "quota", "1000000000"),
				),
			},
		},
	})
}

func TestAccObsBucket_lifecycle(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.name", "rule1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.name", "rule2"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.name", "rule3"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.prefix", ""),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.transition.0.days", "30"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.transition.1.days", "180"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.noncurrent_version_transition.0.days", "60"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.noncurrent_version_transition.1.days", "180"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.abort_incomplete_multipart_upload.0.days", "360"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.abort_incomplete_multipart_upload.0.days", "2147483647"),
				),
			},
		},
	})
}

func TestAccObsBucket_website(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketWebsiteConfigWithRoutingRules(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "website.0.index_document", "index.html"),
					resource.TestCheckResourceAttr(
						resourceName, "website.0.error_document", "error.html"),
				),
			},
		},
	})
}

func TestAccObsBucket_cors(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "cors_rule.0.allowed_origins.0", "https://www.example.com"),
					resource.TestCheckResourceAttr(
						resourceName, "cors_rule.0.allowed_methods.0", "PUT"),
					resource.TestCheckResourceAttr(
						resourceName, "cors_rule.0.allowed_headers.0", "*"),
					resource.TestCheckResourceAttr(
						resourceName, "cors_rule.0.expose_headers.1", "ETag"),
					resource.TestCheckResourceAttr(
						resourceName, "cors_rule.0.max_age_seconds", "3000"),
				),
			},
		},
	})
}

func TestAccObsBucket_userDomainNames(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "huaweicloud_obs_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
			acceptance.TestAccPreCheckOBSUserDomainNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckObsBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObsBucketConfigWithUserDomainNames(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "user_domain_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_domain_names.0",
						acceptance.HW_OBS_USER_DOMAIN_NAME1),
				),
			},
			{
				Config: testAccObsBucketConfigWithUserDomainNames_update1(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "user_domain_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "user_domain_names.0",
						acceptance.HW_OBS_USER_DOMAIN_NAME2),
				),
			},
			{
				Config: testAccObsBucketConfigWithUserDomainNames_update2(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckObsBucketExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "bucket", testAccObsBucketName(rInt)),
					resource.TestCheckResourceAttr(resourceName, "user_domain_names.#", "0"),
				),
			},
		},
	})
}

func testAccCheckObsBucketDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	obsClient, err := conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OBS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_obs_bucket" {
			continue
		}

		_, err := obsClient.HeadBucket(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("OBS Bucket %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckObsBucketExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		obsClient, err := conf.ObjectStorageClient(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OBS client: %s", err)
		}

		_, err = obsClient.HeadBucket(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("OBS Bucket not found: %v", err)
		}
		return nil
	}
}

func testAccCheckObsBucketLogging(name, target, prefix, agency string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		obsClient, err := conf.ObjectStorageClientWithSignature(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating OBS client with signature: %s", err)
		}

		output, err := obsClient.GetBucketLoggingConfiguration(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error getting logging configuration of OBS bucket: %s", err)
		}

		if output.TargetBucket != target {
			return fmt.Errorf("%s.logging: Attribute 'target_bucket' expected %s, got %s",
				name, output.TargetBucket, target)
		}
		if output.TargetPrefix != prefix {
			return fmt.Errorf("%s.logging: Attribute 'target_prefix' expected %s, got %s",
				name, output.TargetPrefix, prefix)
		}

		if output.Agency != agency {
			return fmt.Errorf("%s.logging: Attribute 'agency' expected %s, got %s",
				name, output.Agency, agency)
		}

		return nil
	}
}

// These need a bit of randomness as the name can only be used once globally
func testAccObsBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccObsBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.obs.%s.myhuaweicloud.com", randInt, acceptance.HW_REGION_NAME)
}

func testAccObsBucket_basic(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, randInt)
}

func testAccObsBucket_basic_update(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "WARM"
  acl           = "public-read"

  tags = {
    owner = "terraform"
    key   = "value1"
  }
}
`, randInt)
}

func testAccObsBucket_encryption(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "key_1" {
  key_alias    = "kms-%d"
  pending_days = "7"
}

resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
  encryption    = true
  kms_key_id    = huaweicloud_kms_key.key_1.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, randInt, randInt)
}

func testAccObsBucket_encryptionByDefaultKMSMasterKey(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
  encryption    = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, randInt)
}

func testAccObsBucket_encryptionBySSEObsMode(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
  encryption    = true
  sse_algorithm = "AES256"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, randInt)
}

func testAccObsBucket_updateEncryptionToSSEKMSMode(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "tf-test-bucket-%d"
  storage_class = "STANDARD"
  acl           = "private"
  encryption    = true
  sse_algorithm = "kms"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, randInt)
}

func testAccObsBucket_epsId(randInt int, enterpriseProjectId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket                = "tf-test-bucket-%d"
  storage_class         = "STANDARD"
  acl                   = "private"
  enterprise_project_id = "%s"
}
`, randInt, enterpriseProjectId)
}

func testAccObsBucketConfigMultiAZ(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket   = "tf-test-bucket-%d"
  acl      = "private"
  multi_az = true

  tags = {
    key      = "value"
    multi_az = "3az"
  }
}
`, randInt)
}

func testAccObsBucketConfigParallelFS(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket      = "tf-test-bucket-%d"
  acl         = "private"
  multi_az    = true
  parallel_fs = true

  tags = {
    parallel_fs = "true"
    multi_az    = "3az"
  }
}
`, randInt)
}

func testAccObsBucketConfigWithVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket     = "tf-test-bucket-%d"
  acl        = "private"
  versioning = true
}
`, randInt)
}

func testAccObsBucketConfigWithDisableVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket     = "tf-test-bucket-%d"
  acl        = "private"
  versioning = false
}
`, randInt)
}

func testAccObsBucketConfigWithLogging(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "log_bucket" {
  bucket        = "tf-test-log-bucket-%d"
  acl           = "log-delivery-write"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "private"

  logging {
    target_bucket = huaweicloud_obs_bucket.log_bucket.id
    target_prefix = "log/"
    agency        = "live_to_obs"
  }
}
`, randInt, randInt)
}

func testAccObsBucketConfigWithQuota(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "private"
  quota  = 1000000000
}
`, randInt)
}

func testAccObsBucketConfigWithLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket     = "tf-test-bucket-%d"
  acl        = "private"
  versioning = true

  lifecycle_rule {
    name    = "rule1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
    abort_incomplete_multipart_upload {
      days = 360
    }
  }
  lifecycle_rule {
    name    = "rule2"
    prefix  = "path2/"
    enabled = true

    expiration {
      days = 365
    }

    transition {
      days          = 30
      storage_class = "WARM"
    }
    transition {
      days          = 180
      storage_class = "COLD"
    }
    abort_incomplete_multipart_upload {
      days = 2147483647
    }
  }
  lifecycle_rule {
    name    = "rule3"
    enabled = true

    noncurrent_version_expiration {
      days = 365
    }
    noncurrent_version_transition {
      days          = 60
      storage_class = "WARM"
    }
    noncurrent_version_transition {
      days          = 180
      storage_class = "COLD"
    }
  }
}
`, randInt)
}

func testAccObsBucketWebsiteConfigWithRoutingRules(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  website {
    index_document = "index.html"
    error_document = "error.html"
    routing_rules = <<EOF
[{
  "Condition": {
    "KeyPrefixEquals": "docs/"
  },
  "Redirect": {
    "ReplaceKeyPrefixWith": "documents/"
  }
}]
EOF
  }
}
`, randInt)
}

func testAccObsBucketConfigWithCORS(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT","POST"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["x-amz-server-side-encryption","ETag"]
    max_age_seconds = 3000
  }
}
`, randInt)
}

func testAccObsBucketConfigWithUserDomainNames(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  user_domain_names = ["%s"]
}
`, randInt, acceptance.HW_OBS_USER_DOMAIN_NAME1)
}

func testAccObsBucketConfigWithUserDomainNames_update1(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  user_domain_names = ["%s"]
}
`, randInt, acceptance.HW_OBS_USER_DOMAIN_NAME2)
}

func testAccObsBucketConfigWithUserDomainNames_update2(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket = "tf-test-bucket-%d"
  acl    = "public-read"

  user_domain_names = []
}
`, randInt)
}
