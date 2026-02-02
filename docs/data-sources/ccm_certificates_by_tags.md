---
subcategory: "Cloud Certificate Manager (CCM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ccm_certificates_by_tags"
description: |-
  Use this data source to get the list of CCM SSL certificates by tags.
---

# huaweicloud_ccm_certificates_by_tags

Use this data source to get the list of CCM SSL certificates by tags.

## Example Usage

### Filter the certificate list

```hcl
data "huaweicloud_ccm_certificates_by_tags" "test" {
  resource_instances = "resource_instances"
  action             = "filter"
}
```

### Query the certificate total count

```hcl
data "huaweicloud_ccm_certificates_by_tags" "test" {
  resource_instances = "resource_instances"
  action             = "count"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `resource_instances` - (Required, String) Specifies the resource instances. Only support **resource_instances**.

* `action` - (Required, String) Specifies the action value. Valid values are **filter** and **count**.

* `without_any_tag` - (Optional, Bool) Specifies whether to query the certificates without tags.
  When the value of this field is **true**, all resources without tags are queried, and the fields `tags`, `tags_any`,
  `not_tags`, and `not_tags_any` are ignored.

* `tags` - (Optional, List) Specifies the containing tags.
  It can contain a maximum of `20` keys, with a maximum of `20` values ​​under each key. The value array for each key can
  be empty, but the structure must be complete. Keys cannot be duplicated, and values ​​for the same key cannot be
  duplicated. The result returns a list of resources containing all tags. Keys are related by AND, and values ​​in the
  key-value structure are related by OR. Without tag filtering, the full dataset is returned.

  The [tags](#tags_Tag) structure is documented below.

* `tags_any` - (Optional, List) Specifies any included tags.
  It can contain a maximum of `20` keys, with a maximum of `20` values ​​under each key. The value array for each key can
  be empty, but the structure must be complete. Keys cannot be duplicated, and values ​​for the same key cannot be
  duplicated. The result returns a list of resources containing tags. Keys are related by OR, and values ​​in the
  key-value structure are also related by OR. Without filtering, the full dataset is returned.

  The [tags_any](#tags_Tag) structure is documented below.

* `not_tags` - (Optional, List) Specifies not included tags.
  It can contain a maximum of `20` keys, with a maximum of `20` values ​​under each key. The value array for each key can
  be empty, but the structure must be complete. Keys cannot be duplicated, and values ​​for the same key cannot be
  duplicated. The result returns a list of resources excluding tags. Keys are related by AND, and values ​​in the
  key-value structure are related by OR. Without filtering, the full dataset is returned.

  The [not_tags](#tags_Tag) structure is documented below.

* `not_tags_any` - (Optional, List) Specifies any tags that are not included.
  It can contain a maximum of `20` keys, with a maximum of `20` values ​​under each key. The value array for each key can
  be empty, but the structure must be complete. Keys cannot be duplicated, and values ​​for the same key cannot be
  duplicated. The result returns a list of resources excluding tags. Keys are related by OR, and values ​​in the
  key-value structure are also related by OR. Without filtering, the full dataset is returned.

  The [not_tags_any](#tags_Tag) structure is documented below.

* `matches` - (Optional, List) Specifies the matches condition.
  The key is the field to be matched, such as **resource_name**. The value is the match value.
  The key is a fixed dictionary value and cannot contain duplicate keys or unsupported keys.

  The [matches](#matches_Match) structure is documented below.

<a name="tags_Tag"></a>
The `tags`, `tags_any`, `not_tags`, `not_tags_any` block supports:

* `key` - (Optional, String) Specifies the key.

* `values` - (Optional, List) Specifies the value array.

<a name="matches_Match"></a>
The `matches` block supports:

* `key` - (Optional, String) Specifies the field to be matched.

* `value` - (Optional, String) Specifies the match value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_count` - The total number of resources.

* `resources` - The resource list.
  The [resources](#resources_response_struct) structure is documented below.

<a name="resources_response_struct"></a>
The `resources` block supports:

* `resource_id` - The resource ID of the certificate.

* `tags` - The tags of the certificate.
  The [tags](#tags_response_struct) structure is documented below.

* `resource_name` - The resource name of the certificate.

* `resource_detail` - The resource detail of the certificate.
  The [resource_detail](#resource_detail_struct) structure is documented below.

<a name="tags_response_struct"></a>
The `tags` block supports:

* `key` - The key of the tag.

* `value` - The value of the tag.

<a name="resource_detail_struct"></a>
The `resource_detail` block supports:

* `cert_id` - The ID of the certificate.

* `cert_name` - The name of the certificate.

* `domain` - The domain of the certificate.

* `cert_type` - The type of the certificate. Valid values are **OV_SSL_CERT** and **EV_SSL_CERT**.

* `cert_brand` - The brand of the certificate.

* `domain_type` - The type of the domain. Valid values are **SINGLE_DOMAIN**, **MULTI_DOMAIN**, and **WILDCARD**.

* `purchase_period` - The purchase period of the certificate, the unit is year.

* `expired_time` - The expiration time of the certificate, the unit is milliseconds.

* `order_status` - The order status of the certificate.

* `domain_num` - The number of domains.

* `wildcard_number` - The number of wildcard domains.

* `sans` - The SANs of the certificate.

* `cert_des` - The description of the certificate.

* `signature_algorithm` - The signature algorithm.

* `fail_reason` - The fail reason.

* `partner_order_id` - The order serial number.

* `push_support` - Whether the certificate is supported to push.

* `cert_issued_time` - The certificate issue time, unit is milliseconds.

* `resource_id` - The resource ID of the certificate.

* `unsubscribe_support` - Whether the certificate is supported to unsubscribe.

* `enterprise_project_id` - The enterprise project ID.

* `origin_cert_id` - The origin certificate ID.

* `renewal_cert_id` - The renewal certificate ID.

* `auto_renew_status` - The auto-renewal status.

* `remain_cert_number` - The remaining number of certificates.

* `auto_deploy_support` - Whether the certificate is supported to auto-deploy.
