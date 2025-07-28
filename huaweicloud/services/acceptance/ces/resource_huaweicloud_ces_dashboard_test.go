package ces

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getDashboardFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	var (
		httpUrl = "v2/{project_id}/dashboards?dashboard_id={id}"
		product = "ces"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating CES Client: %s", err)
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{id}", state.Primary.ID)

	getResourceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResourceGroupResp, err := client.Request("GET", path, &getResourceGroupOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CES dashboard: %s", err)
	}
	return utils.FlattenResponse(getResourceGroupResp)
}

func TestAccDashboard_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_dashboard.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDashboardFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboard_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "1"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDashboard_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "2"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dashboard_id",
				},
			},
		},
	})
}

func TestAccDashboard_copy(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_dashboard.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDashboardFunc,
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
				Config: testDashboard_copy_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "1"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestCheckResourceAttrPair(rName, "dashboard_id", "huaweicloud_ces_dashboard.base", "id"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDashboard_copy_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "3"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
				),
			},
		},
	})
}

func TestAccDashboard_extend_info(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_ces_dashboard.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getDashboardFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDashboard_extend_info_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "1"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "true"),
					resource.TestCheckResourceAttr(rName, "extend_info.0.filter", "average"),
					resource.TestCheckResourceAttr(rName, "extend_info.0.period", "60"),
					resource.TestCheckResourceAttrSet(rName, "creator_name"),
					resource.TestMatchResourceAttr(rName,
						"created_at", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			{
				Config: testDashboard_extend_info_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "row_widget_num", "2"),
					resource.TestCheckResourceAttr(rName, "is_favorite", "false"),
					resource.TestCheckResourceAttr(rName, "extend_info.0.filter", "min"),
					resource.TestCheckResourceAttr(rName, "extend_info.0.period", "300"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"dashboard_id",
				},
			},
		},
	})
}

func testDashboard_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s"
  row_widget_num = 1
  is_favorite    = true
}
`, name)
}

func testDashboard_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s-update"
  row_widget_num = 2
  is_favorite    = false
}
`, name)
}

func testDashboard_copy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "base" {
  name           = "%[1]s"
  row_widget_num = 2
}

resource "huaweicloud_ces_dashboard" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
  dashboard_id          = huaweicloud_ces_dashboard.base.id
  row_widget_num        = 1
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDashboard_copy_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "base" {
  name           = "%[1]s"
  row_widget_num = 2
}

resource "huaweicloud_ces_dashboard" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
  dashboard_id          = huaweicloud_ces_dashboard.base.id
  is_favorite           = true
  row_widget_num        = 3
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDashboard_extend_info_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s"
  row_widget_num = 1
  is_favorite    = true

  extend_info {
    filter                  = "average"
    period                  = "60"
    display_time            = 60
    refresh_time            = 60
    from                    = 1753321953000
    to                      = 1753322953000
    screen_color            = "green"
    enable_screen_auto_play = true
    time_interval           = 10000
    enable_legend           = true
    full_screen_widget_num  = 4
  }
}
`, name)
}

func testDashboard_extend_info_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_ces_dashboard" "test" {
  name           = "%[1]s-update"
  row_widget_num = 2
  is_favorite    = false

  extend_info {
    filter                  = "min"
    period                  = "300"
    display_time            = 15
    refresh_time            = 10
    from                    = 1753321953000
    to                      = 1753322953000
    screen_color            = "blue"
    enable_screen_auto_play = false
    time_interval           = 30000
    enable_legend           = false
    full_screen_widget_num  = 9
  }
}
`, name)
}
