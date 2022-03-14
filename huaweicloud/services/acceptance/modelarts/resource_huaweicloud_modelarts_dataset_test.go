package modelarts

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDatesetResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ModelArtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmtp.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	return dataset.Get(client, state.Primary.ID, dataset.GetOpts{})
}

func TestAccResourceDateset_basic(t *testing.T) {
	var instance dataset.CreateOpts
	resourceName := "huaweicloud_modelarts_dataset.test"
	name := acceptance.RandomAccResourceName()
	updateName := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDatesetResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOBS(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDateset_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_format", "Default"),
					resource.TestCheckResourceAttr(resourceName, "output_path", fmt.Sprintf("/%s/%s/", obsName, "output")),
					resource.TestCheckResourceAttr(resourceName, "description", name),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.data_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.path", fmt.Sprintf("/%s/%s/", obsName, "input")),
					resource.TestCheckResourceAttr(resourceName, "labels.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccDateset_basic(updateName, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "type", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_format", "Default"),
					resource.TestCheckResourceAttr(resourceName, "output_path", fmt.Sprintf("/%s/%s/", obsName, "output")),
					resource.TestCheckResourceAttr(resourceName, "description", updateName),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.data_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.path", fmt.Sprintf("/%s/%s/", obsName, "input")),
					resource.TestCheckResourceAttr(resourceName, "labels.0.name", updateName),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
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

func testAccDatesetObs(obsName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true

  lifecycle {
    ignore_changes = [
      cors_rule,
    ]
  }
}

resource "huaweicloud_obs_bucket_object" "input" {
  bucket  = huaweicloud_obs_bucket.bucket.bucket
  key     = "input/t1"
  content = "some_bucket_content"
}

resource "huaweicloud_obs_bucket_object" "output" {
  bucket  = huaweicloud_obs_bucket.bucket.bucket
  key     = "output/t2"
  content = "some_bucket_content"
}
`, obsName)
}

func testAccDateset_basic(rName, obsName string) string {
	obsConfig := testAccDatesetObs(obsName)
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_dataset" "test" {
  name        = "%s"
  type        = 1
  output_path = "/${huaweicloud_obs_bucket.bucket.bucket}/output/"
  description = "%s"
  data_source {
    path = "/${huaweicloud_obs_bucket.bucket.bucket}/input/"
  }

  labels {
    name = "%s"
  }

  depends_on = [
    huaweicloud_obs_bucket_object.input,
    huaweicloud_obs_bucket_object.output
  ]
}
`, obsConfig, rName, rName, rName)
}
