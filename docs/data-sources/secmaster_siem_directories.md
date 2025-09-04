---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_siem_directories"
description: |-
  Use this data source to get the list of directories.
---

# huaweicloud_secmaster_siem_directories

Use this data source to get the list of directories.

## Example Usage

```hcl
variable "workspace_id" {}
variable "category" {}

data "huaweicloud_secmaster_siem_directories" "test" {
  workspace_id = var.workspace_id
  category     = var.category
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `category` - (Required, String) Specifies the directory type.
  The valid values are as follows:
  + **TABLE**
  + **PIPE**
  + **RETRIEVE_SCRIPT**
  + **ANALYSIS_SCRIPT**
  + **DATA_TRANSFORMATION**
  + **ALERT_RULE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workspaceid` - The workspace ID.

* `project_id` - The project ID.

* `directories` - The directory list.

* `directory_i18ns` - The directory I18N list.

  The [directory_i18ns](#directory_i18ns_struct) structure is documented below.

<a name="directory_i18ns_struct"></a>
The `directory_i18ns` block supports:

* `directory` - The directory grouping.

* `directory_en` - The en directory grouping.

* `directory_fr` - The fr directory grouping.
