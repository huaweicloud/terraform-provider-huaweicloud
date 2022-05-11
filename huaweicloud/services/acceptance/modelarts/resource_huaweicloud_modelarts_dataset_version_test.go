package modelarts

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"

	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getDatasetVersionResourceFunc(config *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := config.ModelArtsV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId, versionId, err := modelarts.ParseVersionInfoFromId(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	return version.Get(client, datasetId, versionId)
}

func TestAccDatasetVersionResource_basic(t *testing.T) {
	var instance dataset.CreateOpts
	resourceName := "huaweicloud_modelarts_dataset_version.test"
	name := acceptance.RandomAccResourceName()
	obsName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&instance,
		getDatasetVersionResourceFunc,
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
				Config: testAccDatasetVersion_basic(name, obsName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", name),
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

func testAccDatasetVersion_basic(rName, obsName string) string {
	datasetConfig := testAccDateset_basic(rName, obsName)
	return fmt.Sprintf(`
%s

resource "huaweicloud_modelarts_dataset_version" "test" {
  name        = "%[2]s"
  dataset_id  = huaweicloud_modelarts_dataset.test.id
  description = "%[2]s"
}
`, datasetConfig, rName)
}
