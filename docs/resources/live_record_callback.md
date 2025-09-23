---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_record_callback"
description: |-
  Manages a callback configuration resource within HuaweiCloud.
---

# huaweicloud_live_record_callback

Manages a callback configuration resource within HuaweiCloud.

-> Only one callback configuration can be created for an ingestion domain name.

## Example Usage

### Create a callback configuration for an ingest domain name

```hcl
variable "ingest_domain_name" {}
variable "notify_callback_url"

resource "huaweicloud_live_record_callback" "test" {
  domain_name = var.ingest_domain_name
  url         = var.notify_callback_url
  types       = ["RECORD_NEW_FILE_START"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.
Changing this parameter will create a new resource.

* `url` - (Required, String) Specifies the callback URL for sending recording notifications, which must start with
`http://` or `https://`, and cannot contain message headers or parameters.

* `types` - (Required, List) Specifies the types of recording notifications.
  The valid values are as follows:
  + **RECORD_NEW_FILE_START**: Recording started.
  + **RECORD_FILE_COMPLETE**: Recording file generated.
  + **RECORD_OVER**: Recording completed.
  + **RECORD_FAILED**: Recording failed.

* `sign_type` - (Optional, String) Specifies the sign type.
  The valid values are as follows:
  + **HMACSHA256**
  + **MD5**

  Defaults to **HMACSHA256**.

* `key` - (Optional, String) Specifies the callback key, which is used for authentication. This parameter is configured
  to protect user data security. The value can only contain letters and digits.
  The length cannot be less than `32` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The record callback resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_live_record_callback.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `key`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_live_record_callback" "test" {
  ...
  
  lifecycle {
    ignore_changes = [
      key,
    ]
  }
}
```
