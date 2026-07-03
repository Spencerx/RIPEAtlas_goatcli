/*
  (C) Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package main

import (
	"fmt"
	"os"

	"github.com/robert-kisteleki/goat"
)

// struct to receive/store command line args for probe tag filtering
type listProbTagsFlags struct {
	filterSlug            string
	filterSlugContains    string
	filterSlugStartsWith  string
	filterName            string
	filterNameContains    string
	filterNameStartsWith  string

	limit uint
	count bool
}

// Implementation of the "probetags" subcommand. Parses command line flags
// and interacts with goatAPI to apply those filters+options to fetch results
func commandListProbTags(args []string) {
	flags := parseListProbTagsArgs(args)

	// single tag retrieval shortcut: exact slug with no other filters set
	if flags.filterSlug != "" &&
		flags.filterSlugContains == "" &&
		flags.filterSlugStartsWith == "" &&
		flags.filterName == "" &&
		flags.filterNameContains == "" &&
		flags.filterNameStartsWith == "" {
		tag, err := goat.GetProbeTag(flagVerbose, flags.filterSlug)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%s\t%s\n", tag.Slug, tag.Name)
		return
	}

	filter := parseListProbTagsFlags(flags)

	// counting only
	if flags.count {
		count, err := filter.GetProbeTagCount()
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(count)
		return
	}

	// most of the work is done by goatAPI
	tags := make(chan goat.AsyncTagResult)
	go filter.GetProbeTags(tags)

	for tag := range tags {
		if tag.Err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", tag.Err)
			os.Exit(1)
		}
		fmt.Printf("%s\t%s\n", tag.Tag.Slug, tag.Tag.Name)
	}
}

// Process flags (filters & options), pass them on to goatAPI
func parseListProbTagsFlags(flags *listProbTagsFlags) *goat.ProbeTagFilter {
	filter := goat.NewProbeTagFilter()
	filter.Verbose(flagVerbose)

	if flags.limit > 0 {
		filter.Limit(flags.limit)
	}

	if flags.filterSlug != "" {
		filter.FilterSlug(flags.filterSlug)
	}
	if flags.filterSlugContains != "" {
		filter.FilterSlugContains(flags.filterSlugContains)
	}
	if flags.filterSlugStartsWith != "" {
		filter.FilterSlugStartsWith(flags.filterSlugStartsWith)
	}
	if flags.filterName != "" {
		filter.FilterName(flags.filterName)
	}
	if flags.filterNameContains != "" {
		filter.FilterNameContains(flags.filterNameContains)
	}
	if flags.filterNameStartsWith != "" {
		filter.FilterNameStartsWith(flags.filterNameStartsWith)
	}

	return filter
}

// Define and parse command line args for this subcommand using the flags package
func parseListProbTagsArgs(args []string) *listProbTagsFlags {
	var flags listProbTagsFlags

	flagsListProbTags.StringVar(&flags.filterSlug, "slug", "", "Filter by exact slug match (retrieves a single tag when used alone)")
	flagsListProbTags.StringVar(&flags.filterSlugContains, "slug-contains", "", "Filter for tags with slugs containing a substring")
	flagsListProbTags.StringVar(&flags.filterSlugStartsWith, "slug-startswith", "", "Filter for tags with slugs starting with a prefix")
	flagsListProbTags.StringVar(&flags.filterName, "name", "", "Filter by exact name match")
	flagsListProbTags.StringVar(&flags.filterNameContains, "name-contains", "", "Filter for tags with names containing a substring")
	flagsListProbTags.StringVar(&flags.filterNameStartsWith, "name-startswith", "", "Filter for tags with names starting with a prefix")

	flagsListProbTags.BoolVar(&flags.count, "count", false, "Count only, don't show the actual results")
	flagsListProbTags.UintVar(&flags.limit, "limit", 0, "Maximum number of tags to retrieve (0 = all)")

	_ = flagsListProbTags.Parse(args)

	return &flags
}
