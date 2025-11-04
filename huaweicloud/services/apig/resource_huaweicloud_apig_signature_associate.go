package apig

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceSignatureAssociate is a provider resource of the API signature.
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/sign-bindings
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/unbinded-apis
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/binded-apis
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/{sign_bindings_id}
func ResourceSignatureAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSignatureAssociateCreate,
		ReadContext:   resourceSignatureAssociateRead,
		UpdateContext: resourceSignatureAssociateUpdate,
		DeleteContext: resourceSignatureAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSignatureAssociateImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the signature and the APIs are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the APIs and the signature belong.",
			},
			"signature_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The signature ID for APIs binding.",
			},
			"publish_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The publish IDs corresponding to the APIs bound by the signature.",
			},
		},
	}
}

func signatureBindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, signId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/sign-bindings/unbinded-apis"
			queryUrl = "?sign_id={sign_id}"
			offset   = 0
			result   = make([]interface{}, 0)
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
		queryUrl = strings.ReplaceAll(queryUrl, "{sign_id}", signId)
		listPath += queryUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		for {
			listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
			requestResp, err := client.Request("GET", listPathWithOffset, &opt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving unassociated signatures: %s", err)
			}
			respBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return nil, "ERROR", err
			}
			unbindPublishIds := utils.PathSearch("apis[*].publish_id", respBody, make([]interface{}, 0)).([]interface{})
			if len(unbindPublishIds) < 1 {
				break
			}
			result = append(result, unbindPublishIds...)
			offset += len(unbindPublishIds)
		}

		if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(result), publishIds, false, true) {
			return result, "PENDING", nil
		}
		return result, "COMPLETED", nil
	}
}

func bindSignatureToApis(ctx context.Context, client *golangsdk.ServiceClient, opts signs.BindOpts,
	timeout time.Duration) error {
	var (
		instanceId = opts.InstanceId
		signId     = opts.SignatureId
		publishIds = opts.PublishIds
	)

	_, err := signs.Bind(client, opts)
	if err != nil {
		return fmt.Errorf("error binding signature to the APIs: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: signatureBindingRefreshFunc(client, instanceId, signId, publishIds),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the binding completed: %s", err)
	}
	return nil
}

func resourceSignatureAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}
	var (
		instanceId = d.Get("instance_id").(string)
		signId     = d.Get("signature_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)

		opts = signs.BindOpts{
			InstanceId:  instanceId,
			SignatureId: signId,
			PublishIds:  utils.ExpandToStringListBySet(publishIds),
		}
	)
	err = bindSignatureToApis(ctx, client, opts, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceId, signId))

	return resourceSignatureAssociateRead(ctx, d, meta)
}

func buildSignBindApiListOpts(instanceId, signId string) signs.ListBindOpts {
	return signs.ListBindOpts{
		InstanceId:  instanceId,
		SignatureId: signId,
		Limit:       500,
	}
}

func flattenApiPublishIdsForSignature(apiList []signs.SignBindApiInfo) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishId
	}
	return result
}

func resourceSignatureAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		signId     = d.Get("signature_id").(string)
		opts       = buildSignBindApiListOpts(instanceId, signId)
	)

	resp, err := signs.ListBind(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Signature association")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	return diag.FromErr(d.Set("publish_ids", flattenApiPublishIdsForSignature(resp)))
}

func signatureUnbindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, signId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := buildSignBindApiListOpts(instanceId, signId)
		resp, err := signs.ListBind(client, opts)
		if err != nil {
			// The API returns a 404 error, which means that the instance or signature has been deleted.
			// In this case, there's no need to disassociate API, also this action has been completed.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "instance_or_signature_not_exist", "COMPLETED", nil
			}
			return resp, "", err
		}
		bindPublishIds := flattenApiPublishIdsForSignature(resp)
		if utils.IsSliceContainsAnyAnotherSliceElement(bindPublishIds, publishIds, false, true) {
			return resp, "PENDING", nil
		}
		return resp, "COMPLETED", nil
	}
}

func unbindSignatureFromApis(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	rmList []string, timeout time.Duration) error {
	var (
		instanceId = d.Get("instance_id").(string)
		signId     = d.Get("signature_id").(string)
	)
	opts := buildSignBindApiListOpts(instanceId, signId)
	resp, err := signs.ListBind(client, opts)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return err
		}
		return fmt.Errorf("error getting binding APIs based on signature (%s): %s", signId, err)
	}

	if len(resp) < 1 {
		log.Printf("[DEBUG] All APIs has been disassociated from the signature (%s) under dedicated instance (%s)", signId, instanceId)
		return nil
	}

	for _, val := range resp {
		if utils.StrSliceContains(rmList, val.PublishId) {
			err = signs.Unbind(client, instanceId, val.BindId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					log.Printf("[DEBUG] All APIs has been disassociated from the signature (%s)", signId)
					continue
				}
				return fmt.Errorf("an error occurred during unbind signature: %s", err)
			}
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: signatureUnbindingRefreshFunc(client, instanceId, signId, rmList),
		Timeout: timeout,
		// In most cases, the unbind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the unbind operation completed: %s", err)
	}
	return nil
}

func resourceSignatureAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId     = d.Get("instance_id").(string)
		signId         = d.Get("signature_id").(string)
		oldRaw, newRaw = d.GetChange("publish_ids")

		addSet = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmSet  = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	)

	if rmSet.Len() > 0 {
		err = unbindSignatureFromApis(ctx, client, d, utils.ExpandToStringListBySet(rmSet),
			d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		opts := signs.BindOpts{
			InstanceId:  instanceId,
			SignatureId: signId,
			PublishIds:  utils.ExpandToStringListBySet(addSet),
		}
		// If the target (published) API already has a signature, this update will replace the signature.
		err = bindSignatureToApis(ctx, client, opts, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			diag.FromErr(err)
		}
	}

	return resourceSignatureAssociateRead(ctx, d, meta)
}

func resourceSignatureAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	err = unbindSignatureFromApis(ctx, client, d, utils.ExpandToStringListBySet(d.Get("publish_ids").(*schema.Set)),
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error unbinding APIs from signature")
	}
	return nil
}

func resourceSignatureAssociateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<signature_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("signature_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
