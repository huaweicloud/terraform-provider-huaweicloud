---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_sign_url"
description: |-
  Manages a sign URL resource within HuaweiCloud.
---

# huaweicloud_live_sign_url

Manages a sign URL resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not change the current status.

-> Before creating the resource, you need to cofiguration the URL validation first (creating
  the URL validation resource first).

## Example Usage

```hcl
variable "domain_name" {}
variable "type" {}
variable "stream_name" {}
variable "app_name" {}

resource "huaweicloud_live_sign_url" "test" {
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

* `domain_name` - (Required, String, ForceNew) Specifies the domain name.
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

* `check_level` - (Optional, String, ForceNew) Specifies the check level.
  Changing this parameter will create a new resource.
  The valid values are as follows:
  + **3**: Indicates only checks the LiveID matches or not but not the validity of the signed URL.
  + **5**: Indicates checks the LiveID matches or not and the validity of timestamp.

  -> 1.If the value of the `auth_type` in the URL validation is set to **c_aes**, this parameter is mandatory.
    <br/>2.The value of the LiveID consists of `app_name` and `stream_name`: <app_name>/<stream_name>.

* `start_time` - (Optional, String, ForceNew) Specifies the start time of the valid access time defined by the user.
  Changing this parameter will create a new resource.
  The time is in UTC, the format is **yyyy-mm-ddThh:mm:ssZ**. e.g. **2024-06-01T15:03:01Z**
  If this parameter is not specified or is left empty, the current time is used by default.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, UUID format.

* `key_chain` - The generated signed URLs.
