# Sprint template contract

`tene-workflow document scaffold` generates the canonical documents. Templates retain stable `tene:section:*` markers so headings and free-form content can change without weakening validation.

The required document set is:

- `00-prd/00-prd.md`
- `01-plan/00-plan.md`
- `02-design/00-design.md`
- `03-analysis/00-loop-check.md`
- `04-qa/00-qa-plan.md`
- `05-report/00-report.md`
- `99-archive/archive-manifest.json`

The embedded generator in `internal/document` is the runtime source. Changes to this contract require golden tests and a schema-version review.

