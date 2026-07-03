/*
  (C) 2022, 2023 Robert Kisteleki & RIPE NCC

  See LICENSE file for the license.
*/

package goat

import (
	"testing"
)

// Test that filter methods set the expected query parameters
func TestProbeTagFilterParams(t *testing.T) {
	var filter *ProbeTagFilter

	filter = NewProbeTagFilter()
	filter.FilterName("home")
	if filter.params.Get("name") != "home" {
		t.Errorf("FilterName did not set name param correctly")
	}

	filter = NewProbeTagFilter()
	filter.FilterNameContains("net")
	if filter.params.Get("name__contains") != "net" {
		t.Errorf("FilterNameContains did not set name__contains param correctly")
	}

	filter = NewProbeTagFilter()
	filter.FilterNameStartsWith("ho")
	if filter.params.Get("name__startswith") != "ho" {
		t.Errorf("FilterNameStartsWith did not set name__startswith param correctly")
	}

	filter = NewProbeTagFilter()
	filter.FilterSlug("home")
	if filter.params.Get("slug") != "home" {
		t.Errorf("FilterSlug did not set slug param correctly")
	}

	filter = NewProbeTagFilter()
	filter.FilterSlugContains("net")
	if filter.params.Get("slug__contains") != "net" {
		t.Errorf("FilterSlugContains did not set slug__contains param correctly")
	}

	filter = NewProbeTagFilter()
	filter.FilterSlugStartsWith("ho")
	if filter.params.Get("slug__startswith") != "ho" {
		t.Errorf("FilterSlugStartsWith did not set slug__startswith param correctly")
	}
}

// Test that GetProbeTag rejects an empty slug before making any API call
func TestGetProbeTagEmptySlug(t *testing.T) {
	_, err := GetProbeTag(false, "")
	if err == nil {
		t.Errorf("GetProbeTag with empty slug should return an error")
	}
}
