-- migrate:up
create type item_t as enum ('character', 'lightcone');
create type banner_t as enum ('standard', 'beginner', 'character', 'lightcone');

create table warps (
  id text primary key,
  uid text,
  gacha_id text,
  gacha_type text,
  count smallint,
  item_id text,
  item_type item_t,
  name text,
  rarity smallint,
  time_raw text,
  time timestamp with time zone,
  banner_type banner_t,
  region text
);

create index warps_uid_idx on warps (uid);
create index warps_gacha_id_idx on warps (gacha_id);
create index warps_gacha_type_idx on warps (gacha_type);
create index warps_item_type_idx on warps (item_type);
create index warps_name_idx on warps (name); 
create index warps_time_idx on warps (time);
create index warps_rarity_idx on warps (rarity);
create index warps_banner_type_idx on warps (banner_type);
create index warps_region_idx on warps (region);

-- migrate:down

