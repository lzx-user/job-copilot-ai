create table if not exists public.interview_reports (
  session_id uuid primary key,
  user_id uuid not null references auth.users(id) on delete cascade,
  overall_score smallint not null check (overall_score between 0 and 100),
  technical_score smallint not null check (technical_score between 0 and 100),
  expression_score smallint not null check (expression_score between 0 and 100),
  project_depth_score smallint not null check (project_depth_score between 0 and 100),
  strengths text[] not null check (cardinality(strengths) between 1 and 10),
  weaknesses text[] not null check (cardinality(weaknesses) between 1 and 10),
  recommended_topics text[] not null check (cardinality(recommended_topics) between 1 and 10),
  answer_tips text[] not null check (cardinality(answer_tips) between 1 and 10),
  summary text not null check (char_length(btrim(summary)) between 1 and 3000),
  created_at timestamptz not null default now(),
  constraint interview_reports_owned_session_fk
    foreign key (session_id, user_id)
    references public.interview_sessions (id, user_id)
    on delete cascade
);

create index if not exists interview_reports_user_created_at_idx
  on public.interview_reports (user_id, created_at desc);

alter table public.interview_reports enable row level security;

revoke all on table public.interview_reports from anon, authenticated;
grant select, insert on table public.interview_reports to authenticated;

create policy "Users can read their own interview reports"
  on public.interview_reports
  for select
  to authenticated
  using ((select auth.uid()) = user_id);

create policy "Users can insert their own interview reports"
  on public.interview_reports
  for insert
  to authenticated
  with check ((select auth.uid()) = user_id);

create or replace function public.complete_interview_report(
  p_session_id uuid,
  p_overall_score smallint,
  p_technical_score smallint,
  p_expression_score smallint,
  p_project_depth_score smallint,
  p_strengths text[],
  p_weaknesses text[],
  p_recommended_topics text[],
  p_answer_tips text[],
  p_summary text
)
returns void
language plpgsql
security invoker
set search_path = ''
as $$
declare
  session_round smallint;
  session_max_rounds smallint;
  session_status text;
  interviewer_count integer;
  candidate_count integer;
begin
  if p_overall_score is null or p_overall_score not between 0 and 100
    or p_technical_score is null or p_technical_score not between 0 and 100
    or p_expression_score is null or p_expression_score not between 0 and 100
    or p_project_depth_score is null or p_project_depth_score not between 0 and 100
    or p_strengths is null or cardinality(p_strengths) not between 1 and 10
    or p_weaknesses is null or cardinality(p_weaknesses) not between 1 and 10
    or p_recommended_topics is null or cardinality(p_recommended_topics) not between 1 and 10
    or p_answer_tips is null or cardinality(p_answer_tips) not between 1 and 10
    or exists (
      select 1
      from unnest(p_strengths || p_weaknesses || p_recommended_topics || p_answer_tips) as item
      where item is null or char_length(btrim(item)) not between 1 and 500
    )
    or p_summary is null or char_length(btrim(p_summary)) not between 1 and 3000 then
    raise exception 'invalid interview report';
  end if;

  select status, current_round, max_rounds
  into session_status, session_round, session_max_rounds
  from public.interview_sessions
  where id = p_session_id
    and user_id = (select auth.uid())
  for update;

  select
    count(*) filter (where role = 'interviewer'),
    count(*) filter (where role = 'candidate' and score is not null and feedback is not null)
  into interviewer_count, candidate_count
  from public.interview_messages
  where session_id = p_session_id
    and user_id = (select auth.uid());

  if session_status is null
    or session_status <> 'in_progress'
    or session_round <> 5
    or session_max_rounds <> 5
    or interviewer_count <> 5
    or candidate_count <> 5
    or exists (
      select 1 from public.interview_reports
      where session_id = p_session_id and user_id = (select auth.uid())
    ) then
    raise exception 'interview report conflict';
  end if;

  insert into public.interview_reports (
    session_id, user_id, overall_score, technical_score, expression_score,
    project_depth_score, strengths, weaknesses, recommended_topics, answer_tips, summary
  ) values (
    p_session_id, (select auth.uid()), p_overall_score, p_technical_score, p_expression_score,
    p_project_depth_score, p_strengths, p_weaknesses, p_recommended_topics, p_answer_tips, btrim(p_summary)
  );

  update public.interview_sessions
  set status = 'completed', completed_at = now(), updated_at = now()
  where id = p_session_id
    and user_id = (select auth.uid());
end;
$$;

revoke all on function public.complete_interview_report(uuid, smallint, smallint, smallint, smallint, text[], text[], text[], text[], text) from public;
grant execute on function public.complete_interview_report(uuid, smallint, smallint, smallint, smallint, text[], text[], text[], text[], text) to authenticated;
