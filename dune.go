package switcher

import (
	"encoding/xml"
	"fmt"
	"net/http"
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

func getDuneStatus(ip string) (string, error) {
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

func getIRCodeUrl(ip string, ir_code string) string {
	return fmt.Sprintf("http://%s/cgi-bin/do?cmd=ir_code&ir_code=%s", ip, ir_code)
}

func getDuneIRCode(ip string, ir_code string) (bool, error) {
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

func getDuneOn(ip string) (bool, error) {
	return getDuneIRCode(ip, "A05FBF00")
}

func getDuneOff(ip string) (bool, error) {
	return getDuneIRCode(ip, "A15EBF00")
}
