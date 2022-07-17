package main

import (
	"errors"
	"github.com/hyperledger/aries-framework-go/pkg/client/outofband"
	"github.com/hyperledger/aries-framework-go/pkg/client/vcwallet"
	"github.com/hyperledger/aries-framework-go/pkg/framework/context"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	walletPkg "github.com/hyperledger/aries-framework-go/pkg/wallet"
	"github.com/tryfix/log"
	"time"
)

type store struct {
	userID string
	logger log.Logger
	wallet *vcwallet.Client
}

func newStore(port string, logger log.Logger) *store {
	return &store{
		userID: `agent ` + port,
		logger: logger,
	}
}

func (s *store) init(ctx *context.Provider) {
	pass := `test-passphrase`
	err := vcwallet.CreateProfile(s.userID, ctx, walletPkg.WithPassphrase(pass))
	if err != nil {
		s.logger.Fatal(err)
	}

	wlt, err := vcwallet.New(s.userID, ctx)
	if err != nil {
		s.logger.Fatal(err)
	}
	s.wallet = wlt

	err = s.wallet.Open(walletPkg.WithUnlockByPassphrase(pass))
	if err != nil {
		s.logger.Fatal(err)
	}

	conns, err := s.wallet.GetAll(walletPkg.Connection)
	if err != nil {
		s.logger.Fatal(err)
	}

	s.logger.Info("Connections from wallet", conns)
}

func (s *store) genKeys() error {
	keys, err := s.wallet.CreateKeyPair(kms.ED25519)
	if err != nil {
		return err
	}

	s.logger.Debug("key pair generated and stored", keys.KeyID, keys.PublicKey)
	return nil
}

func (s *store) acceptInv(inv *outofband.Invitation) error {
	connID, err := s.wallet.Connect(inv, walletPkg.WithConnectTimeout(30*time.Second), walletPkg.WithMyLabel("accepting invitation"))
	if err != nil {
		return err
	}

	s.logger.Debug("out-of-band invitation accepted via wallet", connID)
	return nil
}

//func (s *store) getConn() error {
//	msg, err := s.wallet.Get("connection", s.connID)
//	if err != nil {
//		return err
//	}
//
//	s.logger.Debug("wallet connection: ", string(msg))
//	return nil
//}

func (s *store) Close() error {
	if ok := s.wallet.Close(); !ok {
		return errors.New(`terminating kms wallet failed`)
	}
	return nil
}
