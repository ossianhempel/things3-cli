package cli

import "fmt"

const thingsBundleID = "com.culturedcode.ThingsMac"

func ensureThingsLaunched(app *App) {
	if app == nil || app.Launcher == nil {
		return
	}
	if err := app.Launcher.Open("-g", "-b", thingsBundleID); err != nil && app.Debug && app.Err != nil {
		fmt.Fprintf(app.Err, "Note: unable to launch Things in background (%v)\n", err)
	}
}
