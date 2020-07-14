package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/signalsciences/go-sigsci"
)

func resource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTemplatedRuleCreate,
		Update: resourceTemplatedRuleUpdate,
		Read:   resourceTemplatedRuleRead,
		Delete: resourceTemplatedRuleDelete,
		//Importer: &schema.ResourceImporter{ //TODO try importing, make sure it works
		//	State: schema.ImportStatePassthrough, // this only sets the id. Probably a better way
		//},
		Schema: map[string]*schema.Schema{
			"site_short_name": {
				Type:        schema.TypeString,
				Description: "Site short name",
				Required:    true,
				ForceNew:    true,
			},
			"tag_name": {
				Type:        schema.TypeString,
				Description: "The name of the tag whose occurrences the alert is watching. Must match an existing tag",
				Required:    true,
			},
			"long_name": {
				Type:        schema.TypeString,
				Description: "description",
				Optional:    true,
			},
			"interval": {
				Type:        schema.TypeInt,
				Description: "The number of minutes of past traffic to examine. Must be 1, 10 or 60.",
				Optional:    true,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Description: "The number of occurrences of the tag in the interval needed to trigger the alert. Min 1, Max 10000",
				Optional:    true,
			},
			"block_duration_seconds": {
				Type:        schema.TypeInt,
				Description: "The number of seconds this alert is active.",
				Optional:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Description: "A flag to toggle this alert.",
				Optional:    true,
			},
			"action": {
				Type:        schema.TypeString,
				Description: "A flag that describes what happens when the alert is triggered. 'info' creates an incident in the dashboard. 'flagged' creates an incident and blocks traffic for 24 hours. Must be info or flagged.",
				Optional:    true,
			},
		},
	}
}

func resourceTemplatedRuleCreate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.UpdateSiteTemplateRuleByID(pm.Corp, d.Get("site_short_name").(string), sigsci.SiteTemplateRuleBody{
		DetectionAdds:    nil,
		DetectionUpdates: nil,
		DetectionDeletes: nil,
		AlertAdds:        nil,
		AlertUpdates:     nil,
		AlertDeletes:     nil,
	})
	if err != nil {
		return err
	}
	d.SetId(alert.ID)
	return resourceTemplatedRuleRead(d, m)
}

func resourceTemplatedRuleRead(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.GetCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}

	d.SetId(alert.ID)
	err = d.Set("site_short_name", d.Get("site_short_name").(string))
	if err != nil {
		return err
	}
	err = d.Set("tag_name", alert.TagName)
	if err != nil {
		return err
	}
	err = d.Set("long_name", alert.LongName)
	if err != nil {
		return err
	}
	err = d.Set("interval", alert.Interval)
	if err != nil {
		return err
	}
	err = d.Set("threshold", alert.Threshold)
	if err != nil {
		return err
	}
	err = d.Set("enabled", alert.Enabled)
	if err != nil {
		return err
	}
	err = d.Set("action", alert.Action)
	if err != nil {
		return err
	}

	return nil
}

func resourceTemplatedRuleUpdate(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	alert, err := sc.UpdateCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id(), sigsci.CustomAlertBody{
		TagName:   d.Get("tag_name").(string),
		LongName:  d.Get("long_name").(string),
		Interval:  d.Get("interval").(int),
		Threshold: d.Get("threshold").(int),
		Enabled:   d.Get("enabled").(bool),
		Action:    d.Get("action").(string),
	})
	if err != nil {
		return err
	}

	d.SetId(alert.ID)
	return resourceTemplatedRuleRead(d, m)
}

func resourceTemplatedRuleDelete(d *schema.ResourceData, m interface{}) error {
	pm := m.(providerMetadata)
	sc := pm.Client

	err := sc.DeleteCustomAlert(pm.Corp, d.Get("site_short_name").(string), d.Id())
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}