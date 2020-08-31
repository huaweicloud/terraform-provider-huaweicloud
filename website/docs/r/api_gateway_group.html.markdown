---
subcategory: "API Gateway (APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_api_gateway_group"
sidebar_current: "docs-huaweicloud-resource-api-gateway-group"
description: |-
  Provides an API gateway group resource.
---

# huaweicloud\_api\_gateway\_group

Provides an API gateway group resource.

## Example Usage

```hcl
resource "huaweicloud_api_gateway_group" "apigw_group" {
    name = "apigw_group"
    description = "your descpiption"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the API group. An API group name consists of 3–64 characters,
    starting with a letter. Only letters, digits, and underscores (_) are allowed.

* `description` - (Optional) Specifies the description of the API group.
    The description cannot exceed 255 characters.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the API group.
* `status` - Status of the API group.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
