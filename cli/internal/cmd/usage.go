package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newUsageCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Buyer Market usage commands",
	}
	cmd.AddCommand(
		newUsageListCmd(opts),
		newUsageGetCmd(opts),
	)
	return cmd
}

func newUsageListCmd(opts *rootOptions) *cobra.Command {
	var pageSize int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List my Market SkillBot usage records",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			if pageSize > 0 {
				query.Set("pageSize", fmt.Sprintf("%d", pageSize))
			}

			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/users/me/marketUsages", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newUsageGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <run-transaction-id>",
		Short: "Get one Market SkillBot usage record",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			path := "/users/me/marketUsages/" + url.PathEscape(strings.TrimSpace(args[0]))
			var resp map[string]any
			if err := httpClient.GetProductJSON(ctx, path, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
}
