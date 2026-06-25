package cmd

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

func newMarketCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "market",
		Short: "LoomLoom Market SkillBot commands",
	}
	cmd.AddCommand(
		newMarketListCmd(opts),
		newMarketShowCmd(opts),
		newMarketQuoteCmd(opts),
		newMarketRunCmd(opts),
		newDeprecatedMarketPublishCmd(opts),
		newDeprecatedMarketRelistCmd(opts),
	)
	return cmd
}

func newMarketListCmd(opts *rootOptions) *cobra.Command {
	var (
		keyword   string
		limit     int
		pageSize  int
		pageToken string
		orderBy   string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List published Market SkillBots",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("limit") && cmd.Flags().Changed("page-size") {
				return fmt.Errorf("--limit and --page-size cannot be used together")
			}
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			if strings.TrimSpace(keyword) != "" {
				query.Set("keyword", strings.TrimSpace(keyword))
			}
			if pageSize > 0 {
				query.Set("pageSize", fmt.Sprintf("%d", pageSize))
			} else if limit > 0 {
				query.Set("pageSize", fmt.Sprintf("%d", limit))
			}
			if strings.TrimSpace(pageToken) != "" {
				query.Set("pageToken", strings.TrimSpace(pageToken))
			}
			if strings.TrimSpace(orderBy) != "" {
				query.Set("orderBy", strings.TrimSpace(orderBy))
			}
			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/marketListings", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "Search keyword matched against listing title and description")
	cmd.Flags().IntVar(&limit, "limit", 100, "Maximum number of listings")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Opaque pagination token returned by the previous response")
	cmd.Flags().StringVar(&orderBy, "order-by", "", "Sort expression, for example 'createdAt desc'")
	return cmd
}

func newMarketShowCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "show <listing-id>",
		Short: "Show one Market SkillBot detail",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			var resp map[string]any
			path := "/marketListings/" + url.PathEscape(strings.TrimSpace(args[0]))
			if err := httpClient.GetProductJSON(ctx, path, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
}

func newMarketQuoteCmd(opts *rootOptions) *cobra.Command {
	var inputFile string
	cmd := &cobra.Command{
		Use:   "quote <listing-id>",
		Short: "Quote a Market SkillBot run from a Product API request JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload, err := readJSONFileMap(inputFile)
			if err != nil {
				return err
			}
			removeIdentityFields(payload)

			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			var resp map[string]any
			path := "/marketListings/" + url.PathEscape(strings.TrimSpace(args[0])) + ":quote"
			if err := httpClient.PostProductJSON(ctx, path, payload, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&inputFile, "input-file", "", "Product API quote request JSON file")
	return cmd
}

func newMarketRunCmd(opts *rootOptions) *cobra.Command {
	var (
		inputFile       string
		clientRequestID string
		confirm         bool
	)
	cmd := &cobra.Command{
		Use:   "run <listing-id>",
		Short: "Execute a Market SkillBot from a Product API request JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return fmt.Errorf("--confirm is required")
			}
			payload, err := readJSONFileMap(inputFile)
			if err != nil {
				return err
			}
			removeIdentityFields(payload)
			generatedClientRequestID := false
			if strings.TrimSpace(clientRequestID) != "" {
				payload["clientRequestId"] = strings.TrimSpace(clientRequestID)
			} else if value, ok := payload["clientRequestId"]; !ok || strings.TrimSpace(fmt.Sprint(value)) == "" {
				generatedID, generated := effectiveClientRequestID("")
				payload["clientRequestId"] = generatedID
				generatedClientRequestID = generated
			}
			payload["confirm"] = true

			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			printGeneratedClientRequestID(cmd, fmt.Sprint(payload["clientRequestId"]), generatedClientRequestID)
			var resp map[string]any
			path := "/marketListings/" + url.PathEscape(strings.TrimSpace(args[0])) + ":execute"
			if err := httpClient.PostProductJSON(ctx, path, payload, &resp); err != nil {
				return err
			}
			opts.debugf(
				"market run: submitted listing_id=%s run_id=%s transaction_id=%s",
				strings.TrimSpace(args[0]),
				stringMapValue(resp, "runId"),
				stringMapValue(resp, "runTransactionId"),
			)
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&inputFile, "input-file", "", "Product API execute request JSON file")
	cmd.Flags().StringVar(&clientRequestID, "client-request-id", "", "Stable idempotency key for retrying the same confirmed execution")
	cmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm execution and create the Market run")
	return cmd
}

func newDeprecatedMarketPublishCmd(opts *rootOptions) *cobra.Command {
	var (
		listingID         string
		templateID        string
		templateVersionID string
		displayName       string
		description       string
		taskFixedFeeT     int64
	)
	cmd := &cobra.Command{
		Use:        "publish",
		Short:      "Publish a template version as a Market SkillBot",
		Deprecated: "use 'loomloom listing publish <template-id>'",
		RunE: func(cmd *cobra.Command, args []string) error {
			req := publishMarketListingRequest{
				ListingID:         strings.TrimSpace(listingID),
				DisplayName:       strings.TrimSpace(displayName),
				Description:       strings.TrimSpace(description),
				TaskFixedFeeT:     taskFixedFeeT,
				TemplateID:        strings.TrimSpace(templateID),
				TemplateVersionID: strings.TrimSpace(templateVersionID),
			}

			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			var resp map[string]any
			if err := httpClient.PostProductJSON(ctx, "/marketListings", req, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&listingID, "listing-id", "", "Existing listing ID when publishing a new version")
	cmd.Flags().StringVar(&templateID, "template-id", "", "Template ID to publish")
	cmd.Flags().StringVar(&templateVersionID, "template-version-id", "", "Template version ID to publish")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Market SkillBot display name")
	cmd.Flags().StringVar(&description, "description", "", "Market SkillBot description")
	cmd.Flags().Int64Var(&taskFixedFeeT, "task-fixed-fee-t", 0, "Creator fixed fee per billable task, in API units")
	_ = cmd.MarkFlagRequired("template-id")
	_ = cmd.MarkFlagRequired("template-version-id")
	_ = cmd.MarkFlagRequired("display-name")
	_ = cmd.MarkFlagRequired("task-fixed-fee-t")
	return cmd
}

func newDeprecatedMarketRelistCmd(opts *rootOptions) *cobra.Command {
	cmd := newMarketSaleStatusCmd(opts, "relist", "Restore a Market SkillBot listing", "list")
	cmd.Deprecated = "use 'loomloom listing relist <listing-id>'"
	return cmd
}

func newMarketSaleStatusCmd(opts *rootOptions, use string, short string, action string) *cobra.Command {
	return &cobra.Command{
		Use:   use + " <listing-id>",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			var resp map[string]any
			path := "/marketListings/" + url.PathEscape(strings.TrimSpace(args[0])) + ":" + action
			if err := httpClient.PostProductJSONWithQuery(ctx, path, query, nil, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
}
