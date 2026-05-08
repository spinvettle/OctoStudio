

create table channel(
    id bigint generated always as identity primary key,
    name text not null ,
    status boolean default false,
    base_url text
);

create index idx_channel_name on channel(name);
create table channel_key(
    id  bigint generated always as identity primary key,
    channel_id bigint,
    name text,
    metadata jsonb,
    api_key text,
    status text
);

create index idx_channel_key_name on channel_key(name);
