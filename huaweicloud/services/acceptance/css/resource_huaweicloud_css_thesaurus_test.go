package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/thesaurus"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCssThesaurus_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_css_thesaurus.test"
	bucketName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCssThesaurusDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCssThesaurus_basic(rName, bucketName, "main.txt"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssThesaurusExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "main_object", "main.txt"),
				),
			},
			{
				Config: testAccCssThesaurus_basic(rName, bucketName, "main2.txt"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssThesaurusExists(resourceName),
					resource.TestCheckResourceAttrSet(resourceName, "bucket_name"),
					resource.TestCheckResourceAttr(resourceName, "main_object", "main2.txt"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCssThesaurus_basic(rName string, bucketName string, obsObjectKey string) string {
	cssClusterBasic := testAccCssCluster_basic(rName, "Test@passw0rd", 1, "value")

	return fmt.Sprintf(`
%s

resource "huaweicloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}


resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.bucket
  key          = "%s"
  content      = "123"
  content_type = "text/plain"
}

resource "huaweicloud_css_thesaurus" "test" {
  cluster_id  = huaweicloud_css_cluster.test.id
  bucket_name = huaweicloud_obs_bucket.test.bucket
  main_object = huaweicloud_obs_bucket_object.test.key
}

`, cssClusterBasic, bucketName, obsObjectKey)
}

func testAccCheckCssThesaurusDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	client, err := conf.CssV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating CSS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_css_thesaurus" {
			continue
		}

		resp, getErr := thesaurus.Get(client, rs.Primary.ID)
		if getErr != nil {
			if _, ok := getErr.(golangsdk.ErrDefault404); !ok {
				return fmt.Errorf("get CSS thesaurus failed.error: %s", getErr)
			}
		} else {
			if resp.Bucket != "" {
				return fmt.Errorf("CSS thesaurus still exists, cluster_id: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckCssThesaurusExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		client, err := config.CssV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating CSS client: %s", err)
		}

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("error checking huaweicloud_css_thesaurus exist, err: not found this resource")
		}

		resp, errQueryDetail := thesaurus.Get(client, rs.Primary.ID)
		if errQueryDetail != nil {
			return fmt.Errorf("error checking huaweicloud_css_thesaurus exist, send request failed: %s", errQueryDetail)
		}

		if resp == nil || resp.Bucket == "" {
			return fmt.Errorf("CSS thesaurus don't exists, cluster_id: %s", rs.Primary.ID)
		}

		return nil
	}
}
