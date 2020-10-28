package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/mrs/v1/job"
)

func resourceMRSJobV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceMRSJobV1Create,
		Read:   resourceMRSJobV1Read,
		Delete: resourceMRSJobV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"jar_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"arguments": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"input": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"output": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_log": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hive_script_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"job_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func JobStateRefreshFunc(client *golangsdk.ServiceClient, jobID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		jobGet, err := job.Get(client, jobID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return jobGet, "DELETED", nil
			}
			return nil, "", err
		}
		log.Printf("[DEBUG] JobStateRefreshFunc: %#v", jobGet)
		jobState := "Starting"
		if jobGet.JobState == -1 {
			jobState = "Terminated"
		} else if jobGet.JobState == 1 {
			jobState = "Starting"
		} else if jobGet.JobState == 2 {
			jobState = "Running"
		} else if jobGet.JobState == 3 {
			jobState = "Completed"
		} else if jobGet.JobState == 4 {
			jobState = "Abnormal"
		} else if jobGet.JobState == 5 {
			jobState = "Error"
		}
		return jobGet, jobState, nil
	}
}

func resourceMRSJobV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud MRS client: %s", err)
	}

	createOpts := &job.CreateOpts{
		JobType:        d.Get("job_type").(int),
		JobName:        d.Get("job_name").(string),
		ClusterID:      d.Get("cluster_id").(string),
		JarPath:        d.Get("jar_path").(string),
		Arguments:      d.Get("arguments").(string),
		Input:          d.Get("input").(string),
		Output:         d.Get("output").(string),
		JobLog:         d.Get("job_log").(string),
		HiveScriptPath: d.Get("hive_script_path").(string),
		IsProtected:    d.Get("is_protected").(bool),
		IsPublic:       d.Get("is_public").(bool),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	jobCreate, err := job.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Job: %s", err)
	}

	d.SetId(jobCreate.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Starting", "Running"},
		Target:     []string{"Completed"},
		Refresh:    JobStateRefreshFunc(client, jobCreate.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for job (%s) to become ready: %s ",
			jobCreate.ID, err)
	}

	return resourceMRSJobV1Read(d, meta)
}

func resourceMRSJobV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud  MRS client: %s", err)
	}

	jobGet, err := job.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "MRS Job")
	}
	log.Printf("[DEBUG] Retrieved MRS Job %s: %#v", d.Id(), jobGet)

	d.Set("region", GetRegion(d, config))
	d.SetId(jobGet.ID)
	d.Set("job_type", jobGet.JobType)
	d.Set("job_name", jobGet.JobName)
	d.Set("cluster_id", jobGet.ClusterID)
	d.Set("jar_path", jobGet.JarPath)
	d.Set("arguments", jobGet.Arguments)
	d.Set("input", jobGet.Input)
	d.Set("output", jobGet.Output)
	d.Set("job_log", jobGet.JobLog)
	d.Set("hive_script_path", jobGet.HiveScriptPath)

	jobState := "Starting"
	if jobGet.JobState == -1 {
		jobState = "Terminated"
	} else if jobGet.JobState == 1 {
		jobState = "Starting"
	} else if jobGet.JobState == 2 {
		jobState = "Running"
	} else if jobGet.JobState == 3 {
		jobState = "Completed"
	} else if jobGet.JobState == 4 {
		jobState = "Abnormal"
	} else if jobGet.JobState == 5 {
		jobState = "Error"
	}
	d.Set("job_state", jobState)
	return nil
}

func resourceMRSJobV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.MrsV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud client: %s", err)
	}

	rId := d.Id()
	log.Printf("[DEBUG] Deleting MRS Job %s", rId)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := job.Delete(client, rId).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable MRS Job: %s", rId)
			return nil
		}
		return fmt.Errorf("Error deleting MRS Job %s: %s", rId, err)
	}
	return nil
}
