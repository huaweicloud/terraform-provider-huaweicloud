---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_url_authentication"
description: |-
  Manages an URL authentication resource within HuaweiCloud.
---

# huaweicloud_live_url_authentication

Manages an URL authentication resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

-> Before creating the resource, you need to cofiguration the URL validation first (creating
  the URL validation resource first). Refer to
  [URL Validation](https://support.huaweicloud.com/iLive-live/live_01_0049.html) for more details.

## Example Usage

```hcl
variable "domain_name" {}
variable "type" {}
variable "stream_name" {}
variable "app_name" {}

resource "huaweicloud_live_url_authentication" "test" {
  domain_name = var.domain_name
  type        = var.type
  stream_name = var.stream_name
  app_name    = var.app_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the domain name to which the URL validation belongs.
  Including the ingest domain name and streaming domain name.
  Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the domain name.
  The valid values are as follow:
  + **push**: Indicates ingest domain name.
  + **pull**: Indicates streaming domain name.

  Changing this parameter will create a new resource.

* `stream_name` - (Required, String, ForceNew) Specifies the stream name.
  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name.
  Changing this parameter will create a new resource.

* `check_level` - (Optional, Int, ForceNew) Specifies the check level.
  This parameter is valid and mandatory only when the signing method is **c_aes** in the URL validation.
  The valid values are as follows:
  + `3`: The system checks only LiveID but not the validity of the signed URL.
  + `5`: The system checks LiveID and the validity of timestamp.

  -> The value of the LiveID consists of `app_name` and `stream_name`: <app_name>/<stream_name>.
  Changing this parameter will create a new resource.

* `start_time` - (Optional, String, ForceNew) Specifies the start time of the valid access time defined by the user.
  Changing this parameter will create a new resource.
  The time is in UTC, the format is **yyyy-mm-ddThh:mm:ssZ**, e.g. **2024-06-01T15:03:01Z**. Defaults to current time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

* `key_chain` - The generated signed URLs.
