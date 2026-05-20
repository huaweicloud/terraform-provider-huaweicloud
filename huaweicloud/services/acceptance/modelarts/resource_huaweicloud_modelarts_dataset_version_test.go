package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/version"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getDatasetVersionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ModelArtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId, versionId, err := modelarts.ParseVersionInfoFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	return version.Get(client, datasetId, versionId)
}

func TestAccDatasetVersion_basic(t *testing.T) {
	var instance dataset.CreateOpts
	resourceName := "huaweicloud_modelarts_dataset_version.test"
	name := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDatasetVersionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasetVersion_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "split_ratio", "1.00"),
					resource.TestCheckResourceAttr(resourceName, "hard_example", "false"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttr(resourceName, "verification", "false"),
					resource.TestCheckResourceAttr(resourceName, "labeling_type", "unlabeled"),
					resource.TestCheckResourceAttr(resourceName, "files", "0"),
					resource.TestCheckResourceAttrPair(resourceName, "dataset_id",
						"huaweicloud_modelarts_dataset.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
					resource.TestCheckResourceAttrSet(resourceName, "version_id"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_path"),
					resource.TestCheckResourceAttrSet(resourceName, "is_current"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hard_example"},
			},
		},
	})
}

func testAccDatasetVersion_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
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

resource "huaweicloud_modelarts_dataset" "test" {
	name        = "%[1]s"
	type        = 1
	output_path = "/${huaweicloud_obs_bucket.bucket.bucket}/output/"
	description = "Created by terraform script"
	data_source {
	  path = "/${huaweicloud_obs_bucket.bucket.bucket}/input/"
	}
  
	labels {
	  name = "%[1]s"
	}
  
	depends_on = [
	  huaweicloud_obs_bucket_object.input,
	  huaweicloud_obs_bucket_object.output
	]
  }
`, name)
}

func testAccDatasetVersion_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_dataset_version" "test" {
  name        = "%[2]s"
  dataset_id  = huaweicloud_modelarts_dataset.test.id
  description = "Created by terraform script"
}
`, testAccDatasetVersion_base(name), name)
}
