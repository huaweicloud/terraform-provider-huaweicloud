---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_record_callback"
description: ""
---

# huaweicloud_live_record_callback

Manages a callback configuration within HuaweiCloud Live.

-> Only one callback configuration can be created for an ingestion domain name.

## Example Usage

### Create a callback configuration for an ingest domain name

```hcl
variable "ingest_domain_name" {}

resource "huaweicloud_live_domain" "ingestDomain" {
  name = var.ingest_domain_name
  type = "push"
}

resource "huaweicloud_live_record_callback" "callback" {
  domain_name = var.ingest_domain_name
  url         = "http://mycallback.com.cn/record_notify"
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

* `types` - (Required, List) Specifies the types of recording notifications. The options are as follows:
  + **RECORD_NEW_FILE_START**: Recording started.
  + **RECORD_FILE_COMPLETE**: Recording file generated.
  + **RECORD_OVER**: Recording completed.
  + **RECORD_FAILED**: Recording failed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Callback configurations can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_live_record_callback.test 55534eaa-533a-419d-9b40-ec427ea7195a
```
