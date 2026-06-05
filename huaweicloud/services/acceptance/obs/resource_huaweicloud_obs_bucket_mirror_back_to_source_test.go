package obs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getOBSBucketMirrorBackToSourceResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	obsClient, err := cfg.ObjectStorageClientWithSignature(region)
	if err != nil {
		return nil, fmt.Errorf("error creating OBS Client: %s", err)
	}

	return obsClient.GetBucketMirrorBackToSource(state.Primary.Attributes["bucket"])
}

func TestAccBucketMirrorBackToSource_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_obs_bucket_mirror_back_to_source.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getOBSBucketMirrorBackToSourceResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccBucketMirrorBackToSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", fmt.Sprintf("%s-0", name)),
					resource.TestCheckResourceAttrSet(rName, "rule"),
				),
			},
			{
				Config: testAccBucketMirrorBackToSource_update_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", fmt.Sprintf("%s-0", name)),
					resource.TestCheckResourceAttrSet(rName, "rule"),
				),
			},
			{
				Config: testAccBucketMirrorBackToSource_update_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "bucket", fmt.Sprintf("%s-1", name)),
					resource.TestCheckResourceAttrSet(rName, "rule"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateIdFunc: testAccBucketMirrorBackToSourceImportStateFunc(rName),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_force_new",
				},
			},
		},
	})
}

func testAccBucketMirrorBackToSource_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  count = 2

  bucket        = "%[1]s-${count.index}"
  storage_class = "STANDARD"
  acl           = "private"
}

data "huaweicloud_identity_agencies" "test" {
  name = "live_to_obs"
}
`, name)
}

func testAccBucketMirrorBackToSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_mirror_back_to_source" "test" {
  bucket = huaweicloud_obs_bucket.test[0].bucket
  rule   = jsonencode({
    "id" : "terraformtest123",
    "condition" : {
      "httpErrorCodeReturnedEquals" : "404",
      "objectKeyPrefixEquals" : "test"
    },
    "redirect" : {
      "replaceKeyWith" : "test1$${key}test2",
      "publicSource" : {
        "sourceEndpoint" : {
          "master" : ["https://www.tftest1.com/xxx"],
          "slave" : ["https://www.tftest2.com/yyy"]
        }
      },
      "retryConditions" : ["4XX"],
      "returnBaseErrorConditions" : {
        "IFANY" : {
          "HTTP.STATUS_CODE" : ["4XX"]
        }
      },
      "agency" : try(data.huaweicloud_identity_agencies.test.agencies[0].name, ""),
      "ReturnBaseErrorBodyNewRule": true,
      "privateSource": {},
      "redirectHttpCode": "",
      "redirectOriginServer": "",
      "passQueryString" : true,
      "mirrorFollowRedirect" : true,
      "redirectWithoutReferer" : true,
      "mirrorHttpHeader" : {
        "passAll" : true
      }
    }
  })
}
`, testAccBucketMirrorBackToSource_base(name))
}

func testAccBucketMirrorBackToSource_update_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_mirror_back_to_source" "test" {
  bucket = huaweicloud_obs_bucket.test[0].bucket
  rule   = jsonencode({
    "id" : "terraformtest123",
    "condition" : {
      "httpErrorCodeReturnedEquals" : "404",
      "objectKeyPrefixEquals" : "test_update_step1"
    },
    "redirect" : {
      "replaceKeyWith" : "test2$${key}test1",
      "publicSource" : {
        "sourceEndpoint" : {
          "master" : ["https://www.tftest3.com/zzz"],
          "slave" : ["https://www.tftest4.com/www"]
        }
      },
      "retryConditions" : ["4XX"],
      "returnBaseErrorConditions" : {
        "IFANY" : {
          "HTTP.STATUS_CODE" : ["4XX"]
        }
      },
      "agency" : try(data.huaweicloud_identity_agencies.test.agencies[0].name, ""),
      "ReturnBaseErrorBodyNewRule": true,
      "privateSource": {},
      "redirectHttpCode": "",
      "redirectOriginServer": "",
      "passQueryString" : true,
      "mirrorFollowRedirect" : true,
      "redirectWithoutReferer" : true,
      "mirrorHttpHeader" : {
        "passAll" : true
      }
    }
  })
}
`, testAccBucketMirrorBackToSource_base(name))
}

func testAccBucketMirrorBackToSource_update_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_obs_bucket_mirror_back_to_source" "test" {
  bucket = huaweicloud_obs_bucket.test[1].bucket
  rule   = jsonencode({
    "id" : "terraformtest123",
    "condition" : {
      "httpErrorCodeReturnedEquals" : "404",
      "objectKeyPrefixEquals" : "test_update_step2"
    },
    "redirect" : {
      "replaceKeyWith" : "test2$${key}test1",
      "publicSource" : {
        "sourceEndpoint" : {
          "master" : ["https://www.tftest3.com/zzz"],
          "slave" : ["https://www.tftest4.com/www"]
        }
      },
      "retryConditions" : ["4XX"],
      "returnBaseErrorConditions" : {
        "IFANY" : {
          "HTTP.STATUS_CODE" : ["4XX"]
        }
      },
      "agency" : try(data.huaweicloud_identity_agencies.test.agencies[0].name, ""),
      "ReturnBaseErrorBodyNewRule": true,
      "privateSource": {},
      "redirectHttpCode": "",
      "redirectOriginServer": "",
      "passQueryString" : true,
      "mirrorFollowRedirect" : true,
      "redirectWithoutReferer" : true,
      "mirrorHttpHeader" : {
        "passAll" : true
      }
    }
  })

  enable_force_new = "true"
}
`, testAccBucketMirrorBackToSource_base(name))
}

func testAccBucketMirrorBackToSourceImportStateFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		bucket := rs.Primary.Attributes["bucket"]
		if bucket == "" {
			return "", fmt.Errorf("attribute (bucket) of resource (%s) not found: %s", name, rs)
		}
		return fmt.Sprintf("%s/%s", bucket, rs.Primary.ID), nil
	}
}
