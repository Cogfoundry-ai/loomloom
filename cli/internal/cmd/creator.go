package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newCreatorCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "creator",
		Short: "Creator Market commands",
	}
	cmd.AddCommand(
		newCreatorEarningsCmd(opts),
		newCreatorTransactionsCmd(opts),
		newCreatorReviewCmd(opts),
	)
	return cmd
}

func newCreatorEarningsCmd(opts *rootOptions) *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "earnings",
		Short: "List creator Market earnings",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			if limit > 0 {
				query.Set("pageSize", fmt.Sprintf("%d", limit))
			}

			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/creators/me/earnings", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of earning records")
	return cmd
}

func newCreatorTransactionsCmd(opts *rootOptions) *cobra.Command {
	var pageSize int
	cmd := &cobra.Command{
		Use:   "transactions",
		Short: "List creator Market transactions",
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
			if err := httpClient.GetProductJSONWithQuery(ctx, "/creators/me/marketTransactions", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newCreatorReviewCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Creator Market review request commands",
	}
	cmd.AddCommand(
		newCreatorReviewListCmd(opts),
		newCreatorReviewGetCmd(opts),
		newCreatorReviewWithdrawCmd(opts),
	)
	return cmd
}

func newCreatorReviewListCmd(opts *rootOptions) *cobra.Command {
	var (
		status   string
		pageSize int
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List my Market review requests",
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			if strings.TrimSpace(status) != "" {
				query.Set("status", strings.TrimSpace(status))
			}
			if pageSize > 0 {
				query.Set("pageSize", fmt.Sprintf("%d", pageSize))
			}

			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/creators/me/marketReviewRequests", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "Filter by review status")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newCreatorReviewGetCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "get <review-request-id>",
		Short: "Get one Market review request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			path := "/creators/me/marketReviewRequests/" + url.PathEscape(strings.TrimSpace(args[0]))
			var resp map[string]any
			if err := httpClient.GetProductJSON(ctx, path, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
}

func newCreatorReviewWithdrawCmd(opts *rootOptions) *cobra.Command {
	var reason string
	cmd := &cobra.Command{
		Use:   "withdraw <review-request-id>",
		Short: "Withdraw a pending Market review request",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			payload := map[string]any{}
			if strings.TrimSpace(reason) != "" {
				payload["reason"] = strings.TrimSpace(reason)
			}

			path := "/creators/me/marketReviewRequests/" + url.PathEscape(strings.TrimSpace(args[0])) + ":withdraw"
			var resp map[string]any
			if err := httpClient.PostProductJSON(ctx, path, payload, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&reason, "reason", "", "Optional withdrawal reason")
	return cmd
}
