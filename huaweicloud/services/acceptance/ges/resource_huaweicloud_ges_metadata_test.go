package ges

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getGesMetadataResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getMetadataDetail: Query the GES metadata detail.
	var (
		getMetadataDetailHttpUrl = "v1.0/{project_id}/graphs/metadatas/{id}"
		getMetadataDetailProduct = "ges"
	)
	getMetadataDetailClient, err := cfg.NewServiceClient(getMetadataDetailProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating GES Client: %s", err)
	}

	getMetadataDetailPath := getMetadataDetailClient.Endpoint + getMetadataDetailHttpUrl
	getMetadataDetailPath = strings.ReplaceAll(getMetadataDetailPath, "{project_id}", getMetadataDetailClient.ProjectID)
	getMetadataDetailPath = strings.ReplaceAll(getMetadataDetailPath, "{id}", state.Primary.ID)

	getMetadataDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getMetadataDetailResp, err := getMetadataDetailClient.Request("GET", getMetadataDetailPath, &getMetadataDetailOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GES metadata: %s", err)
	}

	getMetadataDetailRespBody, err := utils.FlattenResponse(getMetadataDetailResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GES metadata: %s", err)
	}

	return getMetadataDetailRespBody, nil
}

func TestAccGesMetadata_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ges_metadata.test"
	bucketName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getGesMetadataResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testGesMetadata_basic(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "ges_metadata.0.labels.0.properties.#", "1"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testGesMetadata_basic_update(name, bucketName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "description", "This is a demo"),
					resource.TestCheckResourceAttr(rName, "ges_metadata.0.labels.0.properties.#", "2"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testGesMetadata_basic(name, bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ges_metadata" "test" {
  name          = "%s"
  description   = "This is a demo"
  metadata_path = "${huaweicloud_obs_bucket.test.bucket}/schema_%s.xml"
  ges_metadata {
    labels {
      name = "user"
      properties = [{
        "dataType"    = "char"
        "name"        = "sex"
        "cardinality" = "single"
        }]
    }
  }

  depends_on = [
    huaweicloud_obs_bucket.test
  ]
}
`, bucketName, name, name)
}

func testGesMetadata_basic_update(name, bucketName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket        = "%s"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_ges_metadata" "test" {
  name          = "%s"
  description   = "This is a demo"
  metadata_path = "${huaweicloud_obs_bucket.test.bucket}/schema_%s.xml"
  ges_metadata {
    labels {
      name = "user"
      properties = [{
        "dataType"    = "char"
        "name"        = "sex"
        "cardinality" = "single"
        },
        {
          "dataType"      = "enum"
          "name"          = "country"
          "cardinality"   = "single"
          "typeNameCount" = "3"
          "typeName1"     = "US"
          "typeName2"     = "EN"
          "typeName3"     = "CN"
        }]
    }
  }

  depends_on = [
    huaweicloud_obs_bucket.test
  ]
}
`, bucketName, name, name)
}
