package dli

import (
	"context"
	"fmt"
	"regexp"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/databases"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceDliSqlDatabaseV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDliSqlDatabaseV1Create,
		ReadContext:   ResourceDliSqlDatabaseV1Read,
		UpdateContext: ResourceDliSqlDatabaseV1Update,
		DeleteContext: ResourceDliSqlDatabaseV1Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabaseImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z0-9][\w_]{0,127}$`),
						"The name consists of 1 to 128 characters, starting with a letter or digit. "+
							"Only letters, digits and underscores (_) are allowed."),
					validation.StringMatch(regexp.MustCompile(`[A-Za-z_]`), "The name cannot be all digits."),
				),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func ResourceDliSqlDatabaseV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	opts := databases.CreateOpts{
		Name:                dbName,
		Description:         d.Get("description").(string),
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, config),
	}
	_, err = databases.Create(c, opts)
	if err != nil {
		return diag.Errorf("error creating DLI database, %s", err)
	}
	// The resource ID (database name) at this time is only used as a mark the resource, and the value will be refreshed
	// in the READ method.
	d.SetId(dbName)

	return ResourceDliSqlDatabaseV1Read(ctx, d, meta)
}

func GetDliSqlDatabaseByName(c *golangsdk.ServiceClient, dbName string) (databases.Database, error) {
	resp, err := databases.List(c, databases.ListOpts{
		Keyword: dbName, // Fuzzy matching.
	})
	if err != nil {
		return databases.Database{}, fmt.Errorf("error getting database: %s", err)
	}

	if len(resp.Databases) < 1 {
		return databases.Database{}, golangsdk.ErrDefault404{}
	}
	for _, db := range resp.Databases {
		if db.Name == dbName {
			return db, nil
		}
	}

	return databases.Database{}, golangsdk.ErrDefault404{}
}

func ResourceDliSqlDatabaseV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	db, err := GetDliSqlDatabaseByName(c, dbName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI database")
	}
	d.SetId(db.ResourceId)

	mErr := multierror.Append(nil,
		d.Set("name", db.Name),
		d.Set("description", db.Description),
		d.Set("enterprise_project_id", db.EnterpriseProjectId),
		d.Set("owner", db.Owner),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDliSqlDatabaseV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	_, err = databases.UpdateDBOwner(c, dbName, databases.UpdateDBOwnerOpts{
		NewOwner: d.Get("owner").(string),
	})
	if err != nil {
		return diag.Errorf("error updating SQL database owner: %s", err)
	}

	return ResourceDliSqlDatabaseV1Read(ctx, d, meta)
}

func ResourceDliSqlDatabaseV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	c, err := config.DliV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DLI v1 client: %s", err)
	}

	dbName := d.Get("name").(string)
	err = databases.Delete(c, dbName).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting SQL database: %s", err)
	}
	return nil
}

func resourceDatabaseImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("name", d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving resource name of the DLI database: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
