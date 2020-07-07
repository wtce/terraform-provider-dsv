package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/thycotic/dsv-sdk-go/vault"
)

func setResourceSecretAttributes(secret *vault.Secret, d *schema.ResourceData) {
	secret.Description = d.Get("description").(string)
	secret.Data = d.Get("data").(map[string]interface{})
	secret.Attributes = d.Get("attributes").(map[string]interface{})
}

func resourceSecretRead(d *schema.ResourceData, meta interface{}) error {
	path := d.Get("path").(string)
	dsv, err := vault.New(meta.(vault.Configuration))

	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}
	log.Printf("[DEBUG] getting secret %s", path)

	secret, err := dsv.Secret(path)

	if err != nil {
		log.Printf("[DEBUG] unable to get secret: %s", err)
		return err
	}

	d.SetId(secret.ID)

	setResourceDataAttributes(secret, d)

	return nil
}

func resourceSecretCreate(d *schema.ResourceData, meta interface{}) error {
	dsv, err := vault.New(meta.(vault.Configuration))
	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}

	secret := new(vault.Secret)
	secret.Path = d.Get("path").(string)

	setResourceSecretAttributes(secret, d)

	err = dsv.NewSecret(secret)
	if err != nil {
		log.Printf("[ERROR] unable to create secret: %s", err)
		return err
	}

	return resourceSecretRead(d, meta)
}

func resourceSecretUpdate(d *schema.ResourceData, meta interface{}) error {
	dsv, err := vault.New(meta.(vault.Configuration))
	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}

	secret, err := dsv.Secret(d.Get("path").(string))
	if err != nil {
		log.Printf("[DEBUG] unable to get secret: %s", err)
		return err
	}

	setResourceSecretAttributes(secret, d)

	err = secret.Update(true)
	if err != nil {
		log.Printf("[ERROR] unable to update secret: %s", err)
		return err
	}

	return resourceSecretRead(d, meta)
}

func resourceSecretDelete(d *schema.ResourceData, meta interface{}) error {
	dsv, err := vault.New(meta.(vault.Configuration))
	if err != nil {
		log.Printf("[DEBUG] configuration error: %s", err)
		return err
	}

	secret, err := dsv.Secret(d.Get("path").(string))
	if err != nil {
		log.Printf("[DEBUG] unable to get secret: %s", err)
		return err
	}

	err = secret.Delete(true)
	if err != nil {
		log.Printf("[DEBUG] unable to delete secret: %s", err)
		return err
	}

	return nil
}

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Read:   resourceSecretRead,

		Schema: map[string]*schema.Schema{
			"data": {
				Description: "the data of the secret",
				Required:    true,
				Sensitive:   true,
				Type:        schema.TypeMap,
			},
			"attributes": {
				Description: "the attributes of the secret",
				Optional:    true,
				Type:        schema.TypeMap,
			},
			"description": {
				Description: "the description of the secret",
				Optional:    true,
				Type:        schema.TypeString,
			},
			"path": {
				Description: "the path of the secret",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},
			"version": {
				Computed:    true,
				Description: "the version of the secret",
				Type:        schema.TypeInt,
			},
		},
	}
}
