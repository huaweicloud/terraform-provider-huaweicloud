package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getDatasetResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.ModelArtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts v1 client, err=%s", err)
	}

	return dataset.Get(client, state.Primary.ID, dataset.GetOpts{})
}

func TestAccDataset_basic(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_modelarts_dataset.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getDatasetResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
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
				Config: testAccDataset_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "type", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_format", "Default"),
					resource.TestCheckResourceAttr(resourceName, "output_path", fmt.Sprintf("/%s/%s/", name, "output")),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.data_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.path", fmt.Sprintf("/%s/%s/", name, "input")),
					resource.TestCheckResourceAttr(resourceName, "labels.0.name", name),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				Config: testAccDataset_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-update", name)),
					resource.TestCheckResourceAttr(resourceName, "type", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttr(resourceName, "data_format", "Default"),
					resource.TestCheckResourceAttr(resourceName, "output_path", fmt.Sprintf("/%s/%s/", name, "output")),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.data_type", "0"),
					resource.TestCheckResourceAttr(resourceName, "data_source.0.path", fmt.Sprintf("/%s/%s/", name, "input")),
					resource.TestCheckResourceAttr(resourceName, "labels.0.name", fmt.Sprintf("%s-update", name)),
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

func testAccDataset_basic_base(name string) string {
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
`, name)
}

func testAccDataset_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_dataset" "test" {
  name        = "%[2]s"
  type        = 1
  output_path = "/${huaweicloud_obs_bucket.bucket.bucket}/output/"
  description = "Created by terraform script"
  data_source {
    path = "/${huaweicloud_obs_bucket.bucket.bucket}/input/"
  }

  labels {
    name = "%[2]s"
  }

  depends_on = [
    huaweicloud_obs_bucket_object.input,
    huaweicloud_obs_bucket_object.output
  ]
}
`, testAccDataset_basic_base(name), name)
}

func testAccDataset_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_dataset" "test" {
  name        = "%[2]s-update"
  type        = 1
  output_path = "/${huaweicloud_obs_bucket.bucket.bucket}/output/"
  data_source {
    path = "/${huaweicloud_obs_bucket.bucket.bucket}/input/"
  }

  labels {
    name = "%[2]s-update"
  }

  depends_on = [
    huaweicloud_obs_bucket_object.input,
    huaweicloud_obs_bucket_object.output
  ]
}
`, testAccDataset_basic_base(name), name)
}
