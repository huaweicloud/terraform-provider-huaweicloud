---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_agencies"
description: ""
---

# huaweicloud_identity_agencies

Use this data source to query the IAM agency list within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
data "huaweicloud_identity_agencies" "all" {}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of agency. The name is a string of 1 to 64 characters.

* `trust_domain_id` - (Optional, String) Specifies the ID of delegated user domain.

## Attribute Reference

* `id` - The data source ID.

* `agencies` - The details of the queried IAM agencies. The structure is documented below.

The `agencies` block contains:

* `id` - The agency ID.

* `name` - The agency name.

* `description` - The supplementary information about the agency.

* `expired_at` - The expiration time of agency.

* `created_at` - The time when the agency was created.

* `trust_domain_id` - The ID of delegated user domain.

* `trust_domain_name` - The name of delegated user domain.

* `duration` - The validity period of an agency.
