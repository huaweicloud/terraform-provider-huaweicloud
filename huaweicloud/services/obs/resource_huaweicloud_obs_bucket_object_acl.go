package obs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/obs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OBS PUT /{ObjectName}?acl
// @API OBS GET /{ObjectName}?acl
func ResourceOBSBucketObjectAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOBSBucketObjectAclCreate,
		ReadContext:   resourceOBSBucketObjectAclRead,
		UpdateContext: resourceOBSBucketObjectAclCreate,
		DeleteContext: resourceOBSBucketObjectAclDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceOBSBucketObjectAclImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the bucket which the object belongs to.`,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the object to which to set the acl.`,
			},
			"public_permission": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        objectPublicPermissionSchema(),
				Description: `Specifies the object public permission.`,
			},
			"account_permission": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        objectAccountPermissionSchema(),
				Description: `Specifies the object account permissions.`,
			},
			"owner_permission": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     objectOwnerPermissionSchema(),
			},
		},
	}
}

func objectPublicPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_to_object": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to object. Only **READ** supported.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func objectAccountPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the account id to authorize. The account id cannot be the object owner, 
and must be unique.`,
			},
			"access_to_object": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to object. Only **READ** supported.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ",
					}, false),
				},
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the access to acl. Valid values are **READ_ACP** and **WRITE_ACP**.`,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"READ_ACP", "WRITE_ACP",
					}, false),
				},
			},
		},
	}
	return &sc
}

func objectOwnerPermissionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_to_object": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The access to bucket.`,
			},
			"access_to_acl": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The access to acl.`,
			},
		},
	}
	return &sc
}

func buildObjectOwnerPermissionGrants(obsClient *obs.ObsClient, d *schema.ResourceData,
	domainID string) ([]obs.Grant, error) {
	params := &obs.GetObjectAclInput{
		Bucket: d.Get("bucket").(string),
		Key:    d.Get("key").(string),
	}
	output, err := obsClient.GetObjectAcl(params)
	if err != nil {
		return nil, err
	}
	var ownerGrants []obs.Grant
	for _, grant := range output.Grants {
		grantee := grant.Grantee
		if grantee.Type == obs.GranteeUser && grantee.ID == domainID {
			// owner grants
			ownerGrants = append(ownerGrants, grant)
		}
	}
	if len(ownerGrants) > 0 {
		return ownerGrants, nil
	}

	// add default permissions: READ、READ_ACP、WRITE_ACP
	accesses := []string{"READ", "READ_ACP", "WRITE_ACP"}
	return buildUserTypeGrants(accesses, domainID), nil
}

func buildObjectAccessesFromRawMap(rawMap map[string]interface{}) []string {
	var accesses []string
	if accessArray, ok := rawMap["access_to_object"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}
	if accessArray, ok := rawMap["access_to_acl"].([]interface{}); ok {
		accesses = append(accesses, utils.ExpandToStringList(accessArray)...)
	}
	return accesses
}

func buildObsBucketObjectAclGrants(obsClient *obs.ObsClient, d *schema.ResourceData,
	domainID string) ([]obs.Grant, error) {
	var grants []obs.Grant
	ownerPermissions, err := buildObjectOwnerPermissionGrants(obsClient, d, domainID)
	if err != nil {
		return nil, err
	}
	grants = append(grants, ownerPermissions...)

	permissions := d.Get("account_permission").(*schema.Set)
	for _, raw := range permissions.List() {
		if rawMap, rawOk := raw.(map[string]interface{}); rawOk {
			accountID := rawMap["account_id"].(string)
			if accountID == domainID {
				return nil, fmt.Errorf("the account id cannot be the object owner: %s", accountID)
			}
			accesses := buildObjectAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The account permission accesses: %v.", accesses)
			grants = append(grants, buildUserTypeGrants(accesses, accountID)...)
		}
	}

	if rawArray, ok := d.Get("public_permission").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			accesses := buildObjectAccessesFromRawMap(rawMap)
			log.Printf("[DEBUG] The public permission accesses: %v.", accesses)
			grants = append(grants, buildGroupTypeGrants(accesses, obs.GroupAllUsers)...)
		}
	}
	return grants, nil
}

func resourceOBSBucketObjectAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	grantParam, err := buildObsBucketObjectAclGrants(obsClient, d, cfg.DomainID)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	params := &obs.SetObjectAclInput{
		Bucket: bucket,
		Key:    key,
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = grantParam
	_, err = obsClient.SetObjectAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error creating OBS bucket object acl", bucket, err))
	}
	d.SetId(key)
	return resourceOBSBucketObjectAclRead(ctx, d, meta)
}

func flattenObjectAccessesFromGrant(grant obs.Grant) (objectAccesses []string, aclAccesses []string) {
	switch grant.Permission {
	case obs.PermissionRead:
		objectAccesses = []string{"READ"}
	case obs.PermissionReadAcp:
		aclAccesses = []string{"READ_ACP"}
	case obs.PermissionWriteAcp:
		aclAccesses = []string{"WRITE_ACP"}
	case obs.PermissionFullControl:
		// permission type `PermissionFullControl` means domain id has all permission(READ、READ_ACP、WRITE_ACP)
		objectAccesses = []string{"READ"}
		aclAccesses = []string{"READ_ACP", "WRITE_ACP"}
	default:
		log.Printf("[WARN] The grant permission: %s not support.", grant.Permission)
	}
	return
}

func flattenObjectPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	var accessToObject []string
	var accessToAcl []string
	for _, grant := range grants {
		objectAccesses, aclAccesses := flattenObjectAccessesFromGrant(grant)
		accessToObject = append(accessToObject, objectAccesses...)
		accessToAcl = append(accessToAcl, aclAccesses...)
	}
	if len(accessToObject) == 0 && len(accessToAcl) == 0 {
		return nil
	}
	permissionMap := map[string]interface{}{
		"access_to_object": accessToObject,
		"access_to_acl":    accessToAcl,
	}
	return []map[string]interface{}{permissionMap}
}

// flattenObjectAccountPermission return grants group by accountID.
func flattenObjectAccountPermission(grants []obs.Grant) []map[string]interface{} {
	if len(grants) == 0 {
		return nil
	}
	// accountIDSet stores accountID set
	accountIDSet := make(map[string]bool)
	// accessToObjectMap stores object accesses, the key is accountID
	accessToObjectMap := make(map[string][]string)
	// accessToAclMap stores acl accesses, the key is accountID
	accessToAclMap := make(map[string][]string)
	for _, grant := range grants {
		accountID := grant.Grantee.ID
		// append new accesses
		objectAccesses, aclAccesses := flattenObjectAccessesFromGrant(grant)
		accessToObjectMap[accountID] = append(accessToObjectMap[accountID], objectAccesses...)
		accessToAclMap[accountID] = append(accessToAclMap[accountID], aclAccesses...)
		accountIDSet[accountID] = true
	}

	var m []map[string]interface{}
	for id := range accountIDSet {
		m = append(m, map[string]interface{}{
			"access_to_object": accessToObjectMap[id],
			"access_to_acl":    accessToAclMap[id],
			"account_id":       id,
		})
	}
	return m
}

func resourceOBSBucketObjectAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	params := &obs.GetObjectAclInput{
		Bucket: d.Get("bucket").(string),
		Key:    d.Id(),
	}
	output, err := obsClient.GetObjectAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error retrieving OBS bucket object acl", d.Id(), err))
	}

	permissionTypeMap := flattenGrantsByPermissionType(output.Grants, cfg)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("public_permission", flattenObjectPermission(permissionTypeMap[GrantPublic])),
		d.Set("account_permission", flattenObjectAccountPermission(permissionTypeMap[GrantAccount])),
		d.Set("owner_permission", flattenObjectPermission(permissionTypeMap[GrantOwner])),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OBS bucket object acl fields: %s", err)
	}
	return nil
}

func resourceOBSBucketObjectAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	obsClient, err := cfg.ObjectStorageClient(region)
	if err != nil {
		return diag.Errorf("error creating OBS Client: %s", err)
	}

	// the operation to delete acl is to set the owner permission: READ、READ_ACP、WRITE_ACP
	ownerPermissions, err := buildObjectOwnerPermissionGrants(obsClient, d, cfg.DomainID)
	if err != nil {
		return diag.FromErr(err)
	}

	bucket := d.Get("bucket").(string)
	params := &obs.SetObjectAclInput{
		Bucket: bucket,
		Key:    d.Id(),
	}
	params.Owner.ID = cfg.DomainID
	params.Grants = ownerPermissions
	_, err = obsClient.SetObjectAcl(params)
	if err != nil {
		return diag.FromErr(getObsError("Error deleting OBS bucket object acl", d.Id(), err))
	}
	return nil
}

func resourceOBSBucketObjectAclImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for import id, must be <bucket>/<key>")
		return nil, err
	}

	bucket := parts[0]
	key := parts[1]
	d.SetId(key)
	mErr := multierror.Append(nil,
		d.Set("bucket", bucket),
		d.Set("key", key),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import obs bucket object acl, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
