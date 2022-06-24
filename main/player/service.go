package player

import (
	"errors"
	"padel-backend/main/util"
)

func Create(player Player) Player {
	players = append(players, player)

	return player
}

func FindAll() []Player {
	return players
}

func FindById(ID int) (*Player, error) {
	for i, p := range players {
		if p.ID == ID {
			return &players[i], nil
		}
	}

	return nil, errors.New("player not found")
}

func Find(player Player) (*Player, error) {
	var p, err = FindById(player.ID)

	if err != nil {
		return &Player{}, errors.New("player does not exist")
	}

	return p, nil
}

func Update(player *Player, updated Player) (Player, error) {
	player.Name = updated.Name
	player.Email = updated.Email

	return *player, nil
}

func DeleteById(ID int) error {
	var _, err = FindById(ID)

	if err != nil {
		return errors.New("player does not exist")
	}

	// TODO: Remove from database
	// TODO: Remove lines below when we have a database
	for i, p := range players {
		if p.ID == ID {
			players = util.RemoveElementByIndex(players, i)
		}
	}

	return nil
}

func Delete(player *Player) error {
	return DeleteById(player.ID)
}
