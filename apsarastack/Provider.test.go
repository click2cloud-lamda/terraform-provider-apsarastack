package apsarastack

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"os"
	"testing"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var defaultRegionToTest = os.Getenv("APSARASTACK_REGION")

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"apsarastack": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("APSARASTACK_ACCESS_KEY"); v == "" {
		t.Fatal("APSARASTACK_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_SECRET_KEY"); v == "" {
		t.Fatal("APSARASTACK_SECRET_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("APSARASTACK_REGION"); v == "" {
		log.Println("[INFO] Test: Using cn-beijing as test region")
		os.Setenv("APSARASTACK_REGION", "cn-beijing")
	}
}

var providerCommon = `
provider "apsarastack" {
	assume_role {}
}
`
