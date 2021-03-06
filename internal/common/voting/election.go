package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect   baseclient.Role
	votingMethod  ElectionVotingMethod
	islandsToVote []shared.ClientID
	votes         [][]shared.ClientID
}

type ElectionVotingMethod = int

const (
	BordaCount = iota
	Plurality
	Majority
)

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role baseclient.Role, method ElectionVotingMethod) {
	e.roleToElect = role
	e.votingMethod = method
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.islandsToVote = clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for _, island := range e.islandsToVote {
		e.votes = append(e.votes, clientMap[island].GetVoteForElection(e.roleToElect))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() shared.ClientID {

	var result shared.ClientID

	switch e.votingMethod {
	case BordaCount:
		result = e.bordaCountResult()
	case Plurality:
		result = e.pluralityResult()
	case Majority:
		result = e.majorityResult()
	}
	return result
}

func (e *Election) bordaCountResult() shared.ClientID {
	// TODO implement Borda count winner selection method.
	return e.pluralityResult()
}

func (e *Election) pluralityResult() shared.ClientID {

	// How many first place votes did each island get
	votesPerIsland := map[shared.ClientID]int{}
	for _, ranking := range e.votes {
		votesPerIsland[ranking[0]] += 1
	}

	// Who got the most first place votes
	winVote := 0
	winner := shared.ClientID(1)
	for island, votes := range votesPerIsland {
		if votes >= winVote {
			winVote = votes
			winner = island
		}
	}
	return winner
}

func (e *Election) majorityResult() shared.ClientID {
	// TODO implement majority winner selection method.
	return e.pluralityResult()
}
