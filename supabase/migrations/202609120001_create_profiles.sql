-- 该迁移补齐 Profile 的可部署基线，并保持在现有迁移之后执行。
create table if not exists public.profiles (
  user_id uuid primary key references auth.users(id) on delete cascade,
  nickname text not null default '',
  target_roles text[] not null default '{}',
  expected_cities text[] not null default '{}',
  skills text[] not null default '{}',
  project_summary text not null default '',
  strengths text not null default '',
  availability text not null default '',
  graduation_year integer,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint profiles_graduation_year_check
    check (graduation_year is null or graduation_year between 2000 and 2100)
);

alter table public.profiles enable row level security;

revoke all on table public.profiles from anon, authenticated;
grant select, insert, update on public.profiles to authenticated;

do $$
begin
  if not exists (
    select 1 from pg_policies
    where schemaname = 'public' and tablename = 'profiles' and policyname = 'profiles_select_own'
  ) then
    create policy profiles_select_own on public.profiles
      for select to authenticated using ((select auth.uid()) = user_id);
  end if;

  if not exists (
    select 1 from pg_policies
    where schemaname = 'public' and tablename = 'profiles' and policyname = 'profiles_insert_own'
  ) then
    create policy profiles_insert_own on public.profiles
      for insert to authenticated with check ((select auth.uid()) = user_id);
  end if;

  if not exists (
    select 1 from pg_policies
    where schemaname = 'public' and tablename = 'profiles' and policyname = 'profiles_update_own'
  ) then
    create policy profiles_update_own on public.profiles
      for update to authenticated
      using ((select auth.uid()) = user_id)
      with check ((select auth.uid()) = user_id);
  end if;
end
$$;

create or replace function public.set_profiles_updated_at()
returns trigger
language plpgsql
set search_path = ''
as $$
begin
  new.updated_at = now();
  return new;
end;
$$;

drop trigger if exists profiles_set_updated_at on public.profiles;
create trigger profiles_set_updated_at
before update on public.profiles
for each row execute function public.set_profiles_updated_at();
