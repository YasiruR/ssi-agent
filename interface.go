package main

import "github.com/hyperledger/aries-framework-go/pkg/client/didexchange"

type Agent interface {
	Init()
	CreateInvitation() (*didexchange.Invitation, error)
	IssueVC()
	VerifyVP()
	RevokeVC()
	Stop()
}
