package match

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
	"padel-backend/main/db"
	"padel-backend/main/util"
)

var Context *gin.Context

func InitMatchService(c *gin.Context) {
	Context = c
}

func Create(match Match) (Match, error) {
	firestore, _ := db.GetFirestore()
	defer firestore.Close()

	match.ID = uuid.New().String()
	match.Locked = false

	_, err := firestore.Collection(CollectionMatch).NewDoc().Set(Context, match)
	if err != nil {
		log.Printf("%v %v", util.GetLogPrefix("MatchService", "Create"), err)
		errors.New("could not save player, an error occurred")

		return match, err
	}

	return match, nil
}

func FindAll() (*[]Match, error) {
	firestore, _ := db.GetFirestore()
	defer firestore.Close()

	var matches []Match

	iter := firestore.Collection(CollectionMatch).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return &matches, err
		}

		var m Match

		err = doc.DataTo(&m)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "FindAll"), err)
			return &matches, errors.New("match could not be transformed to type")
		}

		matches = append(matches, m)
	}

	return &matches, nil
}

func FindById(ID string) (*Match, error) {
	var match Match
	firestore, _ := db.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionMatch).Where("ID", "==", ID).Documents(Context)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "FindById"), err)
			return &match, errors.New("player not found")
		}

		err = doc.DataTo(&match)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "FindById"), err)
			return &match, errors.New("player could not be transformed to type")
		}
	}

	return &match, nil
}

func UpdateBasicFields(match, updated *Match) (*Match, error) {
	match.Club = updated.Club
	match.Location = updated.Location
	match.Time = updated.Time

	return Update(match.ID, match)
}

func Update(ID string, updated *Match) (*Match, error) {
	firestore, _ := db.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionMatch).Where("ID", "==", ID).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "Update"), err)
			return &Match{}, errors.New("match not found")
		}

		var m Match
		err = doc.DataTo(&m)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "Update"), err)
			return &m, errors.New("match could not be transformed to type")
		}

		_, err = firestore.Collection(CollectionMatch).Doc(doc.Ref.ID).Set(Context, updated)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "Update"), err)
			return &Match{}, errors.New("could not update match, an error occurred")
		}

		return updated, nil
	}

	return &Match{}, errors.New("match not found")
}

func UpdateMatchResultById(ID string, request ResultMatchRequest) (*Match, error) {
	match, fetchErr := FindById(ID)
	if fetchErr != nil {
		return &Match{}, fetchErr
	}

	update, err := UpdateMatchResult(match, request)
	if err != nil {
		log.Printf("%v %v", util.GetLogPrefix("MatchService", "UpdateMatchResultById"), err)
		return match, errors.New("could not update match, an error occurred")
	}

	return Update(match.ID, update)
}

func UpdatePlayerToMatchById(ID string, request JoinMatchRequest) (*Match, error) {
	var match, fetchErr = FindById(ID)
	if fetchErr != nil {
		return &Match{}, fetchErr
	}

	update, err := UpdatePlayerToMatch(match, request)
	if err != nil {
		log.Printf("%v %v", util.GetLogPrefix("MatchService", "UpdatePlayerToMatchById"), err)
		return match, errors.New("could not update match, an error occurred")
	}

	return Update(ID, update)
}

func DeleteById(ID string) error {
	firestore, _ := db.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionMatch).Where("ID", "==", ID).Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "DeleteById"), err)
			return errors.New("match not found")
		}

		_, err = firestore.Collection(CollectionMatch).Doc(doc.Ref.ID).Delete(Context)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "DeleteById"), err)
			return errors.New("could not delete match, an error occurred")
		}
	}

	return nil
}
