---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_function_templates"
description: |-
  Use this data source to query CSMS secret rotation function templates within HuaweiCloud.
---

# huaweicloud_csms_function_templates

Use this data source to query CSMS secret rotation function templates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_csms_function_templates" "test" {
  secret_type     = "GaussDB-FG"
  secret_sub_type = "SingleUser"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `secret_type` - (Required, String) Specifies the secret type. Valid values are **GaussDB-FG** and **RDS-FG**.

* `secret_sub_type` - (Required, String) Specifies the secret rotation account type.
  The options are as follows:
  + **SingleUser**: Single user mode rotation.
  + **MultiUser**: Dual user mode rotation.

* `engine` - (Optional, String) Specifies the database type. This parameter is mandatory when the secret type is
  **RDS-FG**. The options are as follows:
  + **mysql**
  + **postgresql**
  + **sqlserver**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `function_templates` - The secret rotation function templates.
