package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/yi-jiayu/terraform-provider-telegram/telegram"
)

func main() {
	upgradedSdkProvider, err := tf5to6server.UpgradeServer(
		context.Background(),
		telegram.Provider().GRPCProvider,
	)
	if err != nil {
		panic(err)
	}

	err = tf6server.Serve(
		"registry.terraform.io/example/example",
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
	)
	if err != nil {
		panic(err)
	}
}
