-- 第一题消息与会话进入第 1 轮必须在同一事务内完成，避免只保存其中一项。
create or replace function public.start_interview_session(
  p_session_id uuid,
  p_question text
)
returns void
language plpgsql
security invoker
set search_path = public
as $$
declare
  affected_rows integer;
begin
  update public.interview_sessions
  set status = 'in_progress',
      current_round = 1,
      updated_at = now()
  where id = p_session_id
    and user_id = (select auth.uid())
    and status = 'pending'
    and current_round = 0;

  get diagnostics affected_rows = row_count;
  if affected_rows <> 1 then
    raise exception 'interview session cannot be started';
  end if;

  insert into public.interview_messages (
    session_id,
    user_id,
    role,
    round,
    content
  ) values (
    p_session_id,
    (select auth.uid()),
    'interviewer',
    1,
    btrim(p_question)
  );
end;
$$;

revoke all on function public.start_interview_session(uuid, text) from public;
grant execute on function public.start_interview_session(uuid, text) to authenticated;
