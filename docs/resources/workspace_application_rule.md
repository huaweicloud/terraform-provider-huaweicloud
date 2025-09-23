---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_application_rule"
description: |-
  Manages a Workspace application rule resource within HuaweiCloud.
---

# huaweicloud_workspace_application_rule

Manages a Workspace application rule resource within HuaweiCloud.

## Example Usage

### Create a product rule

```hcl
variable "rule_name" {}
variable "rule_description" {}
variable "identify_condition" {}
variable "publisher" {}
variable "product_name" {}
variable "process_name" {}
variable "support_os" {}

resource "huaweicloud_workspace_application_rule" "test" {
  name        = var.rule_name
  description = var.rule_description

  detail {
    scope = "PRODUCT"

    product_rule {
      identify_condition = var.identify_condition
      publisher          = var.publisher
      product_name       = var.product_name
      process_name       = var.process_name
      support_os         = var.support_os
      version            = "1.0"
      product_version    = "2019"
    }
  }
}
```

### Create a Path Rule

```hcl
variable "rule_name" {}
variable "rule_description" {}
variable "install_path" {}

resource "huaweicloud_workspace_application_rule" "test_path" {
  name        = var.rule_name
  description = var.rule_description

  detail {
    scope = "PATH"

    path_rule {
      path = var.install_path
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application rule is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the application rule.  
  The name must be `1` to `64` characters, only letters, digits, and underscores (_) are allowed, and the name
  cannot contain spaces.

* `detail` - (Required, List) Specifies the detail of the application rule.  
  The [detail](#workspace_application_rule_detail) structure is documented below.

* `description` - (Optional, String) Specifies the description of the application rule.

<a name="workspace_application_rule_detail"></a>
The `detail` block supports:

* `scope` - (Required, String) Specifies the scope of the application rule.  
  The valid values are as follows:
  + **PRODUCT**
  + **PATH**

* `product_rule` - (Optional, List) Specifies the detail of the product rule.  
  The [product_rule](#workspace_application_product_rule_config) structure is documented below.

  -> Required if the value of parameter `scope` is **PRODUCT**.

* `path_rule` - (Optional, List) Specifies the detail of the path rule.  
  The [path_rule](#workspace_application_path_rule_config) structure is documented below.

  -> Required if the value of parameter `scope` is **PATH**.

<a name="workspace_application_product_rule_config"></a>
The `product_rule` block supports:

* `identify_condition` - (Required, String) Specifies the identify condition of the product rule.  
  The valid values are as follows:
  + **publisher**
  + **product**
  + **process**

* `publisher` - (Optional, String) Specifies the publisher of the product.  
  Defaults to empty string, also you can configure this value as asterisk (*).

  -> At least one of `publisher`, `product_name` and `process_name` must be provided,  
  and both of them cannot be asterisk (*) or empty.

* `product_name` - (Optional, String) Specifies the name of the product.  
  Defaults to empty string, also you can configure this value as asterisk (*).

  -> At least one of `publisher`, `product_name` and `process_name` must be provided,  
  and both of them cannot be asterisk (*) or empty.

* `process_name` - (Optional, String) Specifies the process name of the product.  
  Defaults to empty string, also you can configure this value as asterisk (*).

  -> At least one of `publisher`, `product_name` and `process_name` must be provided,  
  and both of them cannot be asterisk (*) or empty.

* `support_os` - (Optional, String) Specifies the list of the supported operating system types.  
  Defaults to **Windows**

* `version` - (Optional, String) Specifies the version of the product rule.

* `product_version` - (Optional, String) Specifies the version of the product.

<a name="workspace_application_path_rule_config"></a>
The `path_rule` block supports:

* `path` - (Required, String) Specifies the path where the product is installed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Application rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_application_rule.test <id>
```
