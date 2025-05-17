package entity

// AuthRequest представляет запрос на проверку авторизации
type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	IP       string `json:"ip"`
}

// AuthResponse представляет ответ на запрос авторизации
type AuthResponse struct {
	OK bool `json:"ok"`
}

// ResetBucketRequest представляет запрос на сброс bucket
type ResetBucketRequest struct {
	Login string `json:"login"`
	IP    string `json:"ip"`
}

// IPSubnetRequest представляет запрос для управления white/black листами
type IPSubnetRequest struct {
	Subnet string `json:"subnet"` // в формате CIDR (например, "192.168.1.0/24")
}
