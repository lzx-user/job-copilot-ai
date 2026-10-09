-- 将现有“JD 分析 -> 模拟面试”扩展为可追溯的真实求职主线。
-- 所有关联列均允许为空，确保已有分析和面试记录无需伪造新实体即可继续使用。

create table if not exists public.resumes (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  title varchar(100) not null check (char_length(btrim(title)) between 1 and 100),
  is_default boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (id, user_id)
);

create unique index if not exists resumes_one_default_per_user_idx
  on public.resumes (user_id) where is_default;
create index if not exists resumes_user_updated_at_idx
  on public.resumes (user_id, updated_at desc);

create table if not exists public.resume_versions (
  id uuid primary key default gen_random_uuid(),
  resume_id uuid not null references public.resumes(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  version_number integer not null check (version_number > 0),
  content text not null check (char_length(btrim(content)) between 1 and 20000),
  skills text[] not null default '{}',
  project_summary text not null default '',
  note varchar(200) not null default '',
  created_at timestamptz not null default now(),
  unique (resume_id, version_number),
  unique (id, user_id),
  constraint resume_versions_owned_resume_fk
    foreign key (resume_id, user_id) references public.resumes(id, user_id) on delete cascade
);

create index if not exists resume_versions_user_created_at_idx
  on public.resume_versions (user_id, created_at desc);

create table if not exists public.jobs (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  company_name varchar(100) not null check (char_length(btrim(company_name)) between 1 and 100),
  job_title varchar(100) not null check (char_length(btrim(job_title)) between 1 and 100),
  location varchar(100) not null default '',
  source_url text not null default '',
  recruitment_status text not null default 'open'
    check (recruitment_status in ('open', 'closed', 'unknown')),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (id, user_id)
);

create index if not exists jobs_user_updated_at_idx
  on public.jobs (user_id, updated_at desc);

create table if not exists public.job_jd_versions (
  id uuid primary key default gen_random_uuid(),
  job_id uuid not null references public.jobs(id) on delete cascade,
  user_id uuid not null references auth.users(id) on delete cascade,
  version_number integer not null check (version_number > 0),
  jd_content text not null check (char_length(btrim(jd_content)) between 200 and 8000),
  created_at timestamptz not null default now(),
  unique (job_id, version_number),
  unique (id, user_id),
  constraint job_jd_versions_owned_job_fk
    foreign key (job_id, user_id) references public.jobs(id, user_id) on delete cascade
);

create index if not exists job_jd_versions_user_created_at_idx
  on public.job_jd_versions (user_id, created_at desc);

create table if not exists public.applications (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  job_id uuid not null,
  resume_version_id uuid,
  status text not null default 'planned'
    check (status in ('planned', 'applied', 'screening', 'interview', 'offer', 'rejected', 'withdrawn', 'accepted')),
  source varchar(100) not null default '',
  applied_at timestamptz,
  deadline timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (id, user_id),
  constraint applications_owned_job_fk
    foreign key (job_id, user_id) references public.jobs(id, user_id) on delete cascade,
  constraint applications_owned_resume_version_fk
    foreign key (resume_version_id, user_id) references public.resume_versions(id, user_id) on delete restrict
);

create index if not exists applications_user_updated_at_idx
  on public.applications (user_id, updated_at desc);
create index if not exists applications_user_status_idx
  on public.applications (user_id, status);

create table if not exists public.application_events (
  id uuid primary key default gen_random_uuid(),
  application_id uuid not null,
  user_id uuid not null references auth.users(id) on delete cascade,
  event_type text not null check (event_type in (
    'planned', 'applied', 'screening', 'written_test', 'interview_scheduled',
    'interview_completed', 'offer_received', 'rejected', 'withdrawn', 'offer_accepted', 'note'
  )),
  occurred_at timestamptz not null default now(),
  outcome varchar(100) not null default '',
  notes text not null default '' check (char_length(notes) <= 5000),
  created_at timestamptz not null default now(),
  constraint application_events_owned_application_fk
    foreign key (application_id, user_id) references public.applications(id, user_id) on delete cascade
);

create index if not exists application_events_application_occurred_at_idx
  on public.application_events (application_id, occurred_at desc);

create table if not exists public.real_interviews (
  id uuid primary key default gen_random_uuid(),
  application_id uuid not null,
  user_id uuid not null references auth.users(id) on delete cascade,
  round_name varchar(100) not null check (char_length(btrim(round_name)) between 1 and 100),
  scheduled_at timestamptz,
  duration_minutes integer check (duration_minutes is null or duration_minutes between 1 and 1440),
  format varchar(50) not null default '',
  location_or_link text not null default '',
  result text not null default 'scheduled'
    check (result in ('scheduled', 'completed', 'passed', 'failed', 'cancelled', 'pending')),
  notes text not null default '' check (char_length(notes) <= 10000),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (id, user_id),
  constraint real_interviews_owned_application_fk
    foreign key (application_id, user_id) references public.applications(id, user_id) on delete cascade
);

create index if not exists real_interviews_application_scheduled_at_idx
  on public.real_interviews (application_id, scheduled_at);

create table if not exists public.interview_retrospectives (
  interview_id uuid primary key,
  user_id uuid not null references auth.users(id) on delete cascade,
  questions text[] not null default '{}',
  self_assessment text not null default '' check (char_length(self_assessment) <= 10000),
  strengths text[] not null default '{}',
  weaknesses text[] not null default '{}',
  follow_up_actions text[] not null default '{}',
  ai_analysis text not null default '' check (char_length(ai_analysis) <= 10000),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint interview_retrospectives_owned_interview_fk
    foreign key (interview_id, user_id) references public.real_interviews(id, user_id) on delete cascade
);

create table if not exists public.offers (
  application_id uuid primary key,
  user_id uuid not null references auth.users(id) on delete cascade,
  received_at timestamptz not null default now(),
  status text not null default 'pending'
    check (status in ('pending', 'accepted', 'declined', 'expired')),
  deadline timestamptz,
  salary_summary varchar(500) not null default '',
  notes text not null default '' check (char_length(notes) <= 5000),
  decided_at timestamptz,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  constraint offers_owned_application_fk
    foreign key (application_id, user_id) references public.applications(id, user_id) on delete cascade
);

alter table public.jd_analyses
  add column if not exists job_id uuid,
  add column if not exists job_jd_version_id uuid,
  add column if not exists resume_version_id uuid,
  add column if not exists prompt_version varchar(50) not null default 'v1',
  add column if not exists model_name varchar(100) not null default '';

do $migration$
begin
  if not exists (select 1 from pg_constraint where conname = 'jd_analyses_owned_job_fk') then
    alter table public.jd_analyses add constraint jd_analyses_owned_job_fk
      foreign key (job_id, user_id) references public.jobs(id, user_id) on delete restrict;
  end if;
  if not exists (select 1 from pg_constraint where conname = 'jd_analyses_owned_jd_version_fk') then
    alter table public.jd_analyses add constraint jd_analyses_owned_jd_version_fk
      foreign key (job_jd_version_id, user_id) references public.job_jd_versions(id, user_id) on delete restrict;
  end if;
  if not exists (select 1 from pg_constraint where conname = 'jd_analyses_owned_resume_version_fk') then
    alter table public.jd_analyses add constraint jd_analyses_owned_resume_version_fk
      foreign key (resume_version_id, user_id) references public.resume_versions(id, user_id) on delete restrict;
  end if;
end
$migration$;

alter table public.interview_sessions
  add column if not exists application_id uuid,
  add column if not exists resume_version_id uuid,
  add column if not exists interview_type text not null default 'technical'
    check (interview_type in ('technical', 'project', 'comprehensive'));

do $migration$
begin
  if not exists (select 1 from pg_constraint where conname = 'interview_sessions_owned_application_fk') then
    alter table public.interview_sessions add constraint interview_sessions_owned_application_fk
      foreign key (application_id, user_id) references public.applications(id, user_id) on delete restrict;
  end if;
  if not exists (select 1 from pg_constraint where conname = 'interview_sessions_owned_resume_version_fk') then
    alter table public.interview_sessions add constraint interview_sessions_owned_resume_version_fk
      foreign key (resume_version_id, user_id) references public.resume_versions(id, user_id) on delete restrict;
  end if;
end
$migration$;

create or replace function public.create_interview_session_with_context(
  p_analysis_id uuid,
  p_application_id uuid,
  p_resume_version_id uuid,
  p_interview_type text
)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  analysis_job_id uuid;
  analysis_resume_id uuid;
  application_job_id uuid;
  application_resume_id uuid;
  session_id uuid;
begin
  select job_id, resume_version_id into analysis_job_id, analysis_resume_id
  from public.jd_analyses
  where id = p_analysis_id and user_id = (select auth.uid());
  if not found then raise exception 'analysis not found'; end if;

  if p_application_id is not null then
    select job_id, resume_version_id into application_job_id, application_resume_id
    from public.applications
    where id = p_application_id and user_id = (select auth.uid());
    if not found then raise exception 'application not found'; end if;
    if analysis_job_id is not null and application_job_id <> analysis_job_id then
      raise exception 'analysis and application job mismatch';
    end if;
    if p_resume_version_id is not null and application_resume_id is not null
      and p_resume_version_id <> application_resume_id then
      raise exception 'application resume mismatch';
    end if;
  end if;

  insert into public.interview_sessions (
    user_id, jd_analysis_id, application_id, resume_version_id, interview_type
  ) values (
    (select auth.uid()), p_analysis_id, p_application_id,
    coalesce(p_resume_version_id, application_resume_id, analysis_resume_id), p_interview_type
  ) returning id into session_id;
  return session_id;
end;
$$;

-- 客户端只使用登录用户 JWT，RLS 是最终的数据所有权边界。
alter table public.resumes enable row level security;
alter table public.resume_versions enable row level security;
alter table public.jobs enable row level security;
alter table public.job_jd_versions enable row level security;
alter table public.applications enable row level security;
alter table public.application_events enable row level security;
alter table public.real_interviews enable row level security;
alter table public.interview_retrospectives enable row level security;
alter table public.offers enable row level security;

do $$
declare
  table_name text;
begin
  foreach table_name in array array[
    'resumes', 'resume_versions', 'jobs', 'job_jd_versions', 'applications',
    'application_events', 'real_interviews', 'interview_retrospectives', 'offers'
  ] loop
    execute format('revoke all on table public.%I from anon, authenticated', table_name);
    execute format('grant select, insert, update, delete on table public.%I to authenticated', table_name);
    execute format('drop policy if exists %I on public.%I', table_name || '_select_own', table_name);
    execute format('drop policy if exists %I on public.%I', table_name || '_insert_own', table_name);
    execute format('drop policy if exists %I on public.%I', table_name || '_update_own', table_name);
    execute format('drop policy if exists %I on public.%I', table_name || '_delete_own', table_name);
    execute format(
      'create policy %I on public.%I for select to authenticated using ((select auth.uid()) = user_id)',
      table_name || '_select_own', table_name
    );
    execute format(
      'create policy %I on public.%I for insert to authenticated with check ((select auth.uid()) = user_id)',
      table_name || '_insert_own', table_name
    );
    execute format(
      'create policy %I on public.%I for update to authenticated using ((select auth.uid()) = user_id) with check ((select auth.uid()) = user_id)',
      table_name || '_update_own', table_name
    );
    execute format(
      'create policy %I on public.%I for delete to authenticated using ((select auth.uid()) = user_id)',
      table_name || '_delete_own', table_name
    );
  end loop;
end
$$;

create or replace function public.create_resume_with_version(
  p_title text,
  p_content text,
  p_skills text[],
  p_project_summary text,
  p_note text,
  p_is_default boolean default false
)
returns jsonb
language plpgsql
security invoker
set search_path = ''
as $$
declare
  resume_id uuid;
  version_id uuid;
begin
  if char_length(btrim(coalesce(p_title, ''))) not between 1 and 100
    or char_length(btrim(coalesce(p_content, ''))) not between 1 and 20000 then
    raise exception 'invalid resume';
  end if;

  if coalesce(p_is_default, false) then
    update public.resumes set is_default = false, updated_at = now()
    where user_id = (select auth.uid()) and is_default;
  end if;

  insert into public.resumes (user_id, title, is_default)
  values ((select auth.uid()), btrim(p_title), coalesce(p_is_default, false))
  returning id into resume_id;

  insert into public.resume_versions (
    resume_id, user_id, version_number, content, skills, project_summary, note
  ) values (
    resume_id, (select auth.uid()), 1, btrim(p_content), coalesce(p_skills, '{}'),
    btrim(coalesce(p_project_summary, '')), btrim(coalesce(p_note, ''))
  ) returning id into version_id;

  return jsonb_build_object('resumeId', resume_id, 'versionId', version_id);
end;
$$;

create or replace function public.create_job_with_jd(
  p_company_name text,
  p_job_title text,
  p_location text,
  p_source_url text,
  p_jd_content text
)
returns jsonb
language plpgsql
security invoker
set search_path = ''
as $$
declare
  job_id uuid;
  version_id uuid;
begin
  if char_length(btrim(coalesce(p_company_name, ''))) not between 1 and 100
    or char_length(btrim(coalesce(p_job_title, ''))) not between 1 and 100
    or char_length(btrim(coalesce(p_jd_content, ''))) not between 200 and 8000 then
    raise exception 'invalid job';
  end if;

  insert into public.jobs (user_id, company_name, job_title, location, source_url)
  values (
    (select auth.uid()), btrim(p_company_name), btrim(p_job_title),
    btrim(coalesce(p_location, '')), btrim(coalesce(p_source_url, ''))
  ) returning id into job_id;

  insert into public.job_jd_versions (job_id, user_id, version_number, jd_content)
  values (job_id, (select auth.uid()), 1, btrim(p_jd_content))
  returning id into version_id;

  return jsonb_build_object('jobId', job_id, 'jdVersionId', version_id);
end;
$$;

create or replace function public.create_resume_version(
  p_resume_id uuid,
  p_content text,
  p_skills text[],
  p_project_summary text,
  p_note text
)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  next_version integer;
  version_id uuid;
begin
  if char_length(btrim(coalesce(p_content, ''))) not between 1 and 20000 then
    raise exception 'invalid resume version';
  end if;
  perform id from public.resumes
  where id = p_resume_id and user_id = (select auth.uid()) for update;
  if not found then raise exception 'resume not found'; end if;
  select coalesce(max(version_number), 0) + 1 into next_version
  from public.resume_versions where resume_id = p_resume_id;
  insert into public.resume_versions (resume_id, user_id, version_number, content, skills, project_summary, note)
  values (p_resume_id, (select auth.uid()), next_version, btrim(p_content), coalesce(p_skills, '{}'),
    btrim(coalesce(p_project_summary, '')), btrim(coalesce(p_note, '')))
  returning id into version_id;
  update public.resumes set updated_at = now() where id = p_resume_id and user_id = (select auth.uid());
  return version_id;
end;
$$;

create or replace function public.create_job_jd_version(p_job_id uuid, p_jd_content text)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  next_version integer;
  version_id uuid;
begin
  if char_length(btrim(coalesce(p_jd_content, ''))) not between 200 and 8000 then
    raise exception 'invalid JD version';
  end if;
  perform id from public.jobs
  where id = p_job_id and user_id = (select auth.uid()) for update;
  if not found then raise exception 'job not found'; end if;
  select coalesce(max(version_number), 0) + 1 into next_version
  from public.job_jd_versions where job_id = p_job_id;
  insert into public.job_jd_versions (job_id, user_id, version_number, jd_content)
  values (p_job_id, (select auth.uid()), next_version, btrim(p_jd_content))
  returning id into version_id;
  update public.jobs set updated_at = now() where id = p_job_id and user_id = (select auth.uid());
  return version_id;
end;
$$;

create or replace function public.create_application_with_event(
  p_job_id uuid,
  p_resume_version_id uuid,
  p_status text,
  p_source text,
  p_applied_at timestamptz,
  p_deadline timestamptz
)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  application_id uuid;
  initial_event text;
begin
  if p_status not in ('planned', 'applied', 'screening', 'interview', 'offer', 'rejected', 'withdrawn', 'accepted') then
    raise exception 'invalid application';
  end if;

  insert into public.applications (
    user_id, job_id, resume_version_id, status, source, applied_at, deadline
  ) values (
    (select auth.uid()), p_job_id, p_resume_version_id, p_status,
    btrim(coalesce(p_source, '')), p_applied_at, p_deadline
  ) returning id into application_id;

  initial_event := case when p_status = 'planned' then 'planned' else 'applied' end;
  insert into public.application_events (application_id, user_id, event_type, occurred_at)
  values (application_id, (select auth.uid()), initial_event, coalesce(p_applied_at, now()));

  return application_id;
end;
$$;

create or replace function public.add_application_event(
  p_application_id uuid,
  p_event_type text,
  p_occurred_at timestamptz,
  p_outcome text,
  p_notes text
)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  event_id uuid;
  next_status text;
begin
  if p_event_type not in (
    'planned', 'applied', 'screening', 'written_test', 'interview_scheduled',
    'interview_completed', 'offer_received', 'rejected', 'withdrawn', 'offer_accepted', 'note'
  ) then
    raise exception 'invalid application event';
  end if;

  insert into public.application_events (
    application_id, user_id, event_type, occurred_at, outcome, notes
  ) values (
    p_application_id, (select auth.uid()), p_event_type, coalesce(p_occurred_at, now()),
    btrim(coalesce(p_outcome, '')), btrim(coalesce(p_notes, ''))
  ) returning id into event_id;

  next_status := case p_event_type
    when 'planned' then 'planned'
    when 'applied' then 'applied'
    when 'screening' then 'screening'
    when 'written_test' then 'screening'
    when 'interview_scheduled' then 'interview'
    when 'interview_completed' then 'interview'
    when 'offer_received' then 'offer'
    when 'rejected' then 'rejected'
    when 'withdrawn' then 'withdrawn'
    when 'offer_accepted' then 'accepted'
    else null
  end;

  if next_status is not null then
    update public.applications
    set status = next_status, updated_at = now(),
        applied_at = case when p_event_type = 'applied' then coalesce(applied_at, p_occurred_at, now()) else applied_at end
    where id = p_application_id and user_id = (select auth.uid());
  end if;

  return event_id;
end;
$$;

create or replace function public.record_real_interview(
  p_application_id uuid,
  p_round_name text,
  p_scheduled_at timestamptz,
  p_duration_minutes integer,
  p_format text,
  p_location_or_link text,
  p_result text,
  p_notes text
)
returns uuid
language plpgsql
security invoker
set search_path = ''
as $$
declare
  interview_id uuid;
begin
  insert into public.real_interviews (
    application_id, user_id, round_name, scheduled_at, duration_minutes,
    format, location_or_link, result, notes
  ) values (
    p_application_id, (select auth.uid()), btrim(p_round_name), p_scheduled_at,
    p_duration_minutes, btrim(coalesce(p_format, '')), btrim(coalesce(p_location_or_link, '')),
    p_result, btrim(coalesce(p_notes, ''))
  ) returning id into interview_id;

  insert into public.application_events (application_id, user_id, event_type, occurred_at, notes)
  values (
    p_application_id, (select auth.uid()),
    case when p_result = 'scheduled' then 'interview_scheduled' else 'interview_completed' end,
    coalesce(p_scheduled_at, now()), btrim(coalesce(p_round_name, ''))
  );
  update public.applications set status = 'interview', updated_at = now()
  where id = p_application_id and user_id = (select auth.uid());
  return interview_id;
end;
$$;

create or replace function public.save_application_offer(
  p_application_id uuid,
  p_received_at timestamptz,
  p_status text,
  p_deadline timestamptz,
  p_salary_summary text,
  p_notes text,
  p_decided_at timestamptz
)
returns boolean
language plpgsql
security invoker
set search_path = ''
as $$
declare
  is_new boolean;
  next_status text;
begin
  select not exists (
    select 1 from public.offers where application_id = p_application_id
      and user_id = (select auth.uid())
  ) into is_new;

  insert into public.offers (
    application_id, user_id, received_at, status, deadline, salary_summary, notes, decided_at
  ) values (
    p_application_id, (select auth.uid()), coalesce(p_received_at, now()), p_status,
    p_deadline, btrim(coalesce(p_salary_summary, '')), btrim(coalesce(p_notes, '')), p_decided_at
  ) on conflict (application_id) do update set
    received_at = excluded.received_at, status = excluded.status, deadline = excluded.deadline,
    salary_summary = excluded.salary_summary, notes = excluded.notes,
    decided_at = excluded.decided_at, updated_at = now()
  where public.offers.user_id = (select auth.uid());

  if is_new then
    insert into public.application_events (application_id, user_id, event_type, occurred_at, notes)
    values (p_application_id, (select auth.uid()), 'offer_received', coalesce(p_received_at, now()), '记录 Offer');
  end if;
  next_status := case p_status when 'accepted' then 'accepted' when 'declined' then 'withdrawn'
    when 'expired' then 'rejected' else 'offer' end;
  update public.applications set status = next_status, updated_at = now()
  where id = p_application_id and user_id = (select auth.uid());
  return true;
end;
$$;

create or replace function public.update_real_interview_result(
  p_interview_id uuid,
  p_result text,
  p_notes text
)
returns boolean
language plpgsql
security invoker
set search_path = ''
as $$
declare
  target_application_id uuid;
  previous_result text;
begin
  select application_id, result into target_application_id, previous_result
  from public.real_interviews
  where id = p_interview_id and user_id = (select auth.uid()) for update;
  if not found then raise exception 'interview not found'; end if;

  update public.real_interviews
  set result = p_result, notes = btrim(coalesce(p_notes, '')), updated_at = now()
  where id = p_interview_id and user_id = (select auth.uid());

  if p_result in ('completed', 'passed', 'failed') and previous_result not in ('completed', 'passed', 'failed') then
    insert into public.application_events (application_id, user_id, event_type, occurred_at, outcome, notes)
    values (target_application_id, (select auth.uid()), 'interview_completed', now(), p_result, '更新真实面试结果');
  end if;
  return true;
end;
$$;

revoke all on function public.create_resume_with_version(text, text, text[], text, text, boolean) from public;
revoke all on function public.create_interview_session_with_context(uuid, uuid, uuid, text) from public;
revoke all on function public.create_job_with_jd(text, text, text, text, text) from public;
revoke all on function public.create_resume_version(uuid, text, text[], text, text) from public;
revoke all on function public.create_job_jd_version(uuid, text) from public;
revoke all on function public.create_application_with_event(uuid, uuid, text, text, timestamptz, timestamptz) from public;
revoke all on function public.add_application_event(uuid, text, timestamptz, text, text) from public;
revoke all on function public.record_real_interview(uuid, text, timestamptz, integer, text, text, text, text) from public;
revoke all on function public.save_application_offer(uuid, timestamptz, text, timestamptz, text, text, timestamptz) from public;
revoke all on function public.update_real_interview_result(uuid, text, text) from public;
grant execute on function public.create_resume_with_version(text, text, text[], text, text, boolean) to authenticated;
grant execute on function public.create_interview_session_with_context(uuid, uuid, uuid, text) to authenticated;
grant execute on function public.create_job_with_jd(text, text, text, text, text) to authenticated;
grant execute on function public.create_resume_version(uuid, text, text[], text, text) to authenticated;
grant execute on function public.create_job_jd_version(uuid, text) to authenticated;
grant execute on function public.create_application_with_event(uuid, uuid, text, text, timestamptz, timestamptz) to authenticated;
grant execute on function public.add_application_event(uuid, text, timestamptz, text, text) to authenticated;
grant execute on function public.record_real_interview(uuid, text, timestamptz, integer, text, text, text, text) to authenticated;
grant execute on function public.save_application_offer(uuid, timestamptz, text, timestamptz, text, text, timestamptz) to authenticated;
grant execute on function public.update_real_interview_result(uuid, text, text) to authenticated;

-- Dashboard 同时保留原有准备数据，并补充真实投递结果与转化率。
create or replace function public.get_dashboard_summary()
returns jsonb
language sql
stable
security invoker
set search_path = ''
as $$
  with metrics as (
    select
      (select count(*)::int from public.applications
       where user_id = (select auth.uid()) and status <> 'planned') as application_count,
      (select count(*)::int from public.applications
       where user_id = (select auth.uid()) and status in ('applied', 'screening', 'interview', 'offer')) as active_count,
      (select count(*)::int from public.real_interviews
       where user_id = (select auth.uid())) as real_interview_count,
      (select count(*)::int from public.offers
       where user_id = (select auth.uid())) as offer_count
  )
  select jsonb_build_object(
    'jdCount', (select count(*) from public.jd_analyses where user_id = (select auth.uid())),
    'interviewCount', (select count(*) from public.interview_sessions where user_id = (select auth.uid())),
    'applicationCount', metrics.application_count,
    'activeApplicationCount', metrics.active_count,
    'realInterviewCount', metrics.real_interview_count,
    'offerCount', metrics.offer_count,
    'applicationToInterviewRate', case when metrics.application_count = 0 then 0
      else least(100, round(metrics.real_interview_count * 100.0 / metrics.application_count)::int) end,
    'interviewToOfferRate', case when metrics.real_interview_count = 0 then 0
      else least(100, round(metrics.offer_count * 100.0 / metrics.real_interview_count)::int) end,
    'weeklyRecordCount', (
      (select count(*) from public.jd_analyses where user_id = (select auth.uid()) and created_at >= now() - interval '7 days')
      + (select count(*) from public.interview_sessions where user_id = (select auth.uid()) and created_at >= now() - interval '7 days')
      + (select count(*) from public.application_events where user_id = (select auth.uid()) and created_at >= now() - interval '7 days')
    ),
    'profileCompleteness', coalesce((
      select 25 * (
        case when btrim(nickname) <> '' and cardinality(target_roles) > 0
          and cardinality(expected_cities) > 0 and btrim(availability) <> '' then 1 else 0 end
        + case when cardinality(skills) >= 5 then 1 else 0 end
        + case when char_length(btrim(project_summary)) >= 100 then 1 else 0 end
        + case when char_length(btrim(strengths)) >= 50 then 1 else 0 end
      ) from public.profiles where user_id = (select auth.uid())
    ), 0),
    'recentRecords', coalesce((
      select jsonb_agg(jsonb_build_object(
        'id', recent.id, 'kind', recent.kind, 'companyName', recent.company_name,
        'jobTitle', recent.job_title, 'description', recent.description, 'createdAt', recent.created_at
      ) order by recent.created_at desc)
      from (
        select * from (
          select analysis.id::text as id, 'jd'::text as kind, analysis.company_name,
            analysis.job_title, '匹配度 ' || analysis.match_score::text || ' 分' as description,
            analysis.created_at
          from public.jd_analyses analysis where analysis.user_id = (select auth.uid())
          union all
          select session.id::text, 'interview'::text, analysis.company_name, analysis.job_title,
            case when session.status = 'completed'
              then '已完成 · 综合 ' || coalesce(report.overall_score::text, '--') || ' 分'
              else '进行中 · 第 ' || session.current_round::text || '/' || session.max_rounds::text || ' 轮' end,
            session.updated_at
          from public.interview_sessions session
          join public.jd_analyses analysis on analysis.id = session.jd_analysis_id and analysis.user_id = session.user_id
          left join public.interview_reports report on report.session_id = session.id and report.user_id = session.user_id
          where session.user_id = (select auth.uid())
        ) combined_records order by created_at desc limit 5
      ) recent
    ), '[]'::jsonb)
  )
  from metrics;
$$;
