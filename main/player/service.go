package player

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
	auth2 "padel-backend/main/auth"
	"padel-backend/main/firebase"
	"padel-backend/main/util"
)

var Context *gin.Context

func InitPlayerService(c *gin.Context) {
	Context = c
}

func Create(player Player) (Player, error) {
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	player.ID = uuid.New().String()

	_, err := firestore.Collection(CollectionPlayer).NewDoc().Set(Context, player)
	if err != nil {
		log.Printf("%v %v", util.GetLogPrefix("PlayerService", "Create"), err)
		errors.New("could not save player, an error occurred")

		return player, err
	}

	return player, nil
}

func FindAll() ([]Player, error) {
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	var players []Player

	iter := firestore.Collection(CollectionPlayer).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return players, err
		}

		var p Player
		err = doc.DataTo(&p)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindAll"), err)
			return players, errors.New("player could not be transformed to type")
		}

		players = append(players, p)
	}

	return players, nil
}

func FindMyPlayer() (*Player, error) {
	var player Player
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()
	auth, _ := firebase.GetFirebaseAuth()

	token, err := auth.VerifyIDToken(context.Background(), Context.GetString(auth2.IdTokenName))
	if err != nil {
		return &player, errors.New("player not found, token was invalid")
	}

	user, err := auth.GetUser(context.Background(), token.UID)
	if err != nil {
		return &player, errors.New("player not found, user does not exist anymore")
	}

	iter := firestore.Collection(CollectionPlayer).Where("Email", "==", user.Email).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindById"), err)
			return &player, errors.New("player not found")
		}

		err = doc.DataTo(&player)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindById"), err)
			return &player, errors.New("player could not be transformed to type")
		}

		return &player, nil
	}

	return &player, errors.New("player could not be found")
}

func FindById(ID string) (*Player, error) {
	var player Player
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionPlayer).Where("ID", "==", ID).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindById"), err)
			return &player, errors.New("player not found")
		}

		err = doc.DataTo(&player)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindById"), err)
			return &player, errors.New("player could not be transformed to type")
		}

		return &player, nil
	}

	return &player, errors.New("player could not be found")
}

func FindByEmail(Email string) (*Player, error) {
	var player Player
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionPlayer).Where("Email", "==", Email).Limit(1).Documents(Context)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindByEmail"), err)
			return &player, errors.New("player not found")
		}

		err = doc.DataTo(&player)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "FindByEmail"), err)
			return &player, errors.New("player could not be transformed to type")
		}
	}

	return &player, nil
}

func Update(player *Player, updated Player) (*Player, error) {
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionPlayer).Where("ID", "==", player.ID).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "Update"), err)
			return player, errors.New("player not found")
		}

		var p Player
		err = doc.DataTo(&p)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "Update"), err)
			return &p, errors.New("player could not be transformed to type")
		}

		updated.ID = p.ID

		_, err = firestore.Collection(CollectionPlayer).Doc(doc.Ref.ID).Set(Context, updated)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "Update"), err)
			return player, errors.New("could not update player, an error occurred")
		}

		return &updated, nil
	}

	return player, errors.New("player not found")
}

func DeleteById(ID string) error {
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionPlayer).Where("ID", "==", ID).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "DeleteById"), err)
			return errors.New("player not found")
		}

		_, err = firestore.Collection(CollectionPlayer).Doc(doc.Ref.ID).Delete(Context)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("PlayerService", "DeleteById"), err)
			return errors.New("could not delete player, an error occurred")
		}
	}

	return nil
}
