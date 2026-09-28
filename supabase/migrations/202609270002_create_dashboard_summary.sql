create or replace function public.get_dashboard_summary()
returns jsonb
language sql
stable
security invoker
set search_path = ''
as $$
  select jsonb_build_object(
    'jdCount', (
      select count(*) from public.jd_analyses
      where user_id = (select auth.uid())
    ),
    'interviewCount', (
      select count(*) from public.interview_sessions
      where user_id = (select auth.uid())
    ),
    'weeklyRecordCount', (
      (select count(*) from public.jd_analyses
       where user_id = (select auth.uid()) and created_at >= now() - interval '7 days')
      +
      (select count(*) from public.interview_sessions
       where user_id = (select auth.uid()) and created_at >= now() - interval '7 days')
    ),
    'profileCompleteness', coalesce((
      select 25 * (
        case when btrim(nickname) <> '' and cardinality(target_roles) > 0
          and cardinality(expected_cities) > 0 and btrim(availability) <> '' then 1 else 0 end
        + case when cardinality(skills) >= 5 then 1 else 0 end
        + case when char_length(btrim(project_summary)) >= 100 then 1 else 0 end
        + case when char_length(btrim(strengths)) >= 50 then 1 else 0 end
      )
      from public.profiles
      where user_id = (select auth.uid())
    ), 0),
    'recentRecords', coalesce((
      select jsonb_agg(
        jsonb_build_object(
          'id', recent.id,
          'kind', recent.kind,
          'companyName', recent.company_name,
          'jobTitle', recent.job_title,
          'description', recent.description,
          'createdAt', recent.created_at
        ) order by recent.created_at desc
      )
      from (
        select * from (
          select
            analysis.id::text as id,
            'jd'::text as kind,
            analysis.company_name,
            analysis.job_title,
            '匹配度 ' || analysis.match_score::text || ' 分' as description,
            analysis.created_at
          from public.jd_analyses as analysis
          where analysis.user_id = (select auth.uid())

          union all

          select
            session.id::text as id,
            'interview'::text as kind,
            analysis.company_name,
            analysis.job_title,
            case
              when session.status = 'completed'
                then '已完成 · 综合 ' || coalesce(report.overall_score::text, '--') || ' 分'
              else '进行中 · 第 ' || session.current_round::text || '/' || session.max_rounds::text || ' 轮'
            end as description,
            session.updated_at as created_at
          from public.interview_sessions as session
          join public.jd_analyses as analysis
            on analysis.id = session.jd_analysis_id and analysis.user_id = session.user_id
          left join public.interview_reports as report
            on report.session_id = session.id and report.user_id = session.user_id
          where session.user_id = (select auth.uid())
        ) as combined_records
        order by created_at desc
        limit 5
      ) as recent
    ), '[]'::jsonb)
  );
$$;

revoke all on function public.get_dashboard_summary() from public;
grant execute on function public.get_dashboard_summary() to authenticated;
