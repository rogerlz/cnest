package camera

import "strconv"

func GetMjpgArguments(c Config) (string, []string) {
	cmdPath := "/usr/local/bin/ustreamer"

	cmdArgs := []string{}

	// TODO noproxy
	// host and port
	cmdArgs = append(cmdArgs, "--host", "127.0.0.1", "--port", strconv.Itoa(c.Port))

	// TODO dev_is_raspicam()
	// Raspicam Workaround
	dev_is_raspicam := false
	if dev_is_raspicam {
		cmdArgs = append(cmdArgs, "--format", "MJPEG", "--device-timeout", "5", "--buffers", "3")
	} else {
		cmdArgs = append(cmdArgs, "--device", c.Device, "--device-timeout", "2")

		// TODO detect_mjpeg()
		// Use MJPEG hardware encoder if possible
		detect_mjpeg := false
		if detect_mjpeg {
			cmdArgs = append(cmdArgs, "--format", "MJPEG", "--encoder", "HW")
		}
	}

	cmdArgs = append(cmdArgs, "--resolution", c.Resolution, "--desired-fps", strconv.Itoa(c.MaxFps))
	cmdArgs = append(cmdArgs, "--allow-origin", "*", "./ustreamer-www")

	// TODO c.CustomFlags

	return cmdPath, cmdArgs
}
