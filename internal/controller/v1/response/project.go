package response

import "github.com/Taya-san/portofolio-backend/internal/entity"

type ProjectList struct {
	Projects []entity.Project `json:"projects"`
	Total    int              `json:"total"`
}
