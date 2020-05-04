package models

import (
	"context"
	"log"
	"os"
	"time"
)

// Resource represents a file attached to a Content
type Resource struct {
	Src     string            `db:"src" yaml:"src"`
	Parent  string            `db:"parent" yaml:"-"`
	Title   *string           `db:"title" yaml:"title,omitempty"`
	Params  map[string]string `db:"params" yaml:"params,omitempty"`
	Bs      []byte            `db:"bs" yaml:"-"`
	Length  uint64            `db:"length" yaml:"-"`
	LastMod time.Time         `db:"lastmod" yaml:"-"`
}

// ReadAllResources loads all resources from database
func ReadAllResources() ([]*Resource, error) {
	preload := os.Getenv("PRELOAD") == "true"
	resources := []*Resource{}
	query := queryReadAllResources
	if preload {
		log.Println("Preloading binary content of resources")
		query = queryReadAllResourcesFull
	}
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resource := Resource{}
		resource.Params = map[string]string{}
		if preload {
			err = rows.Scan(&resource.Src, &resource.Parent, &resource.Title, &resource.Params, &resource.LastMod, &resource.Length, &resource.Bs)
		} else {
			err = rows.Scan(&resource.Src, &resource.Parent, &resource.Title, &resource.Params, &resource.LastMod, &resource.Length)
		}
		if err != nil {
			return nil, err
		}
		resources = append(resources, &resource)
	}

	return resources, nil
}

// Load loads bytes from an only resource on database
func (r *Resource) Load() error {
	return db.QueryRow(context.Background(), queryLoadResource, r.Src).Scan(&r.Bs)
}

// Size returns size in bytes for this resource
func (r *Resource) Size() uint64 {
	return r.Length
}
