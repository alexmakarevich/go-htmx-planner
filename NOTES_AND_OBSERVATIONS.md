# Notes & Observations

## HTMX

Having to use separate JS snippets, not integrated with the other code was quite awkward

## SQLC

* No support of bulk inserts
  https://github.com/sqlc-dev/sqlc/discussions/2385
* Awkward 1-n joins (the "1" side is repeated N times)

## Svelte

* doesn't have native binding of HTNL datetime inputs:
  https://github.com/sveltejs/svelte/issues/3937