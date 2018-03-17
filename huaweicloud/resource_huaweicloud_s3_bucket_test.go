package huaweicloud

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	//"github.com/aws/aws-sdk-go/service/ec2"
	"bytes"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

// PASS
func TestAccS3Bucket_basic(t *testing.T) {
	rInt := acctest.RandInt()
	//arnRegexp := regexp.MustCompile("^arn:aws:s3:::")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		/*
			IDRefreshName:   "huaweicloud_s3_bucket.bucket",
			IDRefreshIgnore: []string{"force_destroy"},
		*/
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					/*resource.TestCheckResourceAttr(
					"huaweicloud_s3_bucket.bucket", "hosted_zone_id", HostedZoneIDForRegion("us-west-2")), */
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "region", OS_REGION_NAME),
					resource.TestCheckNoResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint"),
					/*resource.TestMatchResourceAttr(
					"huaweicloud_s3_bucket.bucket", "arn", arnRegexp), */
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "bucket", testAccBucketName(rInt)),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "bucket_domain_name", testAccBucketDomainName(rInt)),
				),
			},
		},
	})
}

func TestAccAWSS3MultiBucket_withTags(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccCheckAWSS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSS3MultiBucketConfigWithTags(rInt),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_namePrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig_namePrefix,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.test"),
					resource.TestMatchResourceAttr(
						"huaweicloud_s3_bucket.test", "bucket", regexp.MustCompile("^tf-test-")),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_generatedName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig_generatedName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.test"),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_region(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithRegion(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr("huaweicloud_s3_bucket.bucket", "region", OS_REGION_NAME),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_Policy(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketPolicy(
						"huaweicloud_s3_bucket.bucket", testAccS3BucketPolicy(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketPolicy(
						"huaweicloud_s3_bucket.bucket", ""),
				),
			},
			{
				Config: testAccS3BucketConfigWithEmptyPolicy(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketPolicy(
						"huaweicloud_s3_bucket.bucket", ""),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_UpdateAcl(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccS3BucketConfigWithAcl, ri)
	postConfig := fmt.Sprintf(testAccS3BucketConfigWithAclUpdate, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "acl", "public-read"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "acl", "private"),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_Website_Simple(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "index.html", "", "", ""),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketWebsiteConfigWithError(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "index.html", "error.html", "", ""),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "", "", "", ""),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_WebsiteRedirect(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfigWithRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "", "", "", "hashicorp.com"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketWebsiteConfigWithHttpsRedirect(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "", "", "https", "hashicorp.com"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "", "", "", ""),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_WebsiteRoutingRules(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketWebsiteConfigWithRoutingRules(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "index.html", "error.html", "", ""),
					testAccCheckS3BucketWebsiteRoutingRules(
						"huaweicloud_s3_bucket.bucket",
						[]*s3.RoutingRule{
							{
								Condition: &s3.Condition{
									KeyPrefixEquals: aws.String("docs/"),
								},
								Redirect: &s3.Redirect{
									ReplaceKeyPrefixWith: aws.String("documents/"),
								},
							},
						},
					),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", testAccWebsiteEndpoint(rInt)),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketWebsite(
						"huaweicloud_s3_bucket.bucket", "", "", "", ""),
					testAccCheckS3BucketWebsiteRoutingRules("huaweicloud_s3_bucket.bucket", nil),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "website_endpoint", ""),
				),
			},
		},
	})
}

// Test TestAccAWSS3Bucket_shouldFailNotFound is designed to fail with a "plan
// not empty" error in Terraform, to check against regresssions.
// See https://github.com/hashicorp/terraform/pull/2925
// PASS
func TestAccS3Bucket_shouldFailNotFound(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketDestroyedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3DestroyBucket("huaweicloud_s3_bucket.bucket"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// PASS
func TestAccS3Bucket_Versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketVersioning(
						"huaweicloud_s3_bucket.bucket", ""),
				),
			},
			{
				Config: testAccS3BucketConfigWithVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketVersioning(
						"huaweicloud_s3_bucket.bucket", s3.BucketVersioningStatusEnabled),
				),
			},
			{
				Config: testAccS3BucketConfigWithDisableVersioning(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketVersioning(
						"huaweicloud_s3_bucket.bucket", s3.BucketVersioningStatusSuspended),
				),
			},
		},
	})
}

// PASS
func TestAccS3Bucket_Cors(t *testing.T) {
	rInt := acctest.RandInt()

	updateBucketCors := func(n string) resource.TestCheckFunc {
		return func(s *terraform.State) error {
			rs, ok := s.RootModule().Resources[n]
			if !ok {
				return fmt.Errorf("Not found: %s", n)
			}

			config := testAccProvider.Meta().(*Config)
			conn, err := config.computeS3conn(OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
			}
			_, err = conn.PutBucketCors(&s3.PutBucketCorsInput{
				Bucket: aws.String(rs.Primary.ID),
				CORSConfiguration: &s3.CORSConfiguration{
					CORSRules: []*s3.CORSRule{
						{
							AllowedHeaders: []*string{aws.String("*")},
							AllowedMethods: []*string{aws.String("GET")},
							AllowedOrigins: []*string{aws.String("https://www.example.com")},
						},
					},
				},
			})
			if err != nil {
				if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() != "NoSuchCORSConfiguration" {
					return err
				}
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketCors(
						"huaweicloud_s3_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
					updateBucketCors("huaweicloud_s3_bucket.bucket"),
				),
				ExpectNonEmptyPlan: true, // TODO: No diff in real life, so maybe a timing problem?
			},
			{
				Config: testAccS3BucketConfigWithCORS(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketCors(
						"huaweicloud_s3_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								ExposeHeaders:  []*string{aws.String("x-amz-server-side-encryption"), aws.String("ETag")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

// This fails occasionally, need to dig more.
/*
func TestAccS3Bucket_Logging(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithLogging(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					testAccCheckS3BucketLogging(
						"huaweicloud_s3_bucket.bucket", "huaweicloud_s3_bucket.log_bucket", "log/"),
				),
			},
		},
	})
}
*/

// PASS
func TestAccS3Bucket_Lifecycle(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckS3BucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccS3BucketConfigWithLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.days", "365"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.date", ""),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.expiration.2613713285.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.date", "2016-01-12"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.days", "0"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.expiration.2855832418.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.3.id", "id4"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.3.prefix", "path4/"),
				),
			},
			{
				Config: testAccS3BucketConfigWithVersioningLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.0.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.1.noncurrent_version_expiration.80908210.days", "365"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(
						"huaweicloud_s3_bucket.bucket", "lifecycle_rule.2.prefix", "path3/"),
				),
			},
			{
				Config: testAccS3BucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckS3BucketExists("huaweicloud_s3_bucket.bucket"),
				),
			},
		},
	})
}

func TestS3BucketName(t *testing.T) {
	validDnsNames := []string{
		"foobar",
		"foo.bar",
		"foo.bar.baz",
		"1234",
		"foo-bar",
		strings.Repeat("x", 63),
	}

	for _, v := range validDnsNames {
		if err := validateS3BucketName(v, "us-west-2"); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidDnsNames := []string{
		"foo..bar",
		"Foo.Bar",
		"192.168.0.1",
		"127.0.0.1",
		".foo",
		"bar.",
		"foo_bar",
		strings.Repeat("x", 64),
	}

	for _, v := range invalidDnsNames {
		if err := validateS3BucketName(v, "us-west-2"); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}

	validEastNames := []string{
		"foobar",
		"foo_bar",
		"127.0.0.1",
		"foo..bar",
		"foo_bar_baz",
		"foo.bar.baz",
		"Foo.Bar",
		strings.Repeat("x", 255),
	}

	for _, v := range validEastNames {
		if err := validateS3BucketName(v, "us-east-1"); err != nil {
			t.Fatalf("%q should be a valid S3 bucket name", v)
		}
	}

	invalidEastNames := []string{
		"foo;bar",
		strings.Repeat("x", 256),
	}

	for _, v := range invalidEastNames {
		if err := validateS3BucketName(v, "us-east-1"); err == nil {
			t.Fatalf("%q should not be a valid S3 bucket name", v)
		}
	}
}

func testAccCheckS3BucketDestroy(s *terraform.State) error {
	// UNDONE: Why instance check?
	//return testAccCheckInstanceDestroyWithProvider(s, testAccProvider)
	return nil
}

func testAccCheckS3BucketExists(n string) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckS3BucketExistsWithProviders(n, &providers)
}

func testAccCheckS3BucketExistsWithProviders(n string, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			config := testAccProvider.Meta().(*Config)
			conn, err := config.computeS3conn(OS_REGION_NAME)
			if err != nil {
				return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
			}
			_, err = conn.HeadBucket(&s3.HeadBucketInput{
				Bucket: aws.String(rs.Primary.ID),
			})

			if err != nil {
				return fmt.Errorf("S3 Bucket error: %v", err)
			}
			return nil
		}

		return fmt.Errorf("Instance not found")
	}
}

func testAccCheckS3DestroyBucket(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}
		_, err = conn.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("Error destroying Bucket (%s) in testAccCheckS3DestroyBucket: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func testAccCheckS3BucketPolicy(n string, policy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if policy == "" {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchBucketPolicy" {
				// expected
				return nil
			}
			if err == nil {
				return fmt.Errorf("Expected no policy, got: %#v", *out.Policy)
			} else {
				return fmt.Errorf("GetBucketPolicy error: %v, expected %s", err, policy)
			}
		}
		if err != nil {
			return fmt.Errorf("GetBucketPolicy error: %v, expected %s", err, policy)
		}

		if v := out.Policy; v == nil {
			if policy != "" {
				return fmt.Errorf("bad policy, found nil, expected: %s", policy)
			}
		} else {
			expected := make(map[string]interface{})
			if err := json.Unmarshal([]byte(policy), &expected); err != nil {
				return err
			}
			actual := make(map[string]interface{})
			if err := json.Unmarshal([]byte(*v), &actual); err != nil {
				return err
			}

			if !reflect.DeepEqual(expected, actual) {
				return fmt.Errorf("bad policy, expected: %#v, got %#v", expected, actual)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketWebsite(n string, indexDoc string, errorDoc string, redirectProtocol string, redirectTo string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if indexDoc == "" {
				// If we want to assert that the website is not there, than
				// this error is expected
				return nil
			} else {
				return fmt.Errorf("S3BucketWebsite error: %v", err)
			}
		}

		if v := out.IndexDocument; v == nil {
			if indexDoc != "" {
				return fmt.Errorf("bad index doc, found nil, expected: %s", indexDoc)
			}
		} else {
			if *v.Suffix != indexDoc {
				return fmt.Errorf("bad index doc, expected: %s, got %#v", indexDoc, out.IndexDocument)
			}
		}

		if v := out.ErrorDocument; v == nil {
			if errorDoc != "" {
				return fmt.Errorf("bad error doc, found nil, expected: %s", errorDoc)
			}
		} else {
			if *v.Key != errorDoc {
				return fmt.Errorf("bad error doc, expected: %s, got %#v", errorDoc, out.ErrorDocument)
			}
		}

		if v := out.RedirectAllRequestsTo; v == nil {
			if redirectTo != "" {
				return fmt.Errorf("bad redirect to, found nil, expected: %s", redirectTo)
			}
		} else {
			if *v.HostName != redirectTo {
				return fmt.Errorf("bad redirect to, expected: %s, got %#v", redirectTo, out.RedirectAllRequestsTo)
			}
			if redirectProtocol != "" && v.Protocol != nil && *v.Protocol != redirectProtocol {
				return fmt.Errorf("bad redirect protocol to, expected: %s, got %#v", redirectProtocol, out.RedirectAllRequestsTo)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketWebsiteRoutingRules(n string, routingRules []*s3.RoutingRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketWebsite(&s3.GetBucketWebsiteInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if routingRules == nil {
				return nil
			}
			return fmt.Errorf("GetBucketWebsite error: %v", err)
		}

		if !reflect.DeepEqual(out.RoutingRules, routingRules) {
			return fmt.Errorf("bad routing rule, expected: %v, got %v", routingRules, out.RoutingRules)
		}

		return nil
	}
}

func testAccCheckS3BucketVersioning(n string, versioningStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketVersioning error: %v", err)
		}

		if v := out.Status; v == nil {
			if versioningStatus != "" {
				return fmt.Errorf("bad error versioning status, found nil, expected: %s", versioningStatus)
			}
		} else {
			if *v != versioningStatus {
				return fmt.Errorf("bad error versioning status, expected: %s, got %s", versioningStatus, *v)
			}
		}

		return nil
	}
}

func testAccCheckS3BucketCors(n string, corsRules []*s3.CORSRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketCors(&s3.GetBucketCorsInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketCors error: %v", err)
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func testAccCheckS3BucketLogging(n, b, p string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		conn, err := config.computeS3conn(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud s3 client: %s", err)
		}

		out, err := conn.GetBucketLogging(&s3.GetBucketLoggingInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketLogging error: %v", err)
		}

		tb, _ := s.RootModule().Resources[b]

		if v := out.LoggingEnabled.TargetBucket; v == nil {
			if tb.Primary.ID != "" {
				return fmt.Errorf("bad target bucket, found nil, expected: %s", tb.Primary.ID)
			}
		} else {
			if *v != tb.Primary.ID {
				return fmt.Errorf("bad target bucket, expected: %s, got %s", tb.Primary.ID, *v)
			}
		}

		if v := out.LoggingEnabled.TargetPrefix; v == nil {
			if p != "" {
				return fmt.Errorf("bad target prefix, found nil, expected: %s", p)
			}
		} else {
			if *v != p {
				return fmt.Errorf("bad target prefix, expected: %s, got %s", p, *v)
			}
		}

		return nil
	}
}

// These need a bit of randomness as the name can only be used once globally
// within AWS
func testAccBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.s3.amazonaws.com", randInt)
}

func testAccWebsiteEndpoint(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.s3-website.%s.amazonaws.com", randInt, OS_REGION_NAME)
}

func testAccS3BucketPolicy(randInt int) string {
	return fmt.Sprintf(`{ "Version": "2008-10-17", "Statement": [ { "Effect": "Allow", "Principal": { "AWS": ["*"] }, "Action": ["s3:GetObject"], "Resource": ["arn:aws:s3:::tf-test-bucket-%d/*"] } ] }`, randInt)
}

func testAccS3BucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccAWSS3MultiBucketConfigWithTags(randInt int) string {
	t := template.Must(template.New("t1").
		Parse(`
resource "huaweicloud_s3_bucket" "bucket1" {
	bucket = "tf-test-bucket-1-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-1-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "huaweicloud_s3_bucket" "bucket2" {
	bucket = "tf-test-bucket-2-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-2-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "huaweicloud_s3_bucket" "bucket3" {
	bucket = "tf-test-bucket-3-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-3-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "huaweicloud_s3_bucket" "bucket4" {
	bucket = "tf-test-bucket-4-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-4-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "huaweicloud_s3_bucket" "bucket5" {
	bucket = "tf-test-bucket-5-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-5-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}

resource "huaweicloud_s3_bucket" "bucket6" {
	bucket = "tf-test-bucket-6-{{.GUID}}"
	acl = "private"
	force_destroy = true
	tags {
		Name = "tf-test-bucket-6-{{.GUID}}"
		Environment = "{{.GUID}}"
	}
}
`))
	var doc bytes.Buffer
	t.Execute(&doc, struct{ GUID int }{GUID: randInt})
	return doc.String()
}

func testAccS3BucketConfigWithRegion(randInt int) string {
	return fmt.Sprintf(`
provider "huaweicloud" {
	alias = "reg1"
	region = "%s"
}

resource "huaweicloud_s3_bucket" "bucket" {
	provider = "huaweicloud.reg1"
	bucket = "tf-test-bucket-%d"
	region = "%s"
}
`, OS_REGION_NAME, randInt, OS_REGION_NAME)
}

func testAccS3BucketWebsiteConfig(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		index_document = "index.html"
	}
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithError(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		index_document = "index.html"
		error_document = "error.html"
	}
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		redirect_all_requests_to = "hashicorp.com"
	}
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithHttpsRedirect(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

	website {
		redirect_all_requests_to = "https://hashicorp.com"
	}
}
`, randInt)
}

func testAccS3BucketWebsiteConfigWithRoutingRules(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"

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

func testAccS3BucketConfigWithPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	policy = %s
}
`, randInt, strconv.Quote(testAccS3BucketPolicy(randInt)))
}

func testAccS3BucketDestroyedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccS3BucketConfigWithEmptyPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	policy = ""
}
`, randInt)
}

func testAccS3BucketConfigWithVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	versioning {
	  enabled = true
	}
}
`, randInt)
}

func testAccS3BucketConfigWithDisableVersioning(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	versioning {
	  enabled = false
	}
}
`, randInt)
}

func testAccS3BucketConfigWithCORS(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
	cors_rule {
			allowed_headers = ["*"]
			allowed_methods = ["PUT","POST"]
			allowed_origins = ["https://www.example.com"]
			expose_headers = ["x-amz-server-side-encryption","ETag"]
			max_age_seconds = 3000
	}
}
`, randInt)
}

var testAccS3BucketConfigWithAcl = `
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`

var testAccS3BucketConfigWithAclUpdate = `
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
}
`

func testAccS3BucketConfigWithLogging(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "log_bucket" {
	bucket = "tf-test-log-bucket-%d"
	acl = "log-delivery-write"
}
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	force_destroy = "true"
	logging {
		target_bucket = "${huaweicloud_s3_bucket.log_bucket.id}"
		target_prefix = "log/"
	}
}
`, randInt, randInt)
}

func testAccS3BucketConfigWithLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	lifecycle_rule {
		id = "id1"
		prefix = "path1/"
		enabled = true

		expiration {
			days = 365
		}
	}
	lifecycle_rule {
		id = "id2"
		prefix = "path2/"
		enabled = true

		expiration {
			date = "2016-01-12"
		}
	}
	lifecycle_rule {
		id = "id3"
		prefix = "path3/"
		enabled = true

		expiration {
			days = "30"
		}
	}
	lifecycle_rule {
		id = "id4"
		prefix = "path4/"
		enabled = true

		expiration {
			date = "2016-01-12"
		}
	}
}
`, randInt)
}

func testAccS3BucketConfigWithVersioningLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "huaweicloud_s3_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
	versioning {
	  enabled = false
	}
	lifecycle_rule {
		id = "id1"
		prefix = "path1/"
		enabled = true

		noncurrent_version_expiration {
			days = 365
		}
	}
	lifecycle_rule {
		id = "id2"
		prefix = "path2/"
		enabled = false

		noncurrent_version_expiration {
			days = 365
		}
	}
	lifecycle_rule {
		id = "id3"
		prefix = "path3/"
		enabled = true

		noncurrent_version_expiration {
			days = 30
		}
	}
}
`, randInt)
}

const testAccS3BucketConfig_namePrefix = `
resource "huaweicloud_s3_bucket" "test" {
	bucket_prefix = "tf-test-"
}
`

const testAccS3BucketConfig_generatedName = `
resource "huaweicloud_s3_bucket" "test" {
	bucket_prefix = "tf-test-"
}
`
