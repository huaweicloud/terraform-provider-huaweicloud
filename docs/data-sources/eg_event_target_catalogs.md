---
subcategory: "EventGrid (EG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_eg_event_target_catalogs"
description: |-
  Use this data source to get list of EG event target catalogs within HuaweiCloud.
---

# huaweicloud_eg_event_target_catalogs

Use this data source to get list of EG event target catalogs within HuaweiCloud.

## Example Usage

### Query all event target catalogs

```hcl
data "huaweicloud_eg_event_target_catalogs" "test" {}
```

### Query event target catalogs with support types

```hcl
data "huaweicloud_eg_event_target_catalogs" "test" {
  support_types = ["SUBSCRIPTION", "FLOW"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the event target catalogs are located.  
  If omitted, the provider-level region will be used.

* `fuzzy_label` - (Optional, String) Specifies the label of the event target catalog to be queried.  
  Fuzzy search is supported.

* `support_types` - (Optional, List) Specifies the support type list of event targets to be queried.  
  The valid values are as follows:
  + **SUBSCRIPTION**
  + **FLOW**

* `sort` - (Optional, String) Specifies the sort order for querying event target catalogs.  
  The format is `field:order`, where `field` is the field name and `order` is `ASC` or `DESC`.
  e.g. `created_time:ASC`. Only `created_time` and `updated_time` fields are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `catalogs` - All event target catalogs that match the filter parameters.  
  The [catalogs](#eg_event_target_catalogs_attr) structure is documented below.

<a name="eg_event_target_catalogs_attr"></a>
The `catalogs` block supports:

* `id` - The ID of the event target catalog.

* `name` - The name of the event target catalog.

* `label` - The display name of the event target catalog.

* `description` - The description of the event target catalog.

* `support_types` - The support type list of event target catalog.

* `provider_type` - The provider type of the event target catalog.
  + **OFFICIAL**
  + **CUSTOM**

* `parameters` - The parameter list of the event target catalog.  
  The [parameters](#eg_event_target_catalogs_parameters_attr) structure is documented below.

* `created_time` - The creation time of the event target catalog, in UTC format.

* `updated_time` - The latest update time of the event target catalog, in UTC format.

<a name="eg_event_target_catalogs_parameters_attr"></a>
The `parameters` block supports:

* `name` - The name of the target parameter.

* `label` - The display name of the target parameter.

* `metadata` - The metadata of the target parameter, in JSON format.
