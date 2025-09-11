package sfsturbo

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sfsturbo"
)

func getSfsTurboResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		product = "sfs-turbo"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS Turbo client: %s", err)
	}

	return sfsturbo.GetTurboDetail(client, state.Primary.ID)
}

func TestAccSFSTurbo_basic(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSFSTurbo_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", rName)),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_update"),
					resource.TestCheckResourceAttrPair(resourceName, "security_group_id",
						"huaweicloud_networking_secgroup.test_update", "id"),
					resource.TestCheckResourceAttr(resourceName, "status", "232"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_crypt(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_crypt(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "STANDARD"),
					resource.TestCheckResourceAttr(resourceName, "enhanced", "false"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "crypt_key_id"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_withEpsId(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccSFSTurbo_prePaid(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_prePaid_step1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
				),
			},
			{
				Config: testAccSFSTurbo_prePaid_step2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "600"),
				),
			},
		},
	})
}

func TestAccSFSTurbo_hpcShareType(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_shareTypeHpc(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_id", ""),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "HPC"),
					resource.TestCheckResourceAttr(resourceName, "hpc_bandwidth", "40M"),
					resource.TestCheckResourceAttr(resourceName, "size", "1228"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
				),
			},
			{
				Config: testAccSFSTurbo_shareTypeHpc_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "size", "4915"),
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

func TestAccSFSTurbo_hpcCacheShareType(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_shareTypeHpcCache(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "HPC_CACHE"),
					resource.TestCheckResourceAttr(resourceName, "hpc_cache_bandwidth", "2G"),
					resource.TestCheckResourceAttr(resourceName, "size", "4096"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
				),
			},
			{
				Config: testAccSFSTurbo_shareTypeHpcCache_update1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "hpc_cache_bandwidth", "2G"),
					resource.TestCheckResourceAttr(resourceName, "size", "8192"),
				),
			},
			{
				Config: testAccSFSTurbo_shareTypeHpcCache_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "hpc_cache_bandwidth", "4G"),
					resource.TestCheckResourceAttr(resourceName, "size", "8192"),
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

func TestAccSFSTurbo_backupId(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSFSTurboBackupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_backupId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "backup_id", acceptance.HW_SFS_TURBO_BACKUP_ID),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "HPC"),
					resource.TestCheckResourceAttr(resourceName, "hpc_bandwidth", "20M"),
					resource.TestCheckResourceAttr(resourceName, "size", "1228"),
					resource.TestCheckResourceAttr(resourceName, "status", "200"),
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

func TestAccSFSTurbo_checkError(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSFSTurbo_checkError1(rName),
				ExpectError: regexp.MustCompile("`hpc_bandwidth` is required when share type is HPC"),
			},
			{
				Config:      testAccSFSTurbo_checkError2(rName),
				ExpectError: regexp.MustCompile("`hpc_cache_bandwidth` is required when share type is HPC_CACHE"),
			},
			{
				Config:      testAccSFSTurbo_checkError3(rName),
				ExpectError: regexp.MustCompile("HPC_CACHE share type only support in postpaid charging mode"),
			},
			{
				Config: testAccSFSTurbo_checkError4(rName),
				ExpectError: regexp.MustCompile("`hpc_bandwidth` and `hpc_cache_bandwidth` cannot be set" +
					" when share type is STANDARD or PERFORMANCE"),
			},
		},
	})
}

func TestAccSFSTurbo_checkUpdateError(t *testing.T) {
	var turbo interface{}
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sfs_turbo.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&turbo,
		getSfsTurboResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSTurbo_checkUpdateErrorBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "share_proto", "NFS"),
					resource.TestCheckResourceAttr(resourceName, "share_type", "PERFORMANCE"),
					resource.TestCheckResourceAttr(resourceName, "size", "500"),
				),
			},
			{
				Config:      testAccSFSTurbo_checkUpdateErrorBasic_update(rName),
				ExpectError: regexp.MustCompile("only `HPC_CACHE` share type support updating HPC cache bandwidth"),
			},
		},
	})
}

func testAccSFSTurbo_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_networking_secgroup" "test_update" {
  name                 = "%[2]s_update"
  delete_default_rules = true
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%[2]s_update"
  size              = 600
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test_update.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar_update"
    key = "value_update"
  }
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_crypt(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_sfs_turbo" "test" {
  name              = "%s"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  crypt_key_id      = huaweicloud_kms_key.test.id
}
`, common.TestBaseNetwork(rName), rName, rName)
}

func testAccSFSTurbo_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  name                   = "%s"
  size                   = 500
  share_proto            = "NFS"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id  = "%s"
}
`, common.TestBaseNetwork(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccSFSTurbo_prePaid_step1(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%[2]s"
  size                   = 500
  share_proto            = "NFS"
  enterprise_project_id  = "0"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccSFSTurbo_prePaid_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  security_group_id      = huaweicloud_networking_secgroup.test.id
  availability_zone      = data.huaweicloud_availability_zones.test.names[0]
  name                   = "%[2]s"
  size                   = 600
  share_proto            = "NFS"
  enterprise_project_id  = "0"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = true
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_shareTypeHpc(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 1228
  share_proto       = "NFS"
  share_type        = "HPC"
  hpc_bandwidth     = "40M"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_shareTypeHpc_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 4915
  share_proto       = "NFS"
  share_type        = "HPC"
  hpc_bandwidth     = "40M"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_backupId(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 1228
  share_proto       = "NFS"
  share_type        = "HPC"
  hpc_bandwidth     = "20M"
  backup_id         = "%[3]s"
}
`, common.TestBaseNetwork(name), name, acceptance.HW_SFS_TURBO_BACKUP_ID)
}

func testAccSFSTurbo_shareTypeHpcCache(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 4096
  share_proto         = "NFS"
  share_type          = "HPC_CACHE"
  hpc_cache_bandwidth = "2G"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_shareTypeHpcCache_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 8192
  share_proto         = "NFS"
  share_type          = "HPC_CACHE"
  hpc_cache_bandwidth = "2G"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_shareTypeHpcCache_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 8192
  share_proto         = "NFS"
  share_type          = "HPC_CACHE"
  hpc_cache_bandwidth = "4G"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkError1(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 1228
  share_proto       = "NFS"
  share_type        = "HPC"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkError2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 1228
  share_proto       = "NFS"
  share_type        = "HPC_CACHE"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkError3(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 4096
  share_proto         = "NFS"
  share_type          = "HPC_CACHE"
  hpc_cache_bandwidth = "2G"

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  auto_renew    = false
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkError4(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 500
  share_proto         = "NFS"
  hpc_cache_bandwidth = "2G"

}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkUpdateErrorBasic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id            = huaweicloud_vpc.test.id
  subnet_id         = huaweicloud_vpc_subnet.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  name              = "%[2]s"
  size              = 500
  share_proto       = "NFS"
  share_type        = "PERFORMANCE"
}
`, common.TestBaseNetwork(name), name)
}

func testAccSFSTurbo_checkUpdateErrorBasic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_sfs_turbo" "test" {
  vpc_id              = huaweicloud_vpc.test.id
  subnet_id           = huaweicloud_vpc_subnet.test.id
  security_group_id   = huaweicloud_networking_secgroup.test.id
  availability_zone   = data.huaweicloud_availability_zones.test.names[0]
  name                = "%[2]s"
  size                = 500
  share_proto         = "NFS"
  share_type          = "PERFORMANCE"
  hpc_cache_bandwidth = "2G"
}
`, common.TestBaseNetwork(name), name)
}
