

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
    status int
);

create index idx_channel_key_name on channel_key(name);


create table model(
    id bigint generated always as identity primary key,
    name text not null ,
    provider text,
    metadata jsonb
)


create table channel_model(
    channel_id bigint,
    model_id bigint,
    primary key (channel_id,model_id)
)
