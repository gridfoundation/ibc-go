package types_test

import (
	"time"

	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"

	clienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	"github.com/cosmos/ibc-go/v3/modules/light-clients/01-furyint/types"
	ibctesting "github.com/cosmos/ibc-go/v3/testing"
)

func (suite *FuryintTestSuite) TestGetHeight() {
	if suite.chainA.TestChainClient.GetSelfClientType() == exported.Furyint {
		header := suite.chainA.TestChainClient.(*ibctesting.TestChainFuryint).LastHeader
		suite.Require().NotEqual(uint64(0), header.GetHeight())
	} else {
		// chainB must be Furyint
		header := suite.chainB.TestChainClient.(*ibctesting.TestChainFuryint).LastHeader
		suite.Require().NotEqual(uint64(0), header.GetHeight())
	}
}

func (suite *FuryintTestSuite) TestGetTime() {
	if suite.chainA.TestChainClient.GetSelfClientType() == exported.Furyint {
		header := suite.chainA.TestChainClient.(*ibctesting.TestChainFuryint).LastHeader
		suite.Require().NotEqual(time.Time{}, header.GetTime())
	} else {
		// chainB must be Furyint
		header := suite.chainB.TestChainClient.(*ibctesting.TestChainFuryint).LastHeader
		suite.Require().NotEqual(time.Time{}, header.GetTime())
	}
}

func (suite *FuryintTestSuite) TestHeaderValidateBasic() {
	var (
		header      *types.Header
		furyintChain *ibctesting.TestChainFuryint
	)
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{"valid header", func() {}, true},
		{"header is nil", func() {
			header.Header = nil
		}, false},
		{"signed header is nil", func() {
			header.SignedHeader = nil
		}, false},
		{"SignedHeaderFromProto failed", func() {
			header.SignedHeader.Commit.Height = -1
		}, false},
		{"signed header failed furyint ValidateBasic", func() {
			header = furyintChain.LastHeader
			header.SignedHeader.Commit = nil
		}, false},
		{"trusted height is equal to header height", func() {
			header.TrustedHeight = header.GetHeight().(clienttypes.Height)
		}, false},
		{"validator set nil", func() {
			header.ValidatorSet = nil
		}, false},
		{"ValidatorSetFromProto failed", func() {
			header.ValidatorSet.Validators[0].PubKey = tmprotocrypto.PublicKey{}
		}, false},
		{"header validator hash does not equal hash of validator set", func() {
			// generated new validator set
			val := tmprototypes.Validator{}
			valSet := tmprototypes.ValidatorSet{
				Validators:       []*tmprototypes.Validator{&val},
				Proposer:         &val,
				TotalVotingPower: 0,
			}
			header.ValidatorSet = &valSet
		}, false},
	}

	suite.Require().Equal(exported.Furyint, suite.header.ClientType())

	for _, tc := range testCases {
		tc := tc

		suite.Run(tc.name, func() {
			suite.SetupTest()

			if suite.chainA.TestChainClient.GetSelfClientType() == exported.Furyint {
				furyintChain = suite.chainA.TestChainClient.(*ibctesting.TestChainFuryint)
			} else {
				// chainB must be Furyint
				furyintChain = suite.chainB.TestChainClient.(*ibctesting.TestChainFuryint)
			}

			header = furyintChain.LastHeader // must be explicitly changed in malleate

			tc.malleate()

			err := header.ValidateBasic()

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
