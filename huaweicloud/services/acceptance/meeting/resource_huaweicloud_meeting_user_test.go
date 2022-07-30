package meeting

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/meeting/v1/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/meeting"
)

func getUserFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	token, err := meeting.NewMeetingToken(conf, state)
	if err != nil {
		return nil, err
	}

	opt := users.GetOpts{
		Token:   token,
		Account: state.Primary.ID,
	}
	return users.Get(meeting.NewMeetingV1Client(conf), opt)
}

func TestAccUser_basic(t *testing.T) {
	var (
		user         users.User
		resourceName = "huaweicloud_meeting_user.test"
		rName        = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "account", rName),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Name"),
					resource.TestCheckResourceAttr(resourceName, "country", "chinaPR"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "email", "123456789@example.com"),
					resource.TestCheckResourceAttr(resourceName, "english_name", "Test English Name"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+8612345678987"),
					resource.TestCheckResourceAttr(resourceName, "hide_phone", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature", "HuaweiSignature"),
					resource.TestCheckResourceAttr(resourceName, "title", "Test Title"),
					resource.TestCheckResourceAttr(resourceName, "sort_level", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "third_account"),
				),
			},
			{
				Config: testAccUser_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "account", rName),
					resource.TestCheckResourceAttr(resourceName, "name", "New Test Name"),
					resource.TestCheckResourceAttr(resourceName, "country", "chinaHKG"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
					resource.TestCheckResourceAttr(resourceName, "email", "987654321@example.com"),
					resource.TestCheckResourceAttr(resourceName, "english_name", "New Test English Name"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+85212345678"),
					resource.TestCheckResourceAttr(resourceName, "hide_phone", "false"),
					resource.TestCheckResourceAttr(resourceName, "is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature", "NewHuaweiSignature"),
					resource.TestCheckResourceAttr(resourceName, "title", "New Test Title"),
					resource.TestCheckResourceAttr(resourceName, "sort_level", "10"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "third_account"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccUserManagementImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"corp_id",
					"is_send_notify",
					"password",
				},
			},
		},
	})
}

func TestAccUser_thirdAccount(t *testing.T) {
	var (
		user         users.User
		resourceName = "huaweicloud_meeting_user.test"
		rName        = acceptance.RandomAccResourceName()
		thirdAccount = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_thirdAccount(rName, thirdAccount),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "account", rName),
					resource.TestCheckResourceAttr(resourceName, "third_account", thirdAccount),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Name"),
					resource.TestCheckResourceAttr(resourceName, "country", "chinaPR"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "email", "123456789@example.com"),
					resource.TestCheckResourceAttr(resourceName, "english_name", "Test English Name"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+8612345678987"),
					resource.TestCheckResourceAttr(resourceName, "hide_phone", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature", "HuaweiSignature"),
					resource.TestCheckResourceAttr(resourceName, "title", "Test Title"),
					resource.TestCheckResourceAttr(resourceName, "sort_level", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccUserManagementImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"corp_id",
					"is_send_notify",
					"password",
				},
			},
		},
	})
}

func TestAccUser_withoutAccount(t *testing.T) {
	var (
		user         users.User
		resourceName = "huaweicloud_meeting_user.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_withoutAccount(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_MEETING_USER_ID),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Name"),
					resource.TestCheckResourceAttr(resourceName, "country", "chinaPR"),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttr(resourceName, "email", "123456789@example.com"),
					resource.TestCheckResourceAttr(resourceName, "english_name", "Test English Name"),
					resource.TestCheckResourceAttr(resourceName, "phone", "+8612345678987"),
					resource.TestCheckResourceAttr(resourceName, "hide_phone", "true"),
					resource.TestCheckResourceAttr(resourceName, "is_send_notify", "true"),
					resource.TestCheckResourceAttr(resourceName, "signature", "HuaweiSignature"),
					resource.TestCheckResourceAttr(resourceName, "title", "Test Title"),
					resource.TestCheckResourceAttr(resourceName, "sort_level", "5"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "account"),
					resource.TestCheckResourceAttrSet(resourceName, "third_account"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccUserManagementImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"corp_id",
					"is_send_notify",
					"password",
				},
			},
		},
	})
}

func TestAccUser_admin(t *testing.T) {
	var (
		user         users.User
		resourceName = "huaweicloud_meeting_user.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&user,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAppAuth(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccUser_admin(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "app_id", acceptance.HW_MEETING_APP_ID),
					resource.TestCheckResourceAttr(resourceName, "app_key", acceptance.HW_MEETING_APP_KEY),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Name"),
					resource.TestCheckResourceAttr(resourceName, "is_admin", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccUserManagementImportStateIdFunc(),
				ImportStateVerifyIgnore: []string{
					"corp_id",
					"user_id",
					"password",
					"is_send_notify",
					"is_admin",
				},
			},
		},
	})
}

func testAccUserManagementImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var (
			accountName, password         string
			appId, appKey, corpId, userId string
			account                       string
		)
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "huaweicloud_meeting_user" {
				accountName = rs.Primary.Attributes["account_name"]
				password = rs.Primary.Attributes["account_password"]
				appId = rs.Primary.Attributes["app_id"]
				appKey = rs.Primary.Attributes["app_key"]
				corpId = rs.Primary.Attributes["corp_id"]
				userId = rs.Primary.Attributes["user_id"]
				account = rs.Primary.ID
			}
		}
		if account != "" && accountName != "" && password != "" {
			return fmt.Sprintf("%s/%s/%s", account, accountName, password), nil
		}
		if account != "" && appId != "" && appKey != "" {
			return fmt.Sprintf("%s/%s/%s/%s/%s", account, appId, appKey, corpId, userId), nil
		}
		return "", fmt.Errorf("resource not found: %s", account)
	}
}

func testAccUser_basic(account string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  account        = "%[4]s"
  name           = "Test Name"
  password       = "HuaweiTest@123"
  country        = "chinaPR"
  description    = "Created by script"
  email          = "123456789@example.com"
  english_name   = "Test English Name"
  phone          = "+8612345678987"
  hide_phone     = true
  is_send_notify = true
  signature      = "HuaweiSignature"
  title          = "Test Title"
  sort_level     = 5
  status         = "0"
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID, account)
}

func testAccUser_update(account string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  account        = "%[4]s"
  name           = "New Test Name"
  password       = "HuaweiTest@123"
  country        = "chinaHKG"
  description    = "Updated by script"
  email          = "987654321@example.com"
  english_name   = "New Test English Name"
  phone          = "+85212345678"
  hide_phone     = false
  is_send_notify = true
  signature      = "NewHuaweiSignature"
  title          = "New Test Title"
  sort_level     = 10
  status         = "1"
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID, account)
}

func testAccUser_thirdAccount(account, thirdAccount string) string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  account        = "%[4]s"
  third_account  = "%[5]s"
  name           = "Test Name"
  password       = "HuaweiTest@123"
  country        = "chinaPR"
  description    = "Created by script"
  email          = "123456789@example.com"
  english_name   = "Test English Name"
  phone          = "+8612345678987"
  hide_phone     = true
  is_send_notify = true
  signature      = "HuaweiSignature"
  title          = "Test Title"
  sort_level     = 5
  status         = "0"
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID, account, thirdAccount)
}

func testAccUser_withoutAccount() string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"
  user_id = "%[3]s"

  name           = "Test Name"
  password       = "HuaweiTest@123"
  country        = "chinaPR"
  description    = "Created by script"
  email          = "123456789@example.com"
  english_name   = "Test English Name"
  phone          = "+8612345678987"
  hide_phone     = true
  is_send_notify = true
  signature      = "HuaweiSignature"
  title          = "Test Title"
  sort_level     = 5
  status         = "0"
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY, acceptance.HW_MEETING_USER_ID)
}

func testAccUser_admin() string {
	return fmt.Sprintf(`
resource "huaweicloud_meeting_user" "test" {
  app_id  = "%[1]s"
  app_key = "%[2]s"

  name     = "Test Name"
  password = "HuaweiTest@123"
  is_admin = true
}
`, acceptance.HW_MEETING_APP_ID, acceptance.HW_MEETING_APP_KEY)
}
