-- 历史页只允许登录用户删除自己拥有的 JD 与面试会话。
-- 删除 JD 时，已有外键会级联删除关联会话、消息与报告。
grant delete on table public.jd_analyses to authenticated;
grant delete on table public.interview_sessions to authenticated;

drop policy if exists "Users can delete their own JD analyses" on public.jd_analyses;
create policy "Users can delete their own JD analyses"
  on public.jd_analyses
  for delete
  to authenticated
  using ((select auth.uid()) = user_id);

drop policy if exists "Users can delete their own interview sessions" on public.interview_sessions;
create policy "Users can delete their own interview sessions"
  on public.interview_sessions
  for delete
  to authenticated
  using ((select auth.uid()) = user_id);

-- 用数据库侧查询一次返回面试历史摘要，避免客户端拼接跨表数据。
create or replace function public.list_interview_history(p_limit integer default 50)
returns table (
  session_id uuid,
  company_name varchar,
  job_title varchar,
  status text,
  current_round smallint,
  max_rounds smallint,
  overall_score smallint,
  created_at timestamptz,
  updated_at timestamptz,
  completed_at timestamptz
)
language sql
stable
security invoker
set search_path = ''
as $$
  select
    session.id,
    analysis.company_name,
    analysis.job_title,
    session.status,
    session.current_round,
    session.max_rounds,
    report.overall_score,
    session.created_at,
    session.updated_at,
    session.completed_at
  from public.interview_sessions as session
  join public.jd_analyses as analysis
    on analysis.id = session.jd_analysis_id
   and analysis.user_id = session.user_id
  left join public.interview_reports as report
    on report.session_id = session.id
   and report.user_id = session.user_id
  where session.user_id = (select auth.uid())
  order by session.created_at desc
  limit greatest(1, least(coalesce(p_limit, 50), 100));
$$;

revoke all on function public.list_interview_history(integer) from public;
grant execute on function public.list_interview_history(integer) to authenticated;
