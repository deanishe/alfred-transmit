//
// Copyright (c) 2016 Dean Jackson <deanishe@deanishe.net>
//
// MIT Licence. See http://opensource.org/licenses/MIT
//
// Created on 2016-11-06
//

package transmit

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	// BookmarksPath is the path to Favorites.xml
	BookmarksPath = filepath.Join(os.Getenv("HOME"),
		"Library/Application Support/Transmit",
		"Favorites/Favorites.xml")
	favourites []*Favourite
)

// Favourite is a Transmit favourite.
type Favourite struct {
	Name       string
	ID         string
	Server     string
	Protocol   string
	Port       uint32
	RemotePath string
}

// Favourites returns all Transmit favourites.
func Favourites() ([]*Favourite, error) {
	if favourites == nil {
		favs, err := parseFavourites(BookmarksPath)
		if err != nil {
			return nil, err
		}
		favourites = favs
	}
	return favourites, nil
}

// ByID returns the favourite with ID.
func ByID(ID string) (*Favourite, error) {
	favs, err := Favourites()
	if err != nil {
		return nil, err
	}
	for _, fav := range favs {
		if fav.ID == ID {
			return fav, nil
		}
	}
	return nil, fmt.Errorf("No favourite found for ID: %s", ID)
}

// xmlDatabase is the top-level XML object.
//
// The following structs model the structure of a Transmit.app
// Favorites.xml file.
type xmlDatabase struct {
	XMLName xml.Name    `xml:"database"`
	Objects []xmlObject `xml:"object"`
}

// xmlObject is the root of each structure in the Favorites file.
type xmlObject struct {
	Type          string            `xml:"type,attr"`
	ID            string            `xml:"id,attr"`
	Attributes    []xmlAttribute    `xml:"attribute"`
	Relationships []xmlRelationship `xml:"relationship"`
}

// Attr returns the value of an attribute.
func (o *xmlObject) Attr(name string) string {
	for _, a := range o.Attributes {
		if a.Key == name {
			return a.Value
		}
	}
	return ""
}

// Relationship returns the idrefs of an relationship.
func (o *xmlObject) Relationship(name string) []string {
	for _, r := range o.Relationships {
		if r.Name == name {
			return strings.Split(r.IDs, " ")
		}
	}
	return []string{}
}

// Relationship is a history->favourite link.
type xmlRelationship struct {
	Name        string `xml:"name,attr"`
	Destination string `xml:"destination,attr"`
	IDs         string `xml:"idrefs,attr"`
}

// Attribute stores the name and value of one of an Object's attributes.
type xmlAttribute struct {
	Key   string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

// parseFavourites reads and parses a Favorites.xml file.
func parseFavourites(p string) ([]*Favourite, error) {

	data, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	db := xmlDatabase{}
	err = xml.Unmarshal(data, &db)
	if err != nil {
		return nil, err
	}

	var faves = []*Favourite{}
	var ignore = map[string]bool{}

	// Find and parse history first
	for _, o := range db.Objects {
		if o.Type == "HISTORYCOLLECTION" {
			for _, id := range o.Relationship("favorites") {
				// log.Printf("Ignoring ID %s", id)
				ignore[id] = true
			}
		}
	}
	for _, o := range db.Objects {
		if o.Type == "FAVORITE" {
			// Ignore objects that are in history
			if ignore[o.ID] == true {
				// log.Printf("Ignoring %s", o.Attr("nickname"))
				continue
			}
			f, err := objectToFavourite(o)
			if err != nil {
				return nil, err
			}
			faves = append(faves, f)
		}
	}
	return faves, nil
}

// Convert a FAVORITE Object into a Bookmark
func objectToFavourite(o xmlObject) (*Favourite, error) {
	port, err := strconv.ParseInt(o.Attr("port"), 10, 32)
	if err != nil {
		return nil, err
	}

	return &Favourite{
		Name:       o.Attr("nickname"),
		ID:         o.Attr("uuidstring"),
		Protocol:   o.Attr("protocol"),
		Server:     o.Attr("server"),
		Port:       uint32(port),
		RemotePath: o.Attr("initialremotepath"),
	}, nil
}
