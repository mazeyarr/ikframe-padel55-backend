package match

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
	"log"
	"padel-backend/main/firebase"
	"padel-backend/main/util"
	"sync"
)

var Context *gin.Context

func InitMatchService(c *gin.Context) {
	Context = c
}

func Create(match Match) (Match, error) {
	firestore, _ := firebase.GetFirestore()
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
	firestore, _ := firebase.GetFirestore()
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
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	iter := firestore.Collection(CollectionMatch).Where("ID", "==", ID).Documents(Context)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "FindById"), err)
			return &match, errors.New("match not found")
		}

		err = doc.DataTo(&match)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "FindById"), err)
			return &match, errors.New("match could not be transformed to type")
		}

		return &match, nil
	}

	return &match, errors.New("match could not be found")
}

func FindPlayerMatchesByPlayerId(ID string) (*[]Match, error) {
	var wg = sync.WaitGroup{}

	wg.Add(4)

	var matches []Match

	go func() {
		ms, err := FindPlayerMatchesByPath("TeamA.Player1.ID", ID)
		if err != nil {
			log.Printf("%v %v", util.
				GetLogPrefix("MatchService", "FindPlayerMatchesByPlayerId"), err)
			wg.Done()

			return
		}

		for _, m := range *ms {
			matches = append(matches, m)
		}

		wg.Done()
	}()

	go func() {
		ms, err := FindPlayerMatchesByPath("TeamA.Player2.ID", ID)
		if err != nil {
			log.Printf("%v %v", util.
				GetLogPrefix("MatchService", "FindPlayerMatchesByPlayerId"), err)
			wg.Done()

			return
		}

		for _, m := range *ms {
			matches = append(matches, m)
		}

		wg.Done()
	}()

	go func() {
		ms, err := FindPlayerMatchesByPath("TeamB.Player1.ID", ID)
		if err != nil {
			log.Printf("%v %v", util.
				GetLogPrefix("MatchService", "FindPlayerMatchesByPlayerId"), err)
			wg.Done()

			return
		}

		for _, m := range *ms {
			matches = append(matches, m)
		}

		wg.Done()
	}()

	go func() {
		ms, err := FindPlayerMatchesByPath("TeamB.Player2.ID", ID)
		if err != nil {
			log.Printf("%v %v", util.
				GetLogPrefix("MatchService", "FindPlayerMatchesByPlayerId"), err)
			wg.Done()

			return
		}

		for _, m := range *ms {
			matches = append(matches, m)
		}

		wg.Done()
	}()

	wg.Wait()

	return &matches, nil
}

func FindPlayerMatchesByPath(path, ID string) (*[]Match, error) {
	firestore, _ := firebase.GetFirestore()
	defer firestore.Close()

	var ms []Match

	colRef := firestore.Collection(CollectionMatch)
	query := colRef.Query.Where(path, "==", ID)

	iter := query.Documents(Context)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "Update"), err)
			return &ms, errors.New("match not found")
		}

		var m Match
		err = doc.DataTo(&m)
		if err != nil {
			log.Printf("%v %v", util.GetLogPrefix("MatchService", "Update"), err)
			return &ms, errors.New("match could not be transformed to type")
		}

		ms = append(ms, m)
	}

	return &ms, nil
}

func UpdateBasicFields(match, updated *Match) (*Match, error) {
	match.Club = updated.Club
	match.Location = updated.Location
	match.Time = updated.Time

	return Update(match.ID, match)
}

func Update(ID string, updated *Match) (*Match, error) {
	firestore, _ := firebase.GetFirestore()
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
	firestore, _ := firebase.GetFirestore()
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
