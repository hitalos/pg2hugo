# PG2HUGO

This project was created to attempt provide a way to [hugo](https://gohugo.io) generate static sites querying a [postgres](https://www.postgresql.org/) database.

## Building

Execute `make` or `go build ./cmd/pg2hugo` to build the program.

## Preparing

You need two tables (or views) on your database to use with `pg2hugo`.
Follow example below:

### Database views

I had these tables:

```sql
CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    body TEXT,
    date TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    "publishedAt" TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    "expiredAt" TIMESTAMPTZ(0),
    "updatedAt" TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    author VARCHAR,
    tags JSONB,
    draft boolean NOT NULL DEFAULT false
);
CREATE TABLE IF NOT EXISTS attachs (
    id SERIAL PRIMARY KEY,
    filename VARCHAR NOT NULL UNIQUE,
    title VARCHAR,
    "sortPosition" SMALLINT,
    source VARCHAR,
    "updatedAt" TIMESTAMPTZ(0) DEFAULT CURRENT_TIMESTAMP,
    bs BYTEA,
    post_id INT NOT NULL
        REFERENCES posts(id)
        ON UPDATE RESTRICT
        ON DELETE RESTRICT
)
```

I created a view "contents" returning these fields:

```sql
path, title, body, date, publishdate, expirydate, lastmod, author, tags, draft
```

And a view "resources" returning these fields:

```sql
parent, title, params, lastmod, bs, length
```

Example:

```sql
CREATE OR REPLACE VIEW contents AS
    SELECT
        posts.id::VARCHAR AS path,
        title,
        body,
        date,
        "publishedAt" AS publishdate,
        "expiredAt" AS expirydate,
        "updatedAt" AS lastmod,
        author,
        tags,
        draft
    FROM posts
    ORDER BY id;
CREATE OR REPLACE VIEW resources AS
    SELECT
        filename AS src,
        post_id::VARCHAR AS parent,
        title,
        CASE source
            WHEN '' THEN NULL
            ELSE (('{"source": "' || source) || '"}')::JSONB
        END AS params,
        "updatedAt" AS lastmod,
        bs,
        LENGTH(bs) AS length
   FROM attachments
   ORDER BY sortPosition, src;
```

## Configuring

To run `pg2hugo` you need to define some environment variables:

* **`DSN`** - a connection string to postgres
* **`PRELOAD`** - a boolean to set preloading binary content of resources on starting application (optional default false)

If you prefer, copy `env_example` and rename to `.env`. Save this file on same folder that will run `pg2hugo`. You also can use `-p` to set preloading.

## Running

Run `./pg2hugo mountpoint`. Where "mountpoint" is a content folder (or subfolder) of your site project that will be built with [hugo](https://gohugo.io).

## Credits

This project is based on [pgfs](https://github.com/crgimenes/pgfs) and make intensive use of [libfuse](https://github.com/libfuse/libfuse) through of lib [bazil.org/fuse](https://bazil.org/fuse).

Thanks to [@crgimenes](https://github.com/crgimenes) for the idea and incentive.
