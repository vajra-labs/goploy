package dtos

type UserResDto struct {
	ID        int64   `json:"id"        doc:"User ID"`
	Email     *string `json:"email"     doc:"User email address"`
	FirstName *string `json:"firstName" doc:"User first name"`
	LastName  *string `json:"lastName"  doc:"User last name"`
	Avatar    *string `json:"avatar"    doc:"User avatar URL"`
	IsOwner   bool    `json:"isOwner"   doc:"True if the user is owner"`
}
