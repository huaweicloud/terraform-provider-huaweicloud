---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_orchestration_rule"
description: |-
  Manages an orchestration rule under the dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_orchestration_rule

Manages an orchestration rule under the dedicated instance within HuaweiCloud.

## Example Usage

### Manages an orchestration rule with list type and no preprocessing

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "list"

  mapped_param = jsonencode({
    "mapped_param_name": "listParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      # First priority
      # All input values ​​in the list are mapped to 'ValueAA'.
      "mapped_param_value": "ValueAA",
      "map_param_list": ["ValueA"]
    }),
    jsonencode({
      # Second priority
      "mapped_param_value": "ValueBB",
      "map_param_list": ["ValueB"]
    })
  ]
}
```

### Manages an preprocessing orchestration rule with list type

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id      = var.instance_id
  name             = var.orchestration_rule_name
  strategy         = "list"
  is_preprocessing = true

  map = [
    jsonencode({
      # First priority
      # All input values ​​in the list are mapped to 'ValueAA'.
      "mapped_param_value": "ValueAA",
      "map_param_list": ["ValueA"]
    }),
    jsonencode({
      # Second priority
      "mapped_param_value": "ValueBB",
      "map_param_list": ["ValueB"]
    })
  ]
}
```

### Manages an orchestration rule with range type and no preprocessing

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "range"

  mapped_param = jsonencode({
    "mapped_param_name": "rangeParam",
    "mapped_param_type": "number",
    "mapped_param_location": "query"
  })

  map = [
    jsonencode({
      # First priority
      # All input values ​​in the range 1000 to 1999 are mapped to 10001.
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      # Second priority
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    }),
    jsonencode({
      # Third priority
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    })
  ]
}
```

### Manages an orchestration rule with hash type and no preprocessing

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "hash"

  mapped_param = jsonencode({
    # The value of the request header is directly mapped to the new request header after hash calculation.
    "mapped_param_name": "hashParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}
```

### Manages an orchestration rule with hash range type and no preprocessing

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "hash_range"

  mapped_param = jsonencode({
    "mapped_param_name": "hashParam",
    "mapped_param_type": "number",
    "mapped_param_location": "header"
  })

  # Use the request parameter to generate a hash value, and then use the hash value to perform range arrangement.
  map = [
    jsonencode({
      # First priority
      # If the hash value is ​​in the range 1000 to 1999, then mapped to 10001.
      "mapped_param_value": "10001",
      "map_param_range": {
        "range_start": "1000",
        "range_end": "1999"
      }
    }),
    jsonencode({
      # Second priority
      "mapped_param_value": "10002",
      "map_param_range": {
        "range_start": "2000",
        "range_end": "2999"
      }
    }),
    jsonencode({
      # Third priority
      "mapped_param_value": "10003",
      "map_param_range": {
        "range_start": "3000",
        "range_end": "3999"
      }
    })
  ]
}
```

### Manages an orchestration rule with none value type

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "none_value"

  mapped_param = jsonencode({
    "mapped_param_name": "noneValueParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      # This orchestration value will be returned when the request parameter is empty.
      "mapped_param_value": "NoneValueReturnedForExample"
    })
  ]
}
```

### Manages an orchestration rule with default type

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "default"

  mapped_param = jsonencode({
    "mapped_param_name": "defaultParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      # When the request parameters exist but no orchestration rule can match them, the orchestration mapping value of
      # the default rule is returned.
      "mapped_param_value": "DefaultValueReturnedForExample"
    })
  ]
}
```

### Manages an orchestration rule with head N type

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "head_n"

  mapped_param = jsonencode({
    "mapped_param_name": "headNParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      # Try to intercept the first N characters of the string as the new value.
      "intercept_length": 5,
      "mapped_param_value": "" # Note that the absence of this empty value configuration will cause script changes.
    })
  ]
}
```

### Manages an orchestration rule with tail N type

```hcl
variable "instance_id" {}
variable "orchestration_rule_name" {}

resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = var.instance_id
  name        = var.orchestration_rule_name
  strategy    = "tail_n"

  mapped_param = jsonencode({
    "mapped_param_name": "tailNParam",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })

  map = [
    jsonencode({
      # Try to intercept the last N characters of the string as the new value.
      "intercept_length": 5,
      "mapped_param_value": "" # Note that the absence of this empty value configuration will cause script changes.
    })
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the orchestration rule is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the orchestration
  rule belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the name of the orchestration rule.  
  The valid length is limited from `3` to `64`, only letters, digits and underscores (_) are allowed.  
  The name must start with a letter and must be unique.

* `strategy` - (Required, String) Specifies the type of the orchestration rule.  
  The valid values are as follows:
  + **list**: Maps the values ​​in the list to new values.
  + **range**: Maps the values ​​in the range to new values.
  + **hash**: The value of the request header is directly mapped to the new request header after hash calculation.
  + **hash_range**: Use the request parameter to generate a hash value, and then use the hash value to perform range
    arrangement.
  + **none_value**: Value returned when the request parameter is empty.
  + **default**: When the request parameters exist but no orchestration rule can match them, the orchestration
    mapping value of the default rule is returned.
  + **head_n**: Try to intercept the first N characters of the string as the new value.
  + **tail_n**: Try to intercept the last N characters of the string as the new value.

* `is_preprocessing` - (Optional, Bool, ForceNew) Specifies whether rule is a preprocessing rule.  
  Defaults to **false**. Changing this will create a new resource.

  -> Type **none_value** and type **default** do not have this configuration.

* `mapped_param` - (Optional, String, ForceNew) Specifies the parameter configuration after orchestration, in JSON
  format.  
  Changing this will create a new resource.  
  For the parameter configuration, please refer to the [document](https://support.huaweicloud.com/intl/en-us/api-apig/CreateOrchestration.html#CreateOrchestration__request_OrchestrationMappedParam).

* `map` - (Optional, List) Specifies the list of orchestration mapping rules, each item should be in JSON format.
  For the parameter configuration of list items, please refer to the [document](https://support.huaweicloud.com/intl/en-us/api-apig/CreateOrchestration.html#CreateOrchestration__request_OrchestrationMap).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the orchestration rule.

* `created_at` - The creation time of the orchestration rule, in RFC3339 format.

* `updated_at` - The latest update time of the orchestration rule, in RFC3339 format.

## Import

Orchestration rules be imported using related `instance_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_apig_orchestration_rule.test <instance_id>/<id>
```
