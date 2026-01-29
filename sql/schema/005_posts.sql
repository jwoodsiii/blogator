-- +goose Up
--id - a unique identifier for the post
-- created_at - the time the record was created
-- updated_at - the time the record was last updated
-- title - the title of the post
-- url - the URL of the post (this should be unique)
-- description - the description of the post
-- published_at - the time the post was published
-- feed_id - the ID of the feed that the post came from

create table posts (
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp,
    title text not null,
    url text unique not null,
    description text,
    published_at timestamp,
    feed_id uuid not null references feeds(id) on delete cascade
);

-- +goose Down
drop table posts;
