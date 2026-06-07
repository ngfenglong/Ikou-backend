# Agent Rules

## Safety (Non-Negotiable)

Always ask before:
- `git commit`, `git push`, any `git` command that changes history or remote state
- Any `gh` command (`pr create`, `issue`, `repo`, etc.)
- Any `wrangler` command
- `rm`, `rm -rf`, or any destructive file operation
- Installing dependencies (`npm install`, `pip install`, `brew install`, etc.)

Before running any gated command, present:
- Reason the command is needed
- Exact command
- Expected effect
- Risk level

Never commit directly to `main`, `master`, or `production`. Always use a feature branch.

## Task Tracking

Use beads (`bd`) for all task tracking. No markdown TODO lists, no FIXME comments.

Run `bd ready` to see what to work on next.

## Session End

Run `/end-current-session` to wrap up every session.

Report context window % after each major task. Flag when above 50%.

## Plans & Notes

- Ephemeral AI plans → `history/`
- Permanent documentation → `README.md`
