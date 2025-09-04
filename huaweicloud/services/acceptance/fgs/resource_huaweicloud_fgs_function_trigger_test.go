package fgs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/trigger"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/fgs"
)

func getFunctionTriggerFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("fgs", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating FunctionGraph client: %s", err)
	}

	return fgs.GetTriggerById(client, state.Primary.Attributes["function_urn"], state.Primary.Attributes["type"],
		state.Primary.ID)
}

func TestAccFunctionTrigger_basic(t *testing.T) {
	var (
		relatedFunc      function.Function
		timeTrigger      trigger.Trigger
		randName         = acceptance.RandomAccResourceName()
		resNameFunc      = "huaweicloud_fgs_function.test"
		resNameTimerRate = "huaweicloud_fgs_function_trigger.timer_rate"
		resNameTimerCron = "huaweicloud_fgs_function_trigger.timer_cron"

		rcFunc      = acceptance.InitResourceCheck(resNameFunc, &relatedFunc, getFunction)
		rcTimerRate = acceptance.InitResourceCheck(resNameTimerRate, &timeTrigger, getFunctionTriggerFunc)
		rcTimerCron = acceptance.InitResourceCheck(resNameTimerCron, &timeTrigger, getFunctionTriggerFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rcFunc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTimingTrigger_basic_step1(randName),
				Check: resource.ComposeTestCheckFunc(
					// Timing trigger (with rate schedule type)
					rcTimerRate.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resNameTimerRate, "function_urn",
						"${huaweicloud_fgs_function.test.urn}"),
					resource.TestCheckResourceAttr(resNameTimerRate, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "ACTIVE"),
					// Timing trigger (with cron schedule type)
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "type", "TIMER"),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "ACTIVE"),
				),
			},
			{
				Config: testAccFunctionTimingTrigger_basic_step2(randName),
				Check: resource.ComposeTestCheckFunc(
					// Timing trigger (with rate schedule type)
					rcTimerRate.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerRate, "status", "DISABLED"),
					// Timing trigger (with cron schedule type)
					rcTimerCron.CheckResourceExists(),
					resource.TestCheckResourceAttr(resNameTimerCron, "status", "DISABLED"),
				),
			},
			{
				ResourceName:      resNameTimerRate,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerRate),
			},
			{
				ResourceName:      resNameTimerCron,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccFunctionTriggerImportStateFunc(resNameTimerCron),
			},
		},
	})
}

func testAccFunctionTriggerImportStateFunc(rsName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var functionUrn, triggerType, triggerId string
		rs, ok := s.RootModule().Resources[rsName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) of function trigger is not found in the tfstate", rsName)
		}
		functionUrn = rs.Primary.Attributes["function_urn"]
		triggerType = rs.Primary.Attributes["type"]
		triggerId = rs.Primary.ID
		if functionUrn == "" || triggerType == "" || triggerId == "" {
			return "", fmt.Errorf("the function trigger is not exist or related function URN is missing")
		}
		return fmt.Sprintf("%s/%s/%s", functionUrn, triggerType, triggerId), nil
	}
}

func testAccFunctionTrigger_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 128
  timeout     = 10
  runtime     = "Python2.7"
  code_type   = "inline"
  func_code   = "aW1wb3J0IGpzb24KZGVmIGhhbmRsZXIgW1wcyhldybiBvdXRwdXQ="
}`, name)
}

// Test triggers with a limited number (except for Kafka triggers, when released, the elastic network card will be
// locked for one hour and the subnet cannot be deleted).
// The current quantity constraint rules for function triggers are as follows:
//   - The total number of DDS, GAUSSMONGO, DIS, LTS, Kafka and TIMER triggers that can be created under one function
//     version is up to 10.
//   - The maximum number of CTS triggers that can be created under one project is 10.
//   - There is no limit to the number of other triggers.
func testAccFunctionTimingTrigger_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_rate" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_cron" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "ACTIVE"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}

func testAccFunctionTimingTrigger_basic_step2(name string) string {
	return fmt.Sprintf(`
%s

// Timing trigger (with rate schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_rate" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_rate",
    "schedule_type": "Rate",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "3m"
  })
}

// Timing trigger (with cron schedule type)
resource "huaweicloud_fgs_function_trigger" "timer_cron" {
  function_urn = huaweicloud_fgs_function.test.urn
  type         = "TIMER"
  status       = "DISABLED"
  event_data   = jsonencode({
    "name": "%[2]s_cron",
    "schedule_type": "Cron",
    "sync_execution": false,
    "user_event": "Created by acc test",
    "schedule": "@every 1h30m"
  })
}
`, testAccFunctionTrigger_base(name), name)
}

func TestAccFunctionTrigger_eventGrid(t *testing.T) {
	var (
		obj interface{}

		resourceName = "huaweicloud_fgs_function_trigger.test"
		rc           = acceptance.InitResourceCheck(resourceName, &obj, getFunctionTriggerFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckFgsAgency(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionTrigger_eventGrid_step1(name),
				Check: resource.ComposeTestCheckFunc(
					// OBS trigger
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "function_urn",
						"huaweicloud_fgs_function.test", "urn"),
					resource.TestCheckResourceAttr(resourceName, "type", "EVENTGRID"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "cascade_delete_eg_subscription", "true"),
				),
			},
		},
	})
}

func testAccFunctionTrigger_eventGrid_base(name string) string {
	return fmt.Sprintf(`
variable "script_content" {
  default = <<EOT
# -*-coding:utf-8 -*-
import os
import string
import random
import urllib.parse
import shutil
import contextlib
from PIL import Image
from obs import ObsClient

LOCAL_MOUNT_PATH = '/tmp/'


def handler(event, context):
    ak = context.getSecurityAccessKey()
    sk = context.getSecuritySecretKey()
    st = context.getSecurityToken()
    if ak == "" or sk == "" or st == "":
        context.getLogger().error('Failed to access OBS because no temporary '
                                  'AK, SK, or token has been obtained. Please '
                                  'set an agency.')
        return 'Failed to access OBS because no temporary AK, SK, or token ' \
               'has been obtained. Please set an agency. '

    obs_endpoint = context.getUserData('obs_endpoint')
    if not obs_endpoint:
        return 'obs_endpoint is not configured'

    output_bucket = context.getUserData('output_bucket')
    if not output_bucket:
        return 'output_bucket is not configured'

    compress_handler = ThumbnailHandler(context)
    with contextlib.ExitStack() as stack:
        # After upload the thumbnail image to obs, remove the local image
        stack.callback(shutil.rmtree,compress_handler.download_dir)
        data = event.get("data", None)
        return compress_handler.run(data)


class ThumbnailHandler:

    def __init__(self, context):
        self.logger = context.getLogger()
        obs_endpoint = context.getUserData("obs_endpoint")
        self.obs_client = new_obs_client(context, obs_endpoint)
        self.download_dir = gen_local_download_path()
        self.image_local_path = None
        self.src_bucket = None
        self.src_object_key = None
        self.output_bucket = context.getUserData("output_bucket")

    def parse_record(self, record):
        # parses the record to get src_bucket and input_object_key
        (src_bucket, src_object_key) = get_obs_obj_info(record)
        src_object_key = urllib.parse.unquote_plus(src_object_key)
        self.logger.info("src bucket name: %%s", src_bucket)
        self.logger.info("src object key: %%s", src_object_key)
        self.src_bucket = src_bucket
        self.src_object_key = src_object_key
        self.image_local_path = self.download_dir + src_object_key

    def run(self, record):
        self.parse_record(record)
        # Download the original image from obs to local tmp dir.
        if not self.download_from_obs():
            return "ERROR"
        # Thumbnail original image
        thumbnail_object_key, thumbnail_file_path = self.thumbnail()
        self.logger.info("thumbnail_object_key: %%s, thumbnail_file_path: %%s",
                        thumbnail_object_key, thumbnail_file_path)
        # Upload thumbnail image to obs
        if not self.upload_file_to_obs(thumbnail_object_key,
                                    thumbnail_file_path):
            return "ERROR"
        return "OK"

    def download_from_obs(self):
        resp = self.obs_client. \
            getObject(self.src_bucket, self.src_object_key,
                      downloadPath=self.image_local_path)
        if resp.status < 300:
            return True
        else:
            self.logger.error('failed to download from obs, '
                              'errorCode: %%s, errorMessage: %%s' %% (
                                  resp.errorCode, resp.errorMessage))
            return False

    def upload_file_to_obs(self, object_key, local_file):
        resp = self.obs_client.putFile(self.output_bucket,
                                       object_key, local_file)
        if resp.status < 300:
            return True
        else:
            self.logger.error('failed to upload file to obs, errorCode: %%s, '
                              'errorMessage: %%s' %% (resp.errorCode,
                                                    resp.errorMessage))
            return False

    def thumbnail(self):
        image = Image.open(self.image_local_path)
        image_w, image_h = image.size
        new_image_w, new_image_h = (int(image_w / 2), int(image_h / 2))
    
        image.thumbnail((new_image_w, new_image_h), resample=Image.LANCZOS)
    
        (path, filename) = os.path.split(self.src_object_key)
        if path != "" and not path.endswith("/"):
            path = path + "/"
        (filename, ext) = os.path.splitext(filename)
        ext_lower = ext.lower()

        thumbnail_object_key = path + filename + \
                           "_thumbnail_h_" + str(new_image_h) + \
                           "_w_" + str(new_image_w) + ext

        if hasattr(image, '_getexif'):
            image.info.pop('exif', None)

        if new_image_w * new_image_h < 500000:
            webp_ext = '.webp'
            thumbnail_object_key = path + filename + \
                                  "_thumbnail_h_" + str(new_image_h) + \
                                  "_w_" + str(new_image_w) + webp_ext
            thumbnail_file_path = self.download_dir + path + filename + \
                                 "_thumbnail_h_" + str(new_image_h) + \
                                 "_w_" + str(new_image_w) + webp_ext
            image.save(thumbnail_file_path, format='WEBP', quality=80, method=6)
        elif ext_lower in ['.jpg', '.jpeg']:
            thumbnail_file_path = self.download_dir + path + filename + \
                                 "_thumbnail_h_" + str(new_image_h) + \
                                 "_w_" + str(new_image_w) + ext
            image.save(thumbnail_file_path, quality=85, optimize=True, progressive=True)
        elif ext_lower == '.png':
            thumbnail_file_path = self.download_dir + path + filename + \
                                 "_thumbnail_h_" + str(new_image_h) + \
                                 "_w_" + str(new_image_w) + ext
            image.save(thumbnail_file_path, optimize=True, compress_level=6)
        else:
            thumbnail_file_path = self.download_dir + path + filename + \
                                 "_thumbnail_h_" + str(new_image_h) + \
                                 "_w_" + str(new_image_w) + ext
            image.save(thumbnail_file_path, quality=85, optimize=True)

        return thumbnail_object_key, thumbnail_file_path


# generate a temporary directory for downloading things
# from OBS and compress them.
def gen_local_download_path():
    letters = string.ascii_letters
    download_dir = LOCAL_MOUNT_PATH + ''.join(
        random.choice(letters) for i in range(16)) + '/'
    os.makedirs(download_dir)
    return download_dir


def new_obs_client(context, obs_server):
    ak = context.getSecurityAccessKey()
    sk = context.getSecuritySecretKey()
    st = context.getSecurityToken()
    return ObsClient(access_key_id=ak, secret_access_key=sk, security_token=st, server=obs_server)


def get_obs_obj_info(record):
    if 's3' in record:
        s3 = record['s3']
        return s3['bucket']['name'], s3['object']['key']
    else:
        obs_info = record['obs']
        return obs_info['bucket']['name'], obs_info['object']['key']

EOT
}

data "huaweicloud_fgs_dependencies" "test" {
  type = "public"
  name = "pillow-7.1.2"
}

data "huaweicloud_fgs_dependency_versions" "test" {
  dependency_id = try(data.huaweicloud_fgs_dependencies.test.packages[0].id, "NOT_FOUND")
  version       = 1
}

resource "huaweicloud_fgs_function" "test" {
  name        = "%[1]s"
  agency      = "%[2]s"
  app         = "default"
  handler     = "index.handler"
  memory_size = 256
  timeout     = 40
  runtime     = "Python3.9"
  code_type   = "inline"
  func_code   = base64encode(var.script_content)
  depend_list = data.huaweicloud_fgs_dependency_versions.test.versions[*].id

  user_data = jsonencode({
    "output_bucket" = huaweicloud_obs_bucket.target.bucket
    "obs_endpoint"  = "obs.%[3]s.myhuaweicloud.com"
  })
}

// OBS buckets for testing
resource "huaweicloud_obs_bucket" "source" {
  bucket        = "%[1]s-source"
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "target" {
  bucket        = "%[1]s-target"
  acl           = "private"
  force_destroy = true
}

data "huaweicloud_eg_event_channels" "test" {
  provider_type = "OFFICIAL"
  name          = "default"
}
`, name, acceptance.HW_FGS_AGENCY_NAME, acceptance.HW_REGION_NAME)
}

func testAccFunctionTrigger_eventGrid_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_fgs_function_trigger" "test" {
  depends_on = [huaweicloud_obs_bucket.source]

  function_urn                   = huaweicloud_fgs_function.test.urn
  type                           = "EVENTGRID"
  cascade_delete_eg_subscription = true
  status                         = "ACTIVE"
  event_data                     = jsonencode({
    "channel_id"   = try(data.huaweicloud_eg_event_channels.test.channels[0].id, "")
    "channel_name" = try(data.huaweicloud_eg_event_channels.test.channels[0].name, "")
    "source_name"  = "HC.OBS.DWR"
    "trigger_name" = "demo" # Just the name suffix
    "agency"       = "fgs_to_eg"
    "bucket"       = huaweicloud_obs_bucket.source.bucket
    "event_types"  = ["OBS:DWR:ObjectCreated:PUT", "OBS:DWR:ObjectCreated:POST"]
    "Key_encode"   = true
  })

  lifecycle {
    ignore_changes = [
      event_data,
    ]
  }
}
`, testAccFunctionTrigger_eventGrid_base(name))
}
