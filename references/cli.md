# CLI Quick Reference

```text
tene-workflow init --name NAME --profile standard
tene-workflow sprint create --title TITLE [--slug SLUG]
tene-workflow sprint start ID
tene-workflow status --json
tene-workflow phase transition PHASE [--dry-run]
tene-workflow intent capture --statement TEXT --ac TEXT --observable TEXT
tene-workflow intent confirm ID
tene-workflow task add --title TEXT --layer LAYER [--ac IDS]
tene-workflow task link ID [--ac IDS] [--intent IDS] [--replace]
tene-workflow document scaffold|validate
tene-workflow graph providers|build|understand [--changed|--path CSV]|trace ID|validate
tene-workflow context build
tene-workflow loop check|record-gap|resolve-gap
tene-workflow qa capabilities|plan|execute CASE --adapter NAME|observe CASE --input FILE|case|evaluate|status
tene-workflow evidence register|verify|list
tene-workflow report generate|validate
tene-workflow doctor|compact|clear
tene-workflow secret check|list ENV|run ENV -- COMMAND
```

Global options can appear anywhere: `--root`, `--json`, `--expected-revision`, `--request-id`.

Exit codes: 0 success; 2 validation/usage; 3 guard or QA failure; 4 conflict/lock; 5 missing capability; 6 security; 7 corruption/migration; 8 child tool failure; 10 internal.
