package lts

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

func getLtsStreamResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups/{log_group_id}/streams"
	client, err := cfg.NewServiceClient("lts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating LTS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_group_id}", state.Primary.Attributes["group_id"])

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing the log stream: %s", err)
	}

	streamId := state.Primary.ID
	streamResult := utils.PathSearch(fmt.Sprintf("log_streams|[?log_stream_id=='%s']|[0]", streamId), respBody, nil)
	if streamResult == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return streamResult, nil
}

func TestAccLtsStream_basic(t *testing.T) {
	var (
		stream       interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_lts_stream.test"
		rc           = acceptance.InitResourceCheck(resourceName, &stream, getLtsStreamResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLtsStream_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "-1"),
					resource.TestCheckResourceAttr(resourceName, "filter_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrPair(resourceName, "group_id", "huaweicloud_lts_group.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.terraform", ""),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "true"),
				),
			},
			{
				Config: testAccLtsStream_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "ttl_in_days", "60"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStreamImportState(resourceName),
			},
		},
	})
}

func testLtsStreamImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		streamID := rs.Primary.ID
		groupID := rs.Primary.Attributes["group_id"]

		return fmt.Sprintf("%s/%s", groupID, streamID), nil
	}
}

func testAccStream_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_lts_group" "test" {
  group_name  = "%[1]s"
  ttl_in_days = 30

  tags = {
    owner = "terraform"
  }
}
`, rName)
}

func testAccLtsStream_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
  is_favorite = true

  tags = {
    foo       = "bar"
    terraform = ""
  }
}
`, testAccStream_base(rName), rName)
}

func testAccLtsStream_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = "%[2]s"
  ttl_in_days = 60

  tags = {
    owner = "terraform"
  }
}
`, testAccStream_base(rName), rName)
}

func TestAccLtsStream_epsId(t *testing.T) {
	var (
		stream       interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_lts_stream.test"
		rc           = acceptance.InitResourceCheck(resourceName, &stream, getLtsStreamResourceFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccStream_epsId_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "false"),
				),
			},
			{
				Config: testAccStream_epsId_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "stream_name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "is_favorite", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testLtsStreamImportState(resourceName),
			},
		},
	})
}

func testAccStream_epsId_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = "%[2]s"
  enterprise_project_id = "%[3]s"
}
`, testAccStream_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccStream_epsId_step2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = "%[2]s"
  enterprise_project_id = "%[3]s"
  is_favorite           = true
}
`, testAccStream_base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
