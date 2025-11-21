---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_application_providers"
description: |-
  Use this data source to get the Identity Center application providers.
---

# huaweicloud_identitycenter_application_providers

Use this data source to get the Identity Center application providers.

## Example Usage

```hcl
data "huaweicloud_identitycenter_application_providers" "test" {}
```

## Argument Reference

There are no arguments available for this data source.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `application_providers` - The list of the application providers.
  The [application_providers](#application_providers_struct) structure is documented below.

<a name="application_providers_struct"></a>
The `application_providers` block supports:

* `application_provider_urn` - The urn of the application provider.

* `federation_protocol` - The federation protocol of the application provider.

* `application_provider_id` - The ID of the application provider.

* `display_data` - The display data of the application provider.
  The [display_data](#display_data_struct) structure is documented below.

<a name="display_data_struct"></a>
The `display_data` block supports:

* `description` - The description of the application provider.

* `display_name` - The display name of the application provider.

* `icon_url` - The icon url of the application provider.
