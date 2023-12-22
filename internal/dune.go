package internal

import (
	"encoding/xml"
	"fmt"
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

/*
http://192.168.1.148/cgi-bin/do?cmd=status

<?xml version="1.0" ?>
<command_result>
    <param name="protocol_version" value="3"/>
    <param name="player_state" value="standby"/>
    <param name="playback_pip_video_enabled" value="1"/>
    <param name="pip_alpha_level" value="255"/>
    <param name="pip_zorder" value="400"/>
    <param name="video_zorder" value="200"/>
    <param name="osd_zorder" value="500"/>
</command_result>

Remote Control IR CODES from http://dune-hd.com/support/rc/

DISCRETE-POWER-ON
http://192.168.1.148/cgi-bin/do?cmd=ir_code&ir_code=A05FBF00
DISCRETE-POWER-OFF
http://192.168.1.148/cgi-bin/do?cmd=ir_code&ir_code=A15EBF00

<?xml version="1.0" ?>
<command_result>
    <param name="protocol_version" value="3"/>
    <param name="command_status" value="ok"/>
    <param name="player_state" value="loading"/>
    <param name="playback_pip_video_enabled" value="1"/>
    <param name="pip_alpha_level" value="255"/>
    <param name="pip_zorder" value="400"/>
    <param name="video_zorder" value="200"/>
    <param name="osd_zorder" value="500"/>
</command_result>

*/

type Param struct {
	//	XMLName xml.Name `xml:"param"`
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Command struct {
	//	XMLName xml.Name `xml:"command_result"`
	Params []Param `xml:"param"`
}

func getDune(url string) (Command, error) {
	var command Command

	resp, err := http.Get(url)
	if err != nil {
		return command, err
	}

	defer resp.Body.Close()

	decoder := xml.NewDecoder(resp.Body)
	//decoder.CharsetReader = charset.NewReaderLabel

	err = decoder.Decode(&command)
	if err != nil {
		return command, err
	}

	return command, nil
}

func getDuneStatus(ip net.IP) (string, error) {
	url := fmt.Sprintf("http://%s/cgi-bin/do?cmd=status", ip)
	cmd, err := getDune(url)
	if err != nil {
		return "offline", err
	}
	state := "unknown"
	for _, p := range cmd.Params {
		if p.Name == "player_state" {
			state = p.Value
			break
		}
	}
	return state, nil
}

func getIRCodeUrl(ip net.IP, ir_code string) string {
	return fmt.Sprintf("http://%s/cgi-bin/do?cmd=ir_code&ir_code=%s", ip, ir_code)
}

func getDuneIRCode(ip net.IP, ir_code string) (bool, error) {
	url := getIRCodeUrl(ip, ir_code)
	cmd, err := getDune(url)
	if err != nil {
		return false, err
	}
	state := false
	for _, p := range cmd.Params {
		if p.Name == "command_status" {
			state = p.Value == "ok"
			break
		}
	}
	return state, nil
}

func getDuneOn(ip net.IP) (bool, error) {
	return getDuneIRCode(ip, "A05FBF00")
}

func getDuneOff(ip net.IP) (bool, error) {
	return getDuneIRCode(ip, "A15EBF00")
}

func GetDuneStatusById(conf *Config, id string) string {
	status := "unknown"
	if ip := net.ParseIP(conf.GetDuneIP(id)); ip != nil {
		var err error
		if status, err = getDuneStatus(ip); err != nil {
			log.Error().Msgf("getDuneStatus error: %v", err)
		}
	}
	return status
}

func DuneOnById(conf *Config, id string) bool {
	status := false
	if ip := net.ParseIP(conf.GetDuneIP(id)); ip != nil {
		var err error
		if status, err = getDuneOn(ip); err != nil {
			log.Error().Msgf("getDuneOn error: %v", err)
		}
	}
	return status
}

func DuneOffById(conf *Config, id string) bool {
	status := false
	if ip := net.ParseIP(conf.GetDuneIP(id)); ip != nil {
		var err error
		if status, err = getDuneOff(ip); err != nil {
			log.Error().Msgf("getDuneOn error: %v", err)
		}
	}
	return status
}
