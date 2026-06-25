package cmd

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func newListingCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "listing",
		Short: "Creator listing commands",
	}
	cmd.AddCommand(
		newListingPublishCmd(opts),
		newListingUpdateCmd(opts),
		newListingListCmd(opts),
		newListingShowCmd(opts),
		newListingVersionsCmd(opts),
		newListingWithdrawCmd(opts),
		newListingUnlistCmd(opts),
		newListingRelistCmd(opts),
	)
	return cmd
}

func newListingPublishCmd(opts *rootOptions) *cobra.Command {
	var (
		listingID         string
		templateVersionID string
		displayName       string
		description       string
		taskFixedFeeT     int64
	)
	cmd := &cobra.Command{
		Use:   "publish <template-id>",
		Short: "Submit a template version for Market SkillBot review",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			req := publishMarketListingRequest{
				ListingID:         strings.TrimSpace(listingID),
				TemplateID:        strings.TrimSpace(args[0]),
				TemplateVersionID: strings.TrimSpace(templateVersionID),
				DisplayName:       strings.TrimSpace(displayName),
				Description:       strings.TrimSpace(description),
				TaskFixedFeeT:     taskFixedFeeT,
			}

			var resp map[string]any
			if err := httpClient.PostProductJSON(ctx, "/marketListings", req, &resp); err != nil {
				if strings.Contains(err.Error(), "must have a successful run") {
					return fmt.Errorf("%w\nHint: submit at least one successful run (template-spec submit-workbook) before publishing", err)
				}
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&listingID, "listing-id", "", "Existing listing ID when publishing a new version")
	cmd.Flags().StringVar(&templateVersionID, "template-version-id", "", "Template version ID to publish")
	cmd.Flags().StringVar(&displayName, "display-name", "", "Market SkillBot display name")
	cmd.Flags().StringVar(&description, "description", "", "Market SkillBot description")
	cmd.Flags().Int64Var(&taskFixedFeeT, "task-fixed-fee-t", 0, "Creator fixed fee per billable task, in API units")
	_ = cmd.MarkFlagRequired("template-version-id")
	_ = cmd.MarkFlagRequired("display-name")
	_ = cmd.MarkFlagRequired("task-fixed-fee-t")
	return cmd
}

func newListingListCmd(opts *rootOptions) *cobra.Command {
	var (
		keyword   string
		pageSize  int
		pageToken string
		orderBy   string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List my Market SkillBot listings",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			}
			if strings.TrimSpace(pageToken) != "" {
				query.Set("pageToken", strings.TrimSpace(pageToken))
			}
			if strings.TrimSpace(orderBy) != "" {
				query.Set("orderBy", strings.TrimSpace(orderBy))
			}

			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, "/creators/me/marketListings", query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "Search keyword")
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	cmd.Flags().StringVar(&pageToken, "page-token", "", "Pagination token from previous response")
	cmd.Flags().StringVar(&orderBy, "order-by", "", "Sort expression, e.g. 'createdAt desc'")
	return cmd
}

func newListingShowCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "show <listing-id>",
		Short: "Show one creator-owned Market SkillBot listing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			path := "/creators/me/marketListings/" + url.PathEscape(strings.TrimSpace(args[0]))
			var resp map[string]any
			if err := httpClient.GetProductJSON(ctx, path, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
}

func newListingVersionsCmd(opts *rootOptions) *cobra.Command {
	var pageSize int
	cmd := &cobra.Command{
		Use:   "versions <listing-id>",
		Short: "List versions of a creator-owned Market SkillBot listing",
		Args:  cobra.ExactArgs(1),
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

			path := "/creators/me/marketListings/" + url.PathEscape(strings.TrimSpace(args[0])) + "/versions"
			var resp map[string]any
			if err := httpClient.GetProductJSONWithQuery(ctx, path, query, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().IntVar(&pageSize, "page-size", 0, "Page size")
	return cmd
}

func newListingUpdateCmd(opts *rootOptions) *cobra.Command {
	var (
		displayName string
		description string
	)
	cmd := &cobra.Command{
		Use:   "update <listing-id>",
		Short: "Submit a public profile update for review (display name / description)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("display-name") && !cmd.Flags().Changed("description") {
				return fmt.Errorf("at least one of --display-name or --description is required")
			}
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			listingID := strings.TrimSpace(args[0])
			current := creatorListingProfile{}
			if !cmd.Flags().Changed("display-name") || !cmd.Flags().Changed("description") {
				path := "/creators/me/marketListings/" + url.PathEscape(listingID)
				if err := httpClient.GetProductJSON(ctx, path, &current); err != nil {
					return fmt.Errorf("get current listing profile: %w", err)
				}
			}

			nextDisplayName := strings.TrimSpace(displayName)
			if !cmd.Flags().Changed("display-name") {
				nextDisplayName = strings.TrimSpace(current.DisplayName)
			}
			if nextDisplayName == "" {
				return fmt.Errorf("display name is required; pass --display-name or ensure the current listing has one")
			}

			nextDescription := strings.TrimSpace(description)
			if !cmd.Flags().Changed("description") {
				nextDescription = strings.TrimSpace(current.Description)
			}

			payload := map[string]any{
				"displayName": nextDisplayName,
				"description": nextDescription,
			}

			path := "/marketListings/" + url.PathEscape(listingID) + ":updatePublicProfile"
			var resp map[string]any
			if err := httpClient.PostProductJSON(ctx, path, payload, &resp); err != nil {
				return err
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&displayName, "display-name", "", "New display name")
	cmd.Flags().StringVar(&description, "description", "", "New description")
	return cmd
}

type creatorListingProfile struct {
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type creatorReviewSummary struct {
	ID        string `json:"id"`
	ListingID string `json:"listingId"`
	Status    string `json:"status"`
}

type creatorReviewListResponse struct {
	Items []creatorReviewSummary `json:"items"`
}

func newListingWithdrawCmd(opts *rootOptions) *cobra.Command {
	var reason string
	cmd := &cobra.Command{
		Use:   "withdraw <listing-id>",
		Short: "Withdraw the pending review request for a listing",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			httpClient, err := newHTTPClient(opts)
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), opts.timeout)
			defer cancel()

			query := url.Values{}
			query.Set("status", "pending")
			query.Set("pageSize", "500")
			var reviews creatorReviewListResponse
			if err := httpClient.GetProductJSONWithQuery(ctx, "/creators/me/marketReviewRequests", query, &reviews); err != nil {
				return fmt.Errorf("list pending review requests: %w", err)
			}

			listingID := strings.TrimSpace(args[0])
			reviewIDs := make([]string, 0, 1)
			for _, review := range reviews.Items {
				if strings.TrimSpace(review.ListingID) == listingID &&
					strings.EqualFold(strings.TrimSpace(review.Status), "pending") &&
					strings.TrimSpace(review.ID) != "" {
					reviewIDs = append(reviewIDs, strings.TrimSpace(review.ID))
				}
			}
			sort.Strings(reviewIDs)
			switch len(reviewIDs) {
			case 0:
				return fmt.Errorf("listing %s has no pending review request", listingID)
			case 1:
			default:
				return fmt.Errorf(
					"listing %s has multiple pending review requests: %s; withdraw one explicitly with creator review withdraw <review-request-id>",
					listingID,
					strings.Join(reviewIDs, ", "),
				)
			}

			payload := map[string]any{}
			if strings.TrimSpace(reason) != "" {
				payload["reason"] = strings.TrimSpace(reason)
			}
			path := "/creators/me/marketReviewRequests/" + url.PathEscape(reviewIDs[0]) + ":withdraw"
			var resp map[string]any
			if err := httpClient.PostProductJSON(ctx, path, payload, &resp); err != nil {
				return fmt.Errorf("withdraw pending review request %s: %w", reviewIDs[0], err)
			}
			return writeIndentedJSON(cmd.OutOrStdout(), resp)
		},
	}
	cmd.Flags().StringVar(&reason, "reason", "", "Optional withdrawal reason")
	return cmd
}

func newListingUnlistCmd(opts *rootOptions) *cobra.Command {
	return newMarketSaleStatusCmd(opts, "unlist", "Unlist a Market SkillBot listing", "unlist")
}

func newListingRelistCmd(opts *rootOptions) *cobra.Command {
	return newMarketSaleStatusCmd(opts, "relist", "Restore a previously unlisted Market SkillBot listing", "list")
}
