# things 1 "Jul 2021" things "User Manuals"

## SYNOPSIS

`things` [GLOBAL OPTIONS] _COMMAND_ [OPTIONS]

## DESCRIPTION

CLI for Things 3 by Cultured Code (https://culturedcode.com/things/).

https://github.com/ossianhempel/things3-cli#readme

## COMMANDS

*things add*
  Adds new todos to Things.

*things update*
  Update exiting todo

*things delete*
  Delete an existing todo.

*things add-area*
  Add new area.

*things add-project*
  Add new project.

*things update-area*
  Update exiting area.

*things delete-area*
  Delete an existing area.

*things update-project*
  Update exiting project.

*things delete-project*
  Delete an existing project.

*things show*
  Show an area, project, tag, or todo from the Things database.

*things search*
  Search tasks in the Things database.

*things inbox*
  List inbox tasks from the Things database.

*things today*
  List today tasks from the Things database.

*things upcoming*
  List upcoming tasks from the Things database.

*things repeating*
  List repeating tasks from the Things database.

*things anytime*
  List anytime tasks from the Things database.

*things someday*
  List someday tasks from the Things database.

*things logbook*
  List logbook tasks from the Things database.

*things logtoday*
  List tasks completed today from the Things database.

*things createdtoday*
  List tasks created today from the Things database.

*things completed*
  List completed tasks from the Things database.

*things canceled*
  List canceled tasks from the Things database.

*things trash*
  List trashed tasks from the Things database.

*things deadlines*
  List tasks with deadlines from the Things database.

*things all*
  List key sections from the Things database.

*things projects*
  List projects from the Things database.

*things areas*
  List areas from the Things database.

*things tags*
  List tags from the Things database.

*things tasks*
  List todos from the Things database.

*things auth*
  Show Things auth token status and setup help.

*things help [COMMAND]*
  Show documentation for things3-cli and its subcommands.

## GLOBAL OPTIONS

`-V`, `--version`
  Print version information about things3-cli and Things.

`--debug`
  Enable debug mode for things3-cli.

`--foreground`
  Open Things in the foreground.

`--dry-run`
  Print the Things URL without opening it.

## AUTHORIZATION AND DATABASE ACCESS

Ordinary update operations use the Things URL scheme and require an auth token.
Repeat-only updates write directly to the Things database and require writable
database access (normally Full Disk Access), but not a URL token. Commands that
combine ordinary and repeat fields require both.

1. Open Things 3.
2. Settings -> General -> Things URLs.
3. Copy the token (or enable "Allow 'things' CLI to access Things").
4. `export THINGS_AUTH_TOKEN=your_token_here`

Tip: add the export to your shell profile (e.g. `~/.zshrc`) to persist it.
You can run `things auth` to check token status and print these steps.

## things auth

Prints whether `THINGS_AUTH_TOKEN` is set. If missing, prints setup steps for
the Things URL scheme authorization token.

## things add [OPTIONS] [--] [-|TITLE]

Adds new todos to Things.

If `-` is given as a title, it is read from STDIN. When titles have multiple
lines of text, the first is set as the todo's title and the remaining lines
are set as the todo's notes. Notes set this way take precedence over the
`--notes=` option.

Repeating todos are created via the Things database and require a single
explicit title (no `--titles`, `--use-clipboard`, or quick entry).

Use after-completion mode (the default) when the next copy should depend on
completion. Use schedule mode for a fixed calendar cadence. `--repeat-start`
anchors recurrence and is separate from ordinary `--when` scheduling.
`--dry-run` previews repeat writes, and execution re-reads the template before
reporting verified state.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.
  Used for repeat operations.

*--canceled*, *--cancelled*
   Whether or not the todo should be set to canceled. Default: false. Takes
   priority over completed.

*--notes*
  The text to use for the notes field of the todo. Maximum unencoded
  length: 10,000 characters.

*--show-quick-entry*
  Whether or not to show the quick entry dialog (populated with the
  provided data) instead of adding a new todo. Ignored if titles is
  specified. Default: false.

*--checklist-item=ITEM*
  Checklist item to be added to the todo. Can be specified multiple times
  to create additional checklist items (maximum of 100).

*--completed*
  Whether or not the todo should be set to complete. Default: false.
  Ignored if canceled is also set to true.

*--completion-date=<DATE>*
  ISO8601 date time string. The date to set as the completion date for the
  todo in the database. Ignored if the todo is not completed or canceled,
  or if the date is in the future.

*--creation-date=DATE*
  ISO8601 date time string. The date to set as the creation date for the
  todo in the database. Ignored if the date is in the future.

*--deadline=DATE*
  The deadline to apply to the todo.

*--heading=HEADING*
  The title of a heading within a project to add to. Ignored if a project
  is not specified, or if the heading doesn't exist.

*--list=LIST*
  The title of a project or area to add to. Ignored if list-id is present.

*--list-id=ID*
  The ID of a project or area to add to. Takes precedence over list.

*--reveal*
  Whether or not to navigate to and show the newly created todo. If
  multiple todos have been created, the first one will be shown. Ignored
  if show-quick-entry is also set to true. Default: false.

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Does not
  apply a tag if the specified tag doesn't exist.

*--when=DATE|DATETIME*
  Possible values: today, tomorrow, evening, anytime, someday, a date
  string, or a date time string. Using a date time string adds a reminder
  for that time. The time component is ignored if anytime or someday is
  specified.

*--repeat=UNIT*
  Create a repeating template. Units: day, week, month, year.

*--repeat-mode=MODE*
  Repeat mode: after-completion (default) or schedule.

*--repeat-every=N*
  Repeat every N units. Default: 1.

*--repeat-start=DATE*
  Anchor date for the repeat rule (YYYY-MM-DD). Defaults to today.

*--repeat-until=DATE*
  Stop repeating after the given date (YYYY-MM-DD). Optional.

*--repeat-deadline=DAYS*
  Add repeating deadlines; each copy appears in Today DAYS earlier.

*--titles=TITLE1[,TITLE2,TITLE3...]*
  Use instead of title to create multiple todos. Takes priority over title
  and show-quick-entry. The other parameters are applied to all the
  created todos.

*--use-clipboard=VALUE*
  Possible values can be replace-title (newlines overflow into notes,
  replacing them), replace-notes, or replace-checklist-items (newlines
  create multiple checklist rows). Takes priority over title, notes, or
  checklist-items.

**EXAMPLES**

    things add "Finish add to Things script"

    things add "Add a todo with notes
    The first line of text is the note title and the rest of the text is
    notes."

    echo "Create a todo from STDIN" | things add -

    things add -
    Another way to create a todo from STDIN

    I can type a long form note here for my todo, then press ctrl-d...
    ^d

    things add --deadline=2020-08-01 "Ship this script"

    things add --when="2020-08-01 12:30:00" "Lunch time"

    things add --show-quick-entry \
      "Add a pending todo to the quick entry window"

## things add-area [OPTIONS...] [--] [-|TITLE]

Adds a new area using AppleScript. You may be prompted to grant Things
automation permission to your terminal.

If `-` is given as a title, it is read from STDIN. When titles have multiple
lines of text, the first is set as the area's title.

Alias: `create-area`.

**OPTIONS**

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Optional.

**EXAMPLES**

    things add-area "Health"

    things add-area --tags=Personal,Health "Health"

    echo "Area from STDIN" | things add-area -

## things update [OPTIONS...] [--] [-|TITLE]

Updates an existing todo identified by `--id=`.

If `-` is given as a title, it is read from STDIN. When titles have
multiple lines of text, the first is set as the todo's title and the
remaining lines are set as the todo's notes. Notes set this way take
precedence over the `--notes=` option.

Scheduling note: use `--when=someday` for Someday, or `--later` for
This Evening.

Repeat updates require `--id`. Use `--dry-run` to preview the normalized rule
and target. Execution re-reads the template; partial results identify completed
and failed stages and include a trusted UUID when one is known. Repeat-only
updates require writable database access. Ordinary URL updates require the
Things auth token; combined updates require both.

**OPTIONS**

*--auth-token=TOKEN*
   The Things URL scheme authorization token. Required. See below for more
   information on authorization. If not provided, uses THINGS_AUTH_TOKEN.

*--id=ID*
  The ID of the todo to update. Required.

*--notes=NOTES*
  The notes of the todo. This will replace the existing notes. Maximum
  unencoded length: 10,000 characters. Optional.

*--prepend-notes=NOTES*
  Text to add before the existing notes of a todo. Maximum unencoded
  length: 10,000 characters. Optional.

*--append-notes=NOTES*
  Text to add after the existing notes of a todo. Maximum unencoded
  length: 10,000 characters. Optional.

*--when=DATE|DATETIME*
  Set the when field of a todo. Possible values: today, tomorrow,
  evening, someday, a date string, or a date time string. Including a time
  adds a reminder for that time. The time component is ignored if someday
  is specified. An empty string (`--when=""`) clears the scheduled date.
  This field cannot be updated on repeating todo.
  Optional.

*--later*
  Move the todo to This Evening (alias for `--when=evening`). Optional.

*--deadline=DATE*
  The deadline to apply to the todo. An empty string (`--deadline=""`)
  clears the deadline. This field cannot be updated on repeating todo.
  Optional.

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Replaces
  all current tags. An empty string (`--tags=""`) removes all tags. Does
  not apply a tag if the specified tag doesn't exist. Optional.

*--add-tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Adds the
  specified tags to a todo. Does not apply a tag if the specified tag
  doesn't exist. Optional.

*--completed*
  Complete a todo or set a todo to incomplete. Ignored if canceled is also
  set to true. Setting completed=false on a canceled todo will also mark
  it as incomplete. This field cannot be updated on repeating todos.
  Optional.

*--canceled*, *--cancelled*
  Cancel a todo or set a todo to incomplete. Takes priority over
  completed. Setting canceled=false on a completed todo will also mark it
  as incomplete. This field cannot be updated on repeating todos.

*--reveal*
  Whether or not to navigate to and show the updated todo. Default: false.
  Optional.

*--duplicate*
  Set to true to duplicate the todo before updating it, leaving the
  original todo untouched. Repeating todo cannot be duplicated. Default:
  false. Optional.

*--completion-date=DATE*
  ISO8601 date time string. Set the creation date for the todo in the
  database. Ignored if the date is in the future. Optional.

*--creation-date=DATE*
  ISO8601 date time string. Set the completion date for the todo in the
  database. Ignored if the todo is not completed or canceled, or if the
  date is in the future. This field cannot be updated on repeating
  todo. Optional.

*--heading=HEADING*
  The title of a heading within a project to move the todo to. Ignored if
  the todo is not in a project with the specified heading. Can be used
  together with list or list-id.

*--list=LIST*
  The title of a project or area to move the todo into. Ignored if
  `--list-id=` is present.

*--list-id=LISTID*
  The ID of a project or area to move the todo into. Takes precedence
  over `--list=`.

*--checklist-item=ITEM*
  Checklist items of the todo (maximum of 100). Will replace all existing
  checklist items. Can be specified multiple times on the command line.

*--prepend-checklist-item=ITEM*
  Add checklist items to the front of the list of checklist items in the
  todo (maximum of 100). Can be specified multiple times on the command
  line.

*--append-checklist-item=ITEM*
  Add checklist items to the end of the list of checklist items in the
  todo (maximum of 100). Can be specified multiple times on the command
  line.

*--complete-checklist-item=ITEM*
  Mark an existing checklist item complete by exact title. Requires
  `--id` and cannot be combined with other update fields. Can be
  specified multiple times on the command line.

*--incomplete-checklist-item=ITEM*
  Mark an existing checklist item incomplete by exact title. Requires
  `--id` and cannot be combined with other update fields. Can be
  specified multiple times on the command line.

*--repeat=UNIT*
  Set a repeating schedule. Units: day, week, month, year.

*--repeat-mode=MODE*
  Repeat mode: after-completion (default) or schedule.

*--repeat-every=N*
  Repeat every N units. Default: 1.

*--repeat-start=DATE*
  Anchor date for the repeat rule (YYYY-MM-DD). Defaults to today.

*--repeat-until=DATE*
  Stop repeating after the given date (YYYY-MM-DD). Optional.

*--repeat-deadline=DAYS*
  Add repeating deadlines; each copy appears in Today DAYS earlier.

*--repeat-clear*
  Remove the repeating schedule for the todo.

**EXAMPLES**

    things update --id=8TN1bbz946oBsRBGiQ2XBN "Finish add to Things script"

    things update --id=8TN1bbz946oBsRBGiQ2XBN "Add a todo with notes

    The first line of text is the note title and the rest of the text is
    notes."

    echo "Create a todo from STDIN" |
      things update --id=8TN1bbz946oBsRBGiQ2XBN -

    things update --id=8TN1bbz946oBsRBGiQ2XBN -
    Another way to create a todo from STDIN

    I can type a long form note here for my todo, then press ctrl-d...
    ^d

    things update --id=8TN1bbz946oBsRBGiQ2XBN --deadline=2020-08-01 \
      "Ship this script"

    things update --id=8TN1bbz946oBsRBGiQ2XBN --when="2020-08-01 12:30:00" \
      "Lunch time"

**SEE ALSO**

Authorization: https://culturedcode.com/things/support/articles/2803573/#overview-authorization

## things delete [OPTIONS...] [--] [-|TITLE]

Deletes an existing todo using AppleScript. You may be prompted to grant
Things automation permission to your terminal.

When running interactively, you will be prompted to confirm the deletion.
For non-interactive use, pass `--confirm=` with the todo ID or title.

The todo can be identified by `--id=` or by title from the positional
argument/STDIN. If `-` is given as a title, it is read from STDIN.

**OPTIONS**

*--id=ID*
  The ID of the todo to delete. Optional if a title is provided.

*--confirm=VALUE*
  Confirm deletion by typing the todo ID or title. Required in non-interactive
  mode. Optional when prompted.

**EXAMPLES**

    things delete --id=ABC123

    things delete "Pay bills"

## things add-project [OPTIONS...] [-|TITLE]

Adds a new project to Things.

If `-` is given as a title, it is read from STDIN. When titles have
multiple lines of text, the first is set as the todo's title and the
remaining lines are set as the todo's notes. Notes set this way take
precedence over the `--notes=` option.

Repeating projects are not supported. Repeat flags apply only to todos through
`things add` and `things update`.

Alias: `create-project`.

**OPTIONS**

*--area-id=AREAID*
  The ID of an area to add to. Takes precedence over area. Optional.

*--area=AREA*
  The title of an area to add to. Ignored if area-id is present. Optional.

*--canceled*, *--cancelled*
  Whether or not the project should be set to canceled. Default: false.
  Takes priority over completed. Will set all child todos to be canceled.
  Optional.

*--completed*
  Whether or not the project should be set to complete. Default: false.
  Ignored if canceled is also set to true. Will set all child todos to be
  completed. Optional.

*--completion-date=DATE*
  ISO8601 date time string. The date to set as the completion date for the
  project in the database. If the todos parameter is also specified, this
  date is applied to them, too. Ignored if the todo is not completed or
  canceled, or if the date is in the future. Optional.

*--deadline=DATE*
  The deadline to apply to the project. Optional.

*--notes=NOTES*
  The text to use for the notes field of the project. Maximum unencoded
  length: 10,000 characters. Optional.

*--reveal*
  Whether or not to navigate into the newly created project. Default:
  false. Optional.

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Does not
  apply a tag if the specified tag doesn't exist. Optional.

*--when=DATE|DATETIME*
  Possible values: today, tomorrow, evening, anytime, someday, a date
  string, or a date time string. Using a date time string adds a reminder
  for that time. The time component is ignored if anytime or someday is
  specified. Optional.

*--todo=TITLE*
  Title of a todo to add to the project. Can be specified more than once
  to add multiple todos. Optional.

**EXAMPLES**

    things add-project "Take over the world"

## things update-area [OPTIONS...] [--] [-|TITLE]

Updates an existing area using AppleScript. You may be prompted to grant
Things automation permission to your terminal.

The area can be identified by `--id=` or by title from the positional
argument/STDIN. If `-` is given as a title, it is read from STDIN.

**OPTIONS**

*--id=ID*
  The ID of the area to update. Optional if a title is provided.

*--title=TITLE*
  New title for the area. Optional.

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Replaces
  all current tags. Optional.

*--add-tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Adds the
  specified tags to the area. Optional.

**EXAMPLES**

    things update-area --id=ABC123 --tags=Home,Chores

    things update-area --add-tags=Focus "Work"

    things update-area --id=ABC123 --title="New Name"

## things delete-area [OPTIONS...] [--] [-|TITLE]

Deletes an existing area using AppleScript. You may be prompted to grant
Things automation permission to your terminal.

When running interactively, you will be prompted to confirm the deletion.
For non-interactive use, pass `--confirm=` with the area ID or title.

The area can be identified by `--id=` or by title from the positional
argument/STDIN. If `-` is given as a title, it is read from STDIN.

**OPTIONS**

*--id=ID*
  The ID of the area to delete. Optional if a title is provided.

*--confirm=VALUE*
  Confirm deletion by typing the area ID or title. Required in non-interactive
  mode. Optional when prompted.

**EXAMPLES**

    things delete-area --id=ABC123

    things delete-area "Work"

## things update-project [OPTIONS...] [--] [-|TITLE]

Updates an existing project identified by `--id=`.

If `-` is given as a title, it is read from STDIN. When titles have
multiple lines of text, the first is set as the project's title and the
remaining lines are set as the project's notes. Notes set this way take
precedence over the `--notes=` option.

**OPTIONS**

*--auth-token=TOKEN*
   The Things URL scheme authorization token. Required. See below for more
   information on authorization. If not provided, uses THINGS_AUTH_TOKEN.

*--id=ID*
  The ID of the project to update. Required.

*--notes=NOTES*
  The notes of the project. This will replace the existing notes. Maximum
  unencoded length: 10,000 characters. Optional.

*--prepend-notes=NOTES*
  Text to add before the existing notes of a project. Maximum unencoded
  length: 10,000 characters. Optional.

*--append-notes=NOTES*
  Text to add after the existing notes of a project. Maximum unencoded
  length: 10,000 characters. Optional.

*--when=DATE|DATETIME*
  Set the when field of a project. Possible values: today, tomorrow,
  evening, someday, a date string, or a date time string. Including a time
  adds a reminder for that time. The time component is ignored if someday
  is specified. An empty string (`--when=""`) clears the scheduled date.
  Optional.

*--deadline=DATE*
  The deadline to apply to the project. An empty string (`--deadline=""`)
  clears the deadline. Optional.

*--tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Replaces
  all current tags. An empty string (`--tags=""`) removes all tags. Does
  not apply a tag if the specified tag doesn't exist. Optional.

*--add-tags=TAG1[,TAG2,TAG3...]*
  Comma separated strings corresponding to the titles of tags. Adds the
  specified tags to a project. Does not apply a tag if the specified tag
  doesn't exist. Optional.

*--area=AREA*
  The ID of an area to move the project into. Takes precedence over
  `--area=`. Optional.

*--area-id=AREAID*
  The title of an area to move the project into. Ignored if `--area-id=`
  is present. Optional.

*--completed*
  Complete a project or set a project to incomplete. Ignored if canceled
  is also set to true. Setting to true will be ignored unless all child
  todos are completed or canceled and all child headings archived. Setting
  to false on a canceled project will mark it as incomplete. Optional.

*--canceled, --cancelled*
  Cancel a project or set a project to incomplete. Takes priority over
  completed. Setting to true will be ignored unless all child todos are
  completed or canceled and all child headings archived. Setting to false
  on a completed project will mark it as incomplete. Optional.

*--reveal*
  Whether or not to navigate to and show the updated project. Default:
  false. Optional.

*--duplicate*
  Set to true to duplicate the project before updating it, leaving the
  original project untouched. Default: false. Optional.

*--completion-date=DATE*
  ISO8601 date time string. Set the creation date for the project in the
  database. Ignored if the date is in the future. Optional.

*--creation-date=DATE*
  ISO8601 date time string. Set the completion date for the project in the
  database. Ignored if the project is not completed or canceled, or if the
  date is in the future. Optional.

*--todo=TITLE*
  Title of a todo to add to the project. Can be specified more than once
  to add multiple todos. Optional.

**EXAMPLES**

    things update-project --id=8TN1bbz946oBsRBGiQ2XBN "The new project title"

    things update-project --id=8TN1bbz946oBsRBGiQ2XBN "Set Title and add Notes

    The first line of text is the project title and the rest of the text is
    notes."

    echo "Project title from STDIN" |
      things update-project --id=8TN1bbz946oBsRBGiQ2XBN -

    things update-project --id=8TN1bbz946oBsRBGiQ2XBN -
    Another way to set the project title and notes from STDIN

    These are the notes
    ^d

    things update --id=8TN1bbz946oBsRBGiQ2XBN --reveal
      "Ship this project"

**SEE ALSO**

Authorization: https://culturedcode.com/things/support/articles/2803573/#overview-authorization

## things delete-project [OPTIONS...] [--] [-|TITLE]

Deletes an existing project using AppleScript. You may be prompted to grant
Things automation permission to your terminal.

When running interactively, you will be prompted to confirm the deletion.
For non-interactive use, pass `--confirm=` with the project ID or title.

The project can be identified by `--id=` or by title from the positional
argument/STDIN. If `-` is given as a title, it is read from STDIN.

**OPTIONS**

*--id=ID*
  The ID of the project to delete. Optional if a title is provided.

*--confirm=VALUE*
  Confirm deletion by typing the project ID or title. Required in non-interactive
  mode. Optional when prompted.

**EXAMPLES**

    things delete-project --id=ABC123

    things delete-project "Launch"

## things show [OPTIONS...] [--] [-|QUERY]

Looks up a single item in the local Things database. If a query is provided,
it must match exactly (case-insensitive) and return a single result. Use
`things search` for partial matching.

If `-` is given as a query, it is read from STDIN.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--id=ID*
  The ID of an area, project, tag, or todo to show. Takes precedence over
  QUERY. Required if QUERY is not supplied.

*--recursive*
  Include checklist items when showing to-dos.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**EXAMPLES**

    things show --id=1234567890AB

    things show "Project One"

    echo "Home" | things show -

## things search [--] [-|QUERY]

Searches tasks in the local Things database by title or notes. Query is
required.

If `-` is given as a query, it is read from STDIN.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=COUNT*
  Limit the number of results. 0 means no limit. Default: 200.

*--recursive*
  Include checklist items in JSON output.

*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**EXAMPLES**

    things search "Work"

    echo "Home" | things search -

## things projects [OPTIONS...]

Lists projects from the local Things database (read-only). By default only
incomplete, non-trashed projects are shown.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--area=AREA*
  Filter by area title or ID.

*--include-trashed*
  Include trashed projects.

*--all*
  Include completed, canceled, and trashed projects.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

*--recursive*
  Include nested headings/todos.

*--only-projects*
  Only include projects. Implies `--recursive`.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things areas [OPTIONS...]

Lists areas from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

*--recursive*
  Include nested projects/headings/todos.

*--only-projects*
  Only include areas and projects. Implies `--recursive`.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things tags [OPTIONS...]

Lists tags from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things tasks [OPTIONS...]

Lists todos from the local Things database (read-only). By default only
incomplete, non-trashed tasks are shown.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--search=TEXT*
  Case-insensitive substring match on title or notes.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things today [OPTIONS...]

Lists tasks that should appear in Today using the local Things database.
Rows follow the Things Today/This Evening order using the reference date and
Today index. Select `today_index_reference_date` to inspect the raw packed
date value.

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things inbox [OPTIONS...]

Lists inbox tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things upcoming [OPTIONS...]

Lists upcoming tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things repeating [OPTIONS...]

Lists repeating tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things anytime [OPTIONS...]

Lists anytime tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things someday [OPTIONS...]

Lists someday tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things logbook [OPTIONS...]

Lists completed and canceled tasks from the local Things database
(read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: any.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things logtoday [OPTIONS...]

Lists tasks completed or canceled today using the local Things database
(read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: any.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things createdtoday [OPTIONS...]

Lists tasks created today using the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: any.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things completed [OPTIONS...]

Lists completed tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: completed.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things canceled [OPTIONS...]

Lists canceled tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: canceled.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things trash [OPTIONS...]

Lists trashed tasks from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: any.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things deadlines [OPTIONS...]

Lists tasks with deadlines from the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--status=STATUS*
  Filter by status: incomplete, completed, canceled, any. Default: incomplete.

*--project=PROJECT*
  Filter by project title or ID.

*--area=AREA*
  Filter by area title or ID.

*--tag=TAG*
  Filter by tag title or ID.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.


*--include-trashed*
  Include trashed tasks.

*--all*
  Include completed, canceled, and trashed tasks.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things all [OPTIONS...]

Lists Inbox, Today, Upcoming, Repeating, Anytime, Someday, Logbook, No Area, and Areas
sections using the local Things database (read-only).

**OPTIONS**

*--db=PATH*
  Path to the Things database. Overrides the THINGSDB environment variable.

*--limit=N*
  Limit number of results (0 = no limit). Default: 200.

*--recursive*
  Include checklist items in JSON output.

*--json*
  Output JSON.

*--no-header*
  Suppress the header row.

**NOTES**

The database lives in the Things app sandbox. You may need to grant your
terminal Full Disk Access to read it.

## things help [COMMAND]

Prints documentation for things3-cli commands.

## BUG REPORTS

Issues can be reported on GitHub:

<https://github.com/ossianhempel/things3-cli/issues>

## AUTHOR

Joshua Priddle <jpriddle@me.com>

## LICENSE

MIT License

Copyright (c) 2021 Joshua Priddle <jpriddle@me.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

## SEE ALSO

open(1)
