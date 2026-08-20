# CLI Quick Reference

```text
tene-workflow init --name NAME --profile standard
tene-workflow sprint create --title TITLE [--slug SLUG] [--max-iterations N]
tene-workflow sprint start ID
tene-workflow status --json
tene-workflow route --text TEXT [--phase PHASE] [--active auto|true|false]
tene-workflow phase transition PHASE [--dry-run] [--approval ID]
tene-workflow approval request|approve|list
tene-workflow intent capture --statement TEXT --ac TEXT --observable TEXT
tene-workflow intent confirm ID
tene-workflow task add --title TEXT --layer LAYER [--ac IDS]
tene-workflow task link ID [--ac IDS] [--intent IDS] [--replace]
tene-workflow document scaffold|validate
tene-workflow graph providers|build|understand [--changed|--path CSV]|trace ID|impact ID [--depth N] [--call-depth N]|validate
tene-workflow context build [--phase PHASE] [--budget BYTES] [--output PATH]
tene-workflow context validate --input PATH
tene-workflow loop check|iterate|record-gap|resolve-gap|defer-gap
tene-workflow waiver create --gap ID --reason TEXT --approver ID --expires RFC3339|list|revoke ID
tene-workflow qa capabilities|plan|execute CASE --adapter NAME|observe CASE --input FILE|case|evaluate|status
tene-workflow evidence register|verify|list
tene-workflow report generate|validate
tene-workflow migrate status|dry-run|apply
tene-workflow doctor [--repair]|compact|clear
tene-workflow secret check|list ENV|run ENV -- COMMAND
```

Global options can appear anywhere: `--root`, `--json`, `--expected-revision`, `--request-id`.

Exit codes: 0 success; 2 validation/usage; 3 guard or QA failure; 4 conflict/lock; 5 missing capability; 6 security; 7 corruption/migration; 8 child tool failure; 10 internal.

`compact` creates a full, hash-chained projection checkpoint and snapshot. Later events contain compact merge patches. `doctor` compares journal replay with project, active, and master-plan projections; `doctor --repair` backs up divergent files and rebuilds them from replay. It never repairs a corrupt journal.
