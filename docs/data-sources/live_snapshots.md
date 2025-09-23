---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_snapshots"
description: |-
  Use this data source to get a list of the Live snapshots.
---

# huaweicloud_live_snapshots

Use this data source to get a list of the Live snapshots.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_live_snapshots" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Required, String) Specifies the domain name.

* `app_name` - (Optional, String) Specifies the application name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `snapshots` - The snapshot list.

  The [snapshots](#snapshots_struct) structure is documented below.

<a name="snapshots_struct"></a>
The `snapshots` block supports:

* `call_back_url` - The address of the server for receiving callback notifications.

* `domain_name` - The ingest domain name.

* `app_name` - The application name.

* `call_back_auth_key` - The callback authentication key value.

* `frequency` - The snapshot capturing frequency.

* `storage_mode` - The method for storing snapshots in an OBS bucket.
  + `0`: All. A snapshot file name contains the timestamp. All snapshot files of each stream are stored in OBS.
    Example: snapshot/{domain}/{app_name}/{stream_name}/{UnixTimestamp}.jpg.
  + `1`: Latest. A snapshot file name does not contain the timestamp. Only the latest snapshot file of each stream
    will be saved. A new snapshot file overwrites the previous one. Example: snapshot/{domain}/{app_name}/{stream_name}.jpg

* `call_back_enabled` - Whether to enable callback notification.
  + **on**: Enabled.
  + **off**: Disabled.

* `storage_bucket` - The OBS bucket name.

* `storage_location` - The region where the OBS bucket is located.

* `storage_path` - The OBS object path.
