package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/summa-tx/relays/golang/x/relay/types"
)

func (s *KeeperSuite) TestGetHeader() {
	// errors if header is not found
	header := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers[0]
	_, err := s.Keeper.GetHeader(s.Context, header.Hash)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())
}

func (s *KeeperSuite) TestEmitExtension() {
	// tests extension was emitted successfully
	headers := s.Fixtures.HeaderTestCases.ValidateChain[0].Headers
	s.Keeper.emitExtension(s.Context, headers[0], headers[1])

	events := s.Context.EventManager().Events()
	e := events[0]
	s.Equal("extension", e.Type)
}

func (s *KeeperSuite) TestValidateHeaderChain() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		err := validateHeaderChain(tc.Anchor, tc.Headers, tc.Internal, tc.IsMainnet)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.SDKNil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestIngestHeaders() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	// errors if anchor is not found
	err := s.Keeper.ingestHeaders(s.Context, cases[0].Headers, cases[0].Internal)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())

	for _, tc := range cases {
		s.InitTestContext(tc.IsMainnet, false)
		s.Keeper.ingestHeader(s.Context, tc.Anchor)
		err := s.Keeper.ingestHeaders(s.Context, tc.Headers, tc.Internal)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.SDKNil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestIngestHeaderChain() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		if tc.Internal == false {
			s.InitTestContext(tc.IsMainnet, false)
			s.Keeper.ingestHeader(s.Context, tc.Anchor)
			err := s.Keeper.IngestHeaderChain(s.Context, tc.Headers)
			if tc.Output == 0 {
				logIfTestCaseError(tc, err)
				s.SDKNil(err)
			} else {
				s.NotNil(err)
				s.Equal(tc.Output, err.Code())
			}
		}
	}
}

// TestIngestHeader tests ingestHeader, HasHeader, and GetHeader
func (s *KeeperSuite) TestIngestHeader() {
	cases := s.Fixtures.HeaderTestCases.ValidateChain

	for _, tc := range cases {
		s.Keeper.ingestHeader(s.Context, tc.Headers[0])
		hasHeader := s.Keeper.HasHeader(s.Context, tc.Headers[0].Hash)
		s.Equal(true, hasHeader)
		header, err := s.Keeper.GetHeader(s.Context, tc.Headers[0].Hash)
		s.SDKNil(err)
		s.Equal(tc.Headers[0], header)
	}
}

func (s *KeeperSuite) TestValidateDifficultyChange() {
	cases := s.Fixtures.HeaderTestCases.ValidateDiffChange

	for _, tc := range cases {
		err := validateDifficultyChange(tc.Headers, tc.PrevEpochStart, tc.Anchor)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.SDKNil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestIngestDifficultyChange() {
	cases := s.Fixtures.HeaderTestCases.ValidateDiffChange

	// errors if PrevEpochStart is not found
	err := s.Keeper.IngestDifficultyChange(s.Context, cases[0].PrevEpochStart.Hash, cases[0].Headers)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())

	// errors if anchor is not found
	s.Keeper.ingestHeader(s.Context, cases[0].PrevEpochStart)
	err = s.Keeper.IngestDifficultyChange(s.Context, cases[0].PrevEpochStart.Hash, cases[0].Headers)
	s.Equal(sdk.CodeType(types.UnknownBlock), err.Code())

	for _, tc := range cases {
		s.Keeper.ingestHeader(s.Context, tc.PrevEpochStart)
		s.Keeper.ingestHeader(s.Context, tc.Anchor)
		err := s.Keeper.IngestDifficultyChange(s.Context, tc.PrevEpochStart.Hash, tc.Headers)
		if tc.Output == 0 {
			logIfTestCaseError(tc, err)
			s.SDKNil(err)
		} else {
			s.NotNil(err)
			s.Equal(tc.Output, err.Code())
		}
	}
}

func (s *KeeperSuite) TestCompareTargets() {
	cases := s.Fixtures.HeaderTestCases.CompareTargets

	for _, tc := range cases {
		result := compareTargets(tc.Full, tc.Truncated)
		s.Equal(tc.Output, result)
	}
}

func (s *KeeperSuite) TestSetCurrentEpochDiff() {
	val := sdk.NewUint(1000)
	err := s.Keeper.setCurrentEpochDifficulty(s.Context, val)
	s.SDKNil(err)

	d := s.Keeper.getCurrentEpochDifficulty(s.Context)

	s.Equal(d, val)
}

func (s *KeeperSuite) TestSetPrevEpochDiff() {
	val := sdk.NewUint(1000)
	err := s.Keeper.setPrevEpochDifficulty(s.Context, val)
	s.SDKNil(err)

	d := s.Keeper.getPrevEpochDifficulty(s.Context)

	s.Equal(d, val)
}
