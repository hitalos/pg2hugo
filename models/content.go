package models

import (
	"bytes"
	"context"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

// Content represents a content from database
type Content struct {
	Path        string      `db:"path" yaml:"-"`
	Title       string      `db:"title"`
	Body        string      `db:"body" yaml:"-"`
	Date        time.Time   `db:"date" yaml:"date"`
	PublishDate time.Time   `db:"publishdate" yaml:"publishdate,omitempty"`
	ExpiryDate  *time.Time  `db:"expirydate" yaml:"expirydate,omitempty"`
	LastMod     time.Time   `db:"lastmod" yaml:"lastmod,omitempty"`
	Author      string      `db:"author" yaml:"author,omitempty"`
	Tags        []string    `db:"tags" yaml:",omitempty"`
	Resources   []*Resource `yaml:"resources,omitempty"`
	Draft       bool        `db:"draft" yaml:"draft,omitempty"`
}

func (c *Content) String() string {
	bs := new(bytes.Buffer)
	enc := yaml.NewEncoder(bs)
	if err := enc.Encode(c); err != nil {
		log.Println(err)
		return ""
	}
	enc.Close()
	return "---\n" + bs.String() + "---\n" + c.Body + "\n"
}

// Size calculates size in bytes of this content in markdown format
func (c *Content) Size() uint64 {
	return uint64(len([]byte(c.String())))
}

// ReadAllContents loads all contents from database
func ReadAllContents() ([]*Content, error) {
	log.Println("Loading contents and metadata of resources")
	rows, err := db.Query(context.Background(), queryReadAllContents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contents := []*Content{}
	for rows.Next() {
		content := Content{}
		err = rows.Scan(
			&content.Path,
			&content.Title,
			&content.Body,
			&content.Date,
			&content.PublishDate,
			&content.ExpiryDate,
			&content.LastMod,
			&content.Author,
			&content.Tags,
			&content.Draft)
		if err != nil {
			return nil, err
		}
		contents = append(contents, &content)
	}

	resources, err := ReadAllResources()
	if err != nil {
		return nil, err
	}

	for i := range contents {
		contents[i].Resources = []*Resource{}
		for _, r := range resources {
			if contents[i].Path == r.Parent {
				contents[i].Resources = append(contents[i].Resources, r)
			}
		}
	}
	return contents, nil
}
