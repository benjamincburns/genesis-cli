// message contains a collection of messages so they can be changed easily
package message

const (
	PreBrowserAuthShowURL = "You will now be taken to your browser for authentication " +
		"or open the url below in a browser."

	FailedToOpenBrowserTab = "Failed to open browser, you MUST do the manual process."

	AuthenticationTimedOut = "authentication timed out and was cancelled"

	FatalErrorMessage = "If you believe this is a bug, please file a bug report"

	MissingOrgID = "please provide your organization"

	FilePassedValidation = "all checks passed"

	UnknownFormat = "test definition is in an unknown format"

	NoSchema = "unable to retrieve the schema"

	WindowsUpdateAvailable = `There is an update availible, please download the latest version from
https://assets.whiteblock.io/cli/master/bin/windows/amd64/genesis.exe`

	UpdateAvailable = `There is an update availible, please run 'genesis update' to update to the latest version`
)
