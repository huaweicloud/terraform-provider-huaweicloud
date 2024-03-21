package dew

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	kps "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getKpsKeypairResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.HcKmsV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating KMS v3 client: %s", err)
	}

	request := &kps.ListKeypairDetailRequest{
		KeypairName: state.Primary.ID,
	}

	return client.ListKeypairDetail(request)
}

func TestAccKpsKeypair_basic(t *testing.T) {
	var group kps.ListKeypairDetailResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_kps_keypair.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeypair_basic(rName, "created by acc test"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acc test"),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				Config: testKeypair_basic(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
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

func TestAccKpsKeypair_domain(t *testing.T) {
	var group kps.ListKeypairDetailResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_kps_keypair.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeypair_domain(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "account"),
					resource.TestCheckResourceAttr(resourceName, "encryption_type", "kms"),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_name",
						"huaweicloud_kms_key.test", "key_alias"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_type", "kms_key_name"},
			},
		},
	})
}

func TestAccKpsKeypair_publicKey(t *testing.T) {
	var group kps.ListKeypairDetailResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_kps_keypair.test"
	publicKey, _, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")

	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeypair_publicKey(rName, publicKey),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				Config: testKeypair_basic(rName, "updated by acc test"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
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

func TestAccKpsKeypair_privateKey(t *testing.T) {
	var group kps.ListKeypairDetailResponse

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_kps_keypair.test"
	publicKey, privateKeyPEM, _ := acctest.RandSSHKeyPair("Generated-by-AccTest")
	privateKey := strings.ReplaceAll(privateKeyPEM, "\n", " ")
	rc := acceptance.InitResourceCheck(
		resourceName,
		&group,
		getKpsKeypairResourceFunc,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeypair_privateKey(rName, publicKey, privateKey),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKey),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"private_key",
				},
			},
		},
	})
}

func testKeypair_basic(rName, desc string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_keypair" "test" {
  name        = "%s"
  description = "%s"
}
`, rName, desc)
}

func testKeypair_publicKey(rName, key string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_keypair" "test" {
  name       = "%s"
  public_key = "%s"
}
`, rName, key)
}

func testKeypair_privateKey(rName, publicKey, privateKey string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_keypair" "test" {
  name        = "%[1]s"
  public_key  = "%[2]s"
  private_key = "%[3]s"
}
`, rName, publicKey, privateKey)
}

func testKeypair_domain(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}

resource "huaweicloud_kps_keypair" "test" {
  name            = "%s"
  scope           = "account"
  encryption_type = "kms"
  kms_key_name    = huaweicloud_kms_key.test.key_alias
}
`, rName, rName)
}
