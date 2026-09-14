-- 面试会话必须绑定当前用户自己的 JD 分析，避免跨用户引用数据。
create unique index if not exists jd_analyses_id_user_id_idx
  on public.jd_analyses (id, user_id);

create table if not exists public.interview_sessions (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references auth.users(id) on delete cascade,
  jd_analysis_id uuid not null,
  status text not null default 'pending'
    check (status in ('pending', 'in_progress', 'completed')),
  current_round smallint not null default 0
    check (current_round between 0 and 5),
  max_rounds smallint not null default 5
    check (max_rounds = 5),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  completed_at timestamptz,
  constraint interview_sessions_owned_analysis_fk
    foreign key (jd_analysis_id, user_id)
    references public.jd_analyses (id, user_id)
    on delete cascade,
  constraint interview_sessions_state_check check (
    (status = 'pending' and current_round = 0 and completed_at is null)
    or (status = 'in_progress' and current_round between 1 and 5 and completed_at is null)
    or (status = 'completed' and current_round = 5 and completed_at is not null)
  )
);

create unique index if not exists interview_sessions_id_user_id_idx
  on public.interview_sessions (id, user_id);

create index if not exists interview_sessions_user_created_at_idx
  on public.interview_sessions (user_id, created_at desc);

create index if not exists interview_sessions_user_status_idx
  on public.interview_sessions (user_id, status);

create table if not exists public.interview_messages (
  id uuid primary key default gen_random_uuid(),
  session_id uuid not null,
  user_id uuid not null references auth.users(id) on delete cascade,
  role text not null check (role in ('interviewer', 'candidate')),
  round smallint not null check (round between 1 and 5),
  content text not null check (char_length(btrim(content)) between 1 and 10000),
  created_at timestamptz not null default now(),
  constraint interview_messages_owned_session_fk
    foreign key (session_id, user_id)
    references public.interview_sessions (id, user_id)
    on delete cascade
);

create index if not exists interview_messages_session_created_at_idx
  on public.interview_messages (session_id, created_at);

alter table public.interview_sessions enable row level security;
alter table public.interview_messages enable row level security;

revoke all on table public.interview_sessions from anon, authenticated;
revoke all on table public.interview_messages from anon, authenticated;
grant select, insert, update on table public.interview_sessions to authenticated;
grant select, insert on table public.interview_messages to authenticated;

create policy "Users can read their own interview sessions"
  on public.interview_sessions
  for select
  to authenticated
  using ((select auth.uid()) = user_id);

create policy "Users can insert their own interview sessions"
  on public.interview_sessions
  for insert
  to authenticated
  with check ((select auth.uid()) = user_id);

create policy "Users can update their own interview sessions"
  on public.interview_sessions
  for update
  to authenticated
  using ((select auth.uid()) = user_id)
  with check ((select auth.uid()) = user_id);

create policy "Users can read their own interview messages"
  on public.interview_messages
  for select
  to authenticated
  using ((select auth.uid()) = user_id);

create policy "Users can insert their own interview messages"
  on public.interview_messages
  for insert
  to authenticated
  with check ((select auth.uid()) = user_id);
