package match

import (
	"errors"
	"padel-backend/main/player"
	"padel-backend/main/util"
)

func Create(match Match) Match {
	matches = append(matches, match)

	return match
}

func FindAll() []Match {
	return matches
}

func FindById(ID int) (*Match, error) {
	for i, m := range matches {
		if m.ID == ID {
			return &matches[i], nil
		}
	}

	return nil, errors.New("match not found")
}

func IsTeamFull(team *Team) bool {
	return team.Player1 != player.Player{} && team.Player2 != player.Player{}
}

func Find(match Match) (*Match, error) {
	var m, err = FindById(match.ID)

	if err != nil {
		return &Match{}, errors.New("match does not exist")
	}

	return m, nil
}

func Update(match *Match, updated Match) (Match, error) {
	match.Club = updated.Club
	match.Location = updated.Location
	match.Time = updated.Time

	return *match, nil
}

func UpdateAddPlayerToMatchById(id int, request JoinMatchRequest) (*Match, error) {
	var match, fetchErr = FindById(id)

	if fetchErr != nil {
		return &Match{}, fetchErr
	}

	return UpdateAddPlayerToMatch(match, request)
}

func UpdateAddPlayerToMatch(match *Match, request JoinMatchRequest) (*Match, error) {
	var err error

	var player, playerErr = player.FindById(request.PlayerId)

	if playerErr != nil {
		return &Match{}, playerErr
	}

	switch request.Team {
	case TeamA:
		if !IsTeamFull(&match.TeamA) {
			_, err = UpdateAddPlayerToTeam(&match.TeamA, *player)
		}
		break

	case TeamB:
		_, err = UpdateAddPlayerToTeam(&match.TeamB, *player)

		break

	case None:
		_, err = UpdateRemovePlayerFromTeam(&match.TeamB, *player)
	}

	if err != nil {
		return &Match{}, err
	}

	return match, nil
}

func UpdateAddPlayerToTeam(team *Team, newPlayer player.Player) (*Team, error) {
	if (team.Player1 == player.Player{}) {
		team.Player1 = newPlayer

		return team, nil
	}

	if (team.Player2 == player.Player{}) {
		team.Player2 = newPlayer

		return team, nil
	}

	return team, errors.New("could not add player to team, because al spots in team are filled")
}

func UpdateRemovePlayerFromTeam(team *Team, oldPlayer player.Player) (*Team, error) {
	if team.Player1.ID == oldPlayer.ID {
		team.Player1 = player.Player{}

		return team, nil
	}

	if team.Player2.ID == oldPlayer.ID {
		team.Player2 = player.Player{}

		return team, nil
	}

	return team, errors.New("could not remove player to team, because player does not exist in team")
}

func DeleteById(ID int) error {
	var _, err = FindById(ID)

	if err != nil {
		return errors.New("match does not exist")
	}

	// TODO: Remove from database
	// TODO: Remove lines below when we have a database
	for i, m := range matches {
		if m.ID == ID {
			matches = util.RemoveElementByIndex(matches, i)
		}
	}

	return nil
}

func Delete(match *Match) error {
	return DeleteById(match.ID)
}
