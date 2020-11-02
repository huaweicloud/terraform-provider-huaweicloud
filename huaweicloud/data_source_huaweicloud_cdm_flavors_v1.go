package huaweicloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
)

func dataSourceCdmFlavorV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCdmFlavorV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCdmFlavorV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.cdmV11Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating sdk client, err=%s", err)
	}

	dsid, err := getCdmDatastoreV1ID(client)
	if err != nil {
		return err
	}

	version, fs, err := getCdmFlavorV1(client, dsid)
	if err != nil {
		return err
	}

	d.SetId(version)
	d.Set("version", version)
	d.Set("flavors", fs)
	return nil
}

func getCdmDatastoreV1ID(client *golangsdk.ServiceClient) (string, error) {
	url := client.ServiceURL("datastores")
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return "", r.Err
	}

	v, err := navigateValue(r.Body, []string{"datastores"}, nil)
	if err != nil {
		return "", err
	}

	ds, ok := v.([]interface{})
	if !ok {
		return "", fmt.Errorf("can not find datastore")
	}

	for _, item := range ds {
		name, err := navigateValue(item, []string{"name"}, nil)
		if err != nil {
			return "", err
		}
		if "cdm" == name.(string) {
			dsid, err := navigateValue(item, []string{"id"}, nil)
			if err != nil {
				return "", err
			}
			return dsid.(string), nil
		}
	}

	return "", fmt.Errorf("didn't find the datastore id")
}

func getCdmFlavorV1(client *golangsdk.ServiceClient, dsid string) (string, interface{}, error) {
	url := client.ServiceURL("datastores", dsid, "flavors")
	r := golangsdk.Result{}
	_, r.Err = client.Get(url, &r.Body, &golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		}})
	if r.Err != nil {
		return "", nil, r.Err
	}

	v, err := navigateValue(r.Body, []string{"versions"}, nil)
	if err != nil {
		return "", nil, err
	}
	vs, ok := v.([]interface{})
	if !ok {
		return "", nil, fmt.Errorf("can not find flavor")
	}
	for _, item := range vs {
		version, err := navigateValue(item, []string{"name"}, nil)
		if err != nil {
			return "", nil, err
		}
		flavors, err := navigateValue(item, []string{"flavors"}, nil)
		if err != nil {
			return "", nil, err
		}

		fs, ok := flavors.([]interface{})
		if !ok {
			return "", nil, fmt.Errorf("can not find flavor")
		}
		num := len(fs)
		r := make([]interface{}, num, num)
		for i := 0; i < num; i++ {
			item := fs[i]
			name, err := navigateValue(item, []string{"name"}, nil)
			if err != nil {
				return "", nil, err
			}
			fid, err := navigateValue(item, []string{"str_id"}, nil)
			if err != nil {
				return "", nil, err
			}

			r[i] = map[string]interface{}{
				"id":   fid,
				"name": name,
			}
		}
		return version.(string), r, nil
	}

	return "", nil, fmt.Errorf("can not find flavor")
}
