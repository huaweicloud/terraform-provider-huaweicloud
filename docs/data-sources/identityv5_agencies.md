---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_agencies"
description: |-
  Use this data source to get the list of the IAM agencies or trust agencies within HuaweiCloud.
---

# huaweicloud_identityv5_agencies

Use this data source to get the list of the IAM agencies or trust agencies within HuaweiCloud.

## Example Usage

### Query all agencies

```hcl
data "huaweicloud_identityv5_agencies" "test" {}
```

### Query agency by agency ID

```hcl
variable "agency_id" {}

data "huaweicloud_identityv5_agencies" "test" {
  agency_id = var.agency_id
}
```

## Argument Reference

The following arguments are supported:

* `path_prefix` - (Optional, String) Specifies the resourcepath prefix of the agency.  
  It consists of several strings, each string contains one or more letters, digits, special characters (.,+@=_-),
  and ends with slash (/). e.g. `service-linked-agency/service.CBH/`.

* `agency_id` - (Optional, String) Specifies the ID of the agency.  
  This parameter conflicts with `path_prefix`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agencies` - The list of the  agencies.  
  The [agencies](#v5_agencies) structure is documented below.

<a name="v5_agencies"></a>
The `agencies` block supports:

* `agency_id` - The ID of the agency.

* `agency_name` - The name of the agency.

* `trust_policy` - The trust policy of the agency.

* `path` - The path of the agency.

* `trust_domain_id` - The account ID of the trusted domain.  
  Only exists in agencies, not in trust agencies.

* `trust_domain_name` - The account name of the trusted domain.  
  Only exists in agencies, not in trust agencies.

* `urn` - The URN of the agency.

* `created_at` - The creation time of the agency or trust agency.

* `description` - The description of the agency or trust agency.

* `max_session_duration` - The maximum session duration of the agency.
