package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/dli/v2/spark/resources"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/dli"
)

func getPackageResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.DliV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating DLI v2 client: %s", err)
	}

	return dli.GetDliDependentPackageInfo(c, state.Primary.ID)
}

func TestAccDliPackage_basic(t *testing.T) {
	var pkg resources.Resource

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_dli_package.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pkg,
		getPackageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliPackage_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "pyFile"),
					resource.TestCheckResourceAttr(resourceName, "object_path", fmt.Sprintf(
						"https://%s.obs.%s.myhuaweicloud.com/dli/packages/simple_pyspark_test_DLF_refresh.py",
						rName, acceptance.HW_REGION_NAME)),
					resource.TestCheckResourceAttr(resourceName, "object_name", "simple_pyspark_test_DLF_refresh.py"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "status", "READY"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccDliPackage_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "group_name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
				),
			},
		},
	})
}

func testAccDliPackage_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_obs_bucket" "test" {
  bucket = "%s"
  acl    = "private"
}

resource "huaweicloud_obs_bucket_object" "test" {
  bucket  = huaweicloud_obs_bucket.test.bucket
  key     = "dli/packages/simple_pyspark_test_DLF_refresh.py"
  content = <<EOF
#!/usr/bin/env python
# _*_ coding: utf-8 _*_

import sys
import logging
from operator import add
import time

from pyspark.sql import SparkSession
from pyspark.sql import SQLContext

sparkSession = SparkSession.builder.appName("simple pyspark test DLF refresh").getOrCreate()
sc = SQLContext(sparkSession.sparkContext)

logging.basicConfig(format='%%(message)s', level=logging.INFO)
logger = logging.getLogger("Whatever")
logger.info("[DBmethods.py] HELLOOOOOOOOOOO")


sc._jsc.hadoopConfiguration().set("fs.obs.access.key", "%s")
sc._jsc.hadoopConfiguration().set("fs.obs.secret.key", "%s")
sc._jsc.hadoopConfiguration().set("fs.obs.endpoint", "obs.cn-north-4.myhuaweicloud.com")


# Read private bucket with encryption using AK/SK
private_encrypted_file = "obs://dedicated-for-terraform-acc-test/dli/spark/people.csv"

df = sparkSession.read.options(header='True', inferSchema='True', delimiter=',').csv(private_encrypted_file)
df.show()
df.printSchema()
print(df)
print(df.count())
print(time.time())


my_string_to_print = "{} - {}".format(int(time.time()), df.count()/2)
file_name = "my_file-{}-{}".format(int(time.time()), df.count()/2)


print(my_string_to_print)
print(file_name)

private_encrypted_output_folder = "obs://dedicated-for-terraform-acc-test/dli/result/"
# my_string_to_print.write.mode('overwrite').csv(private_encrypted_output_folder)

final_path = "{}{}".format(private_encrypted_output_folder, file_name)
print(final_path)


sparkSession.sparkContext.parallelize([my_string_to_print]).coalesce(1).saveAsTextFile(final_path)


EOF
  content_type = "text/py"
}
`, name, acceptance.HW_ACCESS_KEY, acceptance.HW_SECRET_KEY)
}

func testAccDliPackage_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_package" "test" {
  depends_on  = [huaweicloud_obs_bucket_object.test]
  group_name  = "%s"
  type        = "pyFile"
  object_path = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"

  tags = {
    foo = "bar"
  }
}
`, testAccDliPackage_base(rName), rName)
}

func testAccDliPackage_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_package" "test" {
  depends_on  = [huaweicloud_obs_bucket_object.test]
  group_name  = "%s"
  type        = "pyFile"
  object_path = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"

  tags = {
    owner = "terraform"
  }
}
`, testAccDliPackage_base(rName), rName)
}

func TestAccDliPackage_not_groupName(t *testing.T) {
	var (
		pkg          resources.Resource
		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_dli_package.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pkg,
		getPackageResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliUpdatedOwner(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccDliPackage_not_groupName(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckNoResourceAttr(resourceName, "group_name"),
					resource.TestCheckResourceAttr(resourceName, "type", "modelFile"),
					resource.TestCheckResourceAttr(resourceName, "object_path", fmt.Sprintf(
						"https://%s.obs.%s.myhuaweicloud.com/dli/packages/simple_pyspark_test_DLF_refresh.py",
						name, acceptance.HW_REGION_NAME)),
					resource.TestCheckResourceAttr(resourceName, "object_name", "simple_pyspark_test_DLF_refresh.py"),
					resource.TestCheckResourceAttr(resourceName, "status", "READY"),
					resource.TestCheckResourceAttr(resourceName, "owner", acceptance.HW_DLI_UPDATED_OWNER),
				),
			},
			{
				Config: testAccDliPackage_not_groupName_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckNoResourceAttr(resourceName, "group_name"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "owner", acceptance.HW_DLI_OWNER),
				),
			},
		},
	})
}

func testAccDliPackage_not_groupName(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_package" "test" {
  depends_on  = [huaweicloud_obs_bucket_object.test]
  type        = "modelFile"
  object_path = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"
  owner       = "%s"

  tags = {
    foo = "bar"
  }
}
`, testAccDliPackage_base(name), acceptance.HW_DLI_UPDATED_OWNER)
}

func testAccDliPackage_not_groupName_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_package" "test" {
  depends_on  = [huaweicloud_obs_bucket_object.test]
  type        = "modelFile"
  object_path = "https://${huaweicloud_obs_bucket.test.bucket_domain_name}/dli/packages/simple_pyspark_test_DLF_refresh.py"
  owner       = "%s"

  tags = {
    owner = "terraform"
  }
}
`, testAccDliPackage_base(name), acceptance.HW_DLI_OWNER)
}
