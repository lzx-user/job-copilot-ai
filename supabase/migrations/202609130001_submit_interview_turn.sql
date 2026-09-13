alter table public.interview_messages
  add column if not exists score smallint check (score between 0 and 100),
  add column if not exists feedback text,
  add column if not exists strengths text[],
  add column if not exists improvements text[];

create unique index if not exists interview_messages_one_role_per_round_idx
  on public.interview_messages (session_id, role, round);

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
  affected_rows integer;
begin
  if p_expected_round is null
    or p_expected_round not between 1 and 4
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
    or p_next_question is null
    or char_length(btrim(p_next_question)) not between 1 and 2000 then
    raise exception 'invalid interview turn';
  end if;

  update public.interview_sessions
  set current_round = current_round + 1,
      updated_at = now()
  where id = p_session_id
    and user_id = (select auth.uid())
    and status = 'in_progress'
    and current_round = p_expected_round
    and current_round < max_rounds;

  get diagnostics affected_rows = row_count;
  if affected_rows <> 1 then
    raise exception 'interview turn conflict';
  end if;

  insert into public.interview_messages (
    session_id, user_id, role, round, content,
    score, feedback, strengths, improvements
  ) values (
    p_session_id, (select auth.uid()), 'candidate', p_expected_round, btrim(p_answer),
    p_score, btrim(p_feedback), p_strengths, p_improvements
  );

  insert into public.interview_messages (session_id, user_id, role, round, content)
  values (
    p_session_id, (select auth.uid()), 'interviewer', p_expected_round + 1, btrim(p_next_question)
  );
end;
$$;

revoke all on function public.submit_interview_turn(uuid, smallint, text, smallint, text, text[], text[], text) from public;
grant execute on function public.submit_interview_turn(uuid, smallint, text, smallint, text, text[], text[], text) to authenticated;
