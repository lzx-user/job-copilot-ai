-- 第 5 轮只保存回答与反馈，不生成第 6 题；Session 留待最终报告成功后再标记 completed。
create or replace function public.submit_interview_turn(
  p_session_id uuid,
  p_expected_round smallint,
  p_answer text,
  p_score smallint,
  p_feedback text,
  p_strengths text[],
  p_improvements text[],
  p_next_question text
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
begin
  if p_expected_round is null
    or p_expected_round not between 1 and 5
    or p_answer is null
    or char_length(btrim(p_answer)) not between 1 and 10000
    or p_score is null
    or p_score not between 0 and 100
    or p_feedback is null
    or char_length(btrim(p_feedback)) not between 1 and 2000
    or p_strengths is null
    or p_improvements is null
    or cardinality(p_strengths) not between 1 and 10
    or cardinality(p_improvements) not between 1 and 10
    or exists (
      select 1
      from unnest(p_strengths || p_improvements) as item
      where item is null or char_length(btrim(item)) not between 1 and 500
    )
    or (p_expected_round < 5 and (
      p_next_question is null
      or char_length(btrim(p_next_question)) not between 1 and 2000
    ))
    or (p_expected_round = 5 and p_next_question is not null) then
    raise exception 'invalid interview turn';
  end if;

  select status, current_round, max_rounds
  into session_status, session_round, session_max_rounds
  from public.interview_sessions
  where id = p_session_id
    and user_id = (select auth.uid())
  for update;

  if session_status is null
    or session_status <> 'in_progress'
    or session_round <> p_expected_round
    or session_max_rounds <> 5
    or exists (
      select 1
      from public.interview_messages
      where session_id = p_session_id
        and user_id = (select auth.uid())
        and role = 'candidate'
        and round = p_expected_round
    ) then
    raise exception 'interview turn conflict';
  end if;

  insert into public.interview_messages (
    session_id, user_id, role, round, content,
    score, feedback, strengths, improvements
  ) values (
    p_session_id, (select auth.uid()), 'candidate', p_expected_round, btrim(p_answer),
    p_score, btrim(p_feedback), p_strengths, p_improvements
  );

  if p_expected_round < session_max_rounds then
    insert into public.interview_messages (session_id, user_id, role, round, content)
    values (
      p_session_id, (select auth.uid()), 'interviewer', p_expected_round + 1, btrim(p_next_question)
    );

    update public.interview_sessions
    set current_round = p_expected_round + 1,
        updated_at = now()
    where id = p_session_id
      and user_id = (select auth.uid());
  else
    update public.interview_sessions
    set updated_at = now()
    where id = p_session_id
      and user_id = (select auth.uid());
  end if;
end;
$$;

revoke all on function public.submit_interview_turn(uuid, smallint, text, smallint, text, text[], text[], text) from public;
grant execute on function public.submit_interview_turn(uuid, smallint, text, smallint, text, text[], text[], text) to authenticated;
