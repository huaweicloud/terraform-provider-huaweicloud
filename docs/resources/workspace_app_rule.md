---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_rule"
description: |-
  Manages a Workspace APP rule resource within HuaweiCloud.
---

# huaweicloud_workspace_app_rule

Manages a Workspace APP rule resource within HuaweiCloud.

-> When using `product_rule`, the `identify_condition` field is required. 
  For `product_rule`, the fields `publisher`, `product_name`, and `process_name` cannot all be `*` or empty simultaneously.
  You can specify either `product_rule` or `path_rule`, but not both in the same rule configuration.


## Example Usage

### Usage with Product Rule

```hcl
variable "app_rule_name" {}
variable "identify_condition" {}
variable "publisher" {}
variable "product_name" {}
variable "process_name" {}
variable "support_os" {}
variable "app_rule_description" {}

resource "huaweicloud_workspace_app_rule" "test" {
  name = var.app_rule_name
  rule {
    scope = "PRODUCT"
    product_rule {
      identify_condition = var.identify_condition
      publisher         = var.publisher
      product_name      = var.product_name
      process_name      = var.process_name
      support_os        = var.support_os
      version           = "1.0"
      product_version   = "2019"
    }
  }
  description = var.app_rule_description
}
```

### Usage with Path Rule

```hcl
variable "app_rule_name" {}
variable "app_path" {}
variable "app_rule_description" {}

resource "huaweicloud_workspace_app_rule" "test_path" {
  name = var.app_rule_name
  rule {
    scope = "PATH"
    path_rule {
      path = var.app_path
    }
  }
  description = var.app_rule_description
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the application rule is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.zhe  

* `name` - (Required, String) Specifies the name of the application rule.  
  The name must be `1` to `64` characters, only letters, digits, and underscores (_) are allowed, and the name
  cannot contain spaces.

* `rule` - (Required, List) Specifies the config object list of the application rule.  
  The [rule](#workspace_app_rule_config) structure is documented below.

* `description` - (Optional, String) Specifies the description of the application rule.

<a name="workspace_app_rule_config"></a>
The `rule` block supports:

* `scope` - (Required, String) Specifies the scope of the application rule.  
  The valid values are as follows:
  + **PRODUCT**
  + **PATH**

* `product_rule` - (Optional, List) Specifies the detail of the product rule.  
  The [product_rule](#workspace_app_product_rule_config) structure is documented below.

* `path_rule` - (Optional, List) Specifies the detail of the path rule.  
  The [path_rule](#workspace_app_path_rule_config) structure is documented below.

<a name="workspace_app_product_rule_config"></a>
The `product_rule` block supports:

* `identify_condition` - (Required, String) Specifies the identify condition of the product rule.  
  The valid values are as follows:
  + **publisher**
  + **product**
  + **process**

* `publisher` - (Optional, String) Specifies the publisher of the product.  
  Cannot be `*` or empty if `product_name` and `process_name` are also `*` or empty.

* `product_name` - (Optional, String) Specifies the name of the product.  
  Cannot be `*` or empty if `publisher` and `process_name` are also `*` or empty.

* `process_name` - (Optional, String) Specifies the process name of the product.  
  Cannot be `*` or empty if `publisher` and `product_name` are also `*` or empty.

* `support_os` - (Optional, String) Specifies the list of the supported operating system types.  
  Defaults to **Windows**

* `version` - (Optional, String) Specifies the version of the product rule.

* `product_version` - (Optional, String) Specifies the version of the product.

<a name="workspace_app_path_rule_config"></a>
The `path_rule` block supports:

* `path` - (Required, String) Specifies the path where the product is installed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the application rule ID.

## Import

Application rule can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_workspace_app_rule.test <id>
```

