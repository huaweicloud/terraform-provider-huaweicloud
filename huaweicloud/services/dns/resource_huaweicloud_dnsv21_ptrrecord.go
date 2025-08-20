package dns

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ptrRecordNonUpdatableParams = []string{
	"publicip_id", "enterprise_project_id",
}

// @API DNS POST /v2.1/ptrs
// @API DNS GET /v2.1/ptrs/{ptr_id}
// @API DNS PUT /v2.1/ptrs/{ptr_id}
// @API DNS DELETE /v2.1/ptrs/{ptr_id}
// @API DNS GET /v2/{project_id}/DNS-ptr_record/{resource_id}/tags
// @API DNS POST /v2/{project_id}/DNS-ptr_record/{resource_id}/tags/action
func ResourceDNSV21PtrRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSV21PtrRecordCreate,
		ReadContext:   resourceDNSV21PtrRecordRead,
		UpdateContext: resourceDNSV21PtrRecordUpdate,
		DeleteContext: resourceDNSV21PtrRecordDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(ptrRecordNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"names": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: func(_, _, _ string, d *schema.ResourceData) bool {
					oldNames, newNames := d.GetChange("names")
					o := utils.ExpandToStringList(oldNames.(*schema.Set).List())
					n := utils.ExpandToStringList(newNames.(*schema.Set).List())
					for i, ov := range o {
						o[i] = strings.TrimSuffix(ov, ".")
					}
					for i, nv := range n {
						n[i] = strings.TrimSuffix(nv, ".")
					}

					sort.Strings(o)
					sort.Strings(n)

					return reflect.DeepEqual(o, n)
				},
				Description: `Specifies the domain names of the PTR record.`,
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the EIP.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the PTR record.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the time to live (TTL) of the record set (in seconds).`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID of the PTR record.`,
			},
			"tags": common.TagsSchema(`Specifies the key/value pairs to associate with the PTR record.`),
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The address of the EIP.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the PTR record.`,
			},
		},
	}
}

func resourceDNSV21PtrRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2.1/ptrs"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDNSV21PtrRecordBodyParams(cfg, d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DNS PTR record: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the DNS PTR record ID from the API response")
	}

	d.SetId(id)

	err = waitForDNSV21PtrRecordToBeActive(ctx, client, id, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for PTR record (%s) to be active: %s", id, err)
	}

	return resourceDNSV21PtrRecordRead(ctx, d, meta)
}

func buildCreateDNSV21PtrRecordBodyParams(cfg *config.Config, d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ptrdnames": d.Get("names").(*schema.Set).List(),
		"publicip": map[string]interface{}{
			"region": cfg.GetRegion(d),
			"id":     d.Get("publicip_id"),
		},
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"ttl":                   utils.ValueIgnoreEmpty(d.Get("ttl")),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}

	return bodyParams
}

func waitForDNSV21PtrRecordToBeActive(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			ptrRecord, err := GetDNSV21PtrRecord(client, id)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", ptrRecord, "").(string)
			if status == "ACTIVE" {
				return ptrRecord, "COMPLETED", nil
			}

			return ptrRecord, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        1 * time.Second,
		PollInterval: 3 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceDNSV21PtrRecordRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	getRespBody, err := GetDNSV21PtrRecord(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DNS PTR record")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("names", utils.PathSearch("ptrdnames", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("publicip_id", utils.PathSearch("publicip.id", getRespBody, nil)),
		d.Set("ttl", utils.PathSearch("ttl", getRespBody, nil)),
		d.Set("address", utils.PathSearch("publicip.address", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
	)

	if err := utils.SetResourceTagsToState(d, client, "DNS-ptr_record", d.Id()); err != nil {
		mErr = multierror.Append(mErr, err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDNSV21PtrRecord(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v2.1/ptrs/{ptr_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{ptr_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceDNSV21PtrRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if d.HasChanges("names", "description", "ttl") {
		httpUrl := "v2.1/ptrs/{ptr_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{ptr_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateDNSV21PtrRecordBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		err = waitForDNSV21PtrRecordToBeActive(ctx, client, d.Id(), d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for PTR record (%s) to be active: %s", d.Id(), err)
		}
	}

	if d.HasChanges("tags") {
		tagErr := utils.UpdateResourceTags(client, d, "DNS-ptr_record", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of DNS PTR record %s: %s", d.Id(), tagErr)
		}
	}

	return resourceDNSV21PtrRecordRead(ctx, d, meta)
}

func buildUpdateDNSV21PtrRecordBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ptrdnames":   d.Get("names").(*schema.Set).List(),
		"description": d.Get("description"),
		"ttl":         d.Get("ttl"),
	}

	return bodyParams
}

func resourceDNSV21PtrRecordDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dns_region", region)
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	httpUrl := "v2.1/ptrs/{ptr_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{ptr_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DNS PTR record")
	}

	return nil
}
