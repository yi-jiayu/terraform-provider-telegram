package telegram

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/yi-jiayu/ted"

	"github.com/yi-jiayu/terraform-provider-telegram/telegram/internal"
)

func resourceTelegramBotCommands() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTelegramBotCommandsCreate,
		ReadContext:   resourceTelegramBotCommandsRead,
		UpdateContext: resourceTelegramBotCommandsUpdate,
		Delete:        resourceTelegramBotCommandsDelete,
		Schema: map[string]*schema.Schema{
			"commands": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Required: true,
			},
		},
	}
}

func resourceTelegramBotCommandsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bot := m.(ted.Bot)
	commands := d.Get("commands").([]interface{})
	expanded, err := expandBotCommands(commands)
	if err != nil {
		return diag.FromErr(err)
	}
	setMyCommands := ted.SetMyCommandsRequest{
		Commands: expanded,
	}
	err = internal.Retry(3, func() error {
		_, err := bot.Do(setMyCommands)
		if err != nil {
			return fmt.Errorf("setMyCommands error: %w", err)
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("commands")
	return resourceTelegramBotCommandsRead(ctx, d, m)
}

func resourceTelegramBotCommandsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	bot := m.(ted.Bot)
	myCommands, err := bot.GetMyCommands()
	if err != nil {
		err := fmt.Errorf("getMyCommands error: %w", err)
		return diag.FromErr(err)
	}
	commands := flattenBotCommands(myCommands)
	err = d.Set("commands", commands)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceTelegramBotCommandsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceTelegramBotCommandsCreate(ctx, d, m)
}

func resourceTelegramBotCommandsDelete(d *schema.ResourceData, m interface{}) error {
	bot := m.(ted.Bot)
	setMyCommands := ted.SetMyCommandsRequest{
		Commands: []ted.BotCommand{},
	}
	err := internal.Retry(3, func() error {
		_, err := bot.Do(setMyCommands)
		if err != nil {
			return fmt.Errorf("setMyCommands error: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func flattenBotCommands(commands []ted.BotCommand) []interface{} {
	flattened := make([]interface{}, len(commands))
	for i := range commands {
		flattened[i] = map[string]interface{}{
			"command":     commands[i].Command,
			"description": commands[i].Description,
		}
	}
	return flattened
}

func expandBotCommands(commands []interface{}) ([]ted.BotCommand, error) {
	expanded := make([]ted.BotCommand, len(commands))
	for i, command := range commands {
		command := command.(map[string]interface{})
		cmd, ok := command["command"]
		if !ok {
			return nil, fmt.Errorf("command not set: commands[%d]", i)
		}
		desc, ok := command["description"]
		if !ok {
			return nil, fmt.Errorf("description not set: commands[%d]", i)
		}
		expanded[i] = ted.BotCommand{
			Command:     cmd.(string),
			Description: desc.(string),
		}
	}
	return expanded, nil
}
