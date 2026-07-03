/*
  (C) 2022, 2023 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goat

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// AsyncTagResult is the async result type for probe tag listings
type AsyncTagResult struct {
	Tag Tag
	Err error
}

// tagListingPage is the pagination envelope returned by probes/tags/
type tagListingPage struct {
	Count uint   `json:"count"`
	Next  string `json:"next"`
	Tags  []Tag  `json:"results"`
}

// ProbeTagFilter holds filter parameters for the probes/tags/ endpoint
type ProbeTagFilter struct {
	params  url.Values
	limit   uint
	verbose bool
}

// NewProbeTagFilter creates a new ProbeTagFilter with defaults
func NewProbeTagFilter() *ProbeTagFilter {
	filter := ProbeTagFilter{}
	filter.params = url.Values{}
	return &filter
}

// Verbose sets verbose output for API calls
func (filter *ProbeTagFilter) Verbose(verbose bool) {
	filter.verbose = verbose
}

// Limit sets the maximum number of results to return
func (filter *ProbeTagFilter) Limit(max uint) {
	filter.limit = max
}

// FilterName filters by exact name match
func (filter *ProbeTagFilter) FilterName(name string) {
	filter.params.Set("name", name)
}

// FilterNameContains filters for tags with names containing a substring
func (filter *ProbeTagFilter) FilterNameContains(s string) {
	filter.params.Set("name__contains", s)
}

// FilterNameStartsWith filters for tags with names starting with a prefix
func (filter *ProbeTagFilter) FilterNameStartsWith(s string) {
	filter.params.Set("name__startswith", s)
}

// FilterSlug filters by exact slug match
func (filter *ProbeTagFilter) FilterSlug(slug string) {
	filter.params.Set("slug", slug)
}

// FilterSlugContains filters for tags with slugs containing a substring
func (filter *ProbeTagFilter) FilterSlugContains(s string) {
	filter.params.Set("slug__contains", s)
}

// FilterSlugStartsWith filters for tags with slugs starting with a prefix
func (filter *ProbeTagFilter) FilterSlugStartsWith(s string) {
	filter.params.Set("slug__startswith", s)
}

// verifyFilters checks filter params for obvious errors before making API calls
func (filter *ProbeTagFilter) verifyFilters() error {
	return nil
}

// GetProbeTagCount returns the total number of probe tags matching the filter
func (filter *ProbeTagFilter) GetProbeTagCount() (count uint, err error) {
	err = filter.verifyFilters()
	if err != nil {
		return
	}

	filter.params.Set("page_size", "0")
	query := apiBaseURL + "probes/tags/?" + filter.params.Encode()

	resp, err := apiGetRequest(filter.verbose, query, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var page tagListingPage
	err = json.NewDecoder(resp.Body).Decode(&page)
	if err != nil {
		return 0, err
	}

	return page.Count, nil
}

// GetProbeTags returns probe tags matching the filter via a channel
func (filter *ProbeTagFilter) GetProbeTags(tags chan AsyncTagResult) {
	defer close(tags)

	err := filter.verifyFilters()
	if err != nil {
		tags <- AsyncTagResult{Tag{}, err}
		return
	}

	query := apiBaseURL + "probes/tags/?" + filter.params.Encode()

	resp, err := apiGetRequest(filter.verbose, query, nil)

	// results are paginated with next= (and previous=)
	var total uint = 0
	for {
		if err != nil {
			tags <- AsyncTagResult{Tag{}, err}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			tags <- AsyncTagResult{Tag{}, parseAPIError(resp)}
			return
		}

		var page tagListingPage
		err = json.NewDecoder(resp.Body).Decode(&page)
		if err != nil {
			tags <- AsyncTagResult{Tag{}, err}
			return
		}

		// return items while observing the limit
		for _, tag := range page.Tags {
			tags <- AsyncTagResult{tag, nil}
			total++
			if filter.limit > 0 && total >= filter.limit {
				return
			}
		}

		// no next page => we're done
		if page.Next == "" {
			break
		}

		// just follow the next link
		resp, err = apiGetRequest(filter.verbose, page.Next, nil)
	}
}

// GetProbeTag retrieves a single probe tag by slug
func GetProbeTag(verbose bool, slug string) (*Tag, error) {
	if slug == "" {
		return nil, fmt.Errorf("slug must not be empty")
	}

	query := fmt.Sprintf("%sprobes/tags/%s/", apiBaseURL, slug)

	resp, err := apiGetRequest(verbose, query, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, parseAPIError(resp)
	}

	var tag Tag
	err = json.NewDecoder(resp.Body).Decode(&tag)
	if err != nil {
		return nil, err
	}

	return &tag, nil
}
