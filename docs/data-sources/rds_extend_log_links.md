---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_extend_log_links"
description: |-
  Use this data source to get the list of RDS extend log links.
---

# huaweicloud_rds_extend_log_links

Use this data source to get the list of RDS extend log links.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_extend_log_links" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `file_name` - (Required, String) Specifies the name of the file to be downloaded.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `links` - Indicates the list of extend log links.

  The [links](#links_struct) structure is documented below.

<a name="links_struct"></a>
The `links` block supports:

* `file_name` - Indicates the name of the file.

* `status` - Indicates the status of the link. The value can be one of the following:
  + **SUCCESS**: The download link has been generated.
  + **EXPORTING**: The file is being generated.
  + **FAILED**: The log file fails to be prepared.

* `file_size` - Indicates the file size in KB.

* `file_link` - Indicates the download link.

* `created_at` - Indicates the creation time.

* `updated_at` - Indicates the last update time.
