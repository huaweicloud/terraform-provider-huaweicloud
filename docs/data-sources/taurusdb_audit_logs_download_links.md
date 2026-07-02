---
subcategory: "TaurusDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_taurusdb_audit_logs_download_links"
description: |-
  Use this data source to get the download links for specific audit logs of a TaurusDB instance within HuaweiCloud.
---

# huaweicloud_taurusdb_audit_logs_download_links

Use this data source to get the download links for specific audit logs of a TaurusDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "audit_log_ids" {
  type = list(string)
}

data "huaweicloud_taurusdb_audit_logs_download_links" "test" {
  instance_id = var.instance_id
  ids         = var.audit_log_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the TaurusDB instance.

* `ids` - (Required, List) Specifies the list of audit log IDs.
  A maximum of 50 audit log IDs are allowed in the list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - Indicates the list of audit log download links. The links are valid for 5 minutes.
