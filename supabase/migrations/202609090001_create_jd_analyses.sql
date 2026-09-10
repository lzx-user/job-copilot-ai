create table if not exists public.jd_analyses (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  company_name varchar(100) not null,
  job_title varchar(100) not null,
  jd_content text not null check (char_length(jd_content) between 200 and 8000),
  resume_summary text not null check (char_length(resume_summary) between 1 and 5000),
  skills text[] not null,
  match_score smallint not null check (match_score between 0 and 100),
  job_summary text not null,
  core_requirements text[] not null,
  matched_skills text[] not null,
  missing_skills text[] not null,
  resume_suggestions text[] not null,
  preparation_topics text[] not null,
  greeting_message text not null,
  created_at timestamptz not null default now()
);

create index if not exists jd_analyses_user_created_at_idx
  on public.jd_analyses (user_id, created_at desc);

alter table public.jd_analyses enable row level security;

-- GRANT 决定角色能否执行操作，RLS Policy 再决定角色能操作哪些行。
revoke all on table public.jd_analyses from anon, authenticated;
grant select, insert on table public.jd_analyses to authenticated;

create policy "Users can read their own JD analyses"
  on public.jd_analyses
  for select
  to authenticated
  using ((select auth.uid()) = user_id);

create policy "Users can insert their own JD analyses"
  on public.jd_analyses
  for insert
  to authenticated
  with check ((select auth.uid()) = user_id);
