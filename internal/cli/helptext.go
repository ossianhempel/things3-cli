package cli

const rootHelp = `Usage: things [GLOBAL OPTIONS...] <COMMAND> [ARGS...]

NAME
  things - manage todos with Things 3

SYNOPSIS
  things [GLOBAL OPTIONS...] <COMMAND> [ARGS...]

DESCRIPTION
  CLI interface for Things 3 by Cultured Code.

COMMANDS
  add            - add new todo
  update         - update exiting todo
  delete         - delete an existing todo
  undo           - undo the last bulk action
  add-area       - add new area
  add-project    - add new project
  update-area    - update exiting area
  delete-area    - delete an area
  update-project - update exiting project
  delete-project - delete an existing project
  show           - show an area, project, tag, or todo from the Things database
  search         - search tasks in the Things database
  inbox          - list inbox tasks from the Things database
  today          - list today tasks from the Things database
  upcoming       - list upcoming tasks from the Things database
  anytime        - list anytime tasks from the Things database
  someday        - list someday tasks from the Things database
  logbook        - list logbook tasks from the Things database
  logtoday       - list tasks completed today from the Things database
  createdtoday   - list tasks created today from the Things database
  completed      - list completed tasks from the Things database
  canceled       - list canceled tasks from the Things database
  trash          - list trashed tasks from the Things database
  deadlines      - list tasks with deadlines from the Things database
  all            - list key sections from the Things database
  projects       - list projects from the Things database
  areas          - list areas from the Things database
  tags           - list tags from the Things database
  tasks          - list todos from the Things database
  auth           - show Things auth token status and setup help
  help           - show documentation for the given command

GLOBAL OPTIONS
  -V, --version
    Print version information about things3-cli and Things.

  --debug
    Enable debug mode for things3-cli.

  --foreground
    Open Things in the foreground.

  --dry-run
    Print the Things URL without opening it.

AUTHOR
  Ossian Hempel

REPORTING BUGS
  Issues can be reported on GitHub:

  https://github.com/ossianhempel/things3-cli/issues

LICENSE
  MIT License

SEE ALSO
  https://culturedcode.com/things/support/articles/2803573/
`

const authHelp = `Usage: things auth

NAME
  things auth - show Things auth token status and setup help

SYNOPSIS
  things auth

DESCRIPTION
  Prints whether {{BT}}THINGS_AUTH_TOKEN{{BT}} is set. If missing, prints setup
  steps for the Things URL scheme authorization token.

NOTES
  Token setup:
    1. Open Things 3.
    2. Settings -> General -> Things URLs.
    3. Copy the token (or enable "Allow 'things' CLI to access Things").
    4. export THINGS_AUTH_TOKEN=your_token_here

  Tip: add the export to your shell profile (e.g. ~/.zshrc) to persist it.
`

const projectsHelp = `Usage: things projects [OPTIONS...]

NAME
  things projects - list projects from the Things database

SYNOPSIS
  things projects [OPTIONS...]

DESCRIPTION
  Lists projects from the local Things database (read-only). By default
  only incomplete, non-trashed projects are shown.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --area=AREA
    Filter by area title or ID.

  --include-trashed
    Include trashed projects.

  --all
    Include completed, canceled, and trashed projects.

  --json
    Output JSON.

  --no-header
    Suppress the header row.

  --recursive
    Include nested headings/todos.

  --only-projects
    Only include projects. Implies --recursive.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const areasHelp = `Usage: things areas [OPTIONS...]

NAME
  things areas - list areas from the Things database

SYNOPSIS
  things areas [OPTIONS...]

DESCRIPTION
  Lists areas from the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --json
    Output JSON.

  --no-header
    Suppress the header row.

  --recursive
    Include nested projects/headings/todos.

  --only-projects
    Only include areas and projects. Implies --recursive.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const tagsHelp = `Usage: things tags [OPTIONS...]

NAME
  things tags - list tags from the Things database

SYNOPSIS
  things tags [OPTIONS...]

DESCRIPTION
  Lists tags from the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --json
    Output JSON.

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const tasksHelp = `Usage: things tasks [OPTIONS...]

NAME
  things tasks - list todos from the Things database

SYNOPSIS
  things tasks [OPTIONS...]

DESCRIPTION
  Lists todos from the local Things database (read-only). By default only
  incomplete, non-trashed tasks are shown.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const todayHelp = `Usage: things today [OPTIONS...]

NAME
  things today - list today tasks from the Things database

SYNOPSIS
  things today [OPTIONS...]

DESCRIPTION
  Lists tasks that should appear in Today using the local Things database.
  This mirrors the Things logic for today (including predicted items).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const inboxHelp = `Usage: things inbox [OPTIONS...]

NAME
  things inbox - list inbox tasks from the Things database

SYNOPSIS
  things inbox [OPTIONS...]

DESCRIPTION
  Lists tasks that are in the Inbox list (unfiled) using the local Things
  database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const upcomingHelp = `Usage: things upcoming [OPTIONS...]

NAME
  things upcoming - list upcoming tasks from the Things database

SYNOPSIS
  things upcoming [OPTIONS...]

DESCRIPTION
  Lists tasks scheduled in the future using the local Things database
  (read-only). Tasks with only deadlines are not included.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const anytimeHelp = `Usage: things anytime [OPTIONS...]

NAME
  things anytime - list anytime tasks from the Things database

SYNOPSIS
  things anytime [OPTIONS...]

DESCRIPTION
  Lists tasks in Anytime using the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const somedayHelp = `Usage: things someday [OPTIONS...]

NAME
  things someday - list someday tasks from the Things database

SYNOPSIS
  things someday [OPTIONS...]

DESCRIPTION
  Lists tasks in Someday using the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const logbookHelp = `Usage: things logbook [OPTIONS...]

NAME
  things logbook - list logbook tasks from the Things database

SYNOPSIS
  things logbook [OPTIONS...]

DESCRIPTION
  Lists completed and canceled tasks from the local Things database
  (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: any.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const logtodayHelp = `Usage: things logtoday [OPTIONS...]

NAME
  things logtoday - list tasks completed today from the Things database

SYNOPSIS
  things logtoday [OPTIONS...]

DESCRIPTION
  Lists tasks completed or canceled today using the local Things database
  (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: any.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const createdTodayHelp = `Usage: things createdtoday [OPTIONS...]

NAME
  things createdtoday - list tasks created today from the Things database

SYNOPSIS
  things createdtoday [OPTIONS...]

DESCRIPTION
  Lists tasks created today using the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: any.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const completedHelp = `Usage: things completed [OPTIONS...]

NAME
  things completed - list completed tasks from the Things database

SYNOPSIS
  things completed [OPTIONS...]

DESCRIPTION
  Lists completed tasks from the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: completed.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const canceledHelp = `Usage: things canceled [OPTIONS...]

NAME
  things canceled - list canceled tasks from the Things database

SYNOPSIS
  things canceled [OPTIONS...]

DESCRIPTION
  Lists canceled tasks from the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: canceled.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const trashHelp = `Usage: things trash [OPTIONS...]

NAME
  things trash - list trashed tasks from the Things database

SYNOPSIS
  things trash [OPTIONS...]

DESCRIPTION
  Lists trashed tasks from the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: any.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const deadlinesHelp = `Usage: things deadlines [OPTIONS...]

NAME
  things deadlines - list tasks with deadlines from the Things database

SYNOPSIS
  things deadlines [OPTIONS...]

DESCRIPTION
  Lists tasks with deadlines using the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --search=TEXT
    Case-insensitive substring match on title or notes.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.


  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const allHelp = `Usage: things all [OPTIONS...]

NAME
  things all - list key sections from the Things database

SYNOPSIS
  things all [OPTIONS...]

DESCRIPTION
  Lists Inbox, Today, Upcoming, Anytime, Someday, Logbook, No Area, and Areas
  sections using the local Things database (read-only).

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --recursive
    Include checklist items in JSON output.


  --json
    Output JSON.

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.
`

const addHelp = `Usage: things add [OPTIONS...] [--] [-|TITLE]

NAME
  things add - add new todo

SYNOPSIS
  things add [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Adds new todos to Things.

  If {{BT}}-{{BT}} is given as a title, it is read from STDIN. When titles have
  multiple lines of text, the first is set as the todo's title and the
  remaining lines are set as the todo's notes. Notes set this way take
  precedence over the {{BT}}--notes={{BT}} option.

OPTIONS
  --canceled, --cancelled
    Whether or not the todo should be set to canceled. Default: false. Takes
    priority over completed.

  --notes
    The text to use for the notes field of the todo. Maximum unencoded
    length: 10,000 characters.

  --show-quick-entry
    Whether or not to show the quick entry dialog (populated with the
    provided data) instead of adding a new todo. Ignored if titles is
    specified. Default: false.

  --checklist-item=ITEM
    Checklist item to be added to the todo. Can be specified multiple times
    to create additional checklist items (maximum of 100).

  --completed
    Whether or not the todo should be set to complete. Default: false.
    Ignored if canceled is also set to true.

  --completion-date=<DATE>
    ISO8601 date time string. The date to set as the completion date for the
    todo in the database. Ignored if the todo is not completed or canceled,
    or if the date is in the future.

  --creation-date=DATE
    ISO8601 date time string. The date to set as the creation date for the
    todo in the database. Ignored if the date is in the future.

  --deadline=DATE
    The deadline to apply to the todo.

  --heading=HEADING
    The title of a heading within a project to add to. Ignored if a project
    is not specified, or if the heading doesn't exist.

  --list=LIST
    The title of a project or area to add to. Ignored if list-id is present.

  --list-id=ID
    The ID of a project or area to add to. Takes precedence over list.

  --reveal
    Whether or not to navigate to and show the newly created todo. If
    multiple todos have been created, the first one will be shown. Ignored
    if show-quick-entry is also set to true. Default: false.

  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Does not
    apply a tag if the specified tag doesn't exist.

  --when=DATE|DATETIME
    Possible values: today, tomorrow, evening, anytime, someday, a date
    string, or a date time string. Using a date time string adds a reminder
    for that time. The time component is ignored if anytime or someday is
    specified.

  --titles=TITLE1[,TITLE2,TITLE3...]
    Use instead of title to create multiple todos. Takes priority over title
    and show-quick-entry. The other parameters are applied to all the
    created todos.

  --use-clipboard=VALUE
    Possible values can be replace-title (newlines overflow into notes,
    replacing them), replace-notes, or replace-checklist-items (newlines
    create multiple checklist rows). Takes priority over title, notes, or
    checklist-items.

EXAMPLES
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
`

const addAreaHelp = `Usage: things add-area [OPTIONS...] [-|TITLE]

NAME
  things add-area - add new area

SYNOPSIS
  things add-area [OPTIONS...] [-|TITLE]

DESCRIPTION
  Adds a new area to Things using AppleScript. You may be prompted to grant
  Things automation permission to your terminal.

  If {{BT}}-{{BT}} is given as a title, it is read from STDIN. When titles have
  multiple lines of text, the first is set as the area's title.

OPTIONS
  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Optional.

EXAMPLES
  things add-area "Health"

  things add-area --tags=Personal,Health "Health"

  echo "Area from STDIN" | things add-area -
`

const addProjectHelp = `Usage: things add-project [OPTIONS...] [-|TITLE]

NAME
  things add-project - add new project

SYNOPSIS
  things add-project [OPTIONS...] [-|TITLE]

DESCRIPTION
  Adds a new project to Things.

  If {{BT}}-{{BT}} is given as a title, it is read from STDIN. When titles have
  multiple lines of text, the first is set as the todo's title and the
  remaining lines are set as the todo's notes. Notes set this way take
  precedence over the {{BT}}--notes={{BT}} option.

OPTIONS
  --area-id=AREAID
    The ID of an area to add to. Takes precedence over area. Optional.

  --area=AREA
    The title of an area to add to. Ignored if area-id is present. Optional.

  --canceled, --cancelled
    Whether or not the project should be set to canceled. Default: false.
    Takes priority over completed. Will set all child todos to be canceled.
    Optional.

  --completed
    Whether or not the project should be set to complete. Default: false.
    Ignored if canceled is also set to true. Will set all child todos to be
    completed. Optional.

  --completion-date=DATE
    ISO8601 date time string. The date to set as the completion date for the
    project in the database. If the todos parameter is also specified, this
    date is applied to them, too. Ignored if the todo is not completed or
    canceled, or if the date is in the future. Optional.

  --deadline=DATE
    The deadline to apply to the project. Optional.

  --notes=NOTES
    The text to use for the notes field of the project. Maximum unencoded
    length: 10,000 characters. Optional.

  --reveal
    Whether or not to navigate into the newly created project. Default:
    false. Optional.

  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Does not
    apply a tag if the specified tag doesn't exist. Optional.

  --when=DATE|DATETIME
    Possible values: today, tomorrow, evening, anytime, someday, a date
    string, or a date time string. Using a date time string adds a reminder
    for that time. The time component is ignored if anytime or someday is
    specified. Optional.

  --todo=TITLE
    Title of a todo to add to the project. Can be specified more than once
    to add multiple todos. Optional.

EXAMPLES
  things add-project "Take over the world"
`

const showHelp = `Usage: things show [OPTIONS...] [--] [-|QUERY]

NAME
  things show - show an area, project, tag, or todo from the Things database

SYNOPSIS
  things show --id=ID
  things show [--] <-|QUERY>

DESCRIPTION
  Looks up a single item in the local Things database. If a query is
  provided, it must match exactly (case-insensitive) and return a single
  result. Use {{BT}}things search{{BT}} for partial matching.

  If {{BT}}-{{BT}} is given as a query, it is read from STDIN.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --id=ID
    The ID of an area, project, tag, or todo to show. Takes precedence over
    QUERY. Required if QUERY is not supplied.

  --json
    Output JSON.

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.

EXAMPLES
  things show --id=1234567890AB

  things show "Project One"

  echo "Home" | things show -
`

const searchHelp = `Usage: things search [OPTIONS...] [--] <-|QUERY>

NAME
  things search - search tasks in the Things database

SYNOPSIS
  things search [--] <-|QUERY>

DESCRIPTION
  Searches tasks in the local Things database. The QUERY argument performs a
  case-insensitive substring search on title or notes. Use {{BT}}--query{{BT}}
  for rich queries with boolean ops, fields, and regex.

  Query is required.

  If {{BT}}-{{BT}} is given as a query, it is read from STDIN.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --status=STATUS
    Filter by status: incomplete, completed, canceled, any. Default: incomplete.

  --project=PROJECT
    Filter by project title or ID.

  --area=AREA
    Filter by area title or ID.

  --tag=TAG
    Filter by tag title or ID.

  --query=QUERY
    Rich query with boolean ops, fields, and regex (e.g. title:/regex/ AND tag:work).

  --limit=N
    Limit number of results (0 = no limit). Default: 200.

  --offset=N
    Offset results for pagination.

  --created-after=DATE
    Filter tasks created after (YYYY-MM-DD or RFC3339).

  --created-before=DATE
    Filter tasks created before (YYYY-MM-DD or RFC3339).

  --modified-after=DATE
    Filter tasks modified after (YYYY-MM-DD or RFC3339).

  --modified-before=DATE
    Filter tasks modified before (YYYY-MM-DD or RFC3339).

  --due-before=DATE
    Filter tasks due before (YYYY-MM-DD).

  --start-before=DATE
    Filter tasks starting before (YYYY-MM-DD).

  --has-url
    Filter tasks with URLs in notes.

  --sort=FIELDS
    Sort by fields (e.g. created,-deadline,title).

  --recursive
    Include checklist items in JSON output.

  --include-trashed
    Include trashed tasks.

  --all
    Include completed, canceled, and trashed tasks.

  --format=FORMAT
    Output format: table, json, jsonl, csv.

  --select=FIELDS
    Select fields (comma-separated).

  --json
    Output JSON (alias for --format json).

  --no-header
    Suppress the header row.

NOTES
  The database lives in the Things app sandbox. You may need to grant your
  terminal Full Disk Access to read it.

EXAMPLES
  things search "Work"

  echo "Home" | things search -
`

const updateHelp = `Usage: things update [OPTIONS...] [--] [-|TITLE]

NAME
  things update - updated existing todos

SYNOPSIS
  things update [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Updates an existing todo identified by {{BT}}--id={{BT}}.
  If query filters are provided instead of {{BT}}--id{{BT}}, updates all
  matching todos. Use {{BT}}--dry-run{{BT}} to preview and {{BT}}--yes{{BT}}
  to confirm bulk updates.

  If {{BT}}-{{BT}} is given as a title, it is read from STDIN. When titles have
  multiple lines of text, the first is set as the todo's title and the
  remaining lines are set as the todo's notes. Notes set this way take
  precedence over the {{BT}}--notes={{BT}} option.

  Scheduling note: use {{BT}}--when=someday{{BT}} for Someday, or
  {{BT}}--later{{BT}} for This Evening.

AUTHORIZATION
  Update commands require a Things URL scheme token. Run {{BT}}things auth{{BT}}
  for setup, set {{BT}}THINGS_AUTH_TOKEN{{BT}}, or pass {{BT}}--auth-token{{BT}}.

  Token setup:
    1. Open Things 3.
    2. Settings -> General -> Things URLs.
    3. Copy the token (or enable "Allow 'things' CLI to access Things").

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --auth-token=TOKEN
    The Things URL scheme authorization token. Required. See below for more
    information on authorization. If not provided, uses THINGS_AUTH_TOKEN.

  --id=ID
    The ID of the todo to update. Required for single updates. Optional when
    using query filters for bulk updates.

  --yes
    Confirm bulk update.

  --notes=NOTES
    The notes of the todo. This will replace the existing notes. Maximum
    unencoded length: 10,000 characters. Optional.

  --prepend-notes=NOTES
    Text to add before the existing notes of a todo. Maximum unencoded
    length: 10,000 characters. Optional.

  --append-notes=NOTES
    Text to add after the existing notes of a todo. Maximum unencoded
    length: 10,000 characters. Optional.

  --when=DATE|DATETIME
    Set the when field of a todo. Possible values: today, tomorrow,
    evening, someday, a date string, or a date time string. Including a time
    adds a reminder for that time. The time component is ignored if someday
    is specified. This field cannot be updated on repeating todo.
    Optional.

  --later
    Move the todo to This Evening (alias for {{BT}}--when=evening{{BT}}).
    Optional.

  --deadline=DATE
    The deadline to apply to the todo. This field cannot be updated on
    repeating todo. Optional.

  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Replaces
    all current tags. Does not apply a tag if the specified tag doesn't
    exist. Optional.

  --add-tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Adds the
    specified tags to a todo. Does not apply a tag if the specified tag
    doesn't exist. Optional.

  --completed
    Complete a todo or set a todo to incomplete. Ignored if canceled is also
    set to true. Setting completed=false on a canceled todo will also mark
    it as incomplete. This field cannot be updated on repeating todos.
    Optional.

  --canceled, --cancelled
    Cancel a todo or set a todo to incomplete. Takes priority over
    completed. Setting canceled=false on a completed todo will also mark it
    as incomplete. This field cannot be updated on repeating todos.

  --reveal
    Whether or not to navigate to and show the updated todo. Default: false.
    Optional.

  --duplicate
    Set to true to duplicate the todo before updating it, leaving the
    original todo untouched. Repeating todo cannot be duplicated. Default:
    false. Optional.

  --completion-date=DATE
    ISO8601 date time string. Set the creation date for the todo in the
    database. Ignored if the date is in the future. Optional.

  --creation-date=DATE
    ISO8601 date time string. Set the completion date for the todo in the
    database. Ignored if the todo is not completed or canceled, or if the
    date is in the future. This field cannot be updated on repeating
    todo. Optional.

  --heading=HEADING
    The title of a heading within a project to move the todo to. Ignored if
    the todo is not in a project with the specified heading. Can be used
    together with list or list-id.

  --list=LIST
    The title of a project or area to move the todo into. Ignored if
    {{BT}}--list-id={{BT}} is present.

  --list-id=LISTID
    The ID of a project or area to move the todo into. Takes precedence
    over {{BT}}--list={{BT}}.

  --checklist-item=ITEM
    Checklist items of the todo (maximum of 100). Will replace all existing
    checklist items. Can be specified multiple times on the command line.

  --prepend-checklist-item=ITEM
    Add checklist items to the front of the list of checklist items in the
    todo (maximum of 100). Can be specified multiple times on the command
    line.

  --append-checklist-item=ITEM
    Add checklist items to the end of the list of checklist items in the
    todo (maximum of 100). Can be specified multiple times on the command
    line.

EXAMPLES
  things update --id=8TN1bbz946oBsRBGiQ2XBN "Updated Title"

  things update --id=8TN1bbz946oBsRBGiQ2XBN "Update todo and add notes

  The first line of text is the note title and the rest of the text is
  notes."

  echo "Create a todo from STDIN" |
    things update --id=8TN1bbz946oBsRBGiQ2XBN -

  things update --id=8TN1bbz946oBsRBGiQ2XBN -
  Another way to update a todo from STDIN

  I can type a long form note here for my todo, then press ctrl-d...
  ^d

  things update --id=8TN1bbz946oBsRBGiQ2XBN --deadline=2020-08-01 \
    "Ship this script"

  things update --id=8TN1bbz946oBsRBGiQ2XBN --when="2020-08-01 12:30:00" \
    "Lunch time"

SEE ALSO
  Authorization: https://culturedcode.com/things/support/articles/2803573/#overview-authorization
`

const deleteHelp = `Usage: things delete [OPTIONS...] [--] [-|TITLE]

NAME
  things delete - delete an existing todo

SYNOPSIS
  things delete [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Deletes todos using AppleScript. Provide {{BT}}--id={{BT}} or a title for a
  single todo, or use query filters (same as {{BT}}things tasks{{BT}}) for bulk
  delete. Use {{BT}}--dry-run{{BT}} to preview matches and {{BT}}--yes{{BT}} to
  confirm bulk actions.

  When running interactively, you will be prompted to confirm single deletes.
  For non-interactive use, pass {{BT}}--confirm={{BT}} with the todo ID or title.

  The todo can be identified by {{BT}}--id={{BT}} or by title from the
  positional argument/STDIN. If {{BT}}-{{BT}} is given as a title, it is read
  from STDIN.

OPTIONS
  --db=PATH
    Path to the Things database. Overrides the THINGSDB environment variable.

  --id=ID
    The ID of the todo to delete. Optional if a title is provided.

  --confirm=VALUE
    Confirm deletion by typing the todo ID or title. Required in non-interactive
    mode. Optional when prompted.

  --yes
    Confirm bulk delete.

EXAMPLES
  things delete --id=ABC123

  things delete "Pay bills"
`

const undoHelp = `Usage: things undo [OPTIONS...]

NAME
  things undo - undo the last bulk action

SYNOPSIS
  things undo [OPTIONS...]

DESCRIPTION
  Replays the last bulk update or trash action recorded by things3-cli.
  Undoing updates requires a Things URL scheme token. Undoing trash recreates
  tasks as new items.

OPTIONS
  --auth-token=TOKEN
    The Things URL scheme authorization token. If not provided, uses
    THINGS_AUTH_TOKEN.

  --yes
    Confirm undo for multiple tasks.
`

const updateAreaHelp = `Usage: things update-area [OPTIONS...] [--] [-|TITLE]

NAME
  things update-area - update an existing area

SYNOPSIS
  things update-area [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Updates an existing area using AppleScript. You may be prompted to grant
  Things automation permission to your terminal.

  The area can be identified by {{BT}}--id={{BT}} or by title from the
  positional argument/STDIN. If {{BT}}-{{BT}} is given as a title, it is read
  from STDIN.

OPTIONS
  --id=ID
    The ID of the area to update. Optional if a title is provided.

  --title=TITLE
    New title for the area. Optional.

  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Replaces
    all current tags. Optional.

  --add-tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Adds the
    specified tags to the area. Optional.

EXAMPLES
  things update-area --id=ABC123 --tags=Home,Chores

  things update-area --add-tags=Focus "Work"

  things update-area --id=ABC123 --title="New Name"
`

const deleteAreaHelp = `Usage: things delete-area [OPTIONS...] [--] [-|TITLE]

NAME
  things delete-area - delete an existing area

SYNOPSIS
  things delete-area [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Deletes an existing area using AppleScript. You may be prompted to grant
  Things automation permission to your terminal.

  When running interactively, you will be prompted to confirm the deletion.
  For non-interactive use, pass {{BT}}--confirm={{BT}} with the area ID or title.

  The area can be identified by {{BT}}--id={{BT}} or by title from the
  positional argument/STDIN. If {{BT}}-{{BT}} is given as a title, it is read
  from STDIN.

OPTIONS
  --id=ID
    The ID of the area to delete. Optional if a title is provided.

  --confirm=VALUE
    Confirm deletion by typing the area ID or title. Required in non-interactive
    mode. Optional when prompted.

EXAMPLES
  things delete-area --id=ABC123

  things delete-area "Work"
`

const deleteProjectHelp = `Usage: things delete-project [OPTIONS...] [--] [-|TITLE]

NAME
  things delete-project - delete an existing project

SYNOPSIS
  things delete-project [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Deletes an existing project using AppleScript. You may be prompted to grant
  Things automation permission to your terminal.

  When running interactively, you will be prompted to confirm the deletion.
  For non-interactive use, pass {{BT}}--confirm={{BT}} with the project ID or title.

  The project can be identified by {{BT}}--id={{BT}} or by title from the
  positional argument/STDIN. If {{BT}}-{{BT}} is given as a title, it is read
  from STDIN.

OPTIONS
  --id=ID
    The ID of the project to delete. Optional if a title is provided.

  --confirm=VALUE
    Confirm deletion by typing the project ID or title. Required in
    non-interactive mode. Optional when prompted.

EXAMPLES
  things delete-project --id=ABC123

  things delete-project "Launch"
`

const updateProjectHelp = `Usage: things update-project [OPTIONS...] [--] [-|TITLE]

NAME
  things update-project - update an existing project

SYNOPSIS
  things update-project [OPTIONS...] [--] [-|TITLE]

DESCRIPTION
  Updates an existing project identified by {{BT}}--id={{BT}}.

  If {{BT}}-{{BT}} is given as a title, it is read from STDIN. When titles have
  multiple lines of text, the first is set as the project's title and the
  remaining lines are set as the project's notes. Notes set this way take
  precedence over the {{BT}}--notes={{BT}} option.

AUTHORIZATION
  Update commands require a Things URL scheme token. Run {{BT}}things auth{{BT}}
  for setup, set {{BT}}THINGS_AUTH_TOKEN{{BT}}, or pass {{BT}}--auth-token{{BT}}.

  Token setup:
    1. Open Things 3.
    2. Settings -> General -> Things URLs.
    3. Copy the token (or enable "Allow 'things' CLI to access Things").

OPTIONS
  --auth-token=TOKEN
    The Things URL scheme authorization token. Required. See below for more
    information on authorization. If not provided, uses THINGS_AUTH_TOKEN.

  --id=ID
    The ID of the project to update. Required.

  --notes=NOTES
    The notes of the project. This will replace the existing notes. Maximum
    unencoded length: 10,000 characters. Optional.

  --prepend-notes=NOTES
    Text to add before the existing notes of a project. Maximum unencoded
    length: 10,000 characters. Optional.

  --append-notes=NOTES
    Text to add after the existing notes of a project. Maximum unencoded
    length: 10,000 characters. Optional.

  --when=DATE|DATETIME
    Set the when field of a project. Possible values: today, tomorrow,
    evening, someday, a date string, or a date time string. Including a time
    adds a reminder for that time. The time component is ignored if someday
    is specified. This field cannot be updated on repeating projects.
    Optional.

  --deadline=DATE
    The deadline to apply to the project. This field cannot be updated on
    repeating projects. Optional.

  --tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Replaces
    all current tags. Does not apply a tag if the specified tag doesn't
    exist. Optional.

  --add-tags=TAG1[,TAG2,TAG3...]
    Comma separated strings corresponding to the titles of tags. Adds the
    specified tags to a project. Does not apply a tag if the specified tag
    doesn't exist. Optional.

  --area=AREA
    The ID of an area to move the project into. Takes precedence over
    {{BT}}--area={{BT}}. Optional.

  --area-id=AREAID
    The title of an area to move the project into. Ignored if {{BT}}--area-id={{BT}}
    is present. Optional.

  --completed
    Complete a project or set a project to incomplete. Ignored if canceled
    is also set to true. Setting to true will be ignored unless all child
    todos are completed or canceled and all child headings archived. Setting
    to false on a canceled project will mark it as incomplete. This field
    cannot be updated on repeating projects. Optional.

  --canceled, --cancelled
    Cancel a project or set a project to incomplete. Takes priority over
    completed. Setting to true will be ignored unless all child todos are
    completed or canceled and all child headings archived. Setting to false
    on a completed project will mark it as incomplete. This field cannot be
    updated on repeating projects. Optional.

  --reveal
    Whether or not to navigate to and show the updated project. Default:
    false. Optional.

  --duplicate
    Set to true to duplicate the project before updating it, leaving the
    original project untouched. Repeating projects cannot be duplicated.
    Default: false. Optional.

  --completion-date=DATE
    ISO8601 date time string. Set the creation date for the project in the
    database. Ignored if the date is in the future. Optional.

  --creation-date=DATE
    ISO8601 date time string. Set the completion date for the project in the
    database. Ignored if the project is not completed or canceled, or if the
    date is in the future. This field cannot be updated on repeating
    projects. Optional.

  --todo=TITLE
    Title of a todo to add to the project. Can be specified more than once
    to add multiple todos. Optional.

EXAMPLES
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

SEE ALSO
  Authorization: https://culturedcode.com/things/support/articles/2803573/#overview-authorization
`
