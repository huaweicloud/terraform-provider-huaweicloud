---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_shared_connections"
description: |-
  Use this data source to query DAS shared connections within HuaweiCloud.
---

# huaweicloud_das_shared_connections

Use this data source to query DAS shared connections within HuaweiCloud.

## Example Usage

### Query all shared connections under a specified connection

```hcl
variable "connection_id" {}

data "huaweicloud_das_shared_connections" "test" {
  connection_id = var.connection_id
}
```

### Query shared connections by keyword

```hcl
# Fuzzy search for shared connections where the target fields (e.g., `user_id`, `user_name`) contain this keyword.
variable "connection_id" {}
variable "fuzzy_search_keyword" {}

data "huaweicloud_das_shared_connections" "by_keyword" {
  connection_id = var.connection_id
  keywords      = var.fuzzy_search_keyword 
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the shared connections are located.  
  If omitted, the provider-level region will be used.

* `connection_id` - (Required, String) Specifies the ID of the connection to which the shared connection belongs.

* `keywords` - (Optional, String) Specifies keywords to search for shared connections.  
  The fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, in UUID format.

* `shared_connections` - The list of shared connections that matched filter parameters.  
  The [shared_connections](#das_shared_connections) structure is documented below.

<a name="das_shared_connections"></a>
The `shared_connections` block supports:

* `user_id` - The user ID of the shared connection.

* `user_name` - The user name of the shared connection.

* `expired_at` - The expiration time of the shared connection, in RFC3339 format.

* `shared_at` - The creation time of the shared connection, in RFC3339 format.
