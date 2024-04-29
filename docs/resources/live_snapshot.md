---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_snapshot"
description: ""
---

# huaweicloud_live_snapshot

Manages a Live snapshot resource within HuaweiCloud.

## Example Usage

```hcl
variable "storage_bucket" {}
variable "storage_path" {}
variable "domain_name" {}

resource "huaweicloud_live_snapshot" "test"{
  domain_name    = var.domain_name
  app_name       = "live"
  frequency      = 10
  storage_mode   = "0"
  storage_bucket = var.storage_bucket
  storage_path   = var.storage_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.

  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name.

  Changing this parameter will create a new resource.

* `frequency` - (Required, Int) Specifies the screenshot frequency. Value range: 5-3600. Unit: second.

* `storage_mode` - (Required, Int) Specifies the store mode in OBS bucket. The options are as follows:
  + **0**: real time snapshot. Name the snapshot file with a timestamp and
    save all screenshot files to the OBS bucket.
    For example, snapshot/{domain}/{app_name}/{stream_name}/{UnixTimestamp}.jpg
  + **1**: coverage snapshot. Only the latest snapshot will be saved, old snapshot
    will be covered by new snapshot.
    For example, snapshot/{domain}/{app_name}/{stream_name}.jpg

* `storage_bucket` - (Required, String) Specifies the bucket name of the OBS.

* `storage_path` - (Required, String) Specifies the path of OBS object. Comply with OSS Object Definition.
  + When used to indicate input, it needs to be specified to a specific object.
  + When used to indicate output, only the path to the expected storage of the transcoding
    results needs to be specified.

* `call_back_enabled` - (Optional, String) Specifies whether to enable callback notifications.
  The options are as follows:
  + **on**: enable
  + **off**: no enable

* `call_back_url` - (Optional, String) Specifies the notification server address.
  It must be a legal URL and carry the protocol. The protocol of `http` and `https` are supported.
  The live service will push the status information of the snapshot to this address after the snapshot is completed.
  It is required when `call_back_enabled` is set to `on`.

* `call_back_auth_key` - (Optional, String) Specifies the callback authentication key value.
  Consists of 32 to 128 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The live snapshot can be imported using the `domain_name` and `app_name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_live_snapshot.test <domain_name>/<app_name>
```
