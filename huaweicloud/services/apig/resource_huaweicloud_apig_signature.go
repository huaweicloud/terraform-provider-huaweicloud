package apig

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/signs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// ResourceSignature is a provider resource of the APIG signature.
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/signs/{sign_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/signs/{sign_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/signs
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/signs
func ResourceSignature() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSignatureCreate,
		ReadContext:   resourceSignatureRead,
		UpdateContext: resourceSignatureUpdate,
		DeleteContext: resourceSignatureDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSignatureImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the signature is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the signature belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The signature name.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The signature type.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The signature key.",
			},
			"secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				Description: "The signature secret.",
			},
			"algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The signature algorithm.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the signature.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the signature.",
			},
		},
	}
}

func resourceSignatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	opts := signs.CreateOpts{
		InstanceId:    d.Get("instance_id").(string),
		Name:          d.Get("name").(string),
		SignType:      d.Get("type").(string),
		SignKey:       d.Get("key").(string),
		SignSecret:    d.Get("secret").(string),
		SignAlgorithm: d.Get("algorithm").(string),
	}
	resp, err := signs.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating the signature: %s", err)
	}
	d.SetId(resp.ID)

	return resourceSignatureRead(ctx, d, meta)
}

func resourceSignatureRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		signatureId = d.Id()
		opts        = signs.ListOpts{
			InstanceId: d.Get("instance_id").(string),
			ID:         signatureId,
		}
	)

	resp, err := signs.List(client, opts)
	log.Printf("[DEBUG] The signature result is: %#v", resp)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Signature")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Signature")
	}

	signature := resp[0]
	mErr := multierror.Append(nil,
		d.Set("name", signature.Name),
		d.Set("type", signature.SignType),
		d.Set("key", signature.SignKey),
		d.Set("secret", signature.SignSecret),
		d.Set("algorithm", signature.SignAlgorithm),
		d.Set("created_at", signature.CreatedAt),
		d.Set("updated_at", signature.UpdatedAt),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving signature (%s) fields: %s", signatureId, err)
	}
	return nil
}

func resourceSignatureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		signatureId = d.Id()
		opts        = signs.UpdateOpts{
			InstanceId:    d.Get("instance_id").(string),
			SignatureId:   signatureId,
			Name:          d.Get("name").(string),
			SignType:      d.Get("type").(string),
			SignKey:       d.Get("key").(string),
			SignSecret:    d.Get("secret").(string),
			SignAlgorithm: d.Get("algorithm").(string),
		}
	)
	_, err = signs.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating the signature (%s): %s", signatureId, err)
	}
	return resourceSignatureRead(ctx, d, meta)
}

func resourceSignatureDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId  = d.Get("instance_id").(string)
		signatureId = d.Id()
	)
	err = signs.Delete(client, instanceId, signatureId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting the signature (%s)", signatureId))
	}
	return nil
}

func resourceSignatureImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])
	err := d.Set("instance_id", parts[0])
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving instance ID field: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
