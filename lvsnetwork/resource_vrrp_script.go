package lvsnetwork

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	defaultVrrpScriptRise = 3
	defaultVrrpScriptFall = 3
)

func resourceVrrpScript() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVrrpScriptCreate,
		ReadContext:   resourceVrrpScriptRead,
		UpdateContext: resourceVrrpScriptUpdate,
		DeleteContext: resourceVrrpScriptDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"script": {
				Type:     schema.TypeString,
				Required: true,
			},
			"interval": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"weight_reverse": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"rise": {
				Type:         schema.TypeInt,
				Default:      defaultVrrpScriptRise,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"fall": {
				Type:         schema.TypeInt,
				Default:      defaultVrrpScriptFall,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"init_fail": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceVrrpScriptCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript(ctx, ADD, &vrrpScript)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(d.Get("name").(string))

	return resourceVrrpScriptRead(ctx, d, m)
}

func resourceVrrpScriptRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	vrrpScriptRead, err := client.requestAPIVrrpScript(ctx, CHECK, &vrrpScript)
	if err != nil {
		return diag.FromErr(err)
	}
	if vrrpScriptRead.Name == "" {
		d.SetId("")
	}
	tfErr := d.Set("script", vrrpScriptRead.Script)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("interval", vrrpScriptRead.Interval)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("timeout", vrrpScriptRead.Timeout)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("weight", vrrpScriptRead.Weight)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("weight_reverse", vrrpScriptRead.WeightReverse)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("rise", vrrpScriptRead.Rise)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("fall", vrrpScriptRead.Fall)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("user", vrrpScriptRead.User)
	if tfErr != nil {
		panic(tfErr)
	}
	tfErr = d.Set("init_fail", vrrpScriptRead.InitFail)
	if tfErr != nil {
		panic(tfErr)
	}

	return nil
}

func resourceVrrpScriptUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	d.Partial(true)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript(ctx, CHANGE, &vrrpScript)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceVrrpScriptRead(ctx, d, m)
}

func resourceVrrpScriptDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	vrrpScript := createStrucVrrpScript(d)
	_, err := client.requestAPIVrrpScript(ctx, REMOVE, &vrrpScript)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func createStrucVrrpScript(d *schema.ResourceData) vrrpScript {
	vrrpScript := vrrpScript{
		Name:          d.Get("name").(string),
		Script:        d.Get("script").(string),
		Interval:      d.Get("interval").(int),
		Timeout:       d.Get("timeout").(int),
		Weight:        d.Get("weight").(int),
		WeightReverse: d.Get("weight_reverse").(bool),
		Rise:          d.Get("rise").(int),
		Fall:          d.Get("fall").(int),
		User:          d.Get("user").(string),
		InitFail:      d.Get("init_fail").(bool),
	}

	return vrrpScript
}
