package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/Cogfoundry-ai/loomloom/cli/internal/client"
	"github.com/Cogfoundry-ai/loomloom/cli/internal/version"
	"github.com/spf13/cobra"
)

func newDoctorCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check LoomLoom server reachability and token wiring",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			systemBaseURL, err := productAPISystemBaseURL(opts.server)
			if err != nil {
				return err
			}
			systemClient, err := client.New(client.Config{
				BaseURL: systemBaseURL,
				Token:   opts.token,
				Timeout: opts.timeout,
			})
			if err != nil {
				return err
			}
			var productVersion map[string]any
			versionResp, err := systemClient.GetBinary(ctx, "/version")
			if err != nil {
				return fmt.Errorf("check Product API version: %w", err)
			}
			if strings.Contains(strings.ToLower(versionResp.ContentType), "application/json") {
				_ = json.Unmarshal(versionResp.Body, &productVersion)
			}
			if productVersion == nil {
				productVersion = map[string]any{
					"reachable":   true,
					"contentType": versionResp.ContentType,
				}
			}

			publicQuery := url.Values{}
			publicQuery.Set("pageSize", "1")
			var marketResp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/marketListings", publicQuery, &marketResp); err != nil {
				return fmt.Errorf("check public Market endpoint: %w", err)
			}

			query := url.Values{}
			query.Set("pageSize", "1")
			var probeResp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/users/me/executables", query, &probeResp); err != nil {
				return fmt.Errorf("check authenticated executables endpoint: %w", err)
			}
			healthy := true
			message := "ok"

			versionStatus, versionErr := version.CheckLatest(ctx)
			currentVersion := version.Version
			latestVersion := ""
			currentChannel := version.ReleaseChannel(currentVersion)
			latestChannel := ""
			updateAvailable := false
			upgradeHint := ""
			if versionStatus != nil {
				currentVersion = versionStatus.CurrentVersion
				latestVersion = versionStatus.LatestVersion
				currentChannel = versionStatus.CurrentChannel
				latestChannel = versionStatus.LatestChannel
				updateAvailable = versionStatus.UpdateAvailable
				upgradeHint = versionStatus.UpgradeHint
			}

			if opts.output == "json" {
				payload := map[string]any{
					"server":           opts.server,
					"token_set":        opts.token != "",
					"healthy":          healthy,
					"message":          message,
					"product_version":  productVersion,
					"cli_version":      currentVersion,
					"release_channel":  currentChannel,
					"latest_release":   latestVersion,
					"latest_channel":   latestChannel,
					"update_available": updateAvailable,
					"upgrade_hint":     upgradeHint,
					"base_usage":       "set LOOMLOOM_SERVER and LOOMLOOM_TOKEN before running template commands",
				}
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(payload)
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "server: %s\n", opts.server)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "token: %t\n", opts.token != "")
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "healthy: %t\n", healthy)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "message: %s\n", message)
			if service, ok := productVersion["service"]; ok {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "product_api: %v\n", service)
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "cli_version: %s\n", currentVersion)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "release_channel: %s\n", currentChannel)
			if latestVersion != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "latest_release: %s\n", latestVersion)
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "latest_channel: %s\n", latestChannel)
			}
			if upgradeHint != "" {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "notice: %s\n", upgradeHint)
			} else if versionErr != nil {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "version check unavailable; skipping release notice")
			}
			return nil
		},
	}
}

func productAPISystemBaseURL(server string) (string, error) {
	raw := strings.TrimSpace(server)
	if raw == "" {
		return "", fmt.Errorf("server URL is required; set LOOMLOOM_SERVER or pass --server")
	}
	if !strings.Contains(raw, "://") {
		if strings.Contains(raw, ":") {
			raw = "http://" + raw
		} else {
			raw = "https://" + raw
		}
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("parse server URL: %w", err)
	}
	parsed.Path = strings.TrimSuffix(strings.TrimRight(parsed.Path, "/"), "/loom/v1")
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return strings.TrimRight(parsed.String(), "/"), nil
}
