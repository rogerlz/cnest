package camera

import "errors"

type Config struct {
	Mode        string `ini:"mode"`
	Port        int    `ini:"port"`
	Device      string `ini:"device"`
	Resolution  string `ini:"resolution"`
	MaxFps      int    `ini:"max_fps"`
	CustomFlags string `ini:"custom_flags,omitempty"`
	V4l2ctl     string `ini:"v4l2ctl,omitempty"`
}

func (c Config) Validate() error {
	if c.Mode != "mjpg" && c.Mode != "rtsp" {
		return errors.New("camera mode must be mjpg or rtsp")
	}

	return nil
}
