package public

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
)

func GetDBServer() string {
	url := GetiniValueByKey("db_server")
	if len(url) > 0 {
		return url
	} else {
		return "127.0.0.1"
	}
}

func EnableAccessControlAllowOrigin() bool {
	log := GetiniValueByKey("Access-Control-Allow-Origin")
	//Log("Access-Control-Allow-Origin====" + log)
	if log == "1" {
		return true
	} else {
		return false
	}
}

func EnableMultiplecCert() bool {
	log := GetiniValueByKey("MultiplecCert")
	if log == "1" {
		return true
	} else {
		return false
	}
}
func EnableLog() bool {
	log := GetiniValueByKey("log")
	if log == "1" {
		return true
	} else {
		return false
	}
}

func EnableFmtLog() bool {
	return true
	log := GetiniValueByKey("fmt_log")
	if log == "0" {
		return false
	} else {
		return true
	}
}
func GetAdminName() string {
	url := GetiniValueByKey("config_server_admin_name")
	if len(url) > 0 {
		return url
	} else {
		return "langan"
	}
}
func GetAdminPwd() string {
	url := GetiniValueByKey("config_server_admin_pwd")
	if len(url) > 0 {
		return url
	} else {
		return "la777888"
	}
}
func GetHttpPort() string {
	url := GetiniValueByKey("http_port")
	if len(url) > 0 {
		return url
	} else {
		return "8082"
	}
}
func GetBlockPage() string {
	url := GetiniValueByKey("block_page")
	if len(url) > 3 {
		return url
	} else {
		return "block.html"
	}
}

func GetIp() string {
	url := GetiniValueByKey("ip")
	if len(url) > 7 {
		return url
	} else {
		return ""
	}
}
func GetAppID() string {
	url := GetiniValueByKey("appid")
	if len(url) > 0 {
		return url
	} else {
		return "ocs"
	}
}
func GetRunCMD() string {
	url := GetiniValueByKey("run_cmd")
	if len(url) > 0 {
		return url
	} else {
		return ""
	}
}
func GetCMDByIndex(index string) string {
	url := GetiniValueByKey("run_cmd" + index)
	if len(url) > 0 {
		return url
	} else {
		return ""
	}
}
func GetDockerUpdateCMD() string {
	url := GetiniValueByKey("docker_update_cmd")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetDockerBuildCMD() string {
	url := GetiniValueByKey("docker_build_cmd")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetDockerPushCMD() string {
	url := GetiniValueByKey("docker_push_cmd")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetDockerLoginCMD() string {
	url := GetiniValueByKey("docker_login_cmd")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetDockerTagCMD() string {
	url := GetiniValueByKey("docker_tag_cmd")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetDeployYamlFilePath() string {
	url := GetiniValueByKey("deploy_yaml_path")
	if len(url) > 4 {
		return url
	} else {
		return ""
	}
}
func GetSSLFilePath() string {
	url := GetiniValueByKey("ssl_crt_path")
	if len(url) > 4 {
		return url
	} else {
		return "zizhuxiyiji.com"
	}
}
func GetSSLFileName() string {
	url := GetiniValueByKey("ssl_crt_name")
	if len(url) > 4 {
		return url
	} else {
		return "Cert_Chain.crt"
	}
}
func GetConfigServerKey() string {
	url := GetiniValueByKey("config_server_key")
	if len(url) > 10 {
		return url
	} else {
		return "899MA5FE12AF2DBF2DDEDCFF44B7CCSS"
	}
}
func GetConfigServerName() string {
	url := GetiniValueByKey("config_server")
	if len(url) > 10 {
		return url
	} else {
		return ""
	}
}
func GetLocalConfigServerName() string {
	url := GetiniValueByKey("config_server_parent")
	if len(url) > 10 {
		return url
	} else {
		return "http://127.0.0.1:8099"
	}
}
func GetAuthServerName() string {
	url := GetiniValueByKey("auth_server")
	if len(url) > 10 {
		return url
	} else {
		return "47.75.130.11:8097"
	}
}

var cfg *ini.File

func GetiniValueByKey(key string) string {

	if cfg == nil {

		fmt.Println("GetiniValueByKey key=" + key)
		inipath := GetCurDir() + string(os.PathSeparator) + "config.ini"
		if !ExistsPath(inipath) {
			fmt.Println("CreateFile")
			CreateFile(inipath)
		}
		var err error
		cfg, err = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, inipath)
		if err != nil {
			fmt.Println(err)
			return ""
		}
	}

	return cfg.Section("path").Key(key).String()

}
func ReadConfigIni() {
	fmt.Println("ReadConfigIni")
	inipath := GetCurDir() + string(os.PathSeparator) + "config.ini"
	if !ExistsPath(inipath) {
		CreateFile(inipath)
	}

	var errr error
	cfg, errr = ini.LoadSources(ini.LoadOptions{IgnoreInlineComment: true}, inipath)
	if errr != nil {
		Log(errr)
	}
}

func SetKeyValue(key, value string) {
	fmt.Println("SetKeyValue")
	inipath := GetCurDir() + string(os.PathSeparator) + "config.ini"
	if !ExistsPath(inipath) {
		CreateFile(inipath)
	}
	cfg.Section("path").Key(key).SetValue(value)
	cfg.SaveTo(inipath)
}
