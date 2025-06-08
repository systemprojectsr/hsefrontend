package api

import "core/internal"

type RegisterInfoPost struct {
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PasswordHash string `json:"password"`
	Photo        string `json:"photo"`
	Type         string `json:"type"`
}

type Client struct {
	*RegisterInfoPost `json:"register"`
}

type ClientRegister struct {
	*Client `json:"user"`
}

type GeneralAuth struct {
	*GeneralLogin `json:"user"`
}

type GeneralLogin struct {
	*GeneralLoginAttributes `json:"login"`
}

type GeneralLoginAttributes struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type ClientToken struct {
	Token string `json:"token"`
}

type TokenAccess struct {
	User struct {
		Login struct {
			Token string `json:"token"`
		} `json:"login"`
	} `json:"user"`
}

type TokenCreateCard struct {
	*TokenAccess
	Card struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"card"`
}

type TokenDeleteCard struct {
	*TokenAccess
	Card struct {
		ID int `json:"id"`
	}
}

type TokenListCard struct {
	*TokenAccess
}

type CardResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CompanyID   uint   `json:"company_id"`
}

type CardsWrapper struct {
	Cards []CardResponse `json:"cards"`
}

type UserCompanyRegister struct {
	*CompanyRegister `json:"user"`
}

type CompanyRegister struct {
	*CompanyInfoPost `json:"register"`
}

type CompanyInfoPost struct {
	CompanyName   string   `json:"company_name"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	FullName      string   `json:"full_name"`
	PositionAgent string   `json:"position_agent"`
	IDCompany     string   `json:"id_company"`
	Address       string   `json:"address"`
	TypeService   string   `json:"type_service"`
	PasswordHash  string   `json:"password_hash"`
	Photo         string   `json:"photo"`
	Documents     []string `json:"documents"`
	Type          string   `json:"type"`
}

type AccountInfo struct {
	ID       uint   `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
	Token    string `json:"token"`
	Type     string `json:"type"`
}

type CompanyInfo struct {
	ID            uint     `json:"id"`
	CompanyName   string   `json:"company_name"`
	Email         string   `json:"email"`
	Phone         string   `json:"phone"`
	FullName      string   `json:"full_name"`
	PositionAgent string   `json:"position_agent"`
	IDCompany     string   `json:"id_company"`
	Address       string   `json:"address"`
	TypeService   string   `json:"type_service"`
	PasswordHash  string   `json:"password_hash"`
	Photo         string   `json:"photo"`
	Documents     []string `json:"documents"`
	Type          string   `json:"type"`
}

type ResponseSuccessAccess struct {
	*internal.StatusResponse
	*ResponseUser `json:"user"`
}

type ResponseUser struct {
	ID    uint   `json:"id"`
	Token string `json:"token"`
	Type  string `json:"type"`
}

type ResponseAccount struct {
	*internal.StatusResponse
	User struct {
		Account AccountInfo `json:"account"`
	} `json:"user"`
}

type ResponseAccountCompany struct {
	*internal.StatusResponse
	User struct {
		Account CompanyInfo `json:"account"`
	} `json:"user"`
}
