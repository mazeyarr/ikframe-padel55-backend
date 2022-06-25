package match

import (
	"errors"
	"padel-backend/main/player"
	"padel-backend/main/util"
)

func Create(match Match) Match {
	matches = append(matches, match)

	match.locked = false

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

func IsTeamValid(team *Team) bool {
	return IsTeamFull(team)
}

func isPlayerInTeam(player player.Player, team *Team) bool {
	return team.Player1.ID == player.ID || team.Player2.ID == player.ID
}

func isPlayerInAnyTeam(player player.Player, match *Match) (bool, *Team) {
	if isPlayerInTeam(player, &match.TeamA) {
		return true, &match.TeamA
	}

	if isPlayerInTeam(player, &match.TeamB) {
		return true, &match.TeamB
	}

	return false, &Team{}
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

func UpdateMatchResultById(id int, request ResultMatchRequest) (*Match, error) {
	var match, fetchErr = FindById(id)

	if fetchErr != nil {
		return &Match{}, fetchErr
	}

	return UpdateMatchResult(match, request)
}

func UpdateMatchResult(match *Match, request ResultMatchRequest) (*Match, error) {
	var team *Team
	var playerUpdating, playerErr = player.FindById(request.PlayerId)

	if playerErr != nil {
		return match, playerErr
	}

	if isInAnyTeam, _ := isPlayerInAnyTeam(*playerUpdating, match); !isInAnyTeam {
		return match, errors.New("player is not in team while editing results")
	}

	switch request.Team {
	case TeamA:
		team = &match.TeamA

		break

	case TeamB:
		team = &match.TeamB

		break

	case None:
	default:
		return match, errors.New("no team selected")
	}

	if !IsTeamValid(&match.TeamA) || !IsTeamValid(&match.TeamB) {
		return match, errors.New("one or more team(s) are invalid, please check all players in match")
	}

	_, err := UpdateTeamResult(team, request.TeamResult)

	if err != nil {
		return match, err
	}

	return match, nil
}

func UpdateTeamResult(team *Team, result TeamResult) (*Team, error) {
	var totalSets = len(team.Results)
	var resultSetExists = false

	if !IsTeamValid(team) {
		return team, errors.New("cannot add result, team is invalid, please check all players")
	}

	for i, r := range team.Results {
		if r.Set == result.Set {
			team.Results[i] = result

			resultSetExists = true

			break
		}
	}

	if resultSetExists {
		return team, nil
	}

	if totalSets > 3 {
		return team, errors.New("cannot add result, already 3 sets with results exist")
	}

	if result.Set < 1 || result.Set > 3 {
		return team, errors.New("cannot add result, can only add a set from 1 to 3")
	}

	team.Results = append(team.Results, result)

	return team, nil
}

func UpdatePlayerToMatchById(id int, request JoinMatchRequest) (*Match, error) {
	var match, fetchErr = FindById(id)

	if fetchErr != nil {
		return &Match{}, fetchErr
	}

	return UpdatePlayerToMatch(match, request)
}

func UpdatePlayerToMatch(match *Match, request JoinMatchRequest) (*Match, error) {
	var err error
	var teamPlayerWantsToJoin *Team
	var playerJoining, playerErr = player.FindById(request.PlayerId)

	if match.locked {
		return match, errors.New("cannot join this match anymore, it's started or locked for another reason")
	}

	if playerErr != nil {
		return match, playerErr
	}

	switch request.Team {
	case None:
		var isInAnyTeam, teamPlayerJoined = isPlayerInAnyTeam(*playerJoining, match)

		if isInAnyTeam {
			_, err = UpdateRemovePlayerFromTeam(teamPlayerJoined, *playerJoining)
		}

		if err != nil {
			return match, err
		}

		return match, nil

	case TeamA:
		teamPlayerWantsToJoin = &match.TeamA

		break

	case TeamB:
		teamPlayerWantsToJoin = &match.TeamB

		break

	default:
		return match, errors.New("no team selected")
	}

	if (teamPlayerWantsToJoin != &Team{}) {
		var isInAnyTeam, teamPlayerAlreadyJoined = isPlayerInAnyTeam(*playerJoining, match)

		if isInAnyTeam {
			if teamPlayerAlreadyJoined == teamPlayerWantsToJoin {
				return match, errors.New("player already joined this team")
			}

			if IsTeamFull(teamPlayerWantsToJoin) {
				return match, errors.New("could not join other team because it's already full")
			}

			_, err = UpdateRemovePlayerFromTeam(teamPlayerAlreadyJoined, *playerJoining)
			_, err = UpdateAddPlayerToTeam(teamPlayerWantsToJoin, *playerJoining)

			return match, nil
		}

		_, err = UpdateAddPlayerToTeam(teamPlayerWantsToJoin, *playerJoining)
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
