// Copyright 2012 Ruben Pollan <meskio@sindominio.net>
// Use of this source code is governed by a LGPL licence
// version 3 or later that can be found in the LICENSE file.

package epubgo

import (
	"errors"
)

// Iterator on the epub navigation index tree.
//
// With it is possible to navigate throw the sections, subsections, ...
// of the epub.
type NavigationIterator struct {
	parents []navCursor
	curr    navCursor
}
type navCursor struct {
	navMap []navpoint
	index  int
}

func newNavigationIterator(navMap []navpoint) *NavigationIterator {
	var nav NavigationIterator
	nav.curr.navMap = navMap
	nav.curr.index = 0
	return &nav
}

// Get the title of the item on the iterator
func (nav NavigationIterator) Title() string {
	return nav.item().Title()
}

// Get the url of the item on the iterator
//
// It usually contains a path and a section link after a '#'.
// The path can be open with epub.OpenFile()
func (nav NavigationIterator) Url() string {
	return nav.item().Url()
}

// Do the item has any children sections?
func (nav NavigationIterator) HasChildren() bool {
	return nav.item().Children() != nil
}

// Do the item has any parent sections?
func (nav NavigationIterator) HasParents() bool {
	return nav.parents != nil
}

// Is the item the first of the sections on the same depth level?
func (nav NavigationIterator) IsFirst() bool {
	return nav.curr.index == 0
}

// Is the item the last of the sections on the same depth level?
func (nav NavigationIterator) IsLast() bool {
	return nav.curr.index == len(nav.curr.navMap)-1
}

// Advance the iterator to the next element on the same depth level
//
// Returns an error if is the last
func (nav *NavigationIterator) Next() error {
	if nav.IsLast() {
		return errors.New("It is the last entry")
	}
	nav.curr.index++
	return nil
}

// Step back the iterator to the previous element on the same depth level
//
// Returns an error if is the first
func (nav *NavigationIterator) Previous() error {
	if nav.IsFirst() {
		return errors.New("It is the first entry")
	}
	nav.curr.index--
	return nil
}

// Move the iterator one level in on depth
//
// Returns an error if it don't has children
func (nav *NavigationIterator) In() error {
	if !nav.HasChildren() {
		return errors.New("It has no children")
	}
	nav.parents = append(nav.parents, nav.curr)
	nav.curr.navMap = nav.item().Children()
	nav.curr.index = 0
	return nil
}

// Move the iterator one level out on depth
//
// Returns an error if it don't has parents
func (nav *NavigationIterator) Out() error {
	if !nav.HasParents() {
		return errors.New("It has no parents")
	}
	nav.curr = nav.parents[len(nav.parents)-1]
	nav.parents = nav.parents[:len(nav.parents)-1]
	return nil
}

func (nav NavigationIterator) item() *navpoint {
	return &nav.curr.navMap[nav.curr.index]
}