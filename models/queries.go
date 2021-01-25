package models

const (
	// This query must select fields: path, title, body, date, publishdate, expirydate, lastmod, author and tags.
	// Where "tags" is an array of strings ([]varchar or []text).
	queryReadAllContents = `SELECT path, title, body, date, publishdate, expirydate, lastmod, author, tags, draft FROM contents`

	// This query must select fields: src, parent, title, params, lastmod and length.
	// Where "params" is NULL or JSON(B?) format (will be converted to map[string]string)
	// and "length" is the size in bytes of resource.
	queryReadAllResources = `SELECT src, parent, title, params, lastmod, length FROM resources`

	// Same as "queryReadAllResources" added blob field
	queryReadAllResourcesFull = `SELECT src, parent, title, params, lastmod, length, bs FROM resources`

	// This query returns field "bs" that must be an array of bytes (bytea in postgres).
	queryLoadResource = "SELECT bs FROM resources WHERE src = $1"
)
