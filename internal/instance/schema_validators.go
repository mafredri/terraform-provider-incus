package instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/lxc/terraform-provider-incus/internal/utils"
)

// configKeyValidator ensures config key does not start
// with "volatile." or "image.".
type configKeyValidator struct{}

func (v configKeyValidator) Description(ctx context.Context) string {
	return fmt.Sprintf("config key cannot have %q or %q prefix", "volatile.", "image.")
}

func (v configKeyValidator) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("config key cannot have `%s` or `%s` prefix", "volatile.", "image.")
}

func (v configKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	value := req.ConfigValue.ValueString()

	if utils.HasAnyPrefix(value, []string{"volatile.", "image."}) {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid config key",
			fmt.Sprintf("Config key cannot have %q or %q prefix. Got: %q.", "volatile.", "image.", value),
		)
	}
}

// stopTimeoutValidator ensures the stop timeout is a duration of at least one second.
type stopTimeoutValidator struct{}

func (v stopTimeoutValidator) Description(ctx context.Context) string {
	return "stop timeout must be a duration of at least 1s"
}

func (v stopTimeoutValidator) MarkdownDescription(ctx context.Context) string {
	return "stop timeout must be a duration of at least `1s`"
}

func (v stopTimeoutValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	value := req.ConfigValue.ValueString()

	timeout, err := time.ParseDuration(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid stop timeout", fmt.Sprintf("Stop timeout must be a duration, such as \"30s\" or \"10m\". Got: %q.", value))
		return
	}

	if timeout < time.Second {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid stop timeout", fmt.Sprintf("Stop timeout must be at least 1s, because Incus force-stops on a timeout of 0. Got: %q.", value))
	}
}
