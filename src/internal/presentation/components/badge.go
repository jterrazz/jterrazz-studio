package components

import "github.com/jterrazz/jterrazz-studio/src/internal/presentation/theme"

// Badge renders status badges

// BadgeOK renders a success checkmark badge
func BadgeOK() string {
	return theme.BadgeOK.Render(theme.IconCheck)
}

// BadgeError renders an error cross badge
func BadgeError() string {
	return theme.BadgeError.Render(theme.IconCross)
}

// Badge renders a badge based on condition
func Badge(ok bool) string {
	if ok {
		return BadgeOK()
	}
	return BadgeError()
}

// BadgeLoading renders a loading badge with spinner
func BadgeLoading(spinnerFrame string) string {
	return theme.SpinnerStyle.Render(spinnerFrame)
}

// ServiceBadge renders a service status badge
func ServiceBadge(running bool) string {
	if running {
		return theme.ServiceRunning.Render(theme.IconServiceOn) + " " + theme.Success.Render("running")
	}
	return theme.ServiceStopped.Render(theme.IconServiceOff) + " " + theme.Warning.Render("stopped")
}
