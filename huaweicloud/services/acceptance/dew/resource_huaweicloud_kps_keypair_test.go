package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	kps "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kps/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

const (
	publicKeyValue = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCDXP+kuDyXeZxf3p58VY/WH8yjzMAmNHXXxQKdHsvEy" +
		"3mU4Egb7ANhDHXHFue0aywhY0XSVULVW1O7qwEtHHcJYBnbt7pKPWrE0rzSpvpvGy9BqxRV44AsHK2VTsXoKGEbXsImwgFwt/q" +
		"5FkAqHMzWB3HJr8tb2rPs7mKjTHs9d1lFxHGehPYgjiFtCGcEOwxJnchGRKgPPa8Tcqkdy3poJ9JrGi6IgXcKtnhS/rZQ+naEG" +
		"VdnC3Qi1Onu0ZixApdbplNcHiw/YBovImSf3JTyCn4U32o+dSFmOuUMovTvysfqtgbouJWEqwG1bQhIzGijW93vDIjjHu/vk/aet7sD rsa-key-20240321"
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
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
		},
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
					resource.TestCheckResourceAttr(resourceName, "user_id", acceptance.HW_USER_ID),
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
				Config: testKeypair_privateKey(rName, publicKeyValue),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_id", "huaweicloud_kms_key.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKeyValue),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "fingerprint"),
				),
			},
			{
				Config: testKeypair_updatePrivateKey1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKeyValue),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "false"),
				),
			},
			{
				Config: testKeypair_updatePrivateKey2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "kms_key_name", "huaweicloud_kms_key.test", "key_alias"),
					resource.TestCheckResourceAttr(resourceName, "public_key", publicKeyValue),
					resource.TestCheckResourceAttr(resourceName, "is_managed", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"private_key", "encryption_type", "kms_key_id", "kms_key_name",
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

func testKeypair_privateKey_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%s"
  pending_days = "7"
}
`, name)
}

func testKeypair_privateKey(rName, publicKeyValue string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  encryption_type = "kms"
  kms_key_id      = huaweicloud_kms_key.test.id
  public_key      = "%[3]s"

  private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAg1z/pLg8l3mcX96efFWP1h/Mo8zAJjR118UCnR7LxMt5lOBI
G+wDYQx1xxbntGssIWNF0lVC1VtTu6sBLRx3CWAZ27e6Sj1qxNK80qb6bxsvQasU
VeOALBytlU7F6ChhG17CJsIBcLf6uRZAKhzM1gdxya/LW9qz7O5io0x7PXdZRcRx
noT2II4hbQhnBDsMSZ3IRkSoDz2vE3KpHct6aCfSaxouiIF3CrZ4Uv62UPp2hBlX
Zwt0ItTp7tGYsQKXW6ZTXB4sP2AaLyJkn9yU8gp+FN9qPnUhZjrlDKL078rH6rYG
6LiVhKsBtW0ISMxoo1vd7wyI4x7v75P2nre7AwIDAQABAoIBAHxBcJM3riC96JuK
cTE8ocSx6Zka6Lp6ruk9Mj66zZZFvaiECdFXis62wYVjdiJjqaefRoExIvm73FVM
6Nzp6vMUUwFRJcZpl9+7Ut6TEZodBbNBBwhDHI8dRVhQ3cS+xTPlixKsOj6L2H5Q
vLrY6SyeeBSF0378PWslBmpewsgdATi5hBl0j7/MoxtvHo7l8Fvyt1pMxpGBWDET
7vXERTjmJCqN8zEejR6j6NpT+SY3DU3Xq+wj/1j3k038noBPrM90ECNkv8t9g85g
n3VczkGBlzYkqEmk0Y1/OEJ4CrnZAQ8WlkgwIuIILz5u3vWF069fEr9LmyMQILpX
KAqVyVECgYEA7JK3MqqJIp+5e4Jh6HlPElmekiB5DSAAWonBmMKiRDxwCHCZYD0W
8HC5HaGl8R3opGuNVnGBJnCO3TYyuHeICI5Iv0GZA42yrGyYBmhIj2i8g5N90tAR
Xv6bHKRiutalEgO1IRj67zHSBLBOc3UXUyYUDB8wdVn0RQAzctQ6ylsCgYEAjiaK
x3h1pTVlr29jYByAa0NoK0MJZs0HsiFNZBTcftfkph45Vugaibf2YJWzGISmT9TK
yIuGUP8ApFBBwkj0tBznG/3Kb9etUwRR02daUvDdL+PnAogu1D9j3OZV5w9foyUr
Fq/FiGqH4IsrxAGjVG2jdKM+7O7xyXrifDGKInkCgYA9Rpc7AV753+M8MX5Ip7sq
ZpojAVQ5aROOX+YMOkWrZPgjx36CpfAeISRhn3AK7xNGGzGFtWqdWUQ32gTzMMrE
ZI5FM6l9eSNRc+NArZw1wQwrDHXnt8r4DvyAQ7fq6xPggaNVylGcyQu7+SqozyhW
eiNxLFbx3nXdtXqeAIilxwKBgEuCg8PT5EJ/K+XWOKasXTcdVm9sq8jU7tqbwB2C
y2IB0u6/LVxR7Q7tDs5dlwZWKHZNpe6D1zSdULz3+QZ4dKxckhOXa/qfSe3IZKL0
ytE2K3iuCl+Y8a9DgQutu0IDM51ZOBtUAY0mcclAhF4ZNKa7mtFxihKYFw4c3cR1
GFiZAoGAMIf+7j8AedEnJRiAzbA4ScFdFLGEcN8G0ENQ/VLe9XnwDy27mMiFsVj0
WMzXvWX6JEkuu4bf/oZ16Uz95IkSocvURuRjSpNexC9efQPH38GUY89Rgf9fl2un
elBITIidTQ9uv/yhqiJmEuVDNncL1W+GvHXk599FLXZpYOr24X4=
-----END RSA PRIVATE KEY-----
EOT
}
`, testKeypair_privateKey_base(rName), rName, publicKeyValue)
}

func testKeypair_updatePrivateKey1(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}
`, testKeypair_privateKey_base(rName), rName)
}

func testKeypair_updatePrivateKey2(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[2]s"
  encryption_type = "kms"
  kms_key_name    = huaweicloud_kms_key.test.key_alias
  private_key = <<EOT
-----BEGIN RSA PRIVATE KEY-----
MIIEogIBAAKCAQEAg1z/pLg8l3mcX96efFWP1h/Mo8zAJjR118UCnR7LxMt5lOBI
G+wDYQx1xxbntGssIWNF0lVC1VtTu6sBLRx3CWAZ27e6Sj1qxNK80qb6bxsvQasU
VeOALBytlU7F6ChhG17CJsIBcLf6uRZAKhzM1gdxya/LW9qz7O5io0x7PXdZRcRx
noT2II4hbQhnBDsMSZ3IRkSoDz2vE3KpHct6aCfSaxouiIF3CrZ4Uv62UPp2hBlX
Zwt0ItTp7tGYsQKXW6ZTXB4sP2AaLyJkn9yU8gp+FN9qPnUhZjrlDKL078rH6rYG
6LiVhKsBtW0ISMxoo1vd7wyI4x7v75P2nre7AwIDAQABAoIBAHxBcJM3riC96JuK
cTE8ocSx6Zka6Lp6ruk9Mj66zZZFvaiECdFXis62wYVjdiJjqaefRoExIvm73FVM
6Nzp6vMUUwFRJcZpl9+7Ut6TEZodBbNBBwhDHI8dRVhQ3cS+xTPlixKsOj6L2H5Q
vLrY6SyeeBSF0378PWslBmpewsgdATi5hBl0j7/MoxtvHo7l8Fvyt1pMxpGBWDET
7vXERTjmJCqN8zEejR6j6NpT+SY3DU3Xq+wj/1j3k038noBPrM90ECNkv8t9g85g
n3VczkGBlzYkqEmk0Y1/OEJ4CrnZAQ8WlkgwIuIILz5u3vWF069fEr9LmyMQILpX
KAqVyVECgYEA7JK3MqqJIp+5e4Jh6HlPElmekiB5DSAAWonBmMKiRDxwCHCZYD0W
8HC5HaGl8R3opGuNVnGBJnCO3TYyuHeICI5Iv0GZA42yrGyYBmhIj2i8g5N90tAR
Xv6bHKRiutalEgO1IRj67zHSBLBOc3UXUyYUDB8wdVn0RQAzctQ6ylsCgYEAjiaK
x3h1pTVlr29jYByAa0NoK0MJZs0HsiFNZBTcftfkph45Vugaibf2YJWzGISmT9TK
yIuGUP8ApFBBwkj0tBznG/3Kb9etUwRR02daUvDdL+PnAogu1D9j3OZV5w9foyUr
Fq/FiGqH4IsrxAGjVG2jdKM+7O7xyXrifDGKInkCgYA9Rpc7AV753+M8MX5Ip7sq
ZpojAVQ5aROOX+YMOkWrZPgjx36CpfAeISRhn3AK7xNGGzGFtWqdWUQ32gTzMMrE
ZI5FM6l9eSNRc+NArZw1wQwrDHXnt8r4DvyAQ7fq6xPggaNVylGcyQu7+SqozyhW
eiNxLFbx3nXdtXqeAIilxwKBgEuCg8PT5EJ/K+XWOKasXTcdVm9sq8jU7tqbwB2C
y2IB0u6/LVxR7Q7tDs5dlwZWKHZNpe6D1zSdULz3+QZ4dKxckhOXa/qfSe3IZKL0
ytE2K3iuCl+Y8a9DgQutu0IDM51ZOBtUAY0mcclAhF4ZNKa7mtFxihKYFw4c3cR1
GFiZAoGAMIf+7j8AedEnJRiAzbA4ScFdFLGEcN8G0ENQ/VLe9XnwDy27mMiFsVj0
WMzXvWX6JEkuu4bf/oZ16Uz95IkSocvURuRjSpNexC9efQPH38GUY89Rgf9fl2un
elBITIidTQ9uv/yhqiJmEuVDNncL1W+GvHXk599FLXZpYOr24X4=
-----END RSA PRIVATE KEY-----
EOT
}
`, testKeypair_privateKey_base(rName), rName)
}

func testKeypair_domain(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_kms_key" "test" {
  key_alias    = "%[1]s"
  pending_days = "7"
}

resource "huaweicloud_kps_keypair" "test" {
  name            = "%[1]s"
  scope           = "account"
  encryption_type = "kms"
  user_id         = "%[2]s"
  kms_key_name    = huaweicloud_kms_key.test.key_alias
}
`, rName, acceptance.HW_USER_ID)
}
