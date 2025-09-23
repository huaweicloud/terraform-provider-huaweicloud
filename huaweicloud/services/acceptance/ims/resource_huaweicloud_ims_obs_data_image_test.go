package ims

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getObsDataImageResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "ims"
		httpUrl = "v2/cloudimages"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating IMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += fmt.Sprintf("?id=%s", state.Primary.ID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IMS OBS data image: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	image := utils.PathSearch("images[0]", getRespBody, nil)
	// If the list API return empty, then return `404` error code.
	if image == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return image, nil
}

func TestAccObsDataImage_basic(t *testing.T) {
	var (
		image        interface{}
		rName        = acceptance.RandomAccResourceName()
		rNameUpdate  = rName + "-update"
		resourceName = "huaweicloud_ims_obs_data_image.test"
		defaultEpsId = "0"
		migrateEpsId = acceptance.HW_ENTERPRISE_PROJECT_ID_TEST
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&image,
		getObsDataImageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need setting a non default enterprise project ID.
			acceptance.TestAccPreCheckEpsID(t)
			// This test requires ensuring that there is an image file in the OBS bucket.
			acceptance.TestAccPreCheckImsImageUrl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccObsDataImage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "image_url", acceptance.HW_IMS_IMAGE_URL),
					resource.TestCheckResourceAttr(resourceName, "min_disk", "60"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "description", "terraform description test"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttrSet(resourceName, "file"),
					resource.TestCheckResourceAttrSet(resourceName, "self"),
					resource.TestCheckResourceAttrSet(resourceName, "schema"),
					resource.TestCheckResourceAttrSet(resourceName, "visibility"),
					resource.TestCheckResourceAttrSet(resourceName, "disk_format"),
					resource.TestCheckResourceAttrSet(resourceName, "image_size"),
					resource.TestCheckResourceAttrSet(resourceName, "data_origin"),
					resource.TestMatchResourceAttr(resourceName, "active_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(resourceName, "updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrPair(resourceName, "cmk_id", "huaweicloud_kms_key.test", "id"),
				),
			},
			{
				Config: testAccObsDataImage_update1(rName, rNameUpdate, migrateEpsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", migrateEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				Config: testAccObsDataImage_update2(rName, rNameUpdate, defaultEpsId),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", defaultEpsId),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
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

func testAccObsDataImage_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}
`, rName)
}

func testAccObsDataImage_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_obs_data_image" "test" {
  name        = "%[2]s"
  image_url   = "%[3]s"
  min_disk    = 60
  os_type     = "Windows"
  description = "terraform description test"
  cmk_id      = huaweicloud_kms_key.test.id

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccObsDataImage_base(rName), rName, acceptance.HW_IMS_IMAGE_URL)
}

func testAccObsDataImage_update1(rName, rNameUpdate, migrateEpsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_obs_data_image" "test" {
  name                  = "%[2]s"
  image_url             = "%[3]s"
  min_disk              = 60
  os_type               = "Windows"
  cmk_id                = huaweicloud_kms_key.test.id
  enterprise_project_id = "%[4]s"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccObsDataImage_base(rName), rNameUpdate, acceptance.HW_IMS_IMAGE_URL, migrateEpsId)
}

func testAccObsDataImage_update2(rName, rNameUpdate, defaultEpsId string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_ims_obs_data_image" "test" {
  name                  = "%[2]s"
  image_url             = "%[3]s"
  min_disk              = 60
  os_type               = "Windows"
  cmk_id                = huaweicloud_kms_key.test.id
  enterprise_project_id = "%[4]s"

  tags = {
    foo  = "bar"
    key  = "value1"
    key2 = "value2"
  }
}
`, testAccObsDataImage_base(rName), rNameUpdate, acceptance.HW_IMS_IMAGE_URL, defaultEpsId)
}
