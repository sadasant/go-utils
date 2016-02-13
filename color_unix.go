// +build darwin freebsd linux netbsd openbsd

package utils

// To clear this for terminals that doesn't support colors
// we could set them up on init based on the environment's TERM variable.
// I'm not sure which TERMS allow colors and which doesn't though!
const color_normal = "\x1b[0m"
const color_red = "\x1b[31m"
const color_green = "\x1b[32m"
const color_blue = "\x1b[34m"
const color_cyan = "\x1b[36m"
const color_bold = "\x1b[1m"
const color_cyan_bold = "\x1b[36;1m"
