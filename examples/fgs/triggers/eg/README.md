# Create FunctionGraph OBS Trigger for Image Compression

This example provides best practice code for using Terraform to create a FunctionGraph OBS trigger for automatic image
compression in HuaweiCloud FunctionGraph service. This solution is suitable for compressing single or batch images.
High-quality image files typically occupy large storage space and bandwidth, causing slow loading speeds for websites
and applications.

By combining Object Storage Service (OBS) with FunctionGraph functions, you can build an efficient image compression
solution that automatically processes images stored in OBS buckets, achieving storage space optimization and efficient
resource utilization.

By deploying this practice script, you will get two buckets for compressing images (source is the bucket for uploading
images, and target is the bucket for outputting compressed images), a function trigger for detecting image uploads, and
its underlying execution function.

## Prerequisites

* A HuaweiCloud account (with 'Tenant Administrator' permission)
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* FunctionGraph service enabled in the target region
* OBS service enabled in the target region
* EventGrid service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the FunctionGraph service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `source_bucket_name` - The source bucket name of the OBS service (used to upload images to be compressed)
* `target_bucket_name` - The target bucket name of the OBS service (the location where the compressed images are
  generated)
* `function_name` - The name of the FunctionGraph function
* `function_agency_name` - The agency name that the FunctionGraph function uses
* `trigger_name_suffix` - The suffix name of the OBS trigger (the format of the trigger name is
  `default-{function_name}-{trigger_name_suffix}`)

#### Optional Variables

* `function_memory_size` - The memory size of the function in MB (default: 256)
* `function_timeout` - The timeout of the function in seconds (default: 40)
* `function_runtime` - The runtime of the function (default: "Python3.6")
* `function_code` - The source code of the function (default: provided Python image compression example)
* `function_description` - The description of the function (default: "")
* `trigger_status` - The status of the OBS trigger (default: "ACTIVE")
* `trigger_agency_name` - The agency name that the OBS trigger uses (default: "fgs_to_eg")

## Image Compression Features

The function automatically:

* Downloads the original image from the source OBS bucket
* Compresses the image to 50% of its original size (width and height)
* Saves the compressed image with 85% quality
* Uploads the compressed image to the target OBS bucket
* Names the compressed image with a suffix indicating the new dimensions

## Example Function Code

The default function code is a Python function that automatically compresses images:

```python
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
        self.logger.info("src bucket name: %s", src_bucket)
        self.logger.info("src object key: %s", src_object_key)
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
        self.logger.info("thumbnail_object_key: %s, thumbnail_file_path: %s",
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
                              'errorCode: %s, errorMessage: %s' % (
                                  resp.errorCode, resp.errorMessage))
            return False

    def upload_file_to_obs(self, object_key, local_file):
        resp = self.obs_client.putFile(self.output_bucket,
                                       object_key, local_file)
        if resp.status < 300:
            return True
        else:
            self.logger.error('failed to upload file to obs, errorCode: %s, '
                              'errorMessage: %s' % (resp.errorCode,
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
```

## Dependencies

The function automatically includes the following dependencies:

* `pillow-7.1.2` - Python image processing library for compression operations

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  source_bucket_name   = "tf-test-bucket-source"
  target_bucket_name   = "tf-test-bucket-target"
  function_name        = "tf-test-function-image-thumbnail"
  function_agency_name = "function_all_trust"
  trigger_name_suffix  = "demo"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The OBS trigger is dependent on the FunctionGraph function and EventGrid service
* The function requires the Pillow library for image processing (automatically included via dependencies)
* Both source and target buckets will be created in the specified region
* The trigger supports both ACTIVE and DISABLED status
* The function supports Python3.6 runtime with Pillow image processing capabilities
* Compressed images are automatically named with dimension information
* The solution handles both PUT and POST upload operations

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.78.1 |
