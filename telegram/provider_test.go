package telegram

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProvider  *schema.Provider
	testAccProviders map[string]*schema.Provider
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"telegram": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("TELEGRAM_BOT_TOKEN"); v == "" {
		t.Fatal("TELEGRAM_BOT_TOKEN must be set for acceptance tests")
	}
}
