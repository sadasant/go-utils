package utils

func Red(s string) string {
	return color_red + s + color_normal
}

func Green(s string) string {
	return color_green + s + color_normal
}

func Cyan(s string) string {
	return color_cyan + s + color_normal
}

func CyanBold(s string) string {
	return color_cyan_bold + s + color_normal
}

func Bold(s string) string {
	return color_bold + s + color_normal
}
