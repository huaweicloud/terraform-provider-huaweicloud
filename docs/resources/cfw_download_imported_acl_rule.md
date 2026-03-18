---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_download_imported_acl_rule"
description: |-
  Manages a resource to download CFW imported ACL rule within HuaweiCloud.
---

# huaweicloud_cfw_download_imported_acl_rule

Manages a resource to download CFW imported ACL rule within HuaweiCloud.

-> 1. This resource is a one-time action resource used to download imported ACL rule. Deleting this resource will
  not clear the corresponding request record, but will only remove the resource information from the tf state file.
  <br/>2. Running this resource will generate a file with the suffix **.xlsx** in the current working directory.

## Example Usage

```hcl
variable "object_id" {}

resource "huaweicloud_cfw_download_imported_acl_rule" "test" {
  object_id        = var.object_id
  export_file_name = "temp-download-file.xlsx"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID. This ID is used to distinguish
  between Internet boundary protection and VPC boundary protection after the cloud firewall is created.
  You can get this value from data source `huaweicloud_cfw_firewalls`.

* `export_file_name` - (Required, String, NonUpdatable) Specifies the directory file using to store the downloaded ACL
  rule, and requires the file ending in `.xlsx`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `object_id`.
