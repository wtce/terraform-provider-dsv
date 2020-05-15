package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thycotic/dsv-sdk-go/vault"
)

func providerConfig(d *schema.ResourceData) (interface{}, error) {
	c := vault.Configuration{
		Tenant: d.Get("tenant").(string),
		Credentials: vault.ClientCredential{
			ClientID:     d.Get("client_id").(string),
			ClientSecret: d.Get("client_secret").(string),
		},
	}

	log.Printf("[DEBUG] tenant is set to %s", c.Tenant)

	if tld, ok := d.GetOk("tld"); ok {
		c.TLD = tld.(string)
		log.Printf("[DEBUG] tld is set to %s", c.TLD)
	}

	if ut, ok := d.GetOk("url_template"); ok {
		c.URLTemplate = ut.(string)
		log.Printf("[DEBUG] url_template is set to %s", c.URLTemplate)
	}
	return c, nil
}

// Provider is a Terraform DataSource
func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"dsv_client": dataSourceClient(),
			"dsv_role":   dataSourceRole(),
			"dsv_secret": dataSourceSecret(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsv_client": resourceClient(),
			"dsv_secret": resourceSecret(),
		},
		Schema: map[string]*schema.Schema{
			"tenant": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DevOps Secrets Vault tenant",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The DevOps Secrets Vault client_id",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Sensitive:   true,
				Required:    true,
				Description: "The DevOps Secrets Vault client_secret",
			},
			"tld": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DSV tenant top-level domain",
			},
			"url_template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The DSV SDK API URL template",
			},
		},
		ConfigureFunc: providerConfig,
	}
}
