# TODO's

## Tasks

- create datamodel for timeentry

  - requires day entry
  - time hh:mm
  - tic key
  - comment
  - story number optional

- add timeentry
- track time
- edit timeentry
- delete timeentry

## CLI UX improvemets

- design better control flow

  - dont make user restart from scratch
  - handle errors without necessarily crashing
  - fmt.Errorf and wrap errors with %w

- introduce a logging framework

- make -listDays not scuffed

- visual presentation when editing
  - navigate with arrowkeys/hjkl
  - better user experience
  - use std.out?
