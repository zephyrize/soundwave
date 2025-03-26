package server

// ServiceRegisterRequest 服务注册请求结构
type ServiceRegisterRequest struct {
	Name     string            `json:"name" binding:"required"`
	ID       string            `json:"id" binding:"required"`
	Hostname string            `json:"hostname" binding:"required"`
	IP       string            `json:"ip" binding:"required,ip"`
	Port     int               `json:"port" binding:"required,gt=0,lte=65535"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
}
