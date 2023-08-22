package modelarts

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getModelArtsAuthorizationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME
	// getModelAuth: Query the ModelArts authorization.
	var (
		getModelAuthHttpUrl = "v2/{project_id}/authorizations"
		getModelAuthProduct = "modelarts"
	)
	getModelAuthClient, err := cfg.NewServiceClient(getModelAuthProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelAuthPath := getModelAuthClient.Endpoint + getModelAuthHttpUrl
	getModelAuthPath = strings.ReplaceAll(getModelAuthPath, "{project_id}", getModelAuthClient.ProjectID)

	getModelAuthResp, err := pagination.ListAllItems(
		getModelAuthClient,
		"offset",
		getModelAuthPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, fmt.Errorf("error retrieving ModelArts authorization: %s", err)
	}

	getModelAuthRespJson, err := json.Marshal(getModelAuthResp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ModelArts authorization: %s", err)
	}
	var getModelAuthRespBody interface{}
	err = json.Unmarshal(getModelAuthRespJson, &getModelAuthRespBody)
	if err != nil {
		return nil, fmt.Errorf("error retrieving ModelArts authorization: %s", err)
	}

	getModelAuthRespBody = modelarts.SearchAuthById(getModelAuthRespBody, state.Primary.ID)
	if getModelAuthRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return getModelAuthRespBody, nil
}

func TestAccModelArtsAuthorization_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	password := acceptance.RandomPassword()
	rName := "huaweicloud_modelarts_authorization.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelArtsAuthorizationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelArtsAuthorization_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "user_id", "huaweicloud_identity_user.test", "id"),
					resource.TestCheckResourceAttr(rName, "user_name", name),
					resource.TestCheckResourceAttr(rName, "type", "agency"),
					resource.TestCheckResourceAttr(rName, "agency_name", "ma_agency_tf_acc_test"),
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

func TestAccModelArtsAuthorization_all(t *testing.T) {
	var obj interface{}

	rName := "huaweicloud_modelarts_authorization.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getModelArtsAuthorizationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testModelArtsAuthorization_allUser(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "user_id", "all-users"),
					resource.TestCheckResourceAttr(rName, "user_name", "all-users"),
					resource.TestCheckResourceAttr(rName, "type", "agency"),
					resource.TestCheckResourceAttr(rName, "agency_name", "modelarts_agency"),
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

func testModelArtsAuthorization_basic(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_user" "test" {
  name        = "%[1]s"
  password    = "%[2]s"
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_modelarts_authorization" "test" {
  user_id     = huaweicloud_identity_user.test.id
  type        = "agency"
  agency_name = "ma_agency_tf_acc_test"
}
`, name, password)
}

func testModelArtsAuthorization_allUser() string {
	return `
resource "huaweicloud_modelarts_authorization" "test" {
  user_id     = "all-users"
  type        = "agency"
  agency_name = "modelarts_agency"
}
`
}
