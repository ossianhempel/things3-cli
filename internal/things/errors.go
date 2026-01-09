package things

import "errors"

var errMissingShowTarget = errors.New("Error: Must specify --id=ID or query")
var ErrMissingAuthToken = errors.New("Error: Missing Things auth token. Run `things auth` for setup, set THINGS_AUTH_TOKEN, or pass --auth-token=TOKEN (Things > Settings > General > Things URLs).")
var errMissingID = errors.New("Error: Must specify --id=id")
var errMissingTitle = errors.New("Error: Must specify title")
var errMissingAreaTarget = errors.New("Error: Must specify --id=ID or area title")
var errMissingAreaUpdate = errors.New("Error: Must specify --tags, --add-tags, or --title")
var errMissingTodoTarget = errors.New("Error: Must specify --id=ID or todo title")
var errMissingProjectTarget = errors.New("Error: Must specify --id=ID or project title")
