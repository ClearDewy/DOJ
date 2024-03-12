/**
 * @ Author: ClearDewy
 * @ Desc: 前端页面配置
 **/
package config

type WebConfig struct {
	// BaseUrl string `json:"baseUrl"`
	RecordName    string `json:"recordName"`
	ProjectName   string `json:"projectName"`
	ShortName     string `json:"shortName"`
	RecordUrl     string `json:"recordUrl"`
	ProjectUrl    string `json:"projectUrl"`
	Introduction  string `json:"introduction"`
	AllowRegister bool   `json:"allowRegister"`
}

func NewInitWebConfig() *WebConfig {
	return &WebConfig{
		ProjectName:   "Dewy Online Judge",
		ShortName:     "DOJ",
		ProjectUrl:    "https://github.com/ClearDewy/DOJ",
		AllowRegister: true,
	}
}
