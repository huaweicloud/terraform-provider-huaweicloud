---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_audit_log_download_links"
description: |-
  Use this data source to get the temporary link for downloading full SQL.
---

# huaweicloud_gaussdb_mysql_audit_log_download_links

Use this data source to get the temporary link for downloading full SQL.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_gaussdb_mysql_audit_log_download_links" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB MySQL instance.

* `start_time` - (Required, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Required, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `node_id` - (Optional, String) Specifies the ID of the GaussDB MySQL instance node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - Indicates the list of the full SQL file information.

  The [links](#links_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `name` - Indicates the name of the file.

* `full_name` - Indicates the full name of the file.

* `size` - Indicates the file size, in KB.

* `updated_time` - Indicates the last modification time of the SQL file.

* `download_link` - Indicates the link for downloading the file.

* `link_expired_time` - Indicates the link expiration time.
